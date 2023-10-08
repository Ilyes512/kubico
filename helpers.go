package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tomasen/realip"
)

type defaultKvTableCache struct {
	cache *kvTable
	once  sync.Once
}

var defaultKvTableCaches = map[string]*defaultKvTableCache{
	"Env":     {},
	"General": {},
}

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
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
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

	if td.KvTables == nil {
		td.KvTables = map[string]*kvTable{}
	}

	if _, exists := td.KvTables["General"]; !exists {
		td.KvTables["General"] = func() *kvTable {
			defaultKvTableCaches["General"].once.Do(func() {
				defaultKvTableCaches["General"].cache = &kvTable{
					Title: "General",
					Values: map[string]string{
						"Current Date":   time.Now().Format(time.RFC3339),
						"Remote Address": realip.FromRequest(r),
					},
				}

				if app.config.Timeout > 0 {
					defaultKvTableCaches["General"].cache.Values["Timeout after"] = fmt.Sprintf("%s seconds", strconv.FormatInt(app.config.Timeout, 10))
				}

				if app.config.MaxRequests > 0 {
					defaultKvTableCaches["General"].cache.Values["Max requests"] = strconv.FormatInt(app.config.MaxRequests, 10)
				}
			})
			return defaultKvTableCaches["General"].cache
		}()
	}

	if _, exists := td.KvTables["Env"]; !exists {
		td.KvTables["Env"] = func() *kvTable {
			defaultKvTableCaches["Env"].once.Do(func() {
				defaultKvTableCaches["Env"].cache = &kvTable{
					Title:  "Environment variables",
					Values: map[string]string{},
				}
				for _, env := range os.Environ() {
					if i := strings.IndexByte(env, '='); i >= 0 {
						defaultKvTableCaches["Env"].cache.Values[env[:i]] = env[i+1:]
					}
				}
			})
			return defaultKvTableCaches["Env"].cache
		}()
	}

	return td
}
