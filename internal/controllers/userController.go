package controllers

import (
	"fmt"

	"github.com/Allexsen/Learning-Project/internal/models"
)

func GetUserByEmail(email string) (models.User, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserByEmail()
	if err != nil {
		return models.User{}, err
	}

	u.Records, err = models.RetrieveAllRecordsByUserID(u.ID)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func GetUserIDByEmail(email string) (int64, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserIDByEmail()
	if err != nil {
		return -1, err
	}

	return u.ID, nil
}

func AddNewUser(name, email string) (int64, error) {
	u := models.User{Name: name, Email: email, LogCount: 0}
	var err error
	u.ID, err = u.AddUser()
	if err != nil {
		return u.ID, err
	}

	return u.ID, nil
}

func UpdateUserWorklogInfo(r models.Record) (models.User, error) {
	u := models.User{ID: r.UserID}
	if err := u.RetrieveUserbyID(); err != nil {
		return models.User{}, fmt.Errorf("couldn't retrieve the user: %v", err)
	}

	u.TotalMinutes += r.Minutes
	u.TotalHours += r.Hours + u.TotalMinutes/60
	u.TotalMinutes %= 60
	u.LogCount++

	if err := u.UpdateUserWorklogInfoByID(); err != nil {
		return models.User{}, fmt.Errorf("couldn't update the user worklog info: %v", err)
	}

	return u, nil
}
