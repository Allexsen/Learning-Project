package models

import (
	"fmt"

	"github.com/Allexsen/Learning-Project/internal/db"
)

type Record struct {
	ID      int64 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"user_id"`
	Hours   int   `db:"hours" json:"hours"`
	Minutes int   `db:"minutes" json:"minutes"`
}

func (r *Record) AddRecord() (int64, error) {
	result, err := db.DB.Exec("INSERT INTO practice_db.records (user_id, hours, minutes) VALUES (?, ?, ?)", r.UserID, r.Hours, r.Minutes)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r Record) RemoveRecord() error {
	q := "DELETE FROM practice_db.records WHERE id=?"
	res, err := db.DB.Exec(q, r.ID)
	if err != nil {
		return fmt.Errorf("couldn't delete the record: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't retrieve the rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id: %d", r.ID)
	}

	return nil
}
