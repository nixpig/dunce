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

	mux.HandleFunc("GET /api/tags", tagsController.GetAllHandler)
	mux.HandleFunc("GET /api/tags/{id}", tagsController.GetByIdHandler)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
