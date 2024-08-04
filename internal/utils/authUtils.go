// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/dgrijalva/jwt-go"
)

// JWT key secret loaded from env variables.
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserDTO user.UserDTO
	jwt.StandardClaims
}

// GenerateJWT sets expiration date to 24 hours, generates a new JWT,
// signs and returns the token string.
func GenerateJWT(userDTO *user.UserDTO) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserDTO: *userDTO,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", apperrors.New(
			http.StatusInternalServerError,
			"Failed to generate JWT",
			apperrors.ErrInternalServerError,
			map[string]interface{}{"details": err.Error()},
		)
	}

	return tokenString, nil
}

// ValidateJWT checks the given token string by
// its expiration date, signature, and validity.
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, apperrors.New(
				http.StatusUnauthorized,
				"Invalid signature",
				apperrors.ErrInvalidJWT,
				map[string]interface{}{"detail": err.Error()},
			)
		}
		return nil, apperrors.New(
			http.StatusUnauthorized,
			"Couldn't parse the token",
			apperrors.ErrInvalidJWT,
			map[string]interface{}{"detail": err.Error()},
		)
	}

	if !token.Valid {
		return nil, apperrors.New(
			http.StatusUnauthorized,
			"Invalid token",
			apperrors.ErrInvalidJWT,
			map[string]interface{}{"detail": "The token is not valid"},
		)
	}

	return claims, nil
}

// IsValidEmail is a regex for email validation.
// It returns an error if the email is invalid, nil otherwise.
func IsValidEmail(email string) error {
	log.Println("Validating email: ", email)

	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return apperrors.New(
			http.StatusBadRequest,
			"Invalid email",
			apperrors.ErrInvalidInput,
			map[string]interface{}{"details": fmt.Sprintf("invalid email: %s", email)},
		)
	}

	return nil
}

// IsValidUsername is a regex for username validation
// It returns an error if the username is invalid, nil otherwise.
func IsValidUsername(username string) error {
	log.Printf("Validating username: %s", username)

	usernameRegex := `^[a-zA-Z0-9-_]+$`
	re := regexp.MustCompile(usernameRegex)
	if len(username) <= 3 || !re.MatchString(username) {
		return apperrors.New(
			http.StatusBadRequest,
			"Invalid username",
			apperrors.ErrInvalidInput,
			map[string]interface{}{"details": fmt.Sprintf("invalid username: %s", username)},
		)
	}

	log.Printf("Username %s is valid", username)
	return nil
}
