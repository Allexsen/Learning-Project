// Package handlers defines API endpoints handlers, parsing and validating the request
package handlers

import (
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

// handleError takes error, type asserts it to AppError, and sends for handling.
// If type assertion fails, creates a new unknown error and sends it instead.
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		apperrors.HandleError(c, appErr)
	} else {
		appErr := apperrors.New(http.StatusInternalServerError, "Unknown error occurred", apperrors.ErrInternalServerError, map[string]interface{}{"detail": err.Error()})
		apperrors.HandleError(c, appErr)
	}
}
