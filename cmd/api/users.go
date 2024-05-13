package main

import (
	"errors"
	"net/http"

	"github.com/sparrowsl/greenlight/internal/data"
	"github.com/sparrowsl/greenlight/internal/validator"
)

func (app *application) registerUser(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(writer, request, &input); err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Set and hash the user password
	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	v := validator.New()
	data.ValidateUser(v, user)

	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	if err := app.models.Users.Insert(user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(writer, request, v.Errors)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err := app.writeJSON(writer, http.StatusCreated, map[string]any{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
