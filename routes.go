package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/dist/*
var assets embed.FS

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.maxRequests(app.home))

	assetsFS, err := fs.Sub(assets, "assets/dist")
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assetsFS))))

	return mux
}
