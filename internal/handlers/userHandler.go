package handlers

import (
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
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

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Printf("couldn't bind json: %v", err)
		return
	}

	exists, err := utils.IsExistingEmail(reqData.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking email availability"})
		log.Printf("couldn't check email availability: %v", err)
		return
	}

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "The email already exists"})
		log.Printf("couldn't register a user, existing email: %v", reqData.Email)
		return
	}

	exists, err = utils.IsExistingUsername(reqData.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking username availability"})
		log.Printf("couldn't check username availability: %v", err)
		return
	}

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		log.Printf("couldn't register a user, existing username: %v", reqData.Username)
		return
	}

	u, err := controllers.UserRegister(reqData.Firstname, reqData.Lastname, reqData.Username, reqData.Email, reqData.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't register a new user"})
		log.Print(err)
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

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Printf("couldn't bind json: %v", err)
	}

	if err := controllers.UserLogin(reqData.Email, reqData.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		log.Printf("couldn't authorize a user due to invalid credentials: %v", err)
		return
	}

	tokenString, err := utils.GenerateJWT(reqData.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't generate a JWT"})
		log.Printf("failed to generate a token string: %v", err)
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

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Printf("couldn't bind json: %v", err)
		return
	}

	u := models.User{Email: reqData.Email}
	u, err := controllers.UserGetByEmail(u.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Printf("couldn't retrieve a user: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

func IsAvailableCreds(c *gin.Context) {
	var reqData struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Printf("couldn't bind json: %v", err)
	}

	exists, err := utils.IsExistingEmail(reqData.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking email availability"})
		log.Printf("couldn't check email availability: %v", err)
		return
	}

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "The email already exists"})
		log.Printf("couldn't register a user, existing email: %v", reqData.Email)
		return
	}

	exists, err = utils.IsExistingUsername(reqData.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking username availability"})
		log.Printf("couldn't check username availability: %v", err)
		return
	}

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		log.Printf("couldn't register a user, existing username: %v", reqData.Username)
		return
	}
}
