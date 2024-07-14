// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"net/http"
	"strconv"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
)

// Atoi acts just as strconv.Atoi,
// except it also provides a centrilized error handling
func Atoi(s string) (int, error) {
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0, apperrors.New(
			http.StatusInternalServerError,
			"Failed to convert ASCII to int",
			apperrors.ErrTypeConversion,
			map[string]interface{}{"detail": err.Error()},
		)
	}

	return res, nil
}
