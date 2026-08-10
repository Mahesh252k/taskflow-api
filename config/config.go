package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
