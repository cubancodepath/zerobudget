package env

import (
	"os"

	"github.com/joho/godotenv"
)

func GetString(key string, fallback string) string {
	_ = godotenv.Load()

	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}
