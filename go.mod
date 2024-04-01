module github.com/michalschott/aws-iam-access-key-disabler

go 1.20

require (
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/aws/aws-lambda-go v1.46.0
	github.com/aws/aws-sdk-go v1.51.11
	github.com/michalschott/aws-iam-access-key-disabler/pkg/iam v0.1.0
	github.com/michalschott/aws-iam-access-key-disabler/pkg/lambda v0.1.0
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace (
	github.com/michalschott/aws-iam-access-key-disabler/pkg/env => ./pkg/env
	github.com/michalschott/aws-iam-access-key-disabler/pkg/iam => ./pkg/iam
	github.com/michalschott/aws-iam-access-key-disabler/pkg/lambda => ./pkg/lambda
)
