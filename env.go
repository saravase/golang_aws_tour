package main

import (
	"os"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	log.Info("Environment variables loaded successfully")
}
