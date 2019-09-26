package main

import (
	"net/http"

	"go.uber.org/atomic"
)

var requests = atomic.NewInt64(0)

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

func (app *application) maxRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.config.MaxRequests > 0 {
			if i := requests.Inc(); i >= app.config.MaxRequests {
				close(requestLimitChan)
			}
		}

		h(w, r)
	}
}
