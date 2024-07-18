// Package models declares object models, and provides methods for database interaction.
package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
)

// getLastInsertId handles sql insertion query result,
// and returns the last insert id.
// Returns -1 and error in case of failure.
func getLastInsertId(result sql.Result, q string, data interface{}) (int64, error) {
	dataType := reflect.TypeOf(data).Name()
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return -1, apperrors.New(
			http.StatusInternalServerError,
			"Couldn't retrieve last insert ID",
			err,
			map[string]interface{}{"query": q, dataType: data},
		)
	}

	return lastInsertId, nil
}

// handleUpdateQuery validates sql update query by checking
// update error, affected rows error, or no rows affected error.
func handleUpdateQuery(result sql.Result, err error, q string, data interface{}) error {
	dataType := reflect.TypeOf(data).Name()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Couldn't alter the %s info", dataType),
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, dataType: data, "error": err},
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			"Couldn't retrieve rows affected",
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, dataType: data, "error": err},
		)
	}

	// If no rows were affected, then the requested resource was not found
	if rowsAffected == 0 {
		return apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("No %s found with the given ID", dataType),
			apperrors.ErrNotFound,
			map[string]interface{}{"query": q, dataType: data},
		)
	}

	return nil
}

// getQueryError checks if there was an error in execution.
// Returns AppError if found any, or nil otherwise.
func getQueryError(q, message string, data interface{}, err error) error {
	if err == nil {
		return nil
	}

	code, errType := http.StatusInternalServerError, apperrors.ErrDBQuery
	// If the resource was not found, return 400 and ErrNotFound.
	if err == sql.ErrNoRows {
		code = http.StatusBadRequest
		errType = apperrors.ErrNotFound
	}

	dataType := reflect.TypeOf(data).Name()
	return apperrors.New(
		code,
		message,
		errType,
		map[string]interface{}{"query": q, dataType: data, "details": err.Error()},
	)
}
