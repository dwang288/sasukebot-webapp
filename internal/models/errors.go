package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add error for when user tries an incorrect email/password
	ErrInvalidCredetials = errors.New("models: invalid credentials")
	// Add error for when user tries to sign up with an existing email
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
