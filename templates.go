package main

import (
	"html/template"
	"path/filepath"

	packr "github.com/gobuffalo/packr/v2"
)

type templateData struct {
	CurrentDate string
	Env         map[string]string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	templatesBox := packr.New("templates", "./templates")

	pageCache := map[string]string{}
	layoutCache := map[string]string{}
	partialCache := map[string]string{}
	err := templatesBox.Walk(func(path string, file packr.File) error {
		name := filepath.Base(path)

		if ok, _ := filepath.Match("*.page.tmpl", path); ok {
			pageCache[name] = file.String()
		} else if ok, _ := filepath.Match("*.layout.tmpl", path); ok {
			layoutCache[name] = file.String()
		} else if ok, _ := filepath.Match("*.partial.tmpl", path); ok {
			partialCache[name] = file.String()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for name, file1 := range pageCache {
		ts, err := template.New(name).Parse(file1)
		if err != nil {
			return nil, err
		}

		for _, file2 := range layoutCache {
			ts, err = ts.Parse(file2)
			if err != nil {
				return nil, err
			}
		}

		for _, file3 := range partialCache {
			ts, err = ts.Parse(file3)
			if err != nil {
				return nil, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
