package record

import (
	"database/sql"
	"net/http"
	"testing"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	// Each test uses the same record & q
	q := `INSERT INTO practice_db\.records \(user_id, hours, minutes\) VALUES \(\?, \?, \?\)`
	record := Record{
		UserID:  1,
		Hours:   2,
		Minutes: 30,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.UserID, record.Hours, record.Minutes).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := record.AddRecord(tx)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Insert Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.UserID, record.Hours, record.Minutes).
			WillReturnError(apperrors.ErrDBQuery)

		id, err := record.AddRecord(tx)
		assert.Error(t, err)
		assert.Equal(t, int64(-1), id)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LastInsertID Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.UserID, record.Hours, record.Minutes).
			WillReturnResult(sqlmock.NewErrorResult(apperrors.ErrDBLastInsertId))

		id, err := record.AddRecord(tx)
		assert.Error(t, err)
		assert.Equal(t, int64(-1), id)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBLastInsertId, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRetrieveRecordByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	q := `SELECT user_id, hours, minutes FROM practice_db\.records WHERE id=\?`
	record := &Record{
		ID: 1,
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id", "hours", "minutes"}).
			AddRow(1, 2, 30)

		mock.ExpectQuery(q).WithArgs(record.ID).WillReturnRows(rows)

		err := record.RetrieveRecordByID(db)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), record.UserID)
		assert.Equal(t, 2, record.Hours)
		assert.Equal(t, 30, record.Minutes)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(record.ID).WillReturnError(sql.ErrNoRows)

		err := record.RetrieveRecordByID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrDBNoRows, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(record.ID).WillReturnError(apperrors.ErrDBQuery)

		err := record.RetrieveRecordByID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRetrieveUserIDByRecordID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	q := `SELECT user_id FROM practice_db\.records WHERE id=\?`
	record := &Record{
		ID: 1,
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id"}).
			AddRow(1)

		mock.ExpectQuery(q).WithArgs(record.ID).WillReturnRows(rows)

		err := record.RetrieveUserIDByRecordID(db)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), record.UserID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(record.ID).WillReturnError(sql.ErrNoRows)

		err := record.RetrieveUserIDByRecordID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrDBNoRows, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(record.ID).WillReturnError(apperrors.ErrDBQuery)

		err := record.RetrieveUserIDByRecordID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRemoveRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	q := `DELETE FROM practice_db\.records WHERE id=\?`
	record := Record{
		ID: 1,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := record.RemoveRecord(tx)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.ID).WillReturnResult(sqlmock.NewResult(1, 0))

		err := record.RemoveRecord(tx)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrNotFound, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.ID).WillReturnError(apperrors.ErrDBQuery)

		err := record.RemoveRecord(tx)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Rows Affected Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(record.ID).WillReturnResult(sqlmock.NewErrorResult(apperrors.ErrDBRowsAffected))

		err := record.RemoveRecord(tx)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBRowsAffected, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
