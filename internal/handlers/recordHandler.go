package handlers

import (
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RecordAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Hours   string `json:"hours"`
			Minutes string `json:"minutes"`
		}

		if err := c.ShouldBindJSON(&reqData); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		record, err := controllers.RecordAdd(reqData.Name, reqData.Email, reqData.Hours, reqData.Minutes)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, record)
	}
}
