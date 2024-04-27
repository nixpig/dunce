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
	Tags     map[int]string
	Articles []Article
}

type ArticlesController struct {
	service    ArticleServiceInterface
	tagService tags.TagServiceInterface
}

func NewArticleController(service ArticleServiceInterface, tagsService tags.TagServiceInterface) ArticlesController {
	return ArticlesController{
		service:    service,
		tagService: tagsService,
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
	articles, err := ac.service.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tags, err := ac.tagService.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var viewableTags = map[int]string{}

	for _, tag := range *tags {
		viewableTags[tag.Id] = tag.Name
	}

	articlesScreen := ArticlesScreen{
		Tags:     viewableTags,
		Articles: *articles,
	}

	templates := []string{
		"./web/templates/admin/base.tmpl",
		"./web/templates/admin/articles.tmpl",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", articlesScreen); err != nil {
		fmt.Println("qux >>>", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac *ArticlesController) NewHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := ac.tagService.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var viewableTags = map[int]string{}

	for _, tag := range *tags {
		viewableTags[tag.Id] = tag.Name
	}

	articlesScreen := ArticlesScreen{
		Tags:     viewableTags,
		Articles: nil,
	}

	templates := []string{
		"./web/templates/admin/base.tmpl",
		"./web/templates/admin/new-article.tmpl",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", articlesScreen); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
