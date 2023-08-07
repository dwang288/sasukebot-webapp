package main

import "net/http"

// Move routing logic into its own file/function
func (app *application) routes() *http.ServeMux {
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

	return mux
}
