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
