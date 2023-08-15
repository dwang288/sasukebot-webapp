package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/dwang288/snippetbox/internal/models"
	"github.com/dwang288/snippetbox/internal/validator"

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

	// Initialize data.Form along with any default form values
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// Struct for holding form data
// Fields are exported on purposes because html/template needs them to be
// exported to be read
// Struct tags for mapping HTML form values to struct fields
// `form:"-"` tells decoder to ignore that field
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// Embeds the validator so that snippetCreateForm can use all the fields
	// and methods of the Validator type
	// Validator type contains the FieldsError field so we can access it the same way
	// as we did before
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Adds data in POST request bodies to the r.PostForm map
	// Function also works for PUT and PATCH
	err := r.ParseForm()
	if err != nil {
		// Client is notified of any bad requests
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create form struct and decode parsed content into the struct
	var form snippetCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// If the check is false, then will add the error info the the form errors
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	// Instead of only checking the length, use our Valid method to see if the form is valid
	// If there's form errors, refill form data + reload and send a 422
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	// Pass data to Insert method and receive ID of the inserted method back
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If snippet is successfully added to DB, add key value pair to session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	// Redirect user to the new snippet's view page
	// Use clean URL format
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
