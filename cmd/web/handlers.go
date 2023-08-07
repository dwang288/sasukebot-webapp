package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Undefined paths no longer get routed to /
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Template path slice. Base template must be first in the slice.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Read template file into a template set.
	// If error is present, log error msg and return a generic 500
	// Path either needs to be an absolute path or relative to your current working diretory
	// Passing in template file path slice as variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Write the specified template in the set into response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
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

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
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
