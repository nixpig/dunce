package pkg

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
)

func newTemplateCache(templateDir, pageGlob string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templatePath := path.Join(pwd, templateDir)

	baseTemplate := path.Join(templatePath, "base.tmpl")
	adminPageTemplatePath := path.Join(templatePath, "pages", "admin")
	// partialTemplatePath := path.Join(templatePath, "partials")

	pages, err := filepath.Glob(path.Join(adminPageTemplatePath, pageGlob))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			baseTemplate,
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func NewTemplateCache() (map[string]*template.Template, error) {
	return newTemplateCache("web/templates", "*.tmpl")
}
