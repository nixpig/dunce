package app

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/home"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/internal/user"
	"github.com/nixpig/dunce/pkg"
)

type AppConfig struct {
	Port           string
	Validator      *validator.Validate
	Db             *db.Dbpool
	TemplateCache  map[string]*template.Template
	Logger         pkg.Logger
	SessionManager *scs.SessionManager
}

func Start(appConfig AppConfig) error {
	mux := http.NewServeMux()

	controllerConfig := pkg.ControllerConfig{
		TemplateCache:  appConfig.TemplateCache,
		Log:            appConfig.Logger,
		SessionManager: appConfig.SessionManager,
	}

	static := http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", static))

	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/login", http.StatusPermanentRedirect)
	})

	userRepo := user.NewUserRepository(appConfig.Db.Pool, appConfig.Logger)
	userService := user.NewUserService(userRepo, appConfig.Validator, appConfig.Logger)
	userController := user.NewUserController(userService, controllerConfig)

	mux.HandleFunc("GET /admin/login", pkg.NoSurf(userController.UserLoginGet))
	mux.HandleFunc("POST /admin/login", pkg.NoSurf(userController.UserLoginPost))
	mux.HandleFunc("POST /admin/logout", pkg.NoSurf(userController.UserLogoutPost))
	mux.HandleFunc("GET /admin/users/new", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, userController.CreateUserGet))))
	mux.HandleFunc("GET /admin/users/{slug}", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, userController.UserGet))))
	mux.HandleFunc("GET /admin/users", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, userController.UsersGet))))
	mux.HandleFunc("POST /admin/users", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, userController.CreateUserPost))))
	mux.HandleFunc("POST /admin/users/{username}/delete", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, userController.DeleteUserPost))))

	tagRepository := tag.NewTagRepository(appConfig.Db.Pool, appConfig.Logger)
	tagService := tag.NewTagService(tagRepository, appConfig.Validator, appConfig.Logger)
	tagController := tag.NewTagController(tagService, controllerConfig)

	mux.HandleFunc("POST /admin/tags", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(tagController.PostAdminTagsHandler)))
	mux.HandleFunc("GET /admin/tags", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, tagController.GetAdminTagsHandler))))
	mux.HandleFunc("GET /admin/tags/new", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, tagController.GetAdminTagsNewHandler))))
	mux.HandleFunc("GET /admin/tags/{slug}", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, tagController.GetAdminTagsSlugHandler))))
	mux.HandleFunc("POST /admin/tags/{slug}", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, tagController.PostAdminTagsSlugHandler))))
	mux.HandleFunc("POST /admin/tags/{slug}/delete", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, tagController.DeleteAdminTagsSlugHandler))))

	articleRepository := article.NewArticleRepository(appConfig.Db.Pool, appConfig.Logger)
	articleService := article.NewArticleService(articleRepository, appConfig.Validator, appConfig.Logger)
	articleController := article.NewArticleController(articleService, tagService, controllerConfig)

	mux.HandleFunc("POST /admin/articles", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, articleController.CreateHandler))))
	mux.HandleFunc("GET /admin/articles", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, articleController.GetAllHandler))))
	mux.HandleFunc("GET /admin/articles/new", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, articleController.NewHandler))))
	mux.HandleFunc("GET /admin/articles/{slug}", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, articleController.GetBySlugHander))))
	mux.HandleFunc("POST /admin/articles/{slug}", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, articleController.UpdateHandler))))
	mux.HandleFunc("POST /admin/articles/{slug}/delete", userService.IsAuthenticatedMiddleware(appConfig.SessionManager, pkg.NoSurf(pkg.Protected(appConfig.SessionManager, articleController.AdminArticlesDeleteHandler))))

	api := home.NewApi(articleService, tagService)

	homeController := home.NewHomeController(api, controllerConfig)
	mux.HandleFunc("GET /", homeController.HomeGet)

	mux.HandleFunc("GET /articles/{slug}", articleController.PublicGetArticle)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", appConfig.Port),
		Handler:      appConfig.SessionManager.LoadAndSave(mux),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	appConfig.Logger.Info("starting server on %s", appConfig.Port)

	if err := server.ListenAndServe(); err != nil {
		appConfig.Logger.Error("failed to start server: %s", err)
	}

	return nil
}
