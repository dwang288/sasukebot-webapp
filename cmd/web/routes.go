package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Return http.Handler type instead of *http.ServeMux so we can chain handlers
func (app *application) routes() http.Handler {
	// Initialize the httprouter
	router := httprouter.New()

	// Set httprouter's default notFound handler to our not found function
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a file server for serving static files out of a directory
	// Path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Replace all http.Servemuxes with httprouter, use clean URL pathing
	// Wrap handlers that use session data with session middleware
	router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	router.Handler(http.MethodGet, "/snippet/view/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))
	router.Handler(http.MethodGet, "/snippet/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreate)))
	router.Handler(http.MethodPost, "/snippet/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreatePost)))

	// User authentication routes
	router.Handler(http.MethodGet, "/user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignup)))
	router.Handler(http.MethodPost, "/user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignupPost)))
	router.Handler(http.MethodGet, "/user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogin)))
	router.Handler(http.MethodPost, "/user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLoginPost)))
	router.Handler(http.MethodPost, "/user/logout", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogoutPost)))

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
