package controllers

import (
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	models "github.com/Allexsen/Learning-Project/internal/models"
	"github.com/Allexsen/Learning-Project/internal/utils"
)

func RecordAdd(email, hrStr, minStr string) (models.User, error) {
	hours, err := utils.Atoi(hrStr)
	if err != nil {
		return models.User{}, err
	}

	minutes, err := utils.Atoi(minStr)
	if err != nil {
		return models.User{}, err
	}

	r := models.Record{Hours: hours, Minutes: minutes}
	r.UserID, err = UserGetIDByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return models.User{}, apperrors.New(
			http.StatusInternalServerError,
			"Failed to begin transaction",
			apperrors.ErrDBTransaction,
			map[string]interface{}{"detail": err.Error()},
		)
	}

	// Defer a rollback in case anything fails.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	r.ID, err = r.AddRecord(tx)
	if err != nil {
		return models.User{}, err
	}

	u, err := UserUpdateWorklogInfo(r, 1, tx)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func RecordRemove(rid int) (models.User, error) {
	r := models.Record{ID: int64(rid)}
	if err := r.RetrieveRecordByID(); err != nil {
		return models.User{}, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return models.User{}, apperrors.New(
			http.StatusInternalServerError,
			"Failed to begin transaction",
			apperrors.ErrDBTransaction,
			map[string]interface{}{"detail": err.Error()},
		)
	}

	// Defer a rollback in case anything fails.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := r.RemoveRecord(tx); err != nil {
		return models.User{}, err
	}

	r.Hours *= -1
	r.Minutes *= -1
	u, err := UserUpdateWorklogInfo(r, -1, tx)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
