package tag

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
)

type TagController struct {
	tagService    TagService
	log           logging.Logger
	templates     templates.TemplateCache
	session       session.SessionManager
	csrfToken     func(r *http.Request) string
	errorHandlers errors.ErrorHandlers
}

type TagControllerConfig struct {
	Log            logging.Logger
	TemplateCache  templates.TemplateCache
	SessionManager session.SessionManager
	CsrfToken      func(*http.Request) string
	ErrorHandlers  errors.ErrorHandlers
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

type TagCreateView struct {
	CsrfToken       string
	IsAuthenticated bool
}

func NewTagController(
	tagService TagService,
	config TagControllerConfig,
) TagController {
	return TagController{
		tagService:    tagService,
		log:           config.Log,
		templates:     config.TemplateCache,
		session:       config.SessionManager,
		csrfToken:     config.CsrfToken,
		errorHandlers: config.ErrorHandlers,
	}
}

func (t *TagController) PostAdminTagsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	tag := TagNewRequestDto{
		Name: r.FormValue("name"),
		Slug: r.FormValue("slug"),
	}

	if _, err := t.tagService.Create(&tag); err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}

	t.session.Put(
		r.Context(),
		session.SESSION_KEY_MESSAGE,
		fmt.Sprintf("Created tag '%s'.", tag.Name),
	)

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) DeleteAdminTagsSlugHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		t.errorHandlers.BadRequest(w, r)
		return
	}

	if err := t.tagService.DeleteById(id); err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}

	t.session.Put(
		r.Context(),
		session.SESSION_KEY_MESSAGE,
		fmt.Sprintf("Deleted tag '%s'.", r.FormValue("name")),
	)

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	tags, err := t.tagService.GetAll()
	if err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}

	message := t.session.PopString(r.Context(), session.SESSION_KEY_MESSAGE)

	tagView := TagsView{
		Message:   message,
		Tags:      tags,
		CsrfToken: t.csrfToken(r),
		IsAuthenticated: t.session.Exists(
			r.Context(),
			string(session.IS_LOGGED_IN_CONTEXT_KEY),
		),
	}

	err = t.templates["pages/admin/tags.tmpl"].ExecuteTemplate(
		w,
		"admin",
		tagView,
	)
	if err != nil {
		t.errorHandlers.InternalServerError(w, r)
	}
}

func (t *TagController) GetAdminTagsSlugHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")

	tag, err := t.tagService.GetByAttribute("slug", slug)
	if err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}

	tagView := TagView{
		Tag:       tag,
		CsrfToken: t.csrfToken(r),
		IsAuthenticated: t.session.Exists(
			r.Context(),
			string(session.IS_LOGGED_IN_CONTEXT_KEY),
		),
	}

	if err := t.templates["pages/admin/tag.tmpl"].ExecuteTemplate(w, "admin", tagView); err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (t *TagController) PostAdminTagsSlugHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		t.errorHandlers.BadRequest(w, r)
		return
	}

	tag := TagUpdateRequestDto{
		Id:   id,
		Name: r.FormValue("name"),
		Slug: r.FormValue("slug"),
	}

	if _, err := t.tagService.Update(&tag); err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}

	t.session.Put(
		r.Context(),
		session.SESSION_KEY_MESSAGE,
		fmt.Sprintf("Updated tag '%s'.", tag.Name),
	)

	http.Redirect(w, r, "/admin/tags", http.StatusSeeOther)
}

func (t *TagController) GetAdminTagsNewHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	tagView := TagCreateView{
		CsrfToken: t.csrfToken(r),
		IsAuthenticated: t.session.Exists(
			r.Context(),
			string(session.IS_LOGGED_IN_CONTEXT_KEY),
		),
	}

	if err := t.templates["pages/admin/new-tag.tmpl"].ExecuteTemplate(w, "admin", tagView); err != nil {
		t.errorHandlers.InternalServerError(w, r)
		return
	}
}
