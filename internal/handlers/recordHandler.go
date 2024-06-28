package handlers

import (
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RecordAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		hStr := c.PostForm("hours")
		minStr := c.PostForm("minutes")

		record, err := controllers.RecordAdd(name, email, hStr, minStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "couldn't add the record"})
			return
		}

		c.JSON(http.StatusOK, record)
	}
}
