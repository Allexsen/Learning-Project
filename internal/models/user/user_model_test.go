package user

import (
	"database/sql"
	"net/http"
	"testing"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TODO: Renew Unit Tests - UserDTO has been changed significantly.

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := User{
		UserDTO: UserDTO{Firstname: "fn",
			Lastname: "ln",
			Email:    "em",
			Username: "un",
		},
		Password: "pswdHash",
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
		UserDTO: UserDTO{
			ID: 1,
		},
	}

	q := `SELECT firstname, lastname, email, username
		FROM practice_db\.users
		WHERE id=\?`

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"firstname", "lastname", "email", "username"}).
			AddRow("fn", "ln", "em", "un")

		mock.ExpectQuery(q).WithArgs(user.ID).WillReturnRows(rows)
		err = user.RetrieveUserbyID(db)
		assert.NoError(t, err)

		assert.Equal(t, "fn", user.Firstname)
		assert.Equal(t, "ln", user.Lastname)
		assert.Equal(t, "em", user.Email)
		assert.Equal(t, "un", user.Username)

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
		UserDTO: UserDTO{
			Email: "test@gmail.com",
		},
	}

	q := `SELECT id, firstname, lastname, username
		FROM practice_db\.users
		WHERE email=\?`
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "username"}).
		AddRow(1, "fn", "ln", "un")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(q).WithArgs(user.Email).WillReturnRows(rows)
		err = user.RetrieveUserByEmail(db)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "fn", user.Firstname)
		assert.Equal(t, "ln", user.Lastname)
		assert.Equal(t, "un", user.Username)

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
		UserDTO: UserDTO{
			Email: "test@gmail.com",
		},
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
