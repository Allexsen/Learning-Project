package controllers

import (
	"strconv"

	models "github.com/Allexsen/Learning-Project/internal/models"
)

func RecordAdd(name, email, hStr, minStr string) (models.Record, error) {
	hours, err := strconv.Atoi(hStr)
	if err != nil {
		return models.Record{}, err
	}

	minutes, err := strconv.Atoi(minStr)
	if err != nil {
		return models.Record{}, err
	}

	r := models.Record{Hours: hours, Minutes: minutes}
	r.UserID, err = GetUserIDByEmail(email)
	if err != nil {
		if r.UserID == -1 {
			u := models.User{Name: name, Email: email}
			r.UserID, err = u.AddUser()
			if err != nil {
				return models.Record{}, err
			}
		} else {
			return models.Record{}, nil
		}
	}

	r.ID, err = r.AddRecord()
	if err != nil {
		return models.Record{}, err
	}

	return r, nil
}
