package app

import (
	"fmt"
	"net/http"

	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/articles"
	"github.com/nixpig/dunce/internal/tags"
)

func Start(port string) {
	mux := http.NewServeMux()

	tagsData := tags.NewTagData(db.DB.Conn)
	tagService := tags.NewTagService(tagsData)
	tagsController := tags.NewTagController(tagService)

	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /admin/tags", tagsController.CreateHandler)
	mux.HandleFunc("GET /admin/tags", tagsController.GetAllHandler)
	mux.HandleFunc("GET /admin/tags/new", tagsController.NewHandler)
	mux.HandleFunc("GET /admin/tags/{slug}", tagsController.GetBySlugHandler)
	mux.HandleFunc("POST /admin/tags/{slug}", tagsController.UpdateHandler)
	mux.HandleFunc("DELETE /admin/tags", tagsController.DeleteHandler)

	articlesData := articles.NewArticleData(db.DB.Conn)
	articlesService := articles.NewArticleService(articlesData)
	articlesController := articles.NewArticleController(articlesService, tagService)

	mux.HandleFunc("GET /admin/articles", articlesController.GetAllHandler)
	mux.HandleFunc("GET /admin/articles/new", articlesController.NewHandler)
	mux.HandleFunc("POST /admin/articles", articlesController.CreateHandler)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
