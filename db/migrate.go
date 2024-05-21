package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MigrateUp() error {
	env, err := loadDatabaseEnvironment()
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

type Dbpool struct {
	Pool Dbconn
}

type Dbconn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type databaseEnvironment struct {
	host     string
	port     uint16
	name     string
	username string
	password string
}

func Connect() (*Dbpool, error) {
	env, err := loadDatabaseEnvironment()
	if err != nil {
		return nil, err
	}

	connectionString := buildConnectionString(env)

	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	return &Dbpool{Pool: pool}, nil
}

func loadDatabaseEnvironment() (*databaseEnvironment, error) {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	name := os.Getenv("POSTGRES_DB")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

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
