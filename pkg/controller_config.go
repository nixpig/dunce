package pkg

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type ControllerConfig struct {
	Log            Logger
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
}

func NewControllerConfig(
	log Logger,
	templateCache map[string]*template.Template,
	sessionManager *scs.SessionManager,
) ControllerConfig {
	return ControllerConfig{
		Log:            log,
		TemplateCache:  templateCache,
		SessionManager: sessionManager,
	}
}
