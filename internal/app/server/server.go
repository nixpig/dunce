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

	web.Get("/users/:id", handlers.UserGetHandler)
	web.Get("/users", handlers.UserGetHandler)

	web.Get("/admin/users", handlers.AdminUserGetHandler)
	web.Post("/admin/users", handlers.AdminUserPostHandler)
	web.Put("/admin/users", handlers.AdminUserPutHandler)
	web.Delete("/admin/users", handlers.AdminUserDeleteHander)
	// web.Get("/users/register", handlers.UserRegisterHandler)
	// web.Get("/users/login", handlers.UserLoginHandler)

	web.Get("/admin/tags", handlers.AdminTagGetHandler)
	web.Post("/admin/tags", handlers.AdminTagPostHandler)

	web.Get("/admin/types", handlers.AdminTypeGetHandler)
	web.Post("/admin/types", handlers.AdminTypePostHander)

	web.Get("/admin/articles", handlers.AdminArticleGetHandler)
	web.Post("/admin/articles", handlers.AdminArticlePostHandler)

	// web.Get("/tags", handlers.TagHandler)
	// web.Get("/tags/:name", handlers.TagHandler)
	//
	// web.Get("/:article_type", handlers.ArticleHandler)
	// web.Get("/:article_type/:article_slug", handlers.ArticleHandler)
	//
	// web.Use(handlers.NotFoundHandler)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
