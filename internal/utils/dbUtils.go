// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"database/sql"
	"net/http"

	database "github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

func GetPasswordHashByUsername(db *sql.DB, username string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE username=?`
	var pswdHash string
	err := db.QueryRow(q, username).Scan(&pswdHash)
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

func GetPasswordHashByEmail(db *sql.DB, email string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE email=?`
	var pswdHash string
	err := db.QueryRow(q, email).Scan(&pswdHash)
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
func IsExistingEmail(db *sql.DB, email string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE email=?)`
	err := db.QueryRow(q, email).Scan(&exists)
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
func IsExistingUsername(db *sql.DB, username string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE username=?)`
	err := db.QueryRow(q, username).Scan(&exists)
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
	db := database.DB
	if email != "" {
		if exists, err := IsExistingEmail(db, email); err != nil || exists {
			return true, err
		}
	}

	if username != "" {
		if exists, err := IsExistingUsername(db, username); err != nil || exists {
			return true, err
		}
	}

	return false, nil
}
