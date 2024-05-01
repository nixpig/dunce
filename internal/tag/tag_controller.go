package tag

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/nixpig/dunce/pkg"
)

type TagController struct {
	service       pkg.Service[Tag]
	log           pkg.Logger
	templateCache map[string]*template.Template
}

func NewTagController(
	service pkg.Service[Tag],
	log pkg.Logger,
	templateCache map[string]*template.Template,
) TagController {
	return TagController{
		service:       service,
		log:           log,
		templateCache: templateCache,
	}
}

func (t *TagController) PostAdminTagsHandler(w http.ResponseWriter, r *http.Request) {
	tag := NewTag(r.FormValue("name"), r.FormValue("slug"))

	if _, err := t.service.Create(&tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) DeleteAdminTagsSlugHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := t.service.DeleteById(id); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := t.service.GetAll()
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.templateCache["tags.tmpl"].ExecuteTemplate(w, "base", tags)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (t *TagController) GetAdminTagsSlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := t.service.GetBySlug(slug)
	if err != nil {
		t.log.Error(err.Error())
		w.Write([]byte(err.Error()))
	}

	if err := t.templateCache["tag.tmpl"].ExecuteTemplate(w, "base", tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (t *TagController) PostAdminTagsSlugHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
	}

	tag := NewTagWithId(
		id,
		r.FormValue("name"),
		r.FormValue("slug"),
	)
	if _, err := t.service.Update(&tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Unable to save changes", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsNewHandler(w http.ResponseWriter, r *http.Request) {
	if err := t.templateCache["new-tag.tmpl"].ExecuteTemplate(w, "base", nil); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
