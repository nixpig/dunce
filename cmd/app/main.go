package main

import (
	"log"
	"os"

	"github.com/nixpig/bloggor/internal/app/server"
	"github.com/nixpig/bloggor/internal/pkg/config"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func main() {

	if err := config.Init(); err != nil {
		log.Printf("unable to load config from env '%v' which may not be fatal; continuing...", err)
	}

	if err := models.MigrateUp(); err != nil {
		log.Printf("did not run database migration due to '%v' which may not be fatal; continuing...", err)
	}

	if err := models.Connect(); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	port := config.Get("WEB_PORT")

	app.Start(port)
}
