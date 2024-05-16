package pkg

import (
	"html/template"
	"net/http"
)

type ControllerConfig struct {
	Log            Logger
	TemplateCache  map[string]*template.Template
	SessionManager SessionManager
	CsrfToken      func(*http.Request) string
}

func NewControllerConfig(
	log Logger,
	templateCache map[string]*template.Template,
	sessionManager SessionManager,
	csrfToken func(*http.Request) string,
) ControllerConfig {
	return ControllerConfig{
		Log:            log,
		TemplateCache:  templateCache,
		SessionManager: sessionManager,
		CsrfToken:      csrfToken,
	}
}
