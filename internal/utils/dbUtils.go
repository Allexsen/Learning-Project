// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

func GetPasswordHashByUsername(username string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE username=?`
	var pswdHash string
	err := db.DB.QueryRow(q, username).Scan(&pswdHash)
	if err != nil {
		return "", apperrors.New(
			http.StatusInternalServerError,
			"Failed scanning a row",
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, "username": username, "details": err.Error()},
		)
	}

	return pswdHash, nil
}

func GetPasswordHashByEmail(email string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE email=?`
	var pswdHash string
	err := db.DB.QueryRow(q, email).Scan(&pswdHash)
	if err != nil {
		return "", apperrors.New(
			http.StatusInternalServerError,
			"Failed scanning a row",
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, "email": email, "details": err.Error()},
		)
	}

	return pswdHash, nil
}

// IsExistingEmail checks if the email is present in the db
func IsExistingEmail(email string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE email=?)`
	err := db.DB.QueryRow(q, email).Scan(&exists)
	if err != nil {
		return false, apperrors.New(
			http.StatusInternalServerError,
			"Failed scanning a row",
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, "email": email, "details": err.Error()},
		)
	}

	return exists, nil
}

// IsExistingEmail checks if the username is present in the db
func IsExistingUsername(username string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE username=?)`
	err := db.DB.QueryRow(q, username).Scan(&exists)
	if err != nil {
		return false, apperrors.New(
			http.StatusInternalServerError,
			"Failed scanning a row",
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, "username": username, "details": err.Error()},
		)
	}

	return exists, nil
}

// IsExistingCreds checks if the email and/or username are present in the db
func IsExistingCreds(c *gin.Context, email, username string) (bool, error) {
	if email != "" {
		if exists, err := IsExistingEmail(email); err != nil || exists {
			return true, err
		}
	}

	if username != "" {
		if exists, err := IsExistingUsername(username); err != nil || exists {
			return true, err
		}
	}

	return false, nil
}
