package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/dwang288/snippetbox/internal/models"
)

// Global var to hold the functions we want to pass into our templates
// String to function lookup of our functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// Not every field needs to be filled upon instantiation,
// OK to leave them as nil if they're not being used in the template
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	User        *models.User
	// Form for any default form data
	Form            any
	Flash           string
	IsAuthenticated bool // Mark if the current user is authenticated
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize template cache
	cache := map[string]*template.Template{}

	// Glob grabs all filepaths that match the pattern and sticks them
	// in a slice of strings. Grab all the template filepaths
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Iterate through every page and turn it into a template set.
	// Chain ParseFiles through base template + partials + page to create a template set.
	// Add completed template set to in memory cache.
	for _, page := range pages {
		// Get filename from full filepath
		name := filepath.Base(page)

		// Create a new template set so we can chain .Funcs immediately afterwards
		// .Funcs registers the functions in the FuncMap to the template set
		// Needs to be called first in the chain
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse any partials on the same template set
		// Will add to the template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse the template for this page and add to the same template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add template to the cache with the filename as the key
		cache[name] = ts
	}

	return cache, nil
}

// Formats t as a human readable time string
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
