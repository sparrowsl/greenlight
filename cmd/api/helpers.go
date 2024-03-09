package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (app *application) readJSON(writer http.ResponseWriter, request *http.Request, dest any) error {
	err := json.NewDecoder(request.Body).Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}
