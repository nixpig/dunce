package tags

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/nixpig/dunce/pkg/logging"
)

type TagsController struct {
	service       TagServiceInterface
	log           logging.Logger
	templateCache map[string]*template.Template
}

func NewTagController(
	service TagServiceInterface,
	log logging.Logger,
	templateCache map[string]*template.Template,
) TagsController {
	return TagsController{
		service:       service,
		log:           log,
		templateCache: templateCache,
	}
}

func (tc *TagsController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	tag := NewTag(r.FormValue("name"), r.FormValue("slug"))

	if _, err := tc.service.Create(&tag); err != nil {
		tc.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (tc *TagsController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		tc.log.Error(err.Error())
		w.Write([]byte(err.Error()))
	}

	if err := tc.service.DeleteById(id); err != nil {
		tc.log.Error(err.Error())
		w.Write([]byte(err.Error()))
	}
}

func (tc *TagsController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := tc.service.GetAll()
	if err != nil {
		tc.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tc.templateCache["tags.tmpl"].ExecuteTemplate(w, "base", tags)
	if err != nil {
		tc.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (tc *TagsController) GetBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := tc.service.GetBySlug(slug)
	if err != nil {
		tc.log.Error(err.Error())
		w.Write([]byte(err.Error()))
	}

	if err := tc.templateCache["tag.tmpl"].ExecuteTemplate(w, "base", tag); err != nil {
		tc.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (tc *TagsController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		tc.log.Error(err.Error())
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
	}

	tag := NewTagWithId(
		id,
		r.FormValue("name"),
		r.FormValue("slug"),
	)
	if _, err := tc.service.Update(&tag); err != nil {
		tc.log.Error(err.Error())
		http.Error(w, "Unable to save changes", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (tc *TagsController) NewHandler(w http.ResponseWriter, r *http.Request) {
	if err := tc.templateCache["new-tag.tmpl"].ExecuteTemplate(w, "base", nil); err != nil {
		tc.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
