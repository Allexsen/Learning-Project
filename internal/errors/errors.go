// Package apperrors provides custom errors to centralize
// errors & error logging through the application.
package apperrors

import (
	"fmt"
)

// AppError represents the custom error structure.
type AppError struct {
	Code    int                    // Code sent to the client
	Message string                 // Message sent to the client. Ignored if Code = 500.
	Err     error                  // Error message logged to the console for debugging
	Context map[string]interface{} // Any additional context to aid debugging and provide useful info
}

// AppError.Error() implements the built-in error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, err: %v, context: %v", e.Code, e.Message, e.Err, e.Context)
}

// New() creates a new AppError and returns a reference to it.
func New(code int, message string, err error, context map[string]interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Context: context,
	}
}
