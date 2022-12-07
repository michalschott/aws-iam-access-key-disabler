package main

import (
	"context"
	"flag"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	log "github.com/sirupsen/logrus"

	e "github.com/michalschott/aws-iam-access-key-disabler/pkg/env"
	i "github.com/michalschott/aws-iam-access-key-disabler/pkg/iam"
	l "github.com/michalschott/aws-iam-access-key-disabler/pkg/lambda"
)

type Config struct {
	now               time.Time
	lowThresholdDays  int
	highThresholdDays int
	dryRun            bool
}

const (
	msgKeyDisabled         = "key disabled"
	msgKeyToBeDisabled     = "key to be disabled"
	msgKeyToBeDisabledSoon = "days left to disable key"
)

func parseKey(svc *iam.IAM, config Config, key *iam.AccessKeyMetadata) error {
	disableKey := false
	keyAgeInDays := i.GetAccessKeyAgeInDays(key)
	if int(keyAgeInDays) >= config.highThresholdDays {
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
		}).Info(msgKeyToBeDisabled)
		disableKey = true
	} else if int(keyAgeInDays) < config.highThresholdDays && int(keyAgeInDays) >= config.lowThresholdDays {
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
		}).Info(keyAgeInDays-config.lowThresholdDays, " ", msgKeyToBeDisabledSoon)
	}

	if disableKey && !config.dryRun {
		err := i.DisableAccessKey(svc, key)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
		}).Info(msgKeyDisabled)
	}

	return nil
}

func parseUser(svc *iam.IAM, username string, config Config) error {
	keys, err := i.GetActiveAccessKeys(svc, username)
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = parseKey(svc, config, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleRequest(ctx context.Context) (int, error) {
	var lowThresholdDays = flag.Int("l", e.LookupEnvOrInt("LOWTHRESHOLDDAYS", 90), "Show warning for keys older than X days")
	var highThreshooldDays = flag.Int("h", e.LookupEnvOrInt("HIGHTHRESHOLDDAYS", 180), "Disable keys older than X days")
	var dryRun = flag.Bool("dryrun", e.LookupEnvOrBool("DRYRUN", true), "Dry run")
	var debug = flag.Int("debug", e.LookupEnvOrInt("DEBUG", 0), "Debug")
	var whiteList = flag.String("w", e.LookupEnvOrString("WHITELIST", ""), "Comma separated list of users to skip.")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if *debug == 1 {
		log.SetLevel(log.DebugLevel)
	}

	config := Config{time.Now().UTC(), *lowThresholdDays, *highThreshooldDays, *dryRun}
	whitelistedUsers := strings.Split(*whiteList, ",")

	sort.Strings(whitelistedUsers)
	log.Info("Ignored users: ", whitelistedUsers)

	if *dryRun {
		log.Info("Running in dry run mode.")
	}

	sess, err := session.NewSession()

	if err != nil {
		return 1, err
	}

	svc := iam.New(sess)

	result, err := svc.ListUsers(&iam.ListUsersInput{})

	if err != nil {
		return 1, err
	}

	for _, user := range result.Users {
		if user == nil {
			continue
		}

		// skip whitelistedUsers
		i := sort.Search(len(whitelistedUsers), func(i int) bool { return whitelistedUsers[i] >= *user.UserName })

		if i < len(whitelistedUsers) && whitelistedUsers[i] == *user.UserName {
			continue
		}

		err = parseUser(svc, *user.UserName, config)
		if err != nil {
			return 1, err
		}
	}

	return 0, nil
}

func main() {
	if l.IsLambda() {
		lambda.Start(HandleRequest)
	} else {
		var ret int
		var err error

		ret, err = HandleRequest(context.Background())
		if err != nil {
			log.Error(err)
		}

		os.Exit(ret)
	}
}
