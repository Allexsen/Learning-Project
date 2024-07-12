package utils

import (
	"github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// to be changed to "log" once the full implementation is done and "log" package is not necessary anymore
var logr = logrus.New()

func init() {
	logr.SetFormatter(&logrus.JSONFormatter{})
	logr.SetLevel(logrus.InfoLevel)
}

func HandleError(c *gin.Context, err *errors.CustomError) {
	logError(err)
	c.AbortWithStatusJSON(err.Code, err.Message)
}

func logError(err *errors.CustomError) {
	logr.WithFields(logrus.Fields{
		"code":    err.Code,
		"message": err.Message,
		"error":   err.Err,
		"context": err.Context,
	})
}
