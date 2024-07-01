package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sparrowsl/greenlight/internal/validator"
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
	maxBytes := 1_048_576 // limit for the request body size
	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

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

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError): // optional limit for the request body
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(query url.Values, key string, defaultValue string) string {
	str := query.Get(key)

	if str == "" {
		return defaultValue
	}

	return str
}

func (app *application) readCSV(query url.Values, key string, defaultValue []string) []string {
	csv := query.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(query url.Values, key string, defaultValue int, validator *validator.Validator) int {
	// Extract the value from the query string.
	s := query.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}

	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		validator.AddError(key, "must be an integer value")
		return defaultValue
	}

	// Otherwise, return the converted integer value.
	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		// catch any panic in the background if code fails/crashes
		defer func() {
			if err := recover(); err != nil {
				app.logger.Println(fmt.Errorf("%s", err))
			}
		}()

		fn() // run the function to run in the background
	}()
}
