package main

import (
	"embed"
	"html/template"
	"io/fs"
	"path/filepath"
)

type kvTable struct {
	Title  string
	Values map[string]string
}

type templateData struct {
	KvTables map[string]*kvTable
}

//go:embed templates/*.tmpl
var templates embed.FS

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pageCache := map[string]string{}
	layoutCache := map[string]string{}
	partialCache := map[string]string{}
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := filepath.Base(path)

		if ok, err := filepath.Match("*.page.tmpl", name); ok || err != nil {
			if err != nil {
				return err
			}

			file, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}

			pageCache[name] = string(file)
		} else if ok, err := filepath.Match("*.layout.tmpl", name); ok || err != nil {
			if err != nil {
				return err
			}

			file, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}

			layoutCache[name] = string(file)
		} else if ok, err := filepath.Match("*.partial.tmpl", name); ok || err != nil {
			if err != nil {
				return err
			}

			file, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}

			partialCache[name] = string(file)
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
