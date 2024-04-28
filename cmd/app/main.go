package main

import (
	"log"
	"os"

	"github.com/nixpig/dunce/db"
	app "github.com/nixpig/dunce/internal/app/server"
	"github.com/nixpig/dunce/internal/config"
	"github.com/nixpig/dunce/pkg/templates"
	"github.com/nixpig/dunce/pkg/validation"
)

func main() {
	appConfig := app.AppConfig{}

	if err := config.Init(); err != nil {
		log.Printf("unable to load config from env '%v' which may be fatal; continuing...", err)
	}

	if err := db.MigrateUp(); err != nil {
		log.Printf("did not run database migration due to '%v' which may be fatal; continuing...", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	appConfig.Db = db.DB.Conn

	validate, err := validation.NewValidator()
	if err != nil {
		log.Fatalf("unable to create validation: %v", err)
		os.Exit(1)
	}

	appConfig.Validator = validate

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		log.Fatalf("unable to build template cache: %v", err)
		os.Exit(1)
	}

	appConfig.TemplateCache = templateCache

	appConfig.Port = config.Get("WEB_PORT")

	if err := app.Start(appConfig); err != nil {
		log.Fatalf("unable to start app: %v", err)
		os.Exit(1)
	}
}
