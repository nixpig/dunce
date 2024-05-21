package main

import (
	"log"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/justinas/nosurf"
	"github.com/nixpig/dunce/db"
	app "github.com/nixpig/dunce/internal/app"
	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
	"github.com/nixpig/dunce/pkg/validation"
)

func main() {
	var err error

	appConfig := app.AppConfig{}

	if err := godotenv.Load(".env"); err != nil {
		log.Printf(
			"unable to load config from env due to '%v' which may not be fatal; continuing...",
			err,
		)
	}

	if err := db.MigrateUp(); err != nil {
		log.Printf(
			"did not run database migration due to '%v' which may not be fatal; continuing...",
			err,
		)
	}

	appConfig.Db, err = db.Connect()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	appConfig.Validator, err = validation.NewValidator()
	if err != nil {
		log.Fatalf("unable to create validator: %v", err)
		os.Exit(1)
	}

	appConfig.TemplateCache, err = templates.NewTemplateCache()
	if err != nil {
		log.Fatalf("unable to build template cache: %v", err)
		os.Exit(1)
	}

	appConfig.SessionManager = session.NewSessionManagerImpl(scs.New())

	appConfig.Logger = logging.NewLogger()

	appConfig.CsrfToken = nosurf.Token

	appConfig.ErrorHandlers = errors.NewErrorHandlersImpl(appConfig.TemplateCache)

	appConfig.Port = os.Getenv("WEB_PORT")

	if err := app.Start(appConfig); err != nil {
		log.Fatalf("unable to start app: %v", err)
		os.Exit(1)
	}
}
