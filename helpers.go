package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var cachedEnv map[string]string
var fetchEnvOnce sync.Once

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	app.clientError(w, http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentDate = time.Now().Format(time.RFC3339)
	td.Env = func() map[string]string {
		fetchEnvOnce.Do(func() {
			cachedEnv = map[string]string{}
			for _, env := range os.Environ() {
				if i := strings.IndexByte(env, '='); i >= 0 {
					cachedEnv[env[:i]] = env[i+1:]
				}
			}
		})

		return cachedEnv
	}()

	return td
}
