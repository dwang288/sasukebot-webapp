package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Return http.Handler type instead of *http.ServeMux so we can chain handlers
func (app *application) routes() http.Handler {
	// Initialize the httprouter
	router := httprouter.New()

	// Create a file server for serving static files out of a directory
	// Path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Replace all http.Servemuxes with httprouter, use clean URL pathing
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
