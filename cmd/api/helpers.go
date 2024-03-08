package main

import (
	"encoding/json"
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

func (app *application) writeJSON(writer http.ResponseWriter, status int, data map[string]any, headers http.Header) error {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	result = append(result, '\n') // to display new line on terminal

	for key, value := range headers {
		writer.Header()[key] = value
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(result)

	return nil
}
