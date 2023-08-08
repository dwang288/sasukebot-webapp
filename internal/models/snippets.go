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

// Insert new snippet into DB and return its ID in the db
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Insert SQL statement, use ? as placeholder to prevent SQL injections instead of
	// interpolating values into the string
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Execute the statement along with variables for placeholders
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted record from the returned result
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Convert returned result from int64 to int
	return int(id), nil
}

// Returns snippet based on ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Returns most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
