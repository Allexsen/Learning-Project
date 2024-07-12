package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("resource not found")
	// to be added...
)

type CustomError struct {
	Code    int
	Message string
	Err     error
	Context map[string]interface{}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, err: %v, context: %v", e.Code, e.Message, e.Err, e.Context)
}

func New(code int, message string, err error, context map[string]interface{}) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
		Context: context,
	}
}
