package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	// Define a command line flag with the name addr and a default value
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse value stored in flag and assign to addr. Without parsing, addr will always
	// be set to the default value. Will panic if errors occur during parsing
	flag.Parse()

	// New INFO level logger with output destination, message string prefix, and flags
	// to indicate the additional information to include (local date and time) joined with |
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// New ERROR level logger with all the INFO level logger information + filename/line number
	// of error logged and logged to stderr

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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

	// Specify and initialize a http.Server so we can use our custom errorLog.
	// Otherwise we could just use the http.ListenAndServe shortcut function.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	// Use the ListenAndServe() function on our custom http.Server
	// to start a new web server.
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
