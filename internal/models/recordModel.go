// Package models declares object models, and provides methods for database interaction.
package models

import (
	"database/sql"

	"github.com/Allexsen/Learning-Project/internal/utils"
)

// Record represents an internal object.
// It stores a single work session duration of a user.
type Record struct {
	ID      int64 `db:"id" json:"id"`           // Unique record id
	UserID  int64 `db:"user_id" json:"user_id"` // Corresponding user
	Hours   int   `db:"hours" json:"hours"`     // Hours worked during a session
	Minutes int   `db:"minutes" json:"minutes"` // Minutes worked during a session
}

func (r Record) AddRecord(tx *sql.Tx) (int64, error) {
	q := `INSERT INTO practice_db.records (user_id, hours, minutes) VALUES (?, ?, ?)`
	result, err := tx.Exec(q, r.UserID, r.Hours, r.Minutes)
	if err != nil {
		return -1, utils.GetQueryError(q, "Couldn't add new record", r, err)
	}

	return utils.GetLastInsertId(result, q, r)
}

func (r *Record) RetrieveRecordByID(db *sql.DB) error {
	q := `SELECT user_id, hours, minutes FROM practice_db.records WHERE id=?`
	err := db.QueryRow(q, r.ID).Scan(&r.UserID, &r.Hours, &r.Minutes)

	return utils.GetQueryError(q, "Couldn't find record by id", r, err)
}

func (r *Record) RetrieveUserIDByRecordID(db *sql.DB) error {
	q := `SELECT user_id FROM practice_db.records WHERE id=?`
	err := db.QueryRow(q, r.ID).Scan(&r.UserID)

	return utils.GetQueryError(q, "Couldn't retrieve user id by record id", r, err)
}

func (r Record) RemoveRecord(tx *sql.Tx) error {
	q := `DELETE FROM practice_db.records WHERE id=?`
	result, err := tx.Exec(q, r.ID)

	return utils.HandleUpdateQuery(result, err, q, r)
}
