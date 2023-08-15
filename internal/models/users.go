package models

import (
	"database/sql"
	"time"
)

// User type with fields and types that match the columns in the users table
// Will copy DB data into this model
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel type that wraps a DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert creates a new record in the Users table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate checks if a user exists with this email/password combo and returns
// the user ID if they do
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user with this ID exists.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
