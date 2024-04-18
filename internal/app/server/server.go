package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/nixpig/dunce/internal/app/handlers"
	"github.com/nixpig/dunce/internal/app/middleware"
	"github.com/nixpig/dunce/internal/pkg/config"
	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/nixpig/dunce/internal/pkg/user"
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

	cookieKey := encryptcookie.GenerateKey()

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: cookieKey,
	}))

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
	admin := app.Group("/admin")

	m := middleware.Protected{}
	admin.Use(m.New(middleware.ProtectedConfig{}))

	web.Get("/", handlers.IndexHandler)

	web.Get("/users/:id", handlers.UserGetHandler)
	web.Get("/users", handlers.UserGetHandler)

	// login
	web.Get("/login", handlers.AdminLoginGetHandler)
	web.Post("/login", handlers.AdminLoginPostHandler)
	web.Get("/logout", handlers.AdminLogoutHandler)

	userData := user.NewUserData(models.DB.Conn)
	userService := user.NewUserService(&userData)
	userController := user.NewUserController(&userService)

	// admin -> users
	admin.Get("/users", userController.HandleGetAll)
	admin.Post("/users", userController.HandleSave)
	admin.Get("/users/:id", handlers.AdminUserGetHandler)
	admin.Put("/users/:id", handlers.AdminUserPutHandler)
	admin.Delete("/users/:id", handlers.AdminUserDeleteHander)
	// web.Get("/users/register", handlers.UserRegisterHandler)
	// web.Get("/users/login", handlers.UserLoginHandler)
	// web.Get("/users/logout", handlers.UserLogoutHandler)

	// admin -> tags
	admin.Get("/tags", handlers.AdminTagGetHandler)
	admin.Post("/tags", handlers.AdminTagPostHandler)
	admin.Get("/tags/:id", handlers.AdminTagGetHandler)
	admin.Put("/tags/:id", handlers.AdminTagUpdateHandler)
	admin.Delete("/tags/:id", handlers.AdminTagDeleteHandler)

	// admin -> types
	admin.Get("/types", handlers.AdminTypeGetHandler)
	admin.Post("/types", handlers.AdminTypePostHander)
	admin.Get("/types/:id", handlers.AdminTypeGetHandler)
	admin.Put("/types/:id", handlers.AdminTypePutHandler)
	admin.Delete("/types/:id", handlers.AdminTypeDeleteHandler)

	// admin -> articles
	admin.Get("/articles", handlers.AdminArticleGetHandler)
	admin.Post("/articles", handlers.AdminArticlePostHandler)
	admin.Get("/articles/:id", handlers.AdminArticleGetHandler)
	admin.Put("/articles/:id", handlers.AdminArticlePutHandler)
	admin.Delete("/articles/:id", handlers.AdminArticleDeleteHandler)
	admin.Get("/create", handlers.AdminArticleCreateGetHandler)

	// admin -> site
	admin.Get("/site", handlers.AdminSiteGetHandler)
	admin.Post("/site", handlers.AdminSitePostHandler)

	// web.Get("/:article_type", handlers.ArticleHandler)
	// web.Get("/:article_type/:article_slug", handlers.ArticleHandler)
	//
	// web.Use(handlers.NotFoundHandler)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
