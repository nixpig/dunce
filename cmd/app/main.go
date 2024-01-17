package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nixpig/bloggor/internal/pkg/config"
	"github.com/nixpig/bloggor/internal/pkg/database"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func main() {

	if err := config.Init(); err != nil {
		log.Printf("unable to load config from env '%v' which may not be fatal; continuing...", err)
	}

	if err := database.MigrateUp(); err != nil {
		log.Printf("did not run database migration due to '%v' which may not be fatal; continuing...", err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	newUser := models.NewUser{
		Username: "testname",
		Email:    "test@example.com",
		Link:     "some link",
		Password: "something suprt secret in hree!!!",
		Role:     models.ReaderRole,
	}

	database.DB = &db

	u1, err := models.CreateUser(*database.DB, &newUser)
	if err != nil {
		log.Printf("error: %v", err)
	}
	fmt.Println(u1)

	users, err := models.GetUserById(*database.DB, 3)
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Printf("users: %v", users)

	fmt.Println(("Hello, world!"))
}
