package handlers

import (
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

func RecordAdd(c *gin.Context) {
	var reqData struct {
		Email   string `json:"email"`
		Hours   string `json:"hours"`
		Minutes string `json:"minutes"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

func RecordDelete(c *gin.Context) {
	var reqData struct {
		ID int `json:"id"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

	u, err := controllers.RecordRemove(reqData.ID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}
