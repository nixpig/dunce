package articles

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/nixpig/dunce/internal/tags"
)

type ArticlesScreen struct {
	Tags     []tags.Tag
	Articles []Article
}

type ArticlesController struct {
	service ArticleServiceInterface
	tags    tags.TagServiceInterface
}

func NewArticleController(service ArticleServiceInterface, tagsService tags.TagServiceInterface) ArticlesController {
	return ArticlesController{
		service: service,
		tags:    tagsService,
	}
}

func (ac *ArticlesController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("form data:", r.Form)

	tagsForm := r.Form["tags[]"]

	fmt.Println("tagsForm: ", tagsForm)

	var tagIds []int

	for _, t := range tagsForm {
		tagId, err := strconv.Atoi(t)
		if err != nil {
			http.Error(w, "Unable to parse tags to ints", http.StatusInternalServerError)
			return
		}

		tagIds = append(tagIds, tagId)
	}

	article := NewArticle(
		r.FormValue("title"),
		r.FormValue("subtitle"),
		r.FormValue("slug"),
		r.FormValue("body"),
		time.Now(),
		time.Now(),
		tagIds,
	)

	if _, err := ac.service.create(&article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (ac *ArticlesController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	articles := []Article{}

	tags, err := ac.tags.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articlesScreen := ArticlesScreen{
		Tags:     *tags,
		Articles: articles,
	}

	templates := []string{
		"./web/templates/base.tmpl",
		"./web/templates/articles.tmpl",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if err := ts.ExecuteTemplate(w, "base", articlesScreen); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
