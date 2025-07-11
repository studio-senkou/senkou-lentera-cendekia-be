package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	environmentError := godotenv.Load(".env")

	if environmentError != nil {
		panic("Error loading .env file: " + environmentError.Error())
	}
}

func GetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
