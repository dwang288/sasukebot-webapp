package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

// Deal with duplicated template rending code in the handlers
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
	return &templateData{CurrentYear: time.Now().Year()}
}
