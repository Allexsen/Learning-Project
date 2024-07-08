package controllers

import (
	"fmt"
	"strconv"

	models "github.com/Allexsen/Learning-Project/internal/models"
)

func RecordAdd(email, hrStr, minStr string) (models.User, error) {
	hours, err := strconv.Atoi(hrStr)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to convert hrStr to int: %v", err)
	}

	minutes, err := strconv.Atoi(minStr)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to convert minStr to int: %v", err)
	}

	r := models.Record{Hours: hours, Minutes: minutes}
	r.UserID, err = UserGetIDByEmail(email)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve a user: %v", err)
	}

	r.ID, err = r.AddRecord()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to add a record: %v", err)
	}

	u, err := UserUpdateWorklogInfo(r, 1)
	if err != nil {
		if err2 := r.RemoveRecord(); err2 != nil {
			return models.User{}, fmt.Errorf("failed to update a user worklog: %v, and failed to revert a record back: %v", err, err2)
		}

		return models.User{}, fmt.Errorf("failed to add a record - couldn't update a user worklog: %v", err)
	}

	return u, nil
}

func RecordRemove(rid int) (models.User, error) {
	r := models.Record{ID: int64(rid)}
	if err := r.RetrieveRecordByID(); err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve a record: %v", err)
	}

	if err := r.RemoveRecord(); err != nil {
		return models.User{}, fmt.Errorf("failed to remove a record: %v", err)
	}

	r.Hours *= -1
	r.Minutes *= -1
	u, err := UserUpdateWorklogInfo(r, -1)
	if err != nil {
		if _, err2 := r.AddRecord(); err2 != nil {
			return models.User{}, fmt.Errorf("failed to update a user worklog: %v, and failed to revert a record %d back: %v", err, r.ID, err2)
		}

		return models.User{}, fmt.Errorf("failed to delete a record - couldn't update a user worklog: %v", err)
	}

	return u, nil
}
