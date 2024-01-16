package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
	// _ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/nixpig/bloggor/internal/pkg/config"
)

type databaseEnvironment struct {
	host     string
	port     uint16
	name     string
	username string
	password string
}

func Connect() error {
	env, err := loadEnv()
	if err != nil {
		return err
	}

	connectionString := buildConnectionString(env)

	DB, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	defer DB.Close()

	return nil

}

func loadEnv() (*databaseEnvironment, error) {
	host := config.Get("DATABASE_HOST")
	port := config.Get("DATABASE_PORT")
	name := config.Get("POSTGRES_DB")
	username := config.Get("POSTGRES_USER")
	password := config.Get("POSTGRES_PASSWORD")

	portNumber, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		log.Fatalf("unable to get valid port number from: %s", port)
		return nil, err
	}

	return &databaseEnvironment{
		host:     host,
		port:     uint16(portNumber),
		name:     name,
		username: username,
		password: password,
	}, nil

}

func buildConnectionString(env *databaseEnvironment) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.username, env.password, env.host, env.port, env.name,
	)
}
