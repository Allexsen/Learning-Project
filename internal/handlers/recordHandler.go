package handlers

import (
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if reqData.Hours == "0" && reqData.Minutes == "0" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Hours and minutes can not both be zero"})
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

func RecordDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData struct {
			ID string `json:"id"`
		}

		if err := c.ShouldBindJSON(&reqData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		err := controllers.RecordRemove(reqData.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete the record"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
