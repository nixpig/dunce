package article

import (
	"html/template"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

const longFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

type ArticleController struct {
	service       pkg.Service[Article, ArticleNew]
	tagService    pkg.Service[tag.Tag, tag.TagData]
	log           pkg.Logger
	templateCache map[string]*template.Template
}

func NewArticleController(
	service pkg.Service[Article, ArticleNew],
	tagsService pkg.Service[tag.Tag, tag.TagData],
	log pkg.Logger,
	templateCache map[string]*template.Template,
) ArticleController {
	return ArticleController{
		service:       service,
		tagService:    tagsService,
		log:           log,
		templateCache: templateCache,
	}
}

func (a *ArticleController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tagsForm := r.Form["tags[]"]

	tagIds := make([]int, len(tagsForm))

	for i, t := range tagsForm {
		tagId, err := strconv.Atoi(t)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Unable to parse tags to ints", http.StatusInternalServerError)
			return
		}

		tagIds[i] = tagId
	}

	article := ArticleNew{
		Title:     r.FormValue("title"),
		Subtitle:  r.FormValue("subtitle"),
		Slug:      r.FormValue("slug"),
		Body:      r.FormValue("body"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    tagIds,
	}

	if _, err := a.service.Create(&article); err != nil {
		a.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a *ArticleController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := a.service.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := a.templateCache["articles.tmpl"].ExecuteTemplate(w, "base", articles); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (a *ArticleController) NewHandler(w http.ResponseWriter, r *http.Request) {
	availableTags, err := a.tagService.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := a.templateCache["new-article.tmpl"].ExecuteTemplate(w, "base", availableTags); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (a *ArticleController) GetBySlugHander(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	article, err := a.service.GetBySlug(slug)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	allTags, err := a.tagService.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := a.templateCache["article.tmpl"].ExecuteTemplate(
		w,
		"base",
		struct {
			Article Article
			Tags    []tag.Tag
		}{
			Article: *article,
			Tags:    *allTags,
		},
	); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (a ArticleController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tags := r.Form["tags[]"]

	tagIds := make([]int, len(tags))

	for i, t := range tags {
		id, err := strconv.Atoi(t)
		if err != nil {
			a.log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tagIds[i] = id
	}

	createdAt, err := time.Parse(longFormat, r.FormValue("created_at"))
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articleId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tagList, err := a.tagService.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articleTags := make([]tag.Tag, len(tagIds))

	for index, tagId := range tagIds {
		tagListIndex := slices.IndexFunc(*tagList, func(t tag.Tag) bool {
			return t.Id == tagId
		})

		if tagListIndex > -1 {
			articleTags[index] = (*tagList)[tagListIndex]
		}
	}

	article := NewArticleWithId(
		articleId,
		r.FormValue("title"),
		r.FormValue("subtitle"),
		r.FormValue("slug"),
		r.FormValue("body"),
		createdAt,
		time.Now(),
		articleTags,
	)

	_, err = a.service.Update(&article)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a ArticleController) AdminArticlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := a.service.DeleteById(id); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}
