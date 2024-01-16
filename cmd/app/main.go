package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nixpig/bloggor/internal/pkg/config"
	"github.com/nixpig/bloggor/internal/pkg/database"
)

func main() {

	if err := config.Init(); err != nil {
		log.Printf("unable to load config from env '%v' which may not be fatal; continuing...", err)
	}

	if err := database.MigrateUp(); err != nil {
		log.Printf("did not run database migration due to '%v' which may not be fatal; continuing...", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	fmt.Println(("Hello, world!"))
}
