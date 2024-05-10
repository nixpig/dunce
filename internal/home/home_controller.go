package home

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

type HomeController struct {
	log            pkg.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	api            Api
}

func NewHomeController(api Api, config pkg.ControllerConfig) HomeController {
	return HomeController{
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
		api:            api,
	}
}

func (h *HomeController) HomeGet(w http.ResponseWriter, r *http.Request) {
	if err := h.templateCache["pages/public/index.tmpl"].ExecuteTemplate(w, "public", h.api); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *HomeController) HomeArticlesGet(w http.ResponseWriter, r *http.Request) {
	if err := h.templateCache["pages/public/public-articles.tmpl"].ExecuteTemplate(w, "public", h.api); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *HomeController) HomeTagsGet(w http.ResponseWriter, r *http.Request) {
	if err := h.templateCache["pages/public/public-tags.tmpl"].ExecuteTemplate(w, "public", h.api); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *HomeController) HomeTagGet(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tagFromSlug := h.api.Tag(slug)

	articles, err := h.api.articleService.GetManyByAttribute("tagSlug", slug)
	if err != nil {
	}

	if err := h.templateCache["pages/public/public-tag.tmpl"].ExecuteTemplate(w, "public", struct {
		Tag      *tag.Tag
		Articles *[]article.Article
	}{
		Tag:      tagFromSlug,
		Articles: articles,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

type Api struct {
	articleService article.ArticleService
	tagService     tag.TagService
}

func NewApi(articleService article.ArticleService, tagService tag.TagService) Api {
	return Api{
		articleService: articleService,
		tagService:     tagService,
	}
}

func (a Api) Tags() *[]tag.Tag {
	tags, err := a.tagService.GetAll()
	if err != nil {
		return nil
	}

	return tags
}

func (a Api) Tag(slug string) *tag.Tag {
	tag, err := a.tagService.GetByAttribute("slug", slug)
	if err != nil {
		return nil
	}

	return tag
}

func (a Api) Articles() *[]article.Article {
	articles, err := a.articleService.GetAll()
	if err != nil {
		return nil
	}

	return articles
}
