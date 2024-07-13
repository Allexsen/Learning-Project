package utils

import (
	"net/http"
	"os"
	"regexp"
	"time"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
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

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func IsValidUsername(username string) bool {
	if len(username) <= 3 {
		return false
	}

	usernameRegex := `^[a-zA-Z0-9-_]+$`
	re := regexp.MustCompile(usernameRegex)
	return re.MatchString(username)
}
