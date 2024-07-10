package utils

import (
	"fmt"
	"os"
	"time"

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
		return "", fmt.Errorf("couldn't get the token string: %v", err)
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
			return nil, fmt.Errorf("invalid signature: %v", err)
		}
		return nil, fmt.Errorf("couldn't parse the token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}

func IsAvailableCreds(email, username string) (bool, error) {
	exists, err := IsExistingEmail(email)
	if err != nil {
		return false, fmt.Errorf("couldn't verify email availability: %v", err)
	}

	if exists {
		return false, nil
	}

	exists, err = IsExistingUsername(username)
	if err != nil {
		return false, fmt.Errorf("couldn't verify username availability: %v", err)
	}

	return !exists, nil
}
