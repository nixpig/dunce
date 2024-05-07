package pkg

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"slices"
)

func newTemplateCache(templateDir, pageGlob string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templatePath := path.Join(pwd, templateDir)

	adminBaseTemplate := path.Join(templatePath, "admin.tmpl")
	publicBaseTemplate := path.Join(templatePath, "public.tmpl")
	adminPageTemplatePath := path.Join(templatePath, "pages", "admin")
	publicPageTemplatePath := path.Join(templatePath, "pages", "public")
	// partialTemplatePath := path.Join(templatePath, "partials")

	adminPages, err := filepath.Glob(path.Join(adminPageTemplatePath, pageGlob))
	if err != nil {
		return nil, err
	}

	publicPages, err := filepath.Glob(path.Join(publicPageTemplatePath, pageGlob))
	if err != nil {
		return nil, err
	}

	pages := slices.Concat(adminPages, publicPages)

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			publicBaseTemplate,
			adminBaseTemplate,
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
