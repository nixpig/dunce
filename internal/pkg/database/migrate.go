package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp() error {
	env, err := loadEnv()
	if err != nil {
		return err
	}

	connectionString := buildConnectionString(env)

	log.Print("creating database migration")
	m, err := migrate.New("file://db/migrations", connectionString)
	if err != nil {
		return err
	}

	log.Print("running database migration")
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
