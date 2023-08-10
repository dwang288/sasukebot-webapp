package main

import (
	"errors"
	"fmt"

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

	// Grab latest 10 snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create a new templateData struct and add snippets to the struct
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Replace duplicated rendering logic. Still passing in hardcoded name
	app.render(w, http.StatusOK, "home.tmpl.html", data)
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

	// Create a new templateData struct and add the snippet to the struct
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the render helper. Still passing in hardcoded page name
	app.render(w, http.StatusOK, "view.tmpl.html", data)
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
