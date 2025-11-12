package domain

import "errors"

// Common domain errors
var (
	// User related errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrPermissionDenied   = errors.New("permission denied")

	// Database errors
	ErrDatabase          = errors.New("database error")
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateKey      = errors.New("duplicate key error")

	// Validation errors
	ErrValidation        = errors.New("validation error")
	ErrInvalidInput      = errors.New("invalid input")
)
