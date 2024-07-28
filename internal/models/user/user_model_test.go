package user

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		Firstname: "fn",
		Lastname:  "ln",
		Email:     "em",
		Username:  "un",
		Password:  "pswdHash",
	}

	q := `INSERT INTO practice_db\.users \(firstname, lastname, email, username, password\) VALUES\(\?, \?, \?, \?, \?\)`

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(user.Firstname, user.Lastname, user.Email, user.Username, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = user.AddUser(db)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Insert Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(user.Firstname, user.Lastname, user.Email, user.Username, user.Password).
			WillReturnError(apperrors.ErrDBQuery)

		id, err := user.AddUser(db)
		assert.Error(t, err)
		assert.Equal(t, int64(-1), id)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("LastInsertID Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(user.Firstname, user.Lastname, user.Email, user.Username, user.Password).
			WillReturnResult(sqlmock.NewErrorResult(apperrors.ErrDBLastInsertId))

		id, err := user.AddUser(db)
		assert.Error(t, err)
		assert.Equal(t, int64(-1), id)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBLastInsertId, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRetrieveUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		ID: 1,
	}

	q := `SELECT firstname, lastname, email, username, log_count, total_hours, total_minutes
		FROM practice_db\.users
		WHERE id=\?`

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"firstname", "lastname", "email", "username", "log_count", "total_hours", "total_minutes"}).
			AddRow("fn", "ln", "em", "un", 1, 2, 3)

		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnRows(rows)
		err = user.RetrieveUserbyID(db)
		assert.NoError(t, err)

		assert.Equal(t, "fn", user.Firstname)
		assert.Equal(t, "ln", user.Lastname)
		assert.Equal(t, "em", user.Email)
		assert.Equal(t, "un", user.Username)
		assert.Equal(t, 1, user.LogCount)
		assert.Equal(t, 2, user.TotalHours)
		assert.Equal(t, 3, user.TotalMinutes)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnError(sql.ErrNoRows)
		err = user.RetrieveUserbyID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrDBNoRows, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnError(apperrors.ErrDBQuery)
		err = user.RetrieveUserbyID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRetrieveUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		Email: "test@gmail.com",
	}

	q := `SELECT id, firstname, lastname, username, log_count, total_hours, total_minutes
		FROM practice_db\.users
		WHERE email=\?`
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "username", "log_count", "total_hours", "total_minutes"}).
		AddRow(1, "fn", "ln", "un", 1, 2, 3)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnRows(rows)
		err = user.RetrieveUserByEmail(db)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "fn", user.Firstname)
		assert.Equal(t, "ln", user.Lastname)
		assert.Equal(t, "un", user.Username)
		assert.Equal(t, 1, user.LogCount)
		assert.Equal(t, 2, user.TotalHours)
		assert.Equal(t, 3, user.TotalMinutes)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
		err = user.RetrieveUserByEmail(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrDBNoRows, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnError(apperrors.ErrDBQuery)
		err = user.RetrieveUserByEmail(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRetrieveUserIDByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		Email: "test@gmail.com",
	}

	q := `SELECT id FROM practice_db\.users WHERE email=\?`
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnRows(rows)
		err = user.RetrieveUserIDByEmail(db)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), user.ID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
		err = user.RetrieveUserIDByEmail(db)
		assert.Error(t, err)
		assert.Equal(t, int64(-1), user.ID)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrDBNoRows, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnError(apperrors.ErrDBQuery)
		err = user.RetrieveUserIDByEmail(db)
		assert.Error(t, err)
		assert.Equal(t, int64(-1), user.ID)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUserWorklogInfoByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		ID:           1,
		LogCount:     2,
		TotalHours:   3,
		TotalMinutes: 4,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	q := `UPDATE practice_db\.users
		SET log_count=\?, total_hours=\?, total_minutes=\?
		WHERE id=\?`

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(user.LogCount, user.TotalHours, user.TotalMinutes, user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = user.UpdateUserWorklogInfoByID(tx)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(user.LogCount, user.TotalHours, user.TotalMinutes, user.ID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := user.UpdateUserWorklogInfoByID(tx)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, apperrors.ErrNotFound, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Rows Affected Error", func(t *testing.T) {
		mock.ExpectExec(q).WithArgs(user.LogCount, user.TotalHours, user.TotalMinutes, user.ID).
			WillReturnResult(sqlmock.NewErrorResult(apperrors.ErrDBRowsAffected))

		err := user.UpdateUserWorklogInfoByID(tx)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBRowsAffected, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRetrieveAllRecordsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		ID: 1,
	}

	q := `SELECT id, hours, minutes FROM practice_db\.records WHERE user_id=\?`
	rows := sqlmock.NewRows([]string{"id", "hours", "minutes"}).
		AddRow(1, 5, 30).
		AddRow(2, 10, 45)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnRows(rows)

		err = user.RetrieveAllRecordsByUserID(db)
		assert.NoError(t, err)

		assert.Len(t, user.Records, 2)
		assert.Equal(t, int64(1), user.Records[0].ID)
		assert.Equal(t, 5, user.Records[0].Hours)
		assert.Equal(t, 30, user.Records[0].Minutes)
		assert.Equal(t, int64(2), user.Records[1].ID)
		assert.Equal(t, 10, user.Records[1].Hours)
		assert.Equal(t, 45, user.Records[1].Minutes)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "hours", "minutes"})

		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnRows(rows)

		err = user.RetrieveAllRecordsByUserID(db)
		assert.NoError(t, err)

		assert.Len(t, user.Records, 0)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnError(apperrors.ErrDBQuery)

		err = user.RetrieveAllRecordsByUserID(db)
		assert.Error(t, err)

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Rows Iteration Error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "hours", "minutes"}).
			AddRow(1, 5, 30).
			RowError(0, errors.New("iteration error"))

		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnRows(rows)

		err = user.RetrieveAllRecordsByUserID(db)
		assert.Error(t, err, "An error is expected but got nil")

		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, apperrors.ErrDBQuery, appErr.Err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
