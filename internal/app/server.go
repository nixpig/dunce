package app

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

type AppConfig struct {
	Port          string
	Validator     *validator.Validate
	Db            db.Dbconn
	TemplateCache map[string]*template.Template
	Logger        pkg.Logger
}

func Start(appConfig AppConfig) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {})

	tagRepository := tag.NewTagRepository(appConfig.Db, appConfig.Logger)
	tagService := tag.NewTagService(tagRepository, appConfig.Validator, appConfig.Logger)
	tagController := tag.NewTagController(tagService, appConfig.Logger, appConfig.TemplateCache)

	mux.HandleFunc("POST /admin/tags", tagController.PostAdminTagsHandler)
	mux.HandleFunc("GET /admin/tags", tagController.GetAdminTagsHandler)
	mux.HandleFunc("GET /admin/tags/new", tagController.GetAdminTagsNewHandler)
	mux.HandleFunc("GET /admin/tags/{slug}", tagController.GetAdminTagsSlugHandler)
	mux.HandleFunc("POST /admin/tags/{slug}", tagController.PostAdminTagsSlugHandler)
	mux.HandleFunc("DELETE /admin/tags/{slug}", tagController.DeleteAdminTagsSlugHandler)

	articleRepository := article.NewArticleRepository(appConfig.Db, appConfig.Logger)
	articleService := article.NewArticleService(articleRepository, appConfig.Validator, appConfig.Logger)
	articleController := article.NewArticleController(articleService, tagService, appConfig.Logger, appConfig.TemplateCache)

	// mux.HandleFunc("POST /admin/articles", articlesController.CreateHandler)
	mux.HandleFunc("GET /admin/articles", articleController.GetAllHandler)
	mux.HandleFunc("GET /admin/articles/new", articleController.NewHandler)
	mux.HandleFunc("GET /admin/articles/{slug}", articleController.GetBySlugHander)
	// mux.HandleFunc("POST /admin/articles/{slug}", articlesController.UpdateHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", appConfig.Port),
		Handler: mux,
	}

	appConfig.Logger.Info("starting server on %s", appConfig.Port)

	if err := server.ListenAndServe(); err != nil {
		appConfig.Logger.Error("failed to start server: %s", err)
	}

	return nil
}
