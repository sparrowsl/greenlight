package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) readIDParam(request *http.Request) (int64, error) {
	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil || id < 1 {
		return 0, errors.New("Invalid id parameter")
	}

	return int64(id), nil
}
