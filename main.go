//go:generate packr2
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type config struct {
	envfile      string
	Host         string `env:"KUBICO_HOST" envDefault:"0.0.0.0"`
	Port         string `env:"KUBICO_PORT" envDefault:"8080"`
	Cacheheaders bool   `env:"KUBICO_NO_CACHE"`
	Timeout      int64  `env:"KUBICO_TIMEOUT_SECONDS"`
	MaxRequests  int64  `env:"KUBICO_MAX_REQUESTS"`
}

type application struct {
	config        *config
	errorLog      *log.Logger
	infoLog       *log.Logger
	server        *http.Server
	wg            sync.WaitGroup
	stoppingCh    chan string
	templateCache map[string]*template.Template
}

func main() {
	app := &application{
		config:     &config{},
		stoppingCh: make(chan string, 1),
		errorLog:   log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:    log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime),
	}

	err := app.fetchConfig()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	stopCh := make(chan struct{})

	go func() {
		msg := <-app.stoppingCh
		if msg != "" {
			app.infoLog.Println(msg)
		}
		close(stopCh)
	}()
	app.SetAppTimeout()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-sigCh

		select {
		case app.stoppingCh <- fmt.Sprint("Got a OS interrupt signal"):
		default:
		}
	}()

	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.templateCache = templateCache

	app.Start()
	defer app.Stop()

	<-stopCh
}

func (app *application) Start() {
	app.server = &http.Server{
		Addr:           net.JoinHostPort(app.config.Host, app.config.Port),
		ErrorLog:       app.errorLog,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        app.noCacheHandler(app.routes()),
	}

	app.wg.Add(1)

	go func() {
		app.infoLog.Printf("Starting server on: http://%s:%s\n", app.config.Host, app.config.Port)
		app.server.ListenAndServe()
		app.wg.Done()
	}()
}

func (app *application) Stop() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.infoLog.Println("Trying to stop server gracefully...")

	if err = app.server.Shutdown(ctx); err != nil {
		if err = app.server.Close(); err != nil {
			app.errorLog.Printf("Stopping server with error: %v\n", err)
			return err
		}
	}

	app.wg.Wait()
	app.infoLog.Println("Server stopped")
	return nil
}

func (app *application) fetchConfig() error {
	envfile := flag.String("env-file", ".env", "Read in a file of environment variables")
	flag.Parse()

	app.config.envfile = path.Clean(*envfile)

	err := godotenv.Load(app.config.envfile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			app.infoLog.Print("No .env found.")
		} else {
			return err
		}
	} else {
		app.infoLog.Printf("Loaded env-file from: %s", app.config.envfile)
	}

	if err := env.Parse(app.config); err != nil {
		return err
	}

	return nil
}

func (app *application) SetAppTimeout() {
	if app.config.Timeout > 0 {
		time.AfterFunc(time.Duration(app.config.Timeout)*time.Second, func() {
			select {
			case app.stoppingCh <- fmt.Sprintf("Kubico timeout after %d seconds.", app.config.Timeout):
			default:
			}
		})
	}
}
