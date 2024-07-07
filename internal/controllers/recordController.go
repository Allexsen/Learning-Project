package controllers

import (
	"fmt"
	"strconv"

	models "github.com/Allexsen/Learning-Project/internal/models"
)

func RecordAdd(firstname, lastname, email, hrStr, minStr string) (models.User, error) {
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
		if r.UserID == -1 {
			r.UserID, err = UserAdd(firstname, lastname, email)
			if err != nil {
				return models.User{}, fmt.Errorf("failed to add the record - couldn't create a new user: %v", err)
			}
		} else {
			return models.User{}, fmt.Errorf("failed to retrieve the user: %v", err)
		}
	}

	r.ID, err = r.AddRecord()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to add the record: %v", err)
	}

	u, err := UserUpdateWorklogInfo(r, 1)
	if err != nil {
		if err2 := r.RemoveRecord(); err2 != nil {
			return models.User{}, fmt.Errorf("failed to update the user worklog: %v, and failed to revert the record back: %v", err, err2)
		}

		return models.User{}, fmt.Errorf("failed to add the record - couldn't update the user worklog: %v", err)
	}

	return u, nil
}

func RecordRemove(rid int) (models.User, error) {
	r := models.Record{ID: int64(rid)}
	if err := r.RetrieveRecordByID(); err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve the record: %v", err)
	}

	if err := r.RemoveRecord(); err != nil {
		return models.User{}, fmt.Errorf("failed to remove the record: %v", err)
	}

	r.Hours *= -1
	r.Minutes *= -1
	u, err := UserUpdateWorklogInfo(r, -1)
	if err != nil {
		if _, err2 := r.AddRecord(); err2 != nil {
			return models.User{}, fmt.Errorf("failed to update the user worklog: %v, and failed to revert the record %d back: %v", err, r.ID, err2)
		}

		return models.User{}, fmt.Errorf("failed to delete the record - couldn't update the user worklog: %v", err)
	}

	return u, nil
}
