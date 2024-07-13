package middlewares

import (
	"net/http"
	"strings"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apperrors.HandleError(c, apperrors.New(
				http.StatusUnauthorized,
				"Authorization header required",
				apperrors.ErrMissingAuthHeader,
				map[string]interface{}{"details": "empty jwt in the request header"},
			))
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			apperrors.HandleError(c, apperrors.New(
				http.StatusUnauthorized,
				"Couldn't validate JWT",
				apperrors.ErrInvalidJWT,
				map[string]interface{}{"details": err.Error()},
			))
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
