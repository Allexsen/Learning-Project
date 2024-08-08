// Package handlers defines API endpoints handlers, parsing and validating the request
package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// TODO: Change reqData to UserDTO(or User for registration/login) and remove the struct
// TODO: Ensure every single route returns a JSON response with "success" key

// UserRegister parses & validates input, and
// sends it to controllers for registering a new user
func UserRegister(c *gin.Context) {
	log.Printf("[HANDLER] Handling registration request for %s", c.ClientIP())

	var reqData struct {
		Firstname string `json:"firstName"`
		Lastname  string `json:"lastName"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	if err := utils.IsValidEmail(reqData.Email); err != nil {
		handleError(c, err)
		return
	}

	if err := utils.IsValidName(reqData.Username); err != nil {
		handleError(c, err)
		return
	}

	// Check if the given email and/or username already exist.
	if exists, err := utils.IsExistingCreds(c, reqData.Email, reqData.Username); err != nil || exists {
		handleError(c, err)
		return
	}

	u, err := controllers.UserRegister(reqData.Firstname, reqData.Lastname, reqData.Username, reqData.Email, reqData.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] User %s has been successfully registered", u.Username)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

// UserLogin parses input, and sends data to controllers.
// If the request is successful, generates and adds JWT to headers.
func UserLogin(c *gin.Context) {
	log.Printf("[HANDLER] Handling login request for %s", c.ClientIP())

	var reqData struct {
		Cred     string `json:"cred"`
		Password string `json:"password"`
	}

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	userDTO, err := controllers.UserLogin(reqData.Cred, reqData.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	tokenString, err := utils.GenerateJWT(userDTO)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] User %s has been successfully logged in", userDTO.Username)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
	})
}

// UserGet parses input, queries controllers for retrieving
// user, and sets it as a header if successful.
func UserGet(c *gin.Context) {
	log.Printf("[HANDLER] Handling user retrieval request for %s", c.ClientIP())

	var reqData struct {
		Cred string `json:"cred"`
	}

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	var u *user.User
	var err error
	if strings.Contains(reqData.Cred, "@") {
		u, err = controllers.UserGetByEmail(reqData.Cred)
		if err != nil {
			handleError(c, err)
			return
		}
	} else {
		u, err = controllers.UserGetByUsername(reqData.Cred)
		if err != nil {
			handleError(c, err)
			return
		}
	}

	log.Printf("[HANDLER] User %s has been successfully retrieved", u.Username)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

// IsAvailableEmail parses input, and checks
// if the email is available
func IsAvailableEmail(c *gin.Context) {
	log.Printf("[HANDLER] Handling email availability check request for %s", c.ClientIP())

	var reqData struct {
		Email string `json:"email"`
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	// Querying with empty username to check email only
	if exists, err := utils.IsExistingCreds(c, reqData.Email, ""); err != nil || exists {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Email %s is available", reqData.Email)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// IsAvailableEmail parses input, and checks
// if the username is available
func IsAvailableUsername(c *gin.Context) {
	log.Printf("[HANDLER] Handling username availability check request for %s", c.ClientIP())

	var reqData struct {
		Username string `json:"username"`
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	// Querying with empty email to check username only
	if exists, err := utils.IsExistingCreds(c, "", reqData.Username); err != nil || exists {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Username %s is available", reqData.Username)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
