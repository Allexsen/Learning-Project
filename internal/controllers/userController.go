package controllers

import (
	"github.com/Allexsen/Learning-Project/internal/models"
)

func GetUserByEmail(email string) (models.User, error) {
	u := models.User{Email: email}
	err := u.RetrieveUserByEmail()
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
