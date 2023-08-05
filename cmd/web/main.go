package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	// Define a command line flag with the name addr and a default value
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse value stored in flag and assign to addr. Without parsing, addr will always
	// be set to the default value. Will panic if errors occur during parsing
	flag.Parse()

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
	// Dereference the flag value and pass in the addr and the servemux
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
