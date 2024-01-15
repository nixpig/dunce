package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var envFile string = ".env"

func Init() error {
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("unable to load env file (%s)\nthis may not be a problem if environment variables are already declared by the environment", envFile)
		return err
	}

	return nil
}

func Get(key string) string {
	return os.Getenv(key)
}
