package articles

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/nixpig/dunce/internal/tags"
	"github.com/nixpig/dunce/pkg/logging"
)

const longFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

type ArticlesScreen struct {
	Tags     map[int]string
	Articles []Article
}

type ArticleScreen struct {
	Tags    map[int]string
	Article Article
}

type ArticlesController struct {
	service    ArticleServiceInterface
	tagService tags.TagServiceInterface
	log        logging.Logger
}

func NewArticleController(
	service ArticleServiceInterface,
	tagsService tags.TagServiceInterface,
	log logging.Logger,
) ArticlesController {
	return ArticlesController{
		service:    service,
		tagService: tagsService,
		log:        log,
	}
}

func (ac *ArticlesController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tagsForm := r.Form["tags[]"]

	var tagIds []int

	for _, t := range tagsForm {
		tagId, err := strconv.Atoi(t)
		if err != nil {
			ac.log.Error(err.Error())
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

	if _, err := ac.service.Create(&article); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (ac *ArticlesController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := ac.service.GetAll()
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tags, err := ac.tagService.GetAll()
	if err != nil {
		ac.log.Error(err.Error())
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
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", articlesScreen); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac *ArticlesController) NewHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := ac.tagService.GetAll()
	if err != nil {
		ac.log.Error(err.Error())
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
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", articlesScreen); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac *ArticlesController) GetBySlugHander(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	article, err := ac.service.GetBySlug(slug)
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tags, err := ac.tagService.GetAll()
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var viewableTags = map[int]string{}

	for _, t := range *tags {
		viewableTags[t.Id] = t.Name
	}

	articleScreen := ArticleScreen{
		Tags:    viewableTags,
		Article: *article,
	}

	templates := []string{
		"./web/templates/admin/base.tmpl",
		"./web/templates/admin/article.tmpl",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", articleScreen); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac ArticlesController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tags := r.Form["tags[]"]

	var tagIds []int

	for _, t := range tags {
		id, err := strconv.Atoi(t)
		if err != nil {
			ac.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tagIds = append(tagIds, id)
	}

	createdAt, err := time.Parse(longFormat, r.FormValue("created_at"))
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articleId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	article := NewArticleWithId(
		articleId,
		r.FormValue("title"),
		r.FormValue("subtitle"),
		r.FormValue("slug"),
		r.FormValue("body"),
		createdAt,
		time.Now(),
		tagIds,
	)

	_, err = ac.service.Update(&article)
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}
