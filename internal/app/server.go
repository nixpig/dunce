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

	controllerConfig := pkg.NewControllerConfig(
		appConfig.Logger,
		appConfig.TemplateCache,
		appConfig.SessionManager,
	)

	userRepo := user.NewUserRepository(appConfig.Db.Pool)
	userService := user.NewUserService(userRepo, appConfig.Validator)
	userController := user.NewUserController(userService, controllerConfig)

	tagRepository := tag.NewTagPostgresRepository(appConfig.Db.Pool)
	tagService := tag.NewTagService(tagRepository, appConfig.Validator)
	tagController := tag.NewTagController(tagService, controllerConfig)

	articleRepository := article.NewArticlePostgresRepository(appConfig.Db.Pool)
	articleService := article.NewArticleService(articleRepository, appConfig.Validator)
	articleController := article.NewArticleController(articleService, tagService, appConfig.SessionManager, controllerConfig)

	isAuthenticated := user.NewIsAuthenticatedMiddleware(userService, appConfig.SessionManager)
	protected := pkg.NewProtectedMiddleware(appConfig.SessionManager)
	noSurf := pkg.NewNoSurfMiddleware()

	static := http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", static))

	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		if appConfig.SessionManager.Exists(r.Context(), string(pkg.IsLoggedInContextKey)) {
			http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		}
	})

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

	mux.HandleFunc("GET /", stripTrailingSlashMiddleware(rootHandler(homeController.HomeGet)))

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

func stripTrailingSlashMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isRoot := r.URL.Path == "/" || len(r.URL.Path) == 0

		if isRoot {
			next(w, r)
			return
		}

		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
			http.Redirect(w, r, r.URL.Path, 301)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rootHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.Error(w, "Not Found", 404)
			return
		}

		next(w, r)
	}
}
