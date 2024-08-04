// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"log"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

// ShouldBindJSON acts just as c.ShouldBindJSON,
// except it also provides a centralized error handling.
func ShouldBindJSON(c *gin.Context, obj interface{}) bool {
	log.Printf("[UTILS] Binding JSON for request from %s", c.ClientIP())

	if err := c.ShouldBindJSON(&obj); err != nil {
		apperrors.HandleError(c, apperrors.New(
			http.StatusBadRequest,
			"Invalid JSON",
			apperrors.ErrBindJson,
			map[string]interface{}{"details": err.Error()},
		))
		return false
	}

	log.Printf("[UTILS] JSON bound successfully: %+v", obj)
	return true
}
