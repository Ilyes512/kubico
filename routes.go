package main

import (
	"net/http"

	packr "github.com/gobuffalo/packr/v2"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.maxRequests(app.home))

	staticBox := packr.New("static", "./assets/dist")
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(staticBox)))

	return mux
}
