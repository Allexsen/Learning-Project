package controllers

import (
	"fmt"
	"log"
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
		log.Printf("Hit record add error: %v", err)
		return models.User{}, err
	}

	u, err := UpdateUserWorklogInfo(r)
	if err != nil {
		if err2 := r.RemoveRecord(); err2 != nil {
			log.Printf("[Error]: Failed to update the user worklog: %v, and failed to delete the corresponding record: %v", err, err2)
			return models.User{}, fmt.Errorf("failed to update the user worklog: %v, and failed to delete the record: %v", err, err2)
		}

		return models.User{}, fmt.Errorf("couldn't add the record - failed to update the user worklog: %v", err)
	}

	return u, nil
}
