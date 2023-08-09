package main

import (
	"html/template"
	"path/filepath"

	"github.com/dwang288/snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// Not every field needs to be filled upon instantiation,
// OK to leave them as nil if they're not being used in the template
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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
	// Add template set to in memory cache.
	for _, page := range pages {
		// Get filename from full filepath
		name := filepath.Base(page)

		// Create a slice of filepaths of our base template + partials + the page
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			page,
		}

		// Parse files into a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add template to the cache with the filename as the key
		cache[name] = ts
	}

	return cache, nil
}
