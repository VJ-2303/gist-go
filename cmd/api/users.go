package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/vj-2303/gist-go/internal/data"
	"github.com/vj-2303/gist-go/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hash, err := app.generateHash(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user := &data.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     input.Password,
		PasswordHash: hash,
		Activated:    false,
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.User.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "email is already in USE")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateLoginUser(v, input.Email, input.Password); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrUserNotFound) {
			app.invalidCredentialsResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	match := app.checkPassAndHash(input.Password, user.PasswordHash)

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	token := &data.Token{
		UserID: user.ID,
		Expiry: 24 * time.Hour,
	}
	err = data.GenerateNewToken(token, app.config.jwt.secret)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"auth_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
