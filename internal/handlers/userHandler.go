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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		}

		c.JSON(http.StatusNotFound, u)
	}
}
