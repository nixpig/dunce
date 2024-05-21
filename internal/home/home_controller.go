package home

import (
	"net/http"

	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
)

type HomeController struct {
	tagService     tag.TagService
	articleService article.ArticleService
	log            logging.Logger
	templateCache  templates.TemplateCache
	sessionManager session.SessionManager
	errorHandlers  errors.ErrorHandlers
}

type HomeControllerConfig struct {
	Log            logging.Logger
	TemplateCache  templates.TemplateCache
	SessionManager session.SessionManager
	CsrfToken      func(*http.Request) string
	ErrorHandlers  errors.ErrorHandlers
}

type HomeView struct {
	Tags     *[]tag.TagResponseDto
	Articles *[]article.ArticleResponseDto
}

type TagView struct {
	Tag      *tag.TagResponseDto
	Articles *[]article.ArticleResponseDto
}

func NewHomeController(
	tagService tag.TagService,
	articleService article.ArticleService,
	config HomeControllerConfig,
) HomeController {
	return HomeController{
		tagService:     tagService,
		articleService: articleService,
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
		errorHandlers:  config.ErrorHandlers,
	}
}

func (h *HomeController) HomeGet(w http.ResponseWriter, r *http.Request) {
	articles, err := h.articleService.GetAll()
	if err != nil {
		h.errorHandlers.InternalServerError(w, r)
		return
	}

	tags, err := h.tagService.GetAll()
	if err != nil {
		h.errorHandlers.InternalServerError(w, r)
		return
	}

	if err := h.templateCache["pages/public/index.tmpl"].ExecuteTemplate(w, "public", HomeView{
		Articles: articles,
		Tags:     tags,
	}); err != nil {
		h.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (h *HomeController) HomeArticlesGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := h.templateCache["pages/public/articles.tmpl"].ExecuteTemplate(w, "public", HomeView{}); err != nil {
		h.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (h *HomeController) HomeTagsGet(w http.ResponseWriter, r *http.Request) {
	if err := h.templateCache["pages/public/tags.tmpl"].ExecuteTemplate(w, "public", HomeView{}); err != nil {
		h.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (h *HomeController) HomeTagGet(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := h.tagService.GetByAttribute("slug", slug)
	if err != nil {
		h.errorHandlers.NotFound(w, r)
		return
	}

	articles, err := h.articleService.GetManyByAttribute("tagSlug", slug)
	if err != nil {
		h.errorHandlers.NotFound(w, r)
		return
	}

	if err := h.templateCache["pages/public/tag.tmpl"].ExecuteTemplate(w, "public", TagView{
		Tag:      tag,
		Articles: articles,
	}); err != nil {
		h.errorHandlers.InternalServerError(w, r)
		return
	}
}
