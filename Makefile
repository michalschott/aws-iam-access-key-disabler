.DEFAULT_GOAL=lambda

lambda: aws-iam-access-key-disabler

aws-iam-access-key-disabler: *.go
	GOOS=linux GOARCH=amd64 go build ./...
