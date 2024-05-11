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
	tagService TagService
	log        pkg.Logger
	templates  map[string]*template.Template
	session    *scs.SessionManager
}

type TagView struct {
	Message         string
	Tag             *TagResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

type TagsView struct {
	Message         string
	Tags            *[]TagResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

func NewTagController(
	tagService TagService,
	config struct {
		Log            pkg.Logger
		TemplateCache  map[string]*template.Template
		SessionManager *scs.SessionManager
	},
) TagController {
	return TagController{
		tagService: tagService,
		log:        config.Log,
		templates:  config.TemplateCache,
		session:    config.SessionManager,
	}
}

func (t *TagController) PostAdminTagsHandler(w http.ResponseWriter, r *http.Request) {
	tag := CreateTagRequestDto{
		Name: r.FormValue("name"),
		Slug: r.FormValue("slug"),
	}

	if _, err := t.tagService.Create(&tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.session.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Created tag '%s'.", tag.Name))

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) DeleteAdminTagsSlugHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := t.tagService.DeleteById(id); err != nil {
		t.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	t.session.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Deleted tag '%s'.", r.FormValue("name")))

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := t.tagService.GetAll()
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := t.session.PopString(r.Context(), pkg.SESSION_KEY_MESSAGE)

	tagView := TagsView{
		Message:         message,
		Tags:            tags,
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: t.session.Exists(r.Context(), string(pkg.IsLoggedInContextKey)),
	}

	err = t.templates["pages/admin/admin-tags.tmpl"].ExecuteTemplate(w, "admin", tagView)
	if err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (t *TagController) GetAdminTagsSlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := t.tagService.GetByAttribute("slug", slug)
	if err != nil {
		t.log.Error(err.Error())
		w.Write([]byte(err.Error()))
	}

	tagView := TagView{
		Tag:             tag,
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: t.session.Exists(r.Context(), string(pkg.IsLoggedInContextKey)),
	}

	if err := t.templates["pages/admin/admin-tag.tmpl"].ExecuteTemplate(w, "admin", tagView); err != nil {
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

	tag := UpdateTagRequestDto{
		Id:   id,
		Name: r.FormValue("name"),
		Slug: r.FormValue("slug"),
	}

	if _, err := t.tagService.Update(&tag); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Unable to save changes", http.StatusInternalServerError)
		return
	}

	t.session.Put(r.Context(), pkg.SESSION_KEY_MESSAGE, fmt.Sprintf("Updated tag '%s'.", tag.Name))

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsNewHandler(w http.ResponseWriter, r *http.Request) {
	tagView := TagsView{
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: t.session.Exists(r.Context(), string(pkg.IsLoggedInContextKey)),
	}

	if err := t.templates["pages/admin/admin-new-tag.tmpl"].ExecuteTemplate(w, "admin", tagView); err != nil {
		t.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
