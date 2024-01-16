package config

import (
	"github.com/joho/godotenv"
	"os"
)

var envFile string = ".env"

func Init() error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}

	return nil
}

func Get(key string) string {
	return os.Getenv(key)
}
