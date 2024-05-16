package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/home"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/internal/user"
	"github.com/nixpig/dunce/pkg"
	"github.com/nixpig/dunce/pkg/middleware"
)

type AppConfig struct {
	Port           string
	Validator      *validator.Validate
	Db             *db.Dbpool
	TemplateCache  map[string]pkg.Template
	Logger         pkg.Logger
	SessionManager pkg.SessionManager
	CsrfToken      func(*http.Request) string
}

func Start(appConfig AppConfig) error {
	mux := http.NewServeMux()

	controllerConfig := pkg.NewControllerConfig(
		appConfig.Logger,
		appConfig.TemplateCache,
		appConfig.SessionManager,
		appConfig.CsrfToken,
	)

	crypto := pkg.NewCryptoImpl()
	userRepo := user.NewUserPostgresRepository(appConfig.Db.Pool)
	userService := user.NewUserService(userRepo, appConfig.Validator, crypto)
	userController := user.NewUserController(userService, controllerConfig)

	tagRepository := tag.NewTagPostgresRepository(appConfig.Db.Pool)
	tagService := tag.NewTagService(tagRepository, appConfig.Validator)
	tagController := tag.NewTagController(tagService, controllerConfig)

	articleRepository := article.NewArticlePostgresRepository(appConfig.Db.Pool)
	articleService := article.NewArticleService(articleRepository, appConfig.Validator)
	articleController := article.NewArticleController(articleService, tagService, controllerConfig)

	isAuthenticated := middleware.NewAuthenticatedMiddleware(userService, appConfig.SessionManager)
	protected := middleware.NewProtectedMiddleware(appConfig.SessionManager)
	noSurf := middleware.NewNoSurfMiddleware()
	stripSlash := middleware.NewStripSlashMiddleware()

	static := http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", static))

	mux.HandleFunc("GET /admin", adminRootHandler(appConfig))
	mux.HandleFunc("GET /admin/login", noSurf(userController.UserLoginGet))
	mux.HandleFunc("POST /admin/login", noSurf(userController.UserLoginPost))
	mux.HandleFunc("POST /admin/logout", noSurf(userController.UserLogoutPost))
	mux.HandleFunc("GET /admin/users/new", isAuthenticated(noSurf(protected(userController.CreateUserGet))))
	mux.HandleFunc("GET /admin/users/{slug}", isAuthenticated(noSurf(protected(userController.UserGet))))
	mux.HandleFunc("GET /admin/users", isAuthenticated(noSurf(protected(userController.UsersGet))))
	mux.HandleFunc("POST /admin/users", isAuthenticated(noSurf(protected(userController.CreateUserPost))))
	mux.HandleFunc("POST /admin/users/{username}/delete", isAuthenticated(noSurf(protected(userController.DeleteUserPost))))

	mux.HandleFunc("POST /admin/tags", isAuthenticated(noSurf(tagController.PostAdminTagsHandler)))
	mux.HandleFunc("GET /admin/tags", isAuthenticated(noSurf(protected(tagController.GetAdminTagsHandler))))
	mux.HandleFunc("GET /admin/tags/new", isAuthenticated(noSurf(protected(tagController.GetAdminTagsNewHandler))))
	mux.HandleFunc("GET /admin/tags/{slug}", isAuthenticated(noSurf(protected(tagController.GetAdminTagsSlugHandler))))
	mux.HandleFunc("POST /admin/tags/{slug}", isAuthenticated(noSurf(protected(tagController.PostAdminTagsSlugHandler))))
	mux.HandleFunc("POST /admin/tags/{slug}/delete", isAuthenticated(noSurf(protected(tagController.DeleteAdminTagsSlugHandler))))

	mux.HandleFunc("POST /admin/articles", isAuthenticated(noSurf(protected(articleController.CreateHandler))))
	mux.HandleFunc("GET /admin/articles", isAuthenticated(noSurf(protected(articleController.GetAllHandler))))
	mux.HandleFunc("GET /admin/articles/new", isAuthenticated(noSurf(protected(articleController.NewHandler))))
	mux.HandleFunc("GET /admin/articles/{slug}", isAuthenticated(noSurf(protected(articleController.GetBySlugHander))))
	mux.HandleFunc("POST /admin/articles/{slug}", isAuthenticated(noSurf(protected(articleController.UpdateHandler))))
	mux.HandleFunc("POST /admin/articles/{slug}/delete", isAuthenticated(noSurf(protected(articleController.AdminArticlesDeleteHandler))))

	homeController := home.NewHomeController(tagService, articleService, controllerConfig)

	mux.HandleFunc("GET /articles", homeController.HomeArticlesGet)
	mux.HandleFunc("GET /articles/{slug}", articleController.PublicGetArticle)
	mux.HandleFunc("GET /tags", homeController.HomeTagsGet)
	mux.HandleFunc("GET /tags/{slug}", homeController.HomeTagGet)

	mux.HandleFunc("GET /", stripSlash(publicRootHandler(homeController.HomeGet)))

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

func publicRootHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.Error(w, "Not Found", 404)
			return
		}

		next(w, r)
	}
}

func adminRootHandler(appConfig AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appConfig.SessionManager.Exists(r.Context(), string(pkg.IS_LOGGED_IN_CONTEXT_KEY)) {
			http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		}
	}
}
