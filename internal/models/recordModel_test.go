package models

import (
	"testing"

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

	record := Record{
		UserID:  1,
		Hours:   2,
		Minutes: 30,
	}

	query := `INSERT INTO practice_db\.records \(user_id, hours, minutes\) VALUES \(\?, \?, \?\)`
	mock.ExpectExec(query).WithArgs(record.UserID, record.Hours, record.Minutes).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = record.AddRecord(tx)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrieveRecordByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	record := &Record{
		ID: 1,
	}

	query := `SELECT user_id, hours, minutes FROM practice_db\.records WHERE id=\?`
	rows := sqlmock.NewRows([]string{"user_id", "hours", "minutes"}).
		AddRow(1, 2, 30)

	mock.ExpectQuery(query).WithArgs(record.ID).WillReturnRows(rows)

	err = record.RetrieveRecordByID(db)
	assert.NoError(t, err)

	assert.Equal(t, int64(1), record.UserID)
	assert.Equal(t, 2, record.Hours)
	assert.Equal(t, 30, record.Minutes)

	assert.NoError(t, mock.ExpectationsWereMet())
}
