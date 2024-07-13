package utils

import (
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(&obj); err != nil {
		apperrors.HandleError(c, apperrors.New(
			http.StatusBadRequest,
			"Invalid JSON",
			apperrors.ErrBindJson,
			map[string]interface{}{"details": err.Error()},
		))
		return false
	}

	return true
}
