// Package controllers is responsible for handler-model interaction, and business logic
package controllers

import (
	"log"
	"net/http"

	database "github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/record"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
)

// RecordAdd adds a new record to the database,
// then updates the user and user's worklog.
func RecordAdd(email, hrStr, minStr string) (*user.User, error) {
	log.Printf("[CONTROLLER] Adding record for %s", email)

	hours, err := utils.Atoi(hrStr)
	if err != nil {
		return nil, err
	}

	minutes, err := utils.Atoi(minStr)
	if err != nil {
		return nil, err
	}

	db := database.DB

	r := record.Record{Hours: hours, Minutes: minutes}
	log.Printf("[CONTROLLER] Record: %+v", r)

	// Get corresponding user
	r.UserID, err = UserGetIDByEmail(db, email)
	if err != nil {
		return nil, err
	}

	// Start transaction to ensure record and user
	// are both updated. If any of them fails, transaction
	// must be rolled back to maintain data integrity.
	tx, err := db.Begin()
	if err != nil {
		return nil, apperrors.New(
			http.StatusInternalServerError,
			"Failed to begin transaction",
			apperrors.ErrDBTransaction,
			map[string]interface{}{"details": err.Error()},
		)
	}

	// Defer rollback in case anything fails.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	r.ID, err = r.AddRecord(tx)
	if err != nil {
		return nil, err
	}

	u, err := userUpdateWorklogInfo(db, r, 1, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[CONTROLLER] Couldn't commit the transaction: %s", err)
		log.Printf("[CONTROLLER] Rolling back the transaction")
		err = tx.Rollback()
	}

	return u, err
}

// RecordRemove deletes a record by record id,
// then updates the user and user's worklog.
func RecordRemove(rid int) (*user.User, error) {
	log.Printf("[CONTROLLER] Removing record %d", rid)

	db := database.DB
	r := record.Record{ID: int64(rid)}

	if err := r.RetrieveRecordByID(db); err != nil {
		return nil, err
	}

	// Start transaction to ensure record and user
	// are both updated. If any of them fails, transaction
	// must be rolled back to maintain data integrity.
	tx, err := db.Begin()
	if err != nil {
		return nil, apperrors.New(
			http.StatusInternalServerError,
			"Failed to begin transaction",
			apperrors.ErrDBTransaction,
			map[string]interface{}{"details": err.Error()},
		)
	}

	// Defer a rollback in case anything fails.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := r.RemoveRecord(tx); err != nil {
		return nil, err
	}

	// Removing record decreases user's work time,
	// which is simulated by negative amount while updating
	r.Hours *= -1
	r.Minutes *= -1
	u, err := userUpdateWorklogInfo(db, r, -1, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[CONTROLLER] Couldn't commit the transaction: %s", err)
		log.Printf("[CONTROLLER] Rolling back the transaction")
		err = tx.Rollback()
	}

	return u, err
}
