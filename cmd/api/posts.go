package main

import (
	"errors"
	"net/http"

	"github.com/vj-2303/gist-go/internal/data"
	"github.com/vj-2303/gist-go/internal/validator"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title    string `json:"title"`
		Language string `json:"language"`
		Code     string `json:"code"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	post := &data.Post{
		Title:    input.Title,
		Language: input.Language,
		Code:     input.Code,
		UserID:   user.ID,
	}
	v := validator.New()

	if data.ValidatePosts(v, post); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Post.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readParamsID(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)

	post, err := app.models.Post.GetByID(id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrPostNotFound) {
			app.notFoundResponse(w, r)
			return
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	posts, err := app.models.Post.GetAllByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"posts": posts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
