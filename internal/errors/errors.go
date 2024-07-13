package apperrors

import (
	"fmt"
)

type AppError struct {
	Code    int
	Message string
	Err     error
	Context map[string]interface{}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, err: %v, context: %v", e.Code, e.Message, e.Err, e.Context)
}

func New(code int, message string, err error, context map[string]interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Context: context,
	}
}
