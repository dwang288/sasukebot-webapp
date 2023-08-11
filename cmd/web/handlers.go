package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/dwang288/snippetbox/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Remove exact base URL check for "/" since httprouter does exact matches

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
	// Grab named parameters from request with ParamsFromContext(r.Context())
	params := httprouter.ParamsFromContext(r.Context())

	// Extract the id parameter from the slice and turn the string into an int
	id, err := strconv.Atoi(params.ByName("id"))

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

// Placeholder function
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Removing check for post

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
	// Use clean URL format
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
