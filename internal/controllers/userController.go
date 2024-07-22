// Package controllers is responsible for handler-model interaction, and business logic
package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	database "github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserRegister takes user info, generates bcrypt
// password hash, and adds the user to the database
func UserRegister(firstname, lastname, username, email, pswd string) (models.User, error) {
	pswdHash, err := bcrypt.GenerateFromPassword([]byte(pswd), 10)
	if err != nil {
		return models.User{}, err
	}

	db := database.DB
	u := models.User{Firstname: firstname, Lastname: lastname, Username: username, Email: email, Password: string(pswdHash)}
	// Query adding to the db
	if u.ID, err = u.AddUser(db); err != nil {
		return models.User{}, err
	}

	return u, nil
}

// UserLogin retrieves bcrypt password hash, authenticates
// user credentials and, if successful, logs in the user
func UserLogin(credential, password string) error {
	var pswdHash string
	var err error
	db := database.DB

	// Check if the provided credentail is email or username, and query accordingly
	if strings.Contains(credential, "@") {
		pswdHash, err = utils.GetPasswordHashByEmail(db, credential)
	} else {
		pswdHash, err = utils.GetPasswordHashByUsername(db, credential)
	}

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(pswdHash), []byte(password))
	if err != nil {
		return apperrors.New(
			http.StatusUnauthorized,
			"Invalid credentials",
			apperrors.ErrUnauthorized,
			map[string]interface{}{"details": err.Error()},
		)
	}

	return nil
}

// UserGetByEmail retrieves user from the database by user email
func UserGetByEmail(email string) (models.User, error) {
	db := database.DB
	u := models.User{Email: email}
	if err := u.RetrieveUserByEmail(db); err != nil {
		return models.User{}, err
	}

	// retrieves records associated with the user by user id
	if err := u.RetrieveAllRecordsByUserID(db); err != nil {
		return models.User{}, err
	}

	return u, nil
}

// UserGetIDByEmail retrieves user ID by user email.
// return -1 and error in case of failure
func UserGetIDByEmail(db *sql.DB, email string) (int64, error) {
	u := models.User{Email: email}
	if err := u.RetrieveUserIDByEmail(db); err != nil {
		return -1, err
	}

	return u.ID, nil
}

// UserUpdateWorklogInfo updates the user worklog data by the provided record
func userUpdateWorklogInfo(db *sql.DB, r models.Record, countChange int, tx *sql.Tx) (models.User, error) {
	u := models.User{ID: r.UserID}
	if err := u.RetrieveUserbyID(db); err != nil {
		return models.User{}, err
	}

	u.TotalMinutes += r.Minutes
	// Adjust minutes if below zero, decrease hours by one
	if u.TotalMinutes < 0 {
		u.TotalHours--
		u.TotalMinutes += 60
	}

	u.TotalHours += r.Hours + u.TotalMinutes/60
	u.TotalMinutes %= 60
	u.LogCount += countChange

	// Query update to the db
	if err := u.UpdateUserWorklogInfoByID(tx); err != nil {
		return models.User{}, err
	}

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	err = u.RetrieveAllRecordsByUserID(db)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
