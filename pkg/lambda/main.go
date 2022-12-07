package lambda

import "os"

// Check if we're running in AWS Lambda environment
func IsLambda() bool {
	_, env1 := os.LookupEnv("LAMBDA_TASK_ROOT")
	_, env2 := os.LookupEnv("AWS_EXECUTION_ENV")

	return (env1 && env2) || false
}
