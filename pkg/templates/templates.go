package templates

import (
	"html/template"
	"io"
	"os"
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type TemplateCache map[string]Template

func newTemplateCache(templateDir, pageGlob string) (TemplateCache, error) {
	cache := TemplateCache{}

	glob := path.Join(templateDir, pageGlob)
	pages, err := doublestar.FilepathGlob(glob)
	if err != nil {
		return nil, err
	}

	adminBaseTemplate := path.Join(templateDir, "base", "admin.tmpl")
	publicBaseTemplate := path.Join(templateDir, "base", "public.tmpl")

	for _, page := range pages {
		files := []string{
			publicBaseTemplate,
			adminBaseTemplate,
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		name := strings.ReplaceAll(
			page,
			templateDir+"/",
			"",
		)

		cache[name] = ts
	}

	return cache, nil
}

func NewTemplateCache() (TemplateCache, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return newTemplateCache(path.Join(pwd, "web", "templates"), "**/*.tmpl")
}
