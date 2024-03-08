package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sparrowsl/greenlight/internal/data"
)

func (app *application) showMovie(writer http.ResponseWriter, request *http.Request) {
	movieId, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	newMovie := data.Movie{
		ID:        movieId,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(writer, http.StatusOK, map[string]any{"movie": newMovie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

func (app *application) createMovie(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Creating a new movie...")
}
