// Package record provides record model and its methods for database interaction.
package record

import (
	"database/sql"
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/common"
)

// Record represents an internal object.
// It stores a single work session duration of a user.
type Record struct {
	ID      int64 `db:"id" json:"id"`           // Unique record id
	UserID  int64 `db:"user_id" json:"user_id"` // Corresponding user
	Hours   int   `db:"hours" json:"hours"`     // Hours worked during a session
	Minutes int   `db:"minutes" json:"minutes"` // Minutes worked during a session
}

// AddRecord adds a new record to the database.
func (r Record) AddRecord(tx *sql.Tx) (int64, error) {
	log.Printf("[RECORD] Adding record: %+v", r)

	q := `INSERT INTO practice_db.records (user_id, hours, minutes) VALUES (?, ?, ?)`
	result, err := tx.Exec(q, r.UserID, r.Hours, r.Minutes)
	if err != nil {
		return -1, common.GetQueryError(q, "Couldn't add new record", r, err)
	}

	id, err := common.GetLastInsertId(result, q, r)
	if err != nil {
		return -1, err
	}

	log.Printf("[RECORD] Record has been successfully added with id %d", id)
	return id, nil
}

// RetrieveAllRecordsByUserID retrieves all records associated with the user by user id.
func (r *Record) RetrieveRecordByID(db *sql.DB) error {
	log.Printf("[RECORD] Retrieving record by id %d", r.ID)

	q := `SELECT user_id, hours, minutes FROM practice_db.records WHERE id=?`
	err := db.QueryRow(q, r.ID).Scan(&r.UserID, &r.Hours, &r.Minutes)

	err = common.GetQueryError(q, "Couldn't retrieve record by id", r, err)
	if err != nil {
		return err
	}

	log.Printf("[RECORD] Record has been successfully retrieved: %+v", r)
	return nil
}

// RetrieveAllRecordsByUserID retrieves all records associated with the user by user id.
func (r *Record) RetrieveUserIDByRecordID(db *sql.DB) error {
	log.Printf("[RECORD] Retrieving user id by record id %d", r.ID)

	q := `SELECT user_id FROM practice_db.records WHERE id=?`
	err := db.QueryRow(q, r.ID).Scan(&r.UserID)

	err = common.GetQueryError(q, "Couldn't retrieve user id by record id", r, err)
	if err != nil {
		return err
	}

	log.Printf("[RECORD] User id %d has been successfully retrieved by record id %d", r.UserID, r.ID)
	return nil
}

// RetrieveAllRecordsByUserID retrieves all records associated with the user by user id.
func (r Record) RemoveRecord(tx *sql.Tx) error {
	log.Printf("[RECORD] Removing record %d", r.ID)

	q := `DELETE FROM practice_db.records WHERE id=?`
	result, err := tx.Exec(q, r.ID)

	err = common.HandleUpdateQuery(result, err, q, r)
	if err != nil {
		return err
	}

	log.Printf("[RECORD] Record %d has been successfully removed", r.ID)
	return nil
}
