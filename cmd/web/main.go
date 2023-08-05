package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a file server for serving static files out of a directory
	// Path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the http.NewServeMux() function to initialize a new servemux
	// Register the home function as the handler for the "/" URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Register FileServer as the handler for URL paths that start with /static/
	// Strip /static prefix from URL path before processing request
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Use the http.ListenAndServe() function to start a new web server
	// Pass in the port and the servemux
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
