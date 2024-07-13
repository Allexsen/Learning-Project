package controllers

import (
	"database/sql"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(firstname, lastname, username, email, pswd string) (models.User, error) {
	pswdHash, err := bcrypt.GenerateFromPassword([]byte(pswd), 10)
	if err != nil {
		return models.User{}, err
	}

	u := models.User{Firstname: firstname, Lastname: lastname, Username: username, Email: email, Password: string(pswdHash)}
	if u.ID, err = u.Register(); err != nil {
		return models.User{}, err
	}

	return u, nil
}

func UserLogin(email, password string) error {
	pswdHash, err := utils.GetPasswordHashByEmail(email)
	if err != nil {
		pswdHash, err = utils.GetPasswordHashByUsername(email)
		if err != nil {
			return err
		}
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

func UserGetByEmail(email string) (models.User, error) {
	u := models.User{Email: email}
	if err := u.RetrieveUserByEmail(); err != nil {
		return models.User{}, err
	}

	if err := u.RetrieveAllRecordsByUserID(); err != nil {
		return models.User{}, err
	}

	return u, nil
}

func UserGetIDByEmail(email string) (int64, error) {
	u := models.User{Email: email}
	if err := u.RetrieveUserIDByEmail(); err != nil {
		return -1, err
	}

	return u.ID, nil
}

func UserAdd(firstname, lastname, email string) (int64, error) {
	u := models.User{Firstname: firstname, Lastname: lastname, Email: email, LogCount: 0}
	var err error
	u.ID, err = u.AddUser()
	if err != nil {
		return u.ID, err
	}

	return u.ID, nil
}

func UserUpdateWorklogInfo(r models.Record, logCountChange int, tx *sql.Tx) (models.User, error) {
	u := models.User{ID: r.UserID}
	if err := u.RetrieveUserbyID(); err != nil {
		return models.User{}, err
	}

	u.TotalMinutes += r.Minutes
	if u.TotalMinutes < 0 {
		u.TotalHours--
		u.TotalMinutes += 60
	}

	u.TotalHours += r.Hours + u.TotalMinutes/60
	u.TotalMinutes %= 60
	u.LogCount += logCountChange

	if err := u.UpdateUserWorklogInfoByID(); err != nil {
		return models.User{}, err
	}

	err := u.RetrieveAllRecordsByUserID()
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
