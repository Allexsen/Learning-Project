package models

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/db"
)

type Record struct {
	ID      int64 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"user_id"`
	Hours   int   `db:"hours" json:"hours"`
	Minutes int   `db:"minutes" json:"minutes"`
}

func (r *Record) AddRecord() (int64, error) {
	log.Print("Hit AddRecord")

	result, err := db.DB.Exec("INSERT INTO practice_db.records (user_id, hours, minutes) VALUES (?, ?, ?)", r.UserID, r.Hours, r.Minutes)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
