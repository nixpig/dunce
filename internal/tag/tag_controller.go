package tag

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"github.com/nixpig/dunce/pkg"
)

type TagController struct {
	service        pkg.Service[Tag, TagData]
	log            pkg.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func NewTagController(
	service pkg.Service[Tag, TagData],
	config pkg.ControllerConfig,
) TagController {
	return TagController{
		service:        service,
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
	}
}

func (t *TagController) PostAdminTagsHandler(w http.ResponseWriter, r *http.Request) {
	tag := TagData{
		Name: r.FormValue("name"),
		Slug: r.FormValue("slug"),
	}

	if _, err := t.service.Create(&tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Created tag '%s'.", tag.Name))

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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	t.sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Deleted tag '%s'.", r.FormValue("name")))

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := t.service.GetAll()
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := t.sessionManager.PopString(r.Context(), pkg.SESSION_KEY_MESSAGE)

	type tagTemplateView struct {
		Message   string
		Tags      *[]Tag
		CsrfToken string
	}

	data := tagTemplateView{
		Message:   message,
		Tags:      tags,
		CsrfToken: nosurf.Token(r),
	}

	err = t.templateCache["tags.tmpl"].ExecuteTemplate(w, "base", data)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (t *TagController) GetAdminTagsSlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := t.service.GetByAttribute("slug", slug)
	if err != nil {
		t.log.Error(err.Error())
		w.Write([]byte(err.Error()))
	}

	if err := t.templateCache["tag.tmpl"].ExecuteTemplate(w, "base", struct {
		Tag       *Tag
		CsrfToken string
	}{
		Tag:       tag,
		CsrfToken: nosurf.Token(r),
	}); err != nil {
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

	tag := Tag{
		Id: id,
		TagData: TagData{
			Name: r.FormValue("name"),
			Slug: r.FormValue("slug"),
		},
	}
	if _, err := t.service.Update(&tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Unable to save changes", http.StatusInternalServerError)
		return
	}

	t.sessionManager.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Updated tag '%s'.", tag.Name))

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsNewHandler(w http.ResponseWriter, r *http.Request) {
	if err := t.templateCache["new-tag.tmpl"].ExecuteTemplate(w, "base", struct {
		CsrfToken string
	}{
		CsrfToken: nosurf.Token(r),
	}); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
