package tags

import (
	"fmt"
	"net/http"
)

type TagsController struct {
	service TagServiceInterface
}

func NewTagController(service TagServiceInterface) TagsController {
	return TagsController{service}
}

func (tc *TagsController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tags get all"))
}

func (tc *TagsController) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.Write([]byte(fmt.Sprintf("tags get by id: %s", id)))
}
