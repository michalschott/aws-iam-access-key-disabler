.DEFAULT_GOAL=lambda

lambda: aws-iam-access-key-disabler

aws-iam-access-key-disabler: *.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o aws-iam-access-key-disabler ./...

.PHONY: clean
clean:
	rm -f aws-iam-access-key-disabler
