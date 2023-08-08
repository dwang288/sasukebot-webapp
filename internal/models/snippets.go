package models

import (
	"database/sql"
	"time"
)

// Individual snippet data struct, matches DB table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Wrapper type for the db connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert new snippet into DB
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// Returns snippet based on ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Returns most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
