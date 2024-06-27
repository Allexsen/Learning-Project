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
		var reqData struct {
			Email string
		}

		if err := c.ShouldBindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		u.Email = reqData.Email
		u, err := controllers.GetUserByEmail(u.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"user":    u,
		})
	}
}
