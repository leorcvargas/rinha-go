package env

import (
	"errors"
	"fmt"
	"os"
)

func GetEnvOrDie(key string) string {
	value := os.Getenv(key)

	if value == "" {
		err := errors.New(fmt.Sprintf("Missing environment variable %s", key))
		panic(err)
	}

	return value
}
