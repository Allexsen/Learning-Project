package handlers

import (
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/Allexsen/Learning-Project/internal/models"
	"github.com/gin-gonic/gin"
)

func UserGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u models.User
		u.Email = c.PostForm("email")

		u, err := controllers.GetUserByEmail(u.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, u)
	}
}
