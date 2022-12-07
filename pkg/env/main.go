package env

import (
	"os"
	"strconv"
)

// Returns environment value if set, otherwise returns defined default string value
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// Returns environment value if set, otherwise returns defined default int value
func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, _ := strconv.Atoi(val)
		return v
	}
	return defaultVal
}

// Returns environment value if set, otherwise returns defined default bool value
func LookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		boolValue, _ := strconv.ParseBool(val)
		return boolValue
	}
	return defaultVal
}
