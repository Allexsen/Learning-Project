package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			log.Printf("empty jwt in the request header")
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Couldn't validate a JWT"})
			log.Printf("couldn't validate a jwt: %v", err)
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
