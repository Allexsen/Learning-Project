package handlers

import (
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RecordAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("Hit handler")

		name := c.PostForm("name")
		email := c.PostForm("email")
		hStr := c.PostForm("hours")
		minStr := c.PostForm("minutes")

		record, err := controllers.RecordAdd(name, email, hStr, minStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't add the record"})
		}

		c.JSON(http.StatusOK, record)
	}
}
