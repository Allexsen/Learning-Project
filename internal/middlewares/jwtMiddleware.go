// Package middlewares provides middleware functions for API request validations.
package middlewares

import (
	"log"
	"net/http"
	"strings"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// CheckJWT middleware reads "Authorization" header and validates JWT if present.
func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[MIDDLEWARE] Checking JWT for %s", c.ClientIP())
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
			if token == "" {
				apperrors.HandleError(c, apperrors.New(
					http.StatusUnauthorized,
					"Authorization header required",
					apperrors.ErrMissingAuthHeader,
					map[string]interface{}{"details": "empty jwt in the request header"},
				))
				return
			}
		}

		tokenString := strings.TrimSpace(strings.Replace(token, "Bearer", "", 1))

		log.Printf("[MIDDLEWARE] JWT: %s", tokenString)

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

		c.Set("userDTO", &claims.UserDTO)
		log.Printf("[MIDDLEWARE] JWT Validated: %+v", claims.UserDTO)
		c.Next()
	}
}
