package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens", app.loginUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/testauth", app.requiredAuthenticatedUser(app.testAuth))

	router.HandlerFunc(http.MethodPost, "/v1/posts", app.requiredAuthenticatedUser(app.createPostHandler))
	router.HandlerFunc(http.MethodGet, "/v1/posts/:id", app.requiredAuthenticatedUser(app.getPostHandler))
	router.HandlerFunc(http.MethodGet, "/v1/posts", app.requiredAuthenticatedUser(app.GetAllPostsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/share/:id/create", app.requiredAuthenticatedUser(app.shareLinkCreateHandler))
	router.HandlerFunc(http.MethodGet, "/v1/share/:shareToken", app.viewSharePostHandler)

	return router
}
