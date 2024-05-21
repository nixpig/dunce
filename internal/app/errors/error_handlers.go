package errors

import (
	"net/http"

	"github.com/nixpig/dunce/pkg/templates"
)

type ErrorHandlers interface {
	NotFound(w http.ResponseWriter, r *http.Request)
	InternalServerError(w http.ResponseWriter, r *http.Request)
}

type ErrorHandlersImpl struct {
	templateCache templates.TemplateCache
}

func NewErrorHandlersImpl(templateCache templates.TemplateCache) ErrorHandlers {
	return ErrorHandlersImpl{templateCache}
}

func (e ErrorHandlersImpl) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := e.templateCache["pages/errors/not-found.tmpl"].ExecuteTemplate(w, "public", nil); err != nil {
		e.InternalServerError(w, r)
	}
}

func (e ErrorHandlersImpl) InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := e.templateCache["pages/errors/internal-server-error.tmpl"].ExecuteTemplate(w, "public", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
