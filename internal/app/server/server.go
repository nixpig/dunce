package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/nixpig/bloggor/internal/app/handlers"
	"github.com/nixpig/bloggor/internal/pkg/config"
)

func Start(port string) {
	engine := html.New("./web/templates/", ".tmpl")

	env := config.Get("APP_ENV")

	if env == "development" {
		engine.Reload(true)
		engine.Debug(true)
	}

	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: handlers.ErrorHandler,
	})

	app.Static("/static", "./web/static")
	app.Static("robots.txt", "./web/robots.txt")

	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(compress.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return config.Get("APP_ENV") == "development"
		},
	}))

	web := app.Group(fmt.Sprintf("/%s", "/"))

	web.Get("/", handlers.IndexHandler)

	web.Get("/user/:id", handlers.UserHandler)
	web.Get("/user", handlers.UserHandler)
	// web.Get("/user/register", handlers.UserRegisterHandler)
	// web.Get("/user/login", handlers.UserLoginHandler)
	//
	// web.Get("/tag", handlers.TagHandler)
	// web.Get("/tag/:name", handlers.TagHandler)
	//
	// web.Get("/:article_type", handlers.ArticleHandler)
	// web.Get("/:article_type/:article_slug", handlers.ArticleHandler)
	//
	// web.Use(handlers.NotFoundHandler)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
