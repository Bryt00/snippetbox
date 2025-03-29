package main

import (
	"html/template"
	"path/filepath"
	"time"

	"raven.net/snippetbox/pkg/forms"

	datalayer "raven.net/snippetbox/pkg/data_layer"
)

// struct to hold templates
type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Snippet         *datalayer.Snippet
	Snippets        []*datalayer.Snippet
	Form            *forms.Form
	IsAuthenticated bool
	User            *datalayer.User
}

func humanDate(t time.Time) string {
	return t.UTC().Format("02 Jan 2006 at 15:04")

}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
