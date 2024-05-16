package home

import (
	"net/http"

	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

type HomeController struct {
	tagService     tag.TagService
	articleService article.ArticleService
	log            pkg.Logger
	templateCache  map[string]pkg.Template
	sessionManager pkg.SessionManager
}

type HomeView struct {
	Tags     *[]tag.TagResponseDto
	Articles *[]article.ArticleResponseDto
}

type TagView struct {
	Tag      *tag.TagResponseDto
	Articles *[]article.ArticleResponseDto
}

func NewHomeController(tagService tag.TagService, articleService article.ArticleService, config pkg.ControllerConfig) HomeController {
	return HomeController{
		tagService:     tagService,
		articleService: articleService,
		log:            config.Log,
		templateCache:  config.TemplateCache,
		sessionManager: config.SessionManager,
	}
}

func (h *HomeController) HomeGet(w http.ResponseWriter, r *http.Request) {
	articles, err := h.articleService.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tags, err := h.tagService.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.templateCache["pages/public/index.tmpl"].ExecuteTemplate(w, "public", HomeView{
		Articles: articles,
		Tags:     tags,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *HomeController) HomeArticlesGet(w http.ResponseWriter, r *http.Request) {
	if err := h.templateCache["pages/public/public-articles.tmpl"].ExecuteTemplate(w, "public", HomeView{}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *HomeController) HomeTagsGet(w http.ResponseWriter, r *http.Request) {
	if err := h.templateCache["pages/public/public-tags.tmpl"].ExecuteTemplate(w, "public", HomeView{}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *HomeController) HomeTagGet(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	tag, err := h.tagService.GetByAttribute("slug", slug)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	articles, err := h.articleService.GetManyByAttribute("tagSlug", slug)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err := h.templateCache["pages/public/public-tag.tmpl"].ExecuteTemplate(w, "public", TagView{
		Tag:      tag,
		Articles: articles,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
