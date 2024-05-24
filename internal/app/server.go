package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/home"
	"github.com/nixpig/dunce/internal/site"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/internal/user"
	"github.com/nixpig/dunce/pkg/crypto"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/middleware"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
	"golang.org/x/crypto/bcrypt"
)

type AppConfig struct {
	Port           string
	Validator      *validator.Validate
	Db             *db.Dbpool
	TemplateCache  templates.TemplateCache
	Logger         logging.Logger
	SessionManager session.SessionManager
	CsrfToken      func(*http.Request) string
	ErrorHandlers  errors.ErrorHandlers
}

func Start(appConfig AppConfig) error {
	mux := http.NewServeMux()

	crypt := crypto.NewCryptoImpl(
		bcrypt.GenerateFromPassword,
		bcrypt.CompareHashAndPassword,
	)

	siteRepo := site.NewSitePostgresRepository(appConfig.Db.Pool)
	siteService := site.NewSiteService(siteRepo)
	siteController := site.NewSiteController(siteService, site.SiteControllerConfig{
		Log:            appConfig.Logger,
		TemplateCache:  appConfig.TemplateCache,
		SessionManager: appConfig.SessionManager,
		CsrfToken:      appConfig.CsrfToken,
		ErrorHandlers:  appConfig.ErrorHandlers,
	})

	userRepo := user.NewUserPostgresRepository(appConfig.Db.Pool)
	userService := user.NewUserService(userRepo, appConfig.Validator, crypt)
	userController := user.NewUserController(userService, user.UserControllerConfig{
		Log:            appConfig.Logger,
		TemplateCache:  appConfig.TemplateCache,
		SessionManager: appConfig.SessionManager,
		CsrfToken:      appConfig.CsrfToken,
		ErrorHandlers:  appConfig.ErrorHandlers,
	})

	tagRepository := tag.NewTagPostgresRepository(appConfig.Db.Pool)
	tagService := tag.NewTagService(tagRepository, appConfig.Validator)
	tagController := tag.NewTagController(tagService, tag.TagControllerConfig{
		Log:            appConfig.Logger,
		TemplateCache:  appConfig.TemplateCache,
		SessionManager: appConfig.SessionManager,
		CsrfToken:      appConfig.CsrfToken,
		ErrorHandlers:  appConfig.ErrorHandlers,
	})

	articleRepository := article.NewArticlePostgresRepository(appConfig.Db.Pool)
	articleService := article.NewArticleService(articleRepository, appConfig.Validator)
	articleController := article.NewArticleController(
		articleService,
		tagService,
		article.ArticleControllerConfig{
			Log:            appConfig.Logger,
			TemplateCache:  appConfig.TemplateCache,
			SessionManager: appConfig.SessionManager,
			CsrfToken:      appConfig.CsrfToken,
			ErrorHandlers:  appConfig.ErrorHandlers,
		},
	)

	isAuthenticated := middleware.NewAuthenticatedMiddleware(userService, appConfig.SessionManager, session.LOGGED_IN_USERNAME)
	protected := middleware.NewProtectedMiddleware(appConfig.SessionManager)
	noSurf := middleware.NewNoSurfMiddleware()
	stripSlash := middleware.NewStripSlashMiddleware()

	static := http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", static))

	mux.HandleFunc("GET /admin", adminRootHandler(appConfig))

	mux.HandleFunc("GET /admin/login", applyMiddlewares(
		userController.UserLoginGet,
		noSurf,
	))
	mux.HandleFunc("POST /admin/login", applyMiddlewares(
		userController.UserLoginPost,
		noSurf,
	))
	mux.HandleFunc("POST /admin/logout", applyMiddlewares(
		userController.UserLogoutPost,
		noSurf,
	))
	mux.HandleFunc("GET /admin/users/new", applyMiddlewares(
		userController.CreateUserGet,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/users/{slug}", applyMiddlewares(
		userController.UserGet,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/users", applyMiddlewares(
		userController.UsersGet,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("POST /admin/users", applyMiddlewares(
		userController.CreateUserPost,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("POST /admin/users/{username}/delete", applyMiddlewares(
		userController.DeleteUserPost,
		protected,
		noSurf,
		isAuthenticated,
	))

	mux.HandleFunc("POST /admin/tags", applyMiddlewares(
		tagController.PostAdminTagsHandler,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/tags", applyMiddlewares(
		tagController.GetAdminTagsHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/tags/new", applyMiddlewares(
		tagController.GetAdminTagsNewHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/tags/{slug}", applyMiddlewares(
		tagController.GetAdminTagsSlugHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("POST /admin/tags/{slug}", applyMiddlewares(
		tagController.PostAdminTagsSlugHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("POST /admin/tags/{slug}/delete", applyMiddlewares(
		tagController.DeleteAdminTagsSlugHandler,
		protected,
		noSurf,
		isAuthenticated,
	))

	mux.HandleFunc("POST /admin/articles", applyMiddlewares(
		articleController.CreateHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/articles", applyMiddlewares(
		articleController.GetAllHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/articles/new", applyMiddlewares(
		articleController.NewHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("GET /admin/articles/{slug}", applyMiddlewares(
		articleController.GetBySlugHander,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("POST /admin/articles/{slug}", applyMiddlewares(
		articleController.UpdateHandler,
		protected,
		noSurf,
		isAuthenticated,
	))
	mux.HandleFunc("POST /admin/articles/{slug}/delete", applyMiddlewares(
		articleController.AdminArticlesDeleteHandler,
		protected,
		noSurf,
		isAuthenticated,
	))

	mux.HandleFunc("GET /admin/site", applyMiddlewares(
		siteController.GetCreateSiteItems,
		protected,
		noSurf,
		isAuthenticated,
	))

	homeController := home.NewHomeController(
		tagService,
		articleService,
		home.HomeControllerConfig{
			Log:            appConfig.Logger,
			TemplateCache:  appConfig.TemplateCache,
			SessionManager: appConfig.SessionManager,
			CsrfToken:      appConfig.CsrfToken,
			ErrorHandlers:  appConfig.ErrorHandlers,
		},
	)

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

func applyMiddlewares(handler http.HandlerFunc, middlewares ...func(next http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	h := handler
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
