// Package common provides helper functions for models package.
package common

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
)

// GetLastInsertId handles sql insertion query result,
// and returns the last insert id.
// Returns -1 and error in case of failure.
func GetLastInsertId(result sql.Result, q string, data interface{}) (int64, error) {
	dataType := reflect.TypeOf(data).Name()
	log.Printf("[COMMON] Retrieving last insert ID for %s", dataType)
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

// GetQueryError checks if there was an error in execution.
// Returns AppError if found any, or nil otherwise.
func GetQueryError(q, message string, data interface{}, err error) error {
	log.Printf("[COMMON] Checking query error for %s", q)
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

// HandleUpdateQuery validates sql update query by checking
// update error, affected rows error, or no rows affected error.
func HandleUpdateQuery(result sql.Result, err error, q string, data interface{}) error {
	log.Printf("[COMMON] Handling update query for %s", q)
	dataType := reflect.TypeOf(data).Name()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Couldn't alter the %s info", dataType),
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, dataType: data, "details": err.Error()},
		)
	}

	log.Printf("[COMMON] Checking rows affected for %s", dataType)
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

	log.Printf("[COMMON] Rows affected for %s: %d", dataType, rowsAffected)
	return nil
}
