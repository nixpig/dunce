package article

import (
	"html/template"
	"net/http"

	"github.com/nixpig/dunce/internal/tag"
	"github.com/nixpig/dunce/pkg"
)

const longFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

type ArticlesController struct {
	service       pkg.Service[Article]
	tagService    pkg.Service[tag.Tag]
	log           pkg.Logger
	templateCache map[string]*template.Template
}

func NewArticleController(
	service pkg.Service[Article],
	tagsService pkg.Service[tag.Tag],
	log pkg.Logger,
	templateCache map[string]*template.Template,
) ArticlesController {
	return ArticlesController{
		service:       service,
		tagService:    tagsService,
		log:           log,
		templateCache: templateCache,
	}
}

func (ac *ArticlesController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// if err := r.ParseForm(); err != nil {
	// 	ac.log.Error(err.Error())
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// 	return
	// }
	//
	// tagsForm := r.Form["tags[]"]
	//
	// var tagIds []int
	//
	// for _, t := range tagsForm {
	// 	tagId, err := strconv.Atoi(t)
	// 	if err != nil {
	// 		ac.log.Error(err.Error())
	// 		http.Error(w, "Unable to parse tags to ints", http.StatusInternalServerError)
	// 		return
	// 	}
	//
	// 	tagIds = append(tagIds, tagId)
	// }
	//
	// article := NewArticle(
	// 	r.FormValue("title"),
	// 	r.FormValue("subtitle"),
	// 	r.FormValue("slug"),
	// 	r.FormValue("body"),
	// 	time.Now(),
	// 	time.Now(),
	// 	tagIds,
	// )
	//
	// if _, err := ac.service.Create(&article); err != nil {
	// 	ac.log.Error(err.Error())
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func (ac *ArticlesController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := ac.service.GetAll()
	if err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ac.templateCache["articles.tmpl"].ExecuteTemplate(w, "base", articles); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac *ArticlesController) NewHandler(w http.ResponseWriter, r *http.Request) {
	if err := ac.templateCache["new-article.tmpl"].ExecuteTemplate(w, "base", nil); err != nil {
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

	if err := ac.templateCache["article.tmpl"].ExecuteTemplate(w, "base", article); err != nil {
		ac.log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ac ArticlesController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// if err := r.ParseForm(); err != nil {
	// 	ac.log.Error(err.Error())
	// 	http.Error(w, "Bad Request", http.StatusBadRequest)
	// 	return
	// }
	//
	// tags := r.Form["tags[]"]
	//
	// var tagIds []int
	//
	// for _, t := range tags {
	// 	id, err := strconv.Atoi(t)
	// 	if err != nil {
	// 		ac.log.Error(err.Error())
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}
	//
	// 	tagIds = append(tagIds, id)
	// }
	//
	// createdAt, err := time.Parse(longFormat, r.FormValue("created_at"))
	// if err != nil {
	// 	ac.log.Error(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	//
	// articleId, err := strconv.Atoi(r.FormValue("id"))
	// if err != nil {
	// 	ac.log.Error(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	//
	// article := NewArticleWithId(
	// 	articleId,
	// 	r.FormValue("title"),
	// 	r.FormValue("subtitle"),
	// 	r.FormValue("slug"),
	// 	r.FormValue("body"),
	// 	createdAt,
	// 	time.Now(),
	// 	tagIds,
	// )
	//
	// _, err = ac.service.Update(&article)
	// if err != nil {
	// 	ac.log.Error(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	//
	// http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}
