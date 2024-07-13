package customErrors

import "errors"

// Custom error types for centralized error handling
var (
	// General
	ErrNotFound            = errors.New("resource not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")

	// User-related
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotAuthorized  = errors.New("user not authorized")

	// Record-related
	ErrRecordNotFound    = errors.New("record not found")
	ErrInvalidRecordData = errors.New("invalid record data")

	// DB-related
	ErrDBConnection      = errors.New("database connection error")
	ErrDBQueryFailed     = errors.New("database query failed")
	ErrTransactionFailed = errors.New("transaction failed")

	// Validation-related
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")

	// Auth-related
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
)
