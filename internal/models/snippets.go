package models

import (
	"database/sql"
	"errors"
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

	// Select statement meant to be sent to DB as a prepared statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Query through the db connection pool with the statement and the id for the
	// placeholder param. Returns a pointer to a sql.Row object with the db result
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Copies values from each column in the row into the struct's values.
	// Must have same number of params as columns.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the error is a sql.ErrNoRows error (a known exception for a known valid
		// case) then return our custom error type.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Returns most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
