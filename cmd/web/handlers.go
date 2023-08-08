package main

import (
	"errors"
	"fmt"

	// Temporarily return just the plain text Latest value for this commit
	// "html/template"
	"net/http"
	"strconv"

	"github.com/dwang288/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Undefined paths no longer get routed to /
	if r.URL.Path != "/" {
		// Use the notFound() helper
		app.notFound(w)
		return
	}

	// Temporarily display snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// // Template path slice. Base template must be first in the slice.
	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }

	// // Read template file into a template set.
	// // If error is present, log error msg and return a generic 500
	// // Path either needs to be an absolute path or relative to your current working diretory
	// // Passing in template file path slice as variadic parameter
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err) // Use the serverError() handler
	// 	return
	// }
	// // Write the specified template in the set into response body
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err) // Use the serverError() handler
	// 	return
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the id parameter from the query and turn the string into and int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// If it cannot be converted or is out of the expected range then return 404
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper
		return
	}

	// Retrieve the snippet data from the db with its id. If no record is found,
	// return a 404. If it's some other error, throw a 500.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Write snippet data as a plain text HTTP response body
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not
	if r.Method != http.MethodPost {
		// Send 405 response if it isn't a POST
		// Also let client know which REST methods are allowed by adding a header to
		// the response header map
		w.Header().Set("Allow", http.MethodPost)
		// Return the correct error code and message
		app.clientError(w, http.StatusMethodNotAllowed) // use the clientError() helper
		return
	}

	// Example dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass data to Insert method and receive ID of the inserted method back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect user to the new snippet's view page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
