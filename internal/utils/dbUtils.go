// Package utils provides useful utilities for common functions throughout the app
package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

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

// getLastInsertId handles sql insertion query result,
// and returns the last insert id.
// Returns -1 and error in case of failure.
func GetLastInsertId(result sql.Result, q string, data interface{}) (int64, error) {
	dataType := reflect.TypeOf(data).Name()
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return -1, apperrors.New(
			http.StatusInternalServerError,
			"Couldn't retrieve last insert ID",
			apperrors.ErrDBLastInsertId,
			map[string]interface{}{"query": q, dataType: data, "error": err},
		)
	}

	return lastInsertId, nil
}

// getQueryError checks if there was an error in execution.
// Returns AppError if found any, or nil otherwise.
func GetQueryError(q, message string, data interface{}, err error) error {
	if err == nil {
		return nil
	}

	code, errType := http.StatusInternalServerError, apperrors.ErrDBQuery
	if err == sql.ErrNoRows { // no rows == resource was not found
		code = http.StatusNotFound
		errType = apperrors.ErrDBNoRows
	}

	dataType := reflect.TypeOf(data).Name()
	return apperrors.New(
		code,
		message,
		errType,
		map[string]interface{}{"query": q, dataType: data, "details": err.Error()},
	)
}

// handleUpdateQuery validates sql update query by checking
// update error, affected rows error, and no rows affected error.
func HandleUpdateQuery(result sql.Result, err error, q string, data interface{}) error {
	dataType := reflect.TypeOf(data).Name()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Couldn't alter the %s info", dataType),
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, dataType: data, "details": err.Error()},
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			"Couldn't retrieve rows affected",
			apperrors.ErrDBRowsAffected,
			map[string]interface{}{"query": q, dataType: data, "details": err.Error()},
		)
	}

	if rowsAffected == 0 { // no rows affected == resource was not found
		return apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("No %s found with the given ID", dataType),
			apperrors.ErrNotFound,
			map[string]interface{}{"query": q, dataType: data},
		)
	}

	return nil
}
