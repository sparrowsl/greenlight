package main

import (
	"encoding/json"
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
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		app.errorResponse(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(writer, "%+v\n", input)
}
