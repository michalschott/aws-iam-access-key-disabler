package main

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	log "github.com/sirupsen/logrus"

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
	msgKeyToBeDisabled     = "days beyond high threshold, key to be disabled"
	msgKeyToBeDisabledSoon = "days beyond low threshold, key will be disabled soon"
)

func HandleRequest(ctx context.Context) (int, error) {
	var lowThresholdDays = kingpin.Flag("lowTresholdDays", "Show warning for keys older than X days").Default("90").Envar("LOWTHRESHOLDDAYS").Short('l').Int()
	var highThresholdDays = kingpin.Flag("highTresholdDays", "Disable keys older than X days").Default("180").Envar("HIGHTHRESHOLDDAYS").Short('h').Int()
	var dryRun = kingpin.Flag("dryrun", "").Default("true").Envar("DRYRUN").Bool()
	var debug = kingpin.Flag("debug", "").Default("0").Envar("DEBUG").Int()
	var whiteList = kingpin.Flag("whitelist", "Comma separated list of users to skip").Default("").Envar("WHITELIST").Short('w').String()
	kingpin.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if *debug == 1 {
		log.SetLevel(log.DebugLevel)
	}

	config := Config{time.Now().UTC(), *lowThresholdDays, *highThresholdDays, *dryRun}
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
