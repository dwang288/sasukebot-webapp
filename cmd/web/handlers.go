package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// Undefined paths no longer get routed to /
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the id parameter from the query and turn the string into and int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// If it cannot be converted or is out of the expected range then return 404
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Write id value into response
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not
	if r.Method != http.MethodPost {
		// Send 405 response if it isn't a POST
		// Also let client know which REST methods are allowed by adding a header to
		// the response header map
		w.Header().Set("Allow", http.MethodPost)
		// Use provided Error function to return the correct error code and message
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
