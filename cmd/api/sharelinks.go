package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vj-2303/gist-go/internal/data"
)

func (app *application) shareLinkCreateHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readParamsID(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)
	post, err := app.models.Post.GetByID(id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrPostNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	token, err := app.GenerateURLSafeString(8)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	shareLink := &data.ShareLink{
		PostID: post.ID,
		Token:  token,
	}
	err = app.models.ShareLink.Insert(shareLink)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	sharableLink := "localhost:4000/share/" + shareLink.Token
	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post, "Sharable_link": sharableLink}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) viewSharePostHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	token := params.ByName("shareToken")

	if token == "" || len(token) != 8 {
		app.notFoundResponse(w, r)
		return
	}
	shareToken, err := app.models.ShareLink.GetByShareToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	post, err := app.models.Post.GetByIDOnly(shareToken.PostID)
	if err != nil {
		if errors.Is(err, data.ErrPostNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
