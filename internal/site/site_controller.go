package site

import (
	"net/http"

	"github.com/nixpig/dunce/internal/app/errors"
	"github.com/nixpig/dunce/pkg/logging"
	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
)

type SiteController struct {
	service       SiteService
	log           logging.Logger
	templates     templates.TemplateCache
	session       session.SessionManager
	csrfToken     func(r *http.Request) string
	errorHandlers errors.ErrorHandlers
}

type SiteControllerConfig struct {
	Log            logging.Logger
	TemplateCache  templates.TemplateCache
	SessionManager session.SessionManager
	CsrfToken      func(r *http.Request) string
	ErrorHandlers  errors.ErrorHandlers
}

type SiteItemsView struct {
	Message         string
	SiteItems       *[]SiteItemResponseDto
	CsrfToken       string
	IsAuthenticated bool
}

func NewSiteController(service SiteService, config SiteControllerConfig) SiteController {
	return SiteController{
		service:       service,
		log:           config.Log,
		templates:     config.TemplateCache,
		session:       config.SessionManager,
		csrfToken:     config.CsrfToken,
		errorHandlers: config.ErrorHandlers,
	}
}

func (s *SiteController) GetCreateSiteItems(w http.ResponseWriter, r *http.Request) {
	if err := s.templates["pages/admin/site.tmpl"].ExecuteTemplate(w, "admin", nil); err != nil {
		s.errorHandlers.InternalServerError(w, r)
		return
	}
}
