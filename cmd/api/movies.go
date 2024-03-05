package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) showMovie(writer http.ResponseWriter, request *http.Request) {
	movieId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		http.NotFound(writer, request)
		// fmt.Fprintln(writer, "No movie with that id exists")
		return
	}

	fmt.Fprintf(writer, "Showing movie with an id: %d\n", movieId)
}

func (app *application) createMovie(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Creating a new movie...")
}
