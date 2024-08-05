// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"database/sql"
	"log"
	"net/http"

	database "github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/gin-gonic/gin"
)

// GetUserIDByUsername retrieves the user ID by username
func GetPasswordHashByUsername(db *sql.DB, username string) (string, error) {
	log.Printf("[UTILS] Getting password hash for username %s", username)

	q := `SELECT password FROM practice_db.users WHERE username=?`
	pswdHash := ""
	err := db.QueryRow(q, username).Scan(&pswdHash)
	err = getQueryError(q, map[string]interface{}{"username": username}, err)
	if err != nil {
		return "", err
	}

	log.Printf("[UTILS] Retrieved password hash for username %s", username)
	return pswdHash, nil
}

// GetPasswordHashByEmail retrieves the password hash by email
func GetPasswordHashByEmail(db *sql.DB, email string) (string, error) {
	log.Printf("[UTILS] Getting password hash for email %s", email)

	q := `SELECT password FROM practice_db.users WHERE email=?`
	pswdHash := ""
	err := db.QueryRow(q, email).Scan(&pswdHash)
	err = getQueryError(q, map[string]interface{}{"email": email}, err)
	if err != nil {
		return "", err
	}

	log.Printf("[UTILS] Retrieved password hash for email %s", email)
	return pswdHash, nil
}

// IsExistingEmail checks if the email is present in the db
func IsExistingEmail(db *sql.DB, email string) (bool, error) {
	log.Printf("[UTILS] Checking if the email %s already exists", email)

	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE email=?)`
	err := db.QueryRow(q, email).Scan(&exists)
	err = getQueryError(q, map[string]interface{}{"email": email}, err)
	if err != nil {
		return false, err
	}

	log.Printf("[UTILS] Email %s exists: %t", email, exists)
	return exists, nil
}

// IsExistingEmail checks if the username is present in the db
func IsExistingUsername(db *sql.DB, username string) (bool, error) {
	log.Printf("[UTILS] Checking if the username %s already exists", username)

	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE username=?)`
	err := db.QueryRow(q, username).Scan(&exists)
	err = getQueryError(q, map[string]interface{}{"username": username}, err)
	if err != nil {
		return false, err
	}

	log.Printf("[UTILS] Username %s exists: %t", username, exists)
	return exists, nil
}

// IsExistingCreds checks if the email and/or username are present in the db
func IsExistingCreds(c *gin.Context, email, username string) (bool, error) {
	log.Printf("[UTILS] Checking if the email %s or username %s already exists", email, username)
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

	log.Printf("[UTILS] Email %s and username %s are unique", email, username)
	return false, nil
}

// getQueryError returns an error if the query execution failed, nil otherwise
func getQueryError(q string, context map[string]interface{}, err error) error {
	log.Printf("[UTILS] Executing query: %s", q)

	if err != nil {
		context["query"] = q
		context["error"] = err.Error()

		return apperrors.New(
			http.StatusInternalServerError,
			"Couldn't scan a row",
			apperrors.ErrDBQuery,
			context,
		)
	}

	log.Printf("[UTILS] Query executed successfully")
	return nil
}
