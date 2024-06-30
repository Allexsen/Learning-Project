package handlers

import (
	"fmt"
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

		if reqData.Hours == "0" && reqData.Minutes == "0" {
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Errorf("hours and minutes can not both be zero"))
			return
		}

		u, err := controllers.RecordAdd(reqData.Name, reqData.Email, reqData.Hours, reqData.Minutes)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"user":    u,
		})
	}
}
