package controllers

import (
	"fmt"

	"github.com/Allexsen/Learning-Project/internal/models"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(firstname, lastname, username, email, pswd string) (models.User, error) {
	pswdHash, err := bcrypt.GenerateFromPassword([]byte(pswd), 10)
	if err != nil {
		return models.User{}, fmt.Errorf("couldn't generate a password hash: %v", err)
	}

	u := models.User{Firstname: firstname, Lastname: lastname, Username: username, Email: email, Password: string(pswdHash)}
	if u.ID, err = u.Register(); err != nil {
		return models.User{}, fmt.Errorf("couldn't register a new user: %v", err)
	}

	return u, nil
}

func UserLogin(email, password string) error {
	pswdHash, err := utils.GetPasswordHashByEmail(email)
	if err != nil {
		pswdHash, err = utils.GetPasswordHashByUsername(email)
		if err != nil {
			return fmt.Errorf("couldn't log in a user: %v", err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(pswdHash), []byte(password))
	if err != nil {
		return fmt.Errorf("couldn't log in a user: %v", err)
	}

	return nil
}

func UserGetByEmail(email string) (models.User, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserByEmail()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user: %v", err)
	}

	u.Records, err = utils.RetrieveAllRecordsByUserID(u.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user records: %v", err)
	}

	return u, nil
}

func UserGetIDByEmail(email string) (int64, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserIDByEmail()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve the user: %v", err)
	}

	return u.ID, nil
}

func UserAdd(firstname, lastname, email string) (int64, error) {
	u := models.User{Firstname: firstname, Lastname: lastname, Email: email, LogCount: 0}
	var err error
	u.ID, err = u.AddUser()
	if err != nil {
		return u.ID, fmt.Errorf("failed to add a new user: %v", err)
	}

	return u.ID, nil
}

func UserUpdateWorklogInfo(r models.Record, logCountChange int) (models.User, error) {
	u := models.User{ID: r.UserID}
	if err := u.RetrieveUserbyID(); err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user: %v", err)
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
		return models.User{}, fmt.Errorf("failed to update the user worklog info: %v", err)
	}

	var err error
	u.Records, err = utils.RetrieveAllRecordsByUserID(u.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user worklog info: %v", err)
	}

	return u, nil
}
