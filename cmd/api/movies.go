package main

import (
	"fmt"
	"net/http"
)

func (app *application) showMovie(writer http.ResponseWriter, request *http.Request) {
	movieId, err := app.readIDParam(request)
	if err != nil {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "Showing movie with an id: %d\n", movieId)
}

func (app *application) createMovie(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Creating a new movie...")
}
