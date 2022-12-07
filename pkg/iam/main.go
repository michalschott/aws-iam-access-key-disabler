package iam

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Returns list of active Access Keys metadata objects
func GetActiveAccessKeys(svc *iam.IAM, username string) ([]*iam.AccessKeyMetadata, error) {
	result, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(username),
	})

	if err != nil {
		return nil, err
	}

	returnValue := []*iam.AccessKeyMetadata{}
	for _, key := range result.AccessKeyMetadata {
		if *key.Status == "Active" {
			returnValue = append(returnValue, key)
		}
	}

	return returnValue, nil
}

// Disabled given Access Key assigned to AccessKeyMetadata
func DisableAccessKey(svc *iam.IAM, key *iam.AccessKeyMetadata) error {
	_, err := svc.UpdateAccessKey(&iam.UpdateAccessKeyInput{
		AccessKeyId: aws.String(*key.AccessKeyId),
		Status:      aws.String(iam.StatusTypeInactive),
		UserName:    aws.String(*key.UserName),
	})

	return err
}

// Returns age of Access Key in days
func GetAccessKeyAgeInDays(key *iam.AccessKeyMetadata) int {
	now := time.Now().UTC()

	return int(now.Sub(*key.CreateDate).Hours() / 24)
}
