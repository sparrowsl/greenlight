package main

import (
	"errors"
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(user.ID, time.Hour*24*3, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	app.background(func() {
		// send the welcome email
		err := app.mailer.Send(user.Email, "user_welcome.html", map[string]any{"activationToken": token.PlainText, "userID": user.ID})
		if err != nil {
			app.logger.Println(err)
		}
	})

	err = app.writeJSON(writer, http.StatusAccepted, map[string]any{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

func (app *application) activateUserToken(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Token string `json:"token"`
	}

	if err := app.readJSON(writer, request, &input); err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	v := validator.New()
	data.ValidateTokenPlainText(v, input.Token)

	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	// app.models.Tokens

	app.logger.Println(input)

	err := app.writeJSON(writer, http.StatusOK, map[string]any{"token": input.Token}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
