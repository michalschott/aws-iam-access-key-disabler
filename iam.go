package main

import (
	"github.com/aws/aws-sdk-go/service/iam"
	i "github.com/michalschott/aws-iam-access-key-disabler/pkg/iam"
	log "github.com/sirupsen/logrus"
)

// check key age, compare with thresholds and disable if crossed high threshold
func parseAndDisableKey(svc *iam.IAM, config Config, key *iam.AccessKeyMetadata) error {
	keyAgeInDays := i.GetAccessKeyAgeInDays(key)
	if int(keyAgeInDays) >= config.highThresholdDays {
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
		}).Info(keyAgeInDays-config.highThresholdDays, " ", msgKeyToBeDisabled)
		if !config.dryRun {
			if err := i.DisableAccessKey(svc, key); err != nil {
				return err
			}

			log.WithFields(log.Fields{
				"Username":    *key.UserName,
				"AccessKeyId": *key.AccessKeyId,
			}).Info(msgKeyDisabled)
		}
	} else if int(keyAgeInDays) < config.highThresholdDays && int(keyAgeInDays) >= config.lowThresholdDays {
		log.WithFields(log.Fields{
			"Username":    *key.UserName,
			"AccessKeyId": *key.AccessKeyId,
		}).Info(keyAgeInDays-config.lowThresholdDays, " ", msgKeyToBeDisabledSoon, " - ", config.highThresholdDays-keyAgeInDays, " days left")
	}

	return nil
}

// gets active access keys for given users and calls parseAndDisableKey() for each key
func parseUser(svc *iam.IAM, username string, config Config) error {
	keys, err := i.GetActiveAccessKeys(svc, username)
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = parseAndDisableKey(svc, config, key)
		if err != nil {
			return err
		}
	}

	return nil
}
