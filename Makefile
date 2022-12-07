.DEFAULT_GOAL=lambda

lambda: aws-iam-access-key-disabler

aws-iam-access-key-disabler: main.go
	GOOS=linux GOARCH=amd64 go build .
