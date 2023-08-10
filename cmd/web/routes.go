package main

import "net/http"

// Return http.Handler type instead of *http.ServeMux so we can chain handlers
func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux
	// Register the home function as the handler for the "/" URL pattern
	mux := http.NewServeMux()

	// Create a file server for serving static files out of a directory
	// Path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Register FileServer as the handler for URL paths that start with /static/
	// Strip /static prefix from URL path before processing request
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Wrap security headers middleware with request logging function
	// Wrap servemux with middleware that adds security header
	// Pass the servemux in as the next handler to be called
	return app.logRequest(secureHeaders(mux))
}
