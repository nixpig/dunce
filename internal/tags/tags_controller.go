package tags

import (
	"fmt"
	"net/http"
	"strconv"
)

type TagsController struct {
	service TagServiceInterface
}

func NewTagController(service TagServiceInterface) TagsController {
	return TagsController{service}
}

func (tc *TagsController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	tag := NewTag(r.FormValue("name"), r.FormValue("slug"))

	createdTag, err := tc.service.create(&tag)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d: %s (%s)", createdTag.Id, createdTag.Name, createdTag.Slug)))
	}
}

func (tc *TagsController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	if err := tc.service.deleteById(id); err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (tc *TagsController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	allTags, err := tc.service.getAll()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(fmt.Sprintf("%v", allTags)))
}

func (tc *TagsController) GetBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := tc.service.getBySlug(slug)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(fmt.Sprintf("%v", tag)))
}

func (tc *TagsController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tag := NewTagWithId(
		id,
		r.FormValue("name"),
		r.FormValue("slug"),
	)
	updatedTag, err := tc.service.update(&tag)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(fmt.Sprintf("%v", updatedTag)))
}
