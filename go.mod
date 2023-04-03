module github.com/michalschott/aws-iam-access-key-disabler

go 1.19

require (
	github.com/aws/aws-lambda-go v1.39.1
	github.com/aws/aws-sdk-go v1.44.234
	github.com/michalschott/aws-iam-access-key-disabler/pkg/env v0.1.0
	github.com/michalschott/aws-iam-access-key-disabler/pkg/iam v0.1.0
	github.com/michalschott/aws-iam-access-key-disabler/pkg/lambda v0.1.0
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace (
	github.com/michalschott/aws-iam-access-key-disabler/pkg/env => ./pkg/env
	github.com/michalschott/aws-iam-access-key-disabler/pkg/iam => ./pkg/iam
	github.com/michalschott/aws-iam-access-key-disabler/pkg/lambda => ./pkg/lambda
)
