package env

import (
	"os"
	"strconv"
)

func GetENVString(key, defaultvalue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultvalue
}

func GetENVIntegers(key string, defaultvalue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}

	}

	return defaultvalue
}
