package pkg

import (
	"net/http"
)

type ControllerConfig struct {
	Log            Logger
	TemplateCache  TemplateCache
	SessionManager SessionManager
	CsrfToken      func(*http.Request) string
}

func NewControllerConfig(
	log Logger,
	templateCache TemplateCache,
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
