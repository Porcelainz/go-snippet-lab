package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snipeetbox.porcelain.com/internal/models"
	"snipeetbox.porcelain.com/ui"
)

type TemplatesData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, "html/base.tmpl", "html/partials/*.tmpl", page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
