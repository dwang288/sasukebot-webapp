package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

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

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Removing check for post

	// Adds data in POST request bodies to the r.PostForm map
	// Function also works for PUT and PATCH
	err := r.ParseForm()
	if err != nil {
		// Client is notified of any bad requests
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Grab body data
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Map for holding form validation errors
	fieldErrors := make(map[string]string)

	// Check title value is not more than 100 character long or blank
	// Use RuneCountInString() function to count characters instead of len() to
	// count bytes. Characters with umlauts for example are unicode characters
	// that take 2 bytes instead of 1.
	if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	}

	// Check expires value is a valid option
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1, 7, or 365"
	}

	// If there are field errors, return map in a plain text HTTP response and
	// return early
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

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
