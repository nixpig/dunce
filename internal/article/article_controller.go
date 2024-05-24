package article

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/markdown"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
)

const longTimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

type ArticleController struct {
	articleService ArticleService
	tagService     tag.TagService
	log            logging.Logger
	templates      templates.TemplateCache
	session        session.SessionManager
	csrfToken      func(r *http.Request) string
	errorHandlers  errors.ErrorHandlers
}

type ArticleControllerConfig struct {
	Log            logging.Logger
	TemplateCache  templates.TemplateCache
	SessionManager session.SessionManager
	CsrfToken      func(*http.Request) string
	ErrorHandlers  errors.ErrorHandlers
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
	config ArticleControllerConfig,
) ArticleController {
	return ArticleController{
		articleService: service,
		tagService:     tagsService,
		session:        config.SessionManager,
		log:            config.Log,
		templates:      config.TemplateCache,
		csrfToken:      config.CsrfToken,
		errorHandlers:  config.ErrorHandlers,
	}
}

func (a *ArticleController) CreateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := r.ParseForm(); err != nil {
		a.errorHandlers.BadRequest(w, r)
		return
	}

	tagsForm := r.Form["tags[]"]

	tagIds := make([]int, len(tagsForm))

	for i, t := range tagsForm {
		tagId, err := strconv.Atoi(t)
		if err != nil {
			a.errorHandlers.InternalServerError(w, r)
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
		a.errorHandlers.InternalServerError(w, r)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a *ArticleController) GetAllHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	articles, err := a.articleService.GetAll()
	if err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}

	if err := a.templates["pages/admin/articles.tmpl"].ExecuteTemplate(w, "admin", ArticlesView{
		Articles:        articles,
		CsrfToken:       a.csrfToken(r),
		IsAuthenticated: a.session.Exists(r.Context(), string(session.IS_LOGGED_IN_CONTEXT_KEY)),
	}); err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (a *ArticleController) NewHandler(w http.ResponseWriter, r *http.Request) {
	availableTags, err := a.tagService.GetAll()
	if err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}

	if err := a.templates["pages/admin/new-article.tmpl"].ExecuteTemplate(w, "admin", ArticlePublishView{
		Tags:            availableTags,
		CsrfToken:       a.csrfToken(r),
		IsAuthenticated: a.session.Exists(r.Context(), string(session.IS_LOGGED_IN_CONTEXT_KEY)),
	}); err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (a *ArticleController) GetBySlugHander(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")

	article, err := a.articleService.GetByAttribute("slug", slug)
	if err != nil {
		a.errorHandlers.BadRequest(w, r)
		return
	}

	allTags, err := a.tagService.GetAll()
	if err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}

	if err := a.templates["pages/admin/article.tmpl"].ExecuteTemplate(
		w,
		"admin",
		ArticleView{
			Article:         article,
			Tags:            allTags,
			CsrfToken:       a.csrfToken(r),
			IsAuthenticated: a.session.Exists(r.Context(), string(session.IS_LOGGED_IN_CONTEXT_KEY)),
		},
	); err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}
}

func (a ArticleController) UpdateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := r.ParseForm(); err != nil {
		a.errorHandlers.BadRequest(w, r)
		return
	}

	tags := r.Form["tags[]"]

	tagIds := make([]int, len(tags))

	for i, t := range tags {
		id, err := strconv.Atoi(t)
		if err != nil {
			a.errorHandlers.InternalServerError(w, r)
			return
		}

		tagIds[i] = id
	}

	createdAt, err := time.Parse(longTimeFormat, r.FormValue("created_at"))
	if err != nil {
		a.errorHandlers.BadRequest(w, r)
		return
	}

	articleId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		a.errorHandlers.BadRequest(w, r)
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
		a.errorHandlers.InternalServerError(w, r)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a ArticleController) AdminArticlesDeleteHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		a.errorHandlers.BadRequest(w, r)
		return
	}

	if err := a.articleService.DeleteById(id); err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (a ArticleController) PublicGetArticle(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")

	article, err := a.articleService.GetByAttribute("slug", slug)
	if err != nil {
		a.errorHandlers.NotFound(w, r)
		return
	}

	content, err := markdown.MdToHtml([]byte(article.Body))
	if err != nil {
		a.errorHandlers.BadRequest(w, r)
		return
	}

	if err := a.templates["pages/public/article.tmpl"].ExecuteTemplate(
		w,
		"public",
		ArticleView{
			Article: article,
			Content: template.HTML(content),
		},
	); err != nil {
		a.errorHandlers.InternalServerError(w, r)
		return
	}
}
