package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/sparrowsl/greenlight/internal/data"
	"github.com/sparrowsl/greenlight/internal/validator"
)

func (app *application) createAuthenticationToken(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	val := validator.New()
	data.ValidateEmail(val, input.Email)
	data.ValidatePasswordPlaintext(val, input.Password)

	if !val.Valid() {
		app.failedValidationResponse(writer, request, val.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(writer, request)
		return
	}

	token, err := app.models.Tokens.New(user.ID, time.Hour*24, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusCreated, map[string]any{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
