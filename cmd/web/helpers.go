package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// On error, logs the error trace and writes the status text for internal server error
// along with the error code to the response
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)

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
