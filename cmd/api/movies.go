package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sparrowsl/greenlight/internal/data"
	"github.com/sparrowsl/greenlight/internal/validator"
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
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	val := validator.New()
	if data.ValidateMovie(val, movie); !val.Valid() {
		app.failedValidationResponse(writer, request, val.Errors)
		return
	}

	if err := app.models.Movies.Insert(movie); err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	app.writeJSON(writer, http.StatusOK, map[string]any{"movie": movie}, nil)
}
