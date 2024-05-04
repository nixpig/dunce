package main

import (
	"log"
	"os"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nixpig/dunce/db"
	app "github.com/nixpig/dunce/internal/app"
	"github.com/nixpig/dunce/internal/config"
	"github.com/nixpig/dunce/pkg"
)

func main() {
	appConfig := app.AppConfig{}

	if err := config.Init(); err != nil {
		log.Printf("unable to load config from env '%v' which may be fatal; continuing...", err)
	}

	if err := db.MigrateUp(); err != nil {
		log.Printf("did not run database migration due to '%v' which may be fatal; continuing...", err)
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	appConfig.Db = db

	appConfig.SessionManager = scs.New()
	appConfig.SessionManager.Store = pgxstore.New(appConfig.Db.Pool.(*pgxpool.Pool))

	validate, err := pkg.NewValidator()
	if err != nil {
		log.Fatalf("unable to create validation: %v", err)
		os.Exit(1)
	}

	appConfig.Validator = validate

	templateCache, err := pkg.NewTemplateCache()
	if err != nil {
		log.Fatalf("unable to build template cache: %v", err)
		os.Exit(1)
	}

	appConfig.TemplateCache = templateCache

	appConfig.Logger = pkg.NewLogger()

	appConfig.Port = config.Get("WEB_PORT")

	if err := app.Start(appConfig); err != nil {
		log.Fatalf("unable to start app: %v", err)
		os.Exit(1)
	}
}
