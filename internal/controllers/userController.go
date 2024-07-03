package controllers

import (
	"fmt"

	"github.com/Allexsen/Learning-Project/internal/models"
)

func GetUserByEmail(email string) (models.User, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserByEmail()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user: %v", err)
	}

	u.Records, err = models.RetrieveAllRecordsByUserID(u.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user records: %v", err)
	}

	return u, nil
}

func GetUserIDByEmail(email string) (int64, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserIDByEmail()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve the user: %v", err)
	}

	return u.ID, nil
}

func AddNewUser(name, email string) (int64, error) {
	u := models.User{Name: name, Email: email, LogCount: 0}
	var err error
	u.ID, err = u.AddUser()
	if err != nil {
		return u.ID, fmt.Errorf("failed to add a new user: %v", err)
	}

	return u.ID, nil
}

func UpdateUserWorklogInfo(r models.Record, logCountChange int) (models.User, error) {
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
	u.Records, err = models.RetrieveAllRecordsByUserID(u.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the user worklog info: %v", err)
	}

	return u, nil
}
