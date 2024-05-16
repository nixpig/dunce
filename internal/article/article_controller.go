package article

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/nosurf"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

const longFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

type ArticleController struct {
	articleService ArticleService
	tagService     tag.TagService
	log            pkg.Logger
	templates      map[string]pkg.Template
	session        pkg.SessionManager
}

type ArticleView struct {
	Message         string
	Article         *ArticleResponseDto
	Tags            *[]tag.TagResponseDto
	Content         template.HTML
	CsrfToken       string
	IsAuthenticated bool
}

type ArticlesView struct {
	Message         string
	Articles        *[]ArticleResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

type ArticlePublishView struct {
	Message         string
	Articles        *[]ArticleResponseDto
	Tags            *[]tag.TagResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

func NewArticleController(
	service ArticleService,
	tagsService tag.TagService,
	config pkg.ControllerConfig,
) ArticleController {
	return ArticleController{
		articleService: service,
		tagService:     tagsService,
		session:        config.SessionManager,
		log:            config.Log,
		templates:      config.TemplateCache,
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

	article := ArticleNewRequestDto{
		Title:     r.FormValue("title"),
		Subtitle:  r.FormValue("subtitle"),
		Slug:      r.FormValue("slug"),
		Body:      r.FormValue("body"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TagIds:    tagIds,
	}

	if _, err := a.articleService.Create(&article); err != nil {
		a.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a *ArticleController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := a.articleService.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := a.templates["pages/admin/admin-articles.tmpl"].ExecuteTemplate(w, "admin", ArticlesView{
		Articles:        articles,
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: a.session.Exists(r.Context(), string(pkg.IS_LOGGED_IN_CONTEXT_KEY)),
	}); err != nil {
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

	if err := a.templates["pages/admin/admin-new-article.tmpl"].ExecuteTemplate(w, "admin", ArticlePublishView{
		Tags:            availableTags,
		CsrfToken:       nosurf.Token(r),
		IsAuthenticated: a.session.Exists(r.Context(), string(pkg.IS_LOGGED_IN_CONTEXT_KEY)),
	}); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (a *ArticleController) GetBySlugHander(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	article, err := a.articleService.GetByAttribute("slug", slug)
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

	if err := a.templates["pages/admin/admin-article.tmpl"].ExecuteTemplate(
		w,
		"admin",
		ArticleView{
			Article:         article,
			Tags:            allTags,
			CsrfToken:       nosurf.Token(r),
			IsAuthenticated: a.session.Exists(r.Context(), string(pkg.IS_LOGGED_IN_CONTEXT_KEY)),
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

	article := ArticleUpdateRequestDto{
		Id:        articleId,
		Title:     r.FormValue("title"),
		Subtitle:  r.FormValue("subtitle"),
		Slug:      r.FormValue("slug"),
		Body:      r.FormValue("body"),
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
		TagIds:    tagIds,
	}

	_, err = a.articleService.Update(&article)
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

	if err := a.articleService.DeleteById(id); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a ArticleController) PublicGetArticle(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	article, err := a.articleService.GetByAttribute("slug", slug)
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	content, err := pkg.MdToHtml([]byte(article.Body))
	if err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err := a.templates["pages/public/public-article.tmpl"].ExecuteTemplate(
		w,
		"public",
		ArticleView{
			Article: article,
			Content: template.HTML(content),
		},
	); err != nil {
		a.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
