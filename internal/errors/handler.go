// Package apperrors provides custom errors to centralize
// errors & error logging through the application.
package apperrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Using logr to avoid confusion with the built-in errors package.
var logr = logrus.New()

// init sets the logging parameters
func init() {
	logr.SetFormatter(&logrus.JSONFormatter{})
	logr.SetLevel(logrus.InfoLevel)
}

// HandleError logs the error, and writes response header.
// Changes error message if Code = 500.
func HandleError(c *gin.Context, err *AppError) {
	logError(err)
	if c != nil {
		if err.Code == http.StatusInternalServerError {
			c.AbortWithStatusJSON(err.Code, gin.H{
				"success": false,
				"error":   "Something went wrong, try again later",
			})
		} else {
			c.AbortWithStatusJSON(err.Code, gin.H{
				"success": false,
				"error":   err.Message,
			})
		}
	}
}

// HandleCriticalError logs the error, and then panics.
func HandleCriticalError(err *AppError) {
	logError(err)
	panic(err)
}

// logError logs the error to the console.
func logError(err *AppError) {
	logr.WithFields(logrus.Fields{
		"code":    err.Code,
		"message": err.Message,
		"error":   err.Err,
		"context": err.Context,
	}).Error("An error occurred")
}
