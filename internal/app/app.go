package app

import (
	"fmt"
	"github.com/pistolricks/m-api/internal/api"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (a *Application) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is Available\n")
}
