package models

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/db"
)

type Record struct {
	ID      int64
	UserID  int64
	Hours   int
	Minutes int
}

func (r *Record) AddRecord() (int64, error) {
	log.Print("Hit AddRecord")

	result, err := db.DB.Exec("INSERT INTO practice_db.records (user_id, hours, minutes) VALUES (?, ?, ?)", r.UserID, r.Hours, r.Minutes)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
