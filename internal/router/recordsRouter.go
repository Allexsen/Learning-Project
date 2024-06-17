package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Record struct {
	// Will get an UID later on
	Name          string
	Email         string
	HoursWorked   int
	MinutesWorked int
}

func initRecordsRouter() {
	var Records []Record
	r.POST("/records/add", func(c *gin.Context) {
		// Temporary In-Memory "database"
		name := c.PostForm("name")
		email := c.PostForm("email")
		hStr := c.PostForm("hours")
		minStr := c.PostForm("minutes")

		hours, err := strconv.Atoi(hStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid hours format"})
			return
		}

		minutes, err := strconv.Atoi(minStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid minutes format"})
			return
		}

		record := Record{name, email, hours, minutes}
		Records = append(Records, record)
	})

	r.POST("/records/retrieve", func(c *gin.Context) {
		email := c.PostForm("email")

		for _, record := range Records {
			if record.Email == email {
				c.JSON(http.StatusOK, record)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
	})
}
