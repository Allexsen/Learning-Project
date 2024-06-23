package models

import (
	"database/sql"

	"github.com/Allexsen/Learning-Project/internal/db"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	TotalHours   int
	TotalMinutes int
	LogCount     int
	Records      []Record
}

func (u *User) AddUser() (int64, error) {
	result, err := db.DB.Exec(`INSERT INTO practiceDB.users (name, email) VALUES (?, ?)`, u.Name, u.Email)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (u *User) RetrieveUser() error {
	err := db.DB.QueryRow(`SELECT * FROM practiceDB.users WHERE id=?`, u.ID).Scan(
		&u.ID, &u.Name, &u.Email, &u.TotalHours, &u.TotalMinutes, &u.LogCount)

	return err
}

func (u *User) GetUserIDByEmail() error {
	err := db.DB.QueryRow(`SELECT id FROM practiceDB.users WHERE email=?`, u.Email).Scan(&u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			u.ID, err = u.AddUser()
		}

		return err
	}

	return nil
}
