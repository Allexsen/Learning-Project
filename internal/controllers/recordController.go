package controllers

import (
	"fmt"
	"strconv"

	models "github.com/Allexsen/Learning-Project/internal/models"
)

func RecordAdd(name, email, hrStr, minStr string) (models.User, error) {
	hours, err := strconv.Atoi(hrStr)
	if err != nil {
		return models.User{}, err
	}

	minutes, err := strconv.Atoi(minStr)
	if err != nil {
		return models.User{}, err
	}

	r := models.Record{Hours: hours, Minutes: minutes}
	r.UserID, err = GetUserIDByEmail(email)
	if err != nil {
		if r.UserID == -1 {
			r.UserID, err = AddNewUser(name, email)
			if err != nil {
				return models.User{}, err
			}
		} else {
			return models.User{}, err
		}
	}

	r.ID, err = r.AddRecord()
	if err != nil {
		return models.User{}, err
	}

	u, err := UpdateUserWorklogInfo(r, 1)
	if err != nil {
		if err2 := r.RemoveRecord(); err2 != nil {
			return models.User{}, fmt.Errorf("failed to update the user worklog: %v, and failed to revert the record back: %v", err, err2)
		}

		return models.User{}, fmt.Errorf("couldn't add the record - failed to update the user worklog: %v", err)
	}

	return u, nil
}

func RecordRemove(ridStr string) error {
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		return fmt.Errorf("invalid record id %q: couldn't convert to int", ridStr)
	}

	r := models.Record{ID: int64(rid)}
	err = r.RetrieveRecordByID()
	if err != nil {
		return fmt.Errorf("couldn't retrieve the record by the record id: %v", err)
	}

	if err := r.RemoveRecord(); err != nil {
		return err
	}

	r.Hours *= -1
	r.Minutes *= -1
	_, err = UpdateUserWorklogInfo(r, -1)
	if err != nil {
		if _, err2 := r.AddRecord(); err2 != nil {
			return fmt.Errorf("failed to update the user worklog: %v, and failed to revert the record %d back: %v", err, r.ID, err2)
		}

		return fmt.Errorf("couldn't delete the record - failed to update the user worklog: %v", err)
	}

	return nil
}
