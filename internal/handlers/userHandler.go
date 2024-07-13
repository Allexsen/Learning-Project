package handlers

import (
	"fmt"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var reqData struct {
		Firstname string `json:"firstName"`
		Lastname  string `json:"lastName"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

	if !utils.IsValidEmail(reqData.Email) {
		apperrors.HandleError(c, apperrors.New(
			http.StatusBadRequest,
			"Invalid email",
			apperrors.ErrInvalidInput,
			map[string]interface{}{"details": fmt.Sprintf("invalid email: %s", reqData.Email)},
		))
		return
	}

	if !utils.IsValidUsername(reqData.Username) {
		apperrors.HandleError(c, apperrors.New(
			http.StatusBadRequest,
			"Invalid username",
			apperrors.ErrInvalidInput,
			map[string]interface{}{"details": fmt.Sprintf("invalid username: %s", reqData.Username)},
		))
		return
	}

	if exists, err := utils.IsExistingCreds(c, reqData.Email, reqData.Username); err != nil || exists {
		handleError(c, err)
		return
	}

	u, err := controllers.UserRegister(reqData.Firstname, reqData.Lastname, reqData.Username, reqData.Email, reqData.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

func UserLogin(c *gin.Context) {
	var reqData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

	if err := controllers.UserLogin(reqData.Email, reqData.Password); err != nil {
		handleError(c, err)
		return
	}

	tokenString, err := utils.GenerateJWT(reqData.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
	})
}

func UserGet(c *gin.Context) {
	var reqData struct {
		Email string `json:"email"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

	u := models.User{Email: reqData.Email}
	u, err := controllers.UserGetByEmail(u.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

func IsAvailableEmail(c *gin.Context) {
	var reqData struct {
		Email string `json:"email"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

	if exists, err := utils.IsExistingCreds(c, reqData.Email, ""); err != nil || exists {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func IsAvailableUsername(c *gin.Context) {
	var reqData struct {
		Username string `json:"username"`
	}

	if !utils.BindJSON(c, &reqData) {
		return
	}

	if exists, err := utils.IsExistingCreds(c, "", reqData.Username); err != nil || exists {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
