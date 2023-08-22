package env

import (
	"fmt"
	"os"
)

func GetEnvOrDie(key string) string {
	value := os.Getenv(key)

	if value == "" {
		err := fmt.Errorf("Missing environment variable %s", key)
		panic(err)
	}

	return value
}
