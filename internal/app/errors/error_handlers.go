package errors

import (
	"net/http"

	"github.com/nixpig/dunce/pkg/templates"
)

type ErrorHandlers interface {
	NotFound(w http.ResponseWriter, r *http.Request)
	InternalServerError(w http.ResponseWriter, r *http.Request)
	BadRequest(w http.ResponseWriter, r *http.Request)
}

type ErrorHandlersImpl struct {
	templateCache templates.TemplateCache
}

type ErrorView struct {
	Title   string
	Message string
}

func NewErrorHandlersImpl(templateCache templates.TemplateCache) ErrorHandlers {
	return ErrorHandlersImpl{templateCache}
}

func (e ErrorHandlersImpl) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := e.templateCache["pages/errors/error.tmpl"].
		ExecuteTemplate(w, "public", ErrorView{
			Title:   "404 Not Found",
			Message: "Unable to find the requested resource.",
		}); err != nil {
		e.InternalServerError(w, r)
	}
}

func (e ErrorHandlersImpl) InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)

	if err := e.templateCache["pages/errors/error.tmpl"].
		ExecuteTemplate(w, "public", ErrorView{
			Title:   "500 Internal Server Error",
			Message: "Something went wrong. Please try again.",
		}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (e ErrorHandlersImpl) BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	if err := e.templateCache["pages/errors/error.tmpl"].
		ExecuteTemplate(w, "public", ErrorView{
			Title:   "400 Bad Request",
			Message: "There was something wrong with your request. Please check and try again.",
		}); err != nil {
		e.InternalServerError(w, r)
	}
}
