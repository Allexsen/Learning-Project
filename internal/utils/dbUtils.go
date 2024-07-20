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
	pswdHash := ""
	err := db.QueryRow(q, username).Scan(&pswdHash)
	return pswdHash, getQueryError(q, map[string]interface{}{"username": username}, err)
}

func GetPasswordHashByEmail(db *sql.DB, email string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE email=?`
	pswdHash := ""
	err := db.QueryRow(q, email).Scan(&pswdHash)
	return pswdHash, getQueryError(q, map[string]interface{}{"email": email}, err)
}

// IsExistingEmail checks if the email is present in the db
func IsExistingEmail(db *sql.DB, email string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE email=?)`
	err := db.QueryRow(q, email).Scan(&exists)
	return exists, getQueryError(q, map[string]interface{}{"email": email}, err)
}

// IsExistingEmail checks if the username is present in the db
func IsExistingUsername(db *sql.DB, username string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE username=?)`
	err := db.QueryRow(q, username).Scan(&exists)
	return exists, getQueryError(q, map[string]interface{}{"username": username}, err)
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

func getQueryError(q string, context map[string]interface{}, err error) error {
	if err == nil {
		return nil
	}

	context["query"] = q
	context["error"] = err.Error()
	return apperrors.New(
		http.StatusInternalServerError,
		"Couldn't scan a row",
		apperrors.ErrDBQuery,
		context,
	)
}
