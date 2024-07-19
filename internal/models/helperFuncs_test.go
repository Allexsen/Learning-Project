package models

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/stretchr/testify/assert"
)

type sqlMockResult struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (r sqlMockResult) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

func (r sqlMockResult) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

func TestGetLastInsertId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sqlMockResult := sqlMockResult{
			lastInsertId: 1,
			err:          nil,
		}

		lastInsertId, err := getLastInsertId(sqlMockResult, "INSERT INTO ...", struct{}{})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), lastInsertId)
	})

	t.Run("Failure", func(t *testing.T) {
		sqlMockResult := sqlMockResult{
			lastInsertId: -1,
			err:          errors.New("some error"),
		}

		lastInsertId, err := getLastInsertId(sqlMockResult, "INSERT INTO ...", struct{}{})
		assert.Error(t, err)
		assert.Equal(t, int64(-1), lastInsertId)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, "Couldn't retrieve last insert ID", appErr.Message)
		assert.Equal(t, "some error", appErr.Err.Error())
		assert.Equal(t, "INSERT INTO ...", appErr.Context["query"])
		assert.Equal(t, nil, appErr.Context["dataType"]) // type assert returns nil for unnamed types, don't expect "struct{}"
	})
}

func TestHandleUpdateQuery(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sqlMockResult := sqlMockResult{
			rowsAffected: 1,
			err:          nil,
		}

		err := handleUpdateQuery(sqlMockResult, nil, "UPDATE practice_db.users SET ...", struct{}{})
		assert.NoError(t, err)
	})

	t.Run("Error in execution", func(t *testing.T) {
		sqlMockResult := sqlMockResult{
			rowsAffected: 0,
			err:          nil,
		}

		execErr := errors.New("execution error")
		err := handleUpdateQuery(sqlMockResult, execErr, "UPDATE practice_db.users SET ...", struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)

		// type assert returns nil for unnamed types, causing extra space
		// therefore it's "the  info" and not "the struct{} info"
		assert.Equal(t, "Couldn't alter the  info", appErr.Message)

		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)
		assert.Equal(t, "UPDATE practice_db.users SET ...", appErr.Context["query"])
		assert.Equal(t, nil, appErr.Context["dataType"]) // type assert returns nil for unnamed types, don't expect "struct{}"
		assert.Equal(t, execErr, appErr.Context["error"])
	})

	t.Run("Error in RowsAffected", func(t *testing.T) {
		sqlMockResult := sqlMockResult{
			rowsAffected: 0,
			err:          errors.New("rows affected error"),
		}

		err := handleUpdateQuery(sqlMockResult, nil, "UPDATE practice_db.users SET ...", struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, "Couldn't retrieve rows affected", appErr.Message)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)
		assert.Equal(t, "UPDATE practice_db.users SET ...", appErr.Context["query"])
		assert.Equal(t, nil, appErr.Context["dataType"]) // type assert returns nil for unnamed types, don't expect "struct{}"
		assert.Equal(t, sqlMockResult.err, appErr.Context["error"])
	})

	t.Run("No rows affected", func(t *testing.T) {
		sqlMockResult := sqlMockResult{
			rowsAffected: 0,
			err:          nil,
		}

		err := handleUpdateQuery(sqlMockResult, nil, "UPDATE practice_db.users SET ...", struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)

		// type assert returns nil for unnamed types, causing extra space
		// therefore it's "No  found" and not "No struct{} found"
		assert.Equal(t, "No  found with the given ID", appErr.Message)

		assert.Equal(t, apperrors.ErrNotFound, appErr.Err)
		assert.Equal(t, "UPDATE practice_db.users SET ...", appErr.Context["query"])
		assert.Equal(t, nil, appErr.Context["dataType"]) // type assert returns nil for unnamed types, don't expect "struct{}"
	})
}

func TestGetQueryError(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		err := getQueryError("SELECT * FROM table", "No error", struct{}{}, nil)
		assert.NoError(t, err)
	})

	t.Run("Resource not found", func(t *testing.T) {
		err := getQueryError("SELECT * FROM table WHERE id=?", "Resource not found", struct{}{}, sql.ErrNoRows)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.Code)
		assert.Equal(t, "Resource not found", appErr.Message)
		assert.Equal(t, apperrors.ErrNotFound, appErr.Err)
		assert.Equal(t, "SELECT * FROM table WHERE id=?", appErr.Context["query"])
		assert.Equal(t, nil, appErr.Context["dataType"]) // type assert returns nil for unnamed types, don't expect "struct{}"
		assert.Equal(t, sql.ErrNoRows.Error(), appErr.Context["details"])
	})

	t.Run("Other error", func(t *testing.T) {
		someErr := errors.New("some error")
		err := getQueryError("SELECT * FROM table WHERE id=?", "Database error", struct{}{}, someErr)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, "Database error", appErr.Message)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)
		assert.Equal(t, "SELECT * FROM table WHERE id=?", appErr.Context["query"])
		assert.Equal(t, nil, appErr.Context["dataType"]) // type assert returns nil for unnamed types, don't expect "struct{}"
		assert.Equal(t, someErr.Error(), appErr.Context["details"])
	})
}
