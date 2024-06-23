package models

import (
	"github.com/Allexsen/Learning-Project/internal/db"
)

type Record struct {
	ID      int64
	UserID  string
	Hours   int
	Minutes int
}

func (r *Record) AddRecord() (int64, error) {
	result, err := db.DB.Exec("INSERT INTO practiceDB.records (user_id, hours, minutes) VALUES (?, ?, ?)", r.UserID, r.Hours, r.Minutes)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
