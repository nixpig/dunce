package app

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/articles"
	"github.com/nixpig/dunce/internal/tags"
	"github.com/nixpig/dunce/pkg/logging"
)

type AppConfig struct {
	Port          string
	Validator     *validator.Validate
	Db            db.Dbconn
	TemplateCache map[string]*template.Template
}

func Start(appConfig AppConfig) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {})

	loggers := logging.NewLogger()

	tagsData := tags.NewTagData(appConfig.Db, loggers)
	tagService := tags.NewTagService(tagsData, appConfig.Validator, loggers)
	tagsController := tags.NewTagController(tagService, loggers, appConfig.TemplateCache)

	mux.HandleFunc("POST /admin/tags", tagsController.CreateHandler)
	mux.HandleFunc("GET /admin/tags", tagsController.GetAllHandler)
	mux.HandleFunc("GET /admin/tags/new", tagsController.NewHandler)
	mux.HandleFunc("GET /admin/tags/{slug}", tagsController.GetBySlugHandler)
	mux.HandleFunc("POST /admin/tags/{slug}", tagsController.UpdateHandler)
	mux.HandleFunc("DELETE /admin/tags", tagsController.DeleteHandler)

	articlesData := articles.NewArticleData(appConfig.Db, loggers)
	articlesService := articles.NewArticleService(articlesData, appConfig.Validator, loggers)
	articlesController := articles.NewArticleController(articlesService, tagService, loggers)

	mux.HandleFunc("POST /admin/articles", articlesController.CreateHandler)
	mux.HandleFunc("GET /admin/articles", articlesController.GetAllHandler)
	mux.HandleFunc("GET /admin/articles/new", articlesController.NewHandler)
	mux.HandleFunc("GET /admin/articles/{slug}", articlesController.GetBySlugHander)
	mux.HandleFunc("POST /admin/articles/{slug}", articlesController.UpdateHandler)

	server := &http.Server{
		Addr:     fmt.Sprintf(":%v", appConfig.Port),
		ErrorLog: loggers.ErrorLogger,
		Handler:  mux,
	}

	loggers.Info("starting server on %s", appConfig.Port)

	if err := server.ListenAndServe(); err != nil {
		loggers.Error("failed to start server: %s", err)
	}

	return nil
}
