package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger *log.Logger
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &Application{
		Logger: logger,
	}

	return app, nil
}

func (a *Application) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Available\n")
}
