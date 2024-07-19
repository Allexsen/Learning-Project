// Package apperrors provides custom errors to centralize
// errors & error logging through the application.
package apperrors

import "errors"

// Custom error types for centralized error handling
var (
	// General
	ErrNotFound            = errors.New("resource not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")
	ErrTypeConversion      = errors.New("type conversion failure")

	// .env
	ErrLoadEnv = errors.New("env vars loading failure")

	// User-related
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidUserCredentials = errors.New("invalid credentials")
	ErrUserNotAuthorized      = errors.New("user not authorized")

	// Record-related
	ErrRecordNotFound    = errors.New("record not found")
	ErrInvalidRecordData = errors.New("invalid record data")

	// DB-related
	ErrDBConnection   = errors.New("database connection error")
	ErrDBQuery        = errors.New("invalid database query")
	ErrDBTransaction  = errors.New("invalid transaction")
	ErrDBLastInsertId = errors.New("sql: last insert ID failure")
	ErrDBRowsAffected = errors.New("sql: affected rows inaccsessible")
	ErrDBNoRows       = errors.New("sql: no rows in result set")

	// Validation-related
	ErrInvalidInput         = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrBindJson             = errors.New("invalid json")

	// Auth-related
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrMissingAuthHeader    = errors.New("missing authorization header")
	ErrInvalidJWT           = errors.New("invalid JWT token")
)
