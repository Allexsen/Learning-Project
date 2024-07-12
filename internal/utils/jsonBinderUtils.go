package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Printf("couldn't bind json: %v", err)
		return false
	}

	return true
}
