package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	// Generate a bcrypt hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}

	statement := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(statement, name, email, string(hashedPassword))

	if err != nil {
		// Check if this has the type *mysql.MySQLError. If yes, then check the
		// error code to see if the value is a duplicate value on a unique column.
		// If the error message contains "users_uc_email", we know that it's this
		// column that has a duplicate value
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
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
