package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sparrowsl/greenlight/internal/data"
	"github.com/sparrowsl/greenlight/internal/validator"
)

func (app *application) listAllMovies(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	val := validator.New()
	query := request.URL.Query()

	input.Title = app.readString(query, "title", "")
	input.Genres = app.readCSV(query, "genres", []string{})

	input.Filters.Page = app.readInt(query, "page", 1, val)
	input.Filters.PageSize = app.readInt(query, "page_size", 20, val)

	// Extract the sort query string value, falling back to "id" i
	// by the client (which will imply a ascending sort on movie I
	input.Filters.Sort = app.readString(query, "sort", "id")

	if !val.Valid() {
		app.failedValidationResponse(writer, request, val.Errors)
		return
	}

	fmt.Fprintf(writer, "%+v\n", input)
}

func (app *application) showMovie(writer http.ResponseWriter, request *http.Request) {
	movieId, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	movie, err := app.models.Movies.Get(movieId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.writeJSON(writer, http.StatusOK, map[string]any{"movie": movie}, nil)
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

	app.writeJSON(writer, http.StatusCreated, map[string]any{"movie": movie}, nil)
}

func (app *application) updateMovie(writer http.ResponseWriter, request *http.Request) {
	movieId, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	movie, err := app.models.Movies.Get(movieId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	if err = app.readJSON(writer, request, &input); err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	val := validator.New()
	if data.ValidateMovie(val, movie); !val.Valid() {
		app.failedValidationResponse(writer, request, val.Errors)
		return
	}

	if err := app.models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}

		return
	}

	err = app.writeJSON(writer, http.StatusOK, map[string]any{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

func (app *application) deleteMovie(writer http.ResponseWriter, request *http.Request) {
	movieId, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	if err := app.models.Movies.Delete(movieId); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.writeJSON(writer, http.StatusNoContent, map[string]any{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
