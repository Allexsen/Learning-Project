package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
)

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

func handleUpdateQuery(res sql.Result, err error, q string, data interface{}) error {
	dataType := reflect.TypeOf(data).Name()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			fmt.Sprintf("Couldn't alter the %s info", dataType),
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, dataType: data, "error": err},
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return apperrors.New(
			http.StatusInternalServerError,
			"Couldn't retrieve rows affected",
			apperrors.ErrDBQuery,
			map[string]interface{}{"query": q, dataType: data, "error": err},
		)
	}

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

func getQueryError(q, message string, data interface{}, err error) error {
	if err == nil {
		return nil
	}

	code, errType := http.StatusInternalServerError, apperrors.ErrDBQuery
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
