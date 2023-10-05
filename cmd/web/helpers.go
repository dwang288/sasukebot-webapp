package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// On error, logs the error trace and writes the status text for internal server error
// along with the error code to the response
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Change frame depth to 2 so we see who called this helper
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends a specific status code to the user in the cases where it's the client and not
// the server that has issues, such as sending a bad request
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Convenience wrapper function around clientError for the specific 404 not found response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Deal with duplicated template rendering code in the handlers
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	// Get template set from cache, if it doesn't exist then throw a 500
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize new buffer to test runtime errors against
	buf := new(bytes.Buffer)

	// Execute template into a buffer first to check for any errors.
	// If it errors out then throw a 500 and return early.
	// We can't error out on the real template response because we'll have sent half
	// of the template already before we hit the runtime error and throw the 500
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Respond with correct header (200/500/404 etc)
	w.WriteHeader(status)

	// Write to the response writer directly from the checked buffer
	buf.WriteTo(w)

}

// Used to initialize template data structs, always want to include the year for the footer
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
		// Add authentication status to the template data
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// Decode request body into target dst
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Parse post form body regularly
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Use the form decoder
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// Panic if we pass in an invalid target destination
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

// isAuthenticated returns true if the request is from an authenticated user
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
