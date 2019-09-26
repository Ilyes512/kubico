package main

import (
	"net/http"
)

func (app *application) noCache(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.config.Cacheheaders == true {
			w.Header().Set("Cache-Control", "no-store, must-revalidate")
		}

		h(w, r)
	}
}

func (app *application) noCacheHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.Cacheheaders == true {
			w.Header().Set("Cache-Control", "no-store, must-revalidate")
		}

		h.ServeHTTP(w, r)
	})
}
