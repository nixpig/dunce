package home

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/nixpig/dunce/internal/article"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

type HomeController struct {
	templateCache map[string]*template.Template
	api           Api
}

func NewHomeController(api Api, config pkg.ControllerConfig) HomeController {
	return HomeController{
		templateCache: config.TemplateCache,
		api:           api,
	}
}

func (h *HomeController) HomeGet(w http.ResponseWriter, r *http.Request) {

	if err := h.templateCache["index.tmpl"].ExecuteTemplate(w, "base", h.api); err != nil {
		fmt.Println(" <<< err >>> ")
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

func (a Api) Articles() *[]article.Article {
	articles, err := a.articleService.GetAll()
	if err != nil {
		return nil
	}

	return articles
}
