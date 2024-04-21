package app

import (
	"fmt"
	"net/http"

	"github.com/nixpig/dunce/db"
	"github.com/nixpig/dunce/internal/tags"
)

func Start(port string) {
	mux := http.NewServeMux()

	tagsData := tags.NewTagData(db.DB.Conn)
	tagService := tags.NewTagService(tagsData)
	tagsController := tags.NewTagController(tagService)

	mux.HandleFunc("POST /tags", tagsController.CreateHandler)
	mux.HandleFunc("GET /tags", tagsController.GetAllHandler)
	mux.HandleFunc("GET /tags/{slug}", tagsController.GetBySlugHandler)
	mux.HandleFunc("POST /tags/{slug}", tagsController.UpdateHandler)
	mux.HandleFunc("DELETE /tags", tagsController.DeleteHandler)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
