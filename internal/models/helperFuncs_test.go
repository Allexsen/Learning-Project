package models

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLastInsertId(t *testing.T) {
	q := "INSERT INTO ..."

	t.Run("Success", func(t *testing.T) {
		res := sqlmock.NewResult(1, 1)

		lastInsertId, err := getLastInsertId(res, q, struct{}{})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), lastInsertId)
	})

	t.Run("Failure", func(t *testing.T) {
		res := sqlmock.NewErrorResult(apperrors.ErrDBLastInsertId)

		lastInsertId, err := getLastInsertId(res, q, struct{}{})
		assert.Error(t, err)
		assert.Equal(t, int64(-1), lastInsertId)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBLastInsertId, appErr.Err)
	})
}

func TestHandleUpdateQuery(t *testing.T) {
	q := "UPDATE ... SET ..."
	t.Run("Success", func(t *testing.T) {
		res := sqlmock.NewResult(1, 1)

		err := handleUpdateQuery(res, nil, q, struct{}{})
		assert.NoError(t, err)
	})

	t.Run("Error in execution", func(t *testing.T) {
		res := sqlmock.NewResult(1, 0)

		execErr := errors.New("execution error")
		err := handleUpdateQuery(res, execErr, q, struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)
	})

	t.Run("Error in RowsAffected", func(t *testing.T) {
		errRes := sqlmock.NewErrorResult(apperrors.ErrDBRowsAffected)

		err := handleUpdateQuery(errRes, nil, q, struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBRowsAffected, appErr.Err)
	})

	t.Run("Not Found", func(t *testing.T) {
		res := sqlmock.NewResult(1, 0)

		err := handleUpdateQuery(res, nil, q, struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrNotFound, appErr.Err)
	})
}

func TestGetQueryError(t *testing.T) {
	q := "SELECT ... FROM ..."
	t.Run("Success", func(t *testing.T) {
		err := getQueryError(q, "No error", struct{}{}, nil)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		err := getQueryError(q, "Resource not found", struct{}{}, sql.ErrNoRows)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrDBNoRows, appErr.Err)
	})

	t.Run("Query Error", func(t *testing.T) {
		err := getQueryError(q, "Database error", struct{}{}, apperrors.ErrDBQuery)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)
	})
}
