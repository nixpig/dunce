package pkg

import (
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func newTemplateCache(templateDir, pageGlob string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

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

		name := strings.ReplaceAll(page, strings.Join([]string{templateDir, "/"}, ""), "")

		cache[name] = ts
	}

	return cache, nil
}

func NewTemplateCache() (map[string]*template.Template, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return newTemplateCache(path.Join(pwd, "web", "templates"), "**/*.tmpl")
}
