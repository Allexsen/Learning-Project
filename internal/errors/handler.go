package apperrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// using logr to differentiate from the built-in log package
var logr = logrus.New()

func init() {
	logr.SetFormatter(&logrus.JSONFormatter{})
	logr.SetLevel(logrus.InfoLevel)
}

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

func HandleCriticalError(err *AppError) {
	logError(err)
	panic(err)
}

func logError(err *AppError) {
	logr.WithFields(logrus.Fields{
		"code":    err.Code,
		"message": err.Message,
		"error":   err.Err,
		"context": err.Context,
	}).Error("An error occurred")
}
