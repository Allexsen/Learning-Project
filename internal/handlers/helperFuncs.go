package handlers

import (
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		apperrors.HandleError(c, appErr)
	} else {
		appErr := apperrors.New(http.StatusInternalServerError, "Unknown error occurred", apperrors.ErrInternalServerError, map[string]interface{}{"detail": err.Error()})
		apperrors.HandleError(c, appErr)
	}
}
