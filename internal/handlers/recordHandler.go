// Package handlers defines API endpoints handlers, parsing and validating requests
package handlers

import (
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// RecordAdd parses & validates input, and
// sends it to controllers for adding a record.
func RecordAdd(c *gin.Context) {
	log.Printf("[HANDLER] Handling record addition request for %s", c.ClientIP())

	var reqData struct {
		Email   string `json:"email"`
		Hours   string `json:"hours"`
		Minutes string `json:"minutes"`
	}

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	// hours=minutes=0 is basically an empty record, making it invalid
	if reqData.Hours == "0" && reqData.Minutes == "0" {
		apperrors.HandleError(c, apperrors.New(
			http.StatusBadRequest,
			"Hours and minutes can not both be zero",
			apperrors.ErrInvalidInput,
			map[string]interface{}{"details": "Hours and minutes can not both be zero"},
		))
		return
	}

	u, err := controllers.RecordAdd(reqData.Email, reqData.Hours, reqData.Minutes)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Record has been successfully added for %s", u.Email)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    &u,
	})
}

// RecordDelete parses input, and sends data
// to controllers for deleting a record
func RecordDelete(c *gin.Context) {
	log.Printf("[HANDLER] Handling record deletion request for %s", c.ClientIP())

	var reqData struct {
		ID int `json:"id"`
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	u, err := controllers.RecordRemove(reqData.ID)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Record has been successfully deleted for %s", u.Email)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    &u,
	})
}
