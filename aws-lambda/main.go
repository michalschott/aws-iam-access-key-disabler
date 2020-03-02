package main

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	now               time.Time
	lowThresholdDays  int
	highThresholdDays int
	dryRun            int
}

func parseKey(svc *iam.IAM, config Config, key *iam.AccessKeyMetadata) {
	disableKey := false
	keyAgeInDays := config.now.Sub(*key.CreateDate).Hours() / 24
	if int(keyAgeInDays) >= config.highThresholdDays {
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
			"Threshold":   strconv.Itoa(config.highThresholdDays) + " days",
		}).Info("Access key age above high threshold. Disabling it.")
		disableKey = true
	} else if int(keyAgeInDays) < config.highThresholdDays && int(keyAgeInDays) >= config.lowThresholdDays {
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
			"Threshold":   strconv.Itoa(config.lowThresholdDays) + " days",
		}).Info("Access key age above low threshold.")
	}

	if disableKey {
		if config.dryRun == 0 {
			_, err := svc.UpdateAccessKey(&iam.UpdateAccessKeyInput{
				AccessKeyId: aws.String(*key.AccessKeyId),
				Status:      aws.String(iam.StatusTypeInactive),
				UserName:    aws.String(*key.UserName),
			})

			if err != nil {
				log.Fatal("Error", err)
			}
		}
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
		}).Debug("Access key has been disabled.")
	}
}

func parseUser(svc *iam.IAM, username string, config Config, wg *sync.WaitGroup) {
	defer wg.Done()
	result, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(username),
	})

	if err != nil {
		log.Fatal("Error", err)
	}

	for _, key := range result.AccessKeyMetadata {
		if *key.Status == "Active" {
			parseKey(svc, config, key)
		} else {
			log.WithFields(log.Fields{
				"Username":    username,
				"AccessKeyId": *key.AccessKeyId,
			}).Debug("Key inactive, skipping.")
		}
	}
}

func HandleRequest(ctx context.Context) (int, error) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	//log.SetLevel(log.DebugLevel)

	var lowThresholdDays = LookupEnvOrInt("LOWTHRESHOLDDAYS", 90)
	var highThreshooldDays = LookupEnvOrInt("HIGHTHRESHOLDDAYS", 180)
	var dryRun = LookupEnvOrInt("DRYRUN", 1)
	var whiteList = LookupEnvOrString("WHITELIST", "")

	config := Config{time.Now().UTC(), lowThresholdDays, highThreshooldDays, dryRun}
	whitelistedUsers := strings.Split(whiteList, ",")

	sort.Strings(whitelistedUsers)
	log.Info("Ignored users: ", whitelistedUsers)

	if dryRun == 1 {
		log.Info("Running in dry run mode.")
	}

	sess, err := session.NewSession()

	if err != nil {
		log.Fatal("Error", err)
	}

	svc := iam.New(sess)

	result, err := svc.ListUsers(&iam.ListUsersInput{})

	if err != nil {
		log.Fatal("Error", err)
	}

	var wg sync.WaitGroup
	for _, user := range result.Users {
		if user == nil {
			continue
		}

		// skip whitelistedUsers
		i := sort.Search(len(whitelistedUsers), func(i int) bool { return whitelistedUsers[i] >= *user.UserName })

		if i < len(whitelistedUsers) && whitelistedUsers[i] == *user.UserName {
			continue
		}

		wg.Add(1)
		go parseUser(svc, *user.UserName, config, &wg)
	}
	wg.Wait()

	return 0, nil
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func main() {
	lambda.Start(HandleRequest)
}
