package models

import (
	"log"

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
	log.Print("Hit AddUser")
	result, err := db.DB.Exec(`INSERT INTO practice_db.users (name, email) VALUES (?, ?)`, u.Name, u.Email)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (u *User) RetrieveUserbyID() error {
	err := db.DB.QueryRow(`SELECT * FROM practice_db.users WHERE id=?`, u.ID).Scan(
		&u.ID, &u.Name, &u.Email, &u.TotalHours, &u.TotalMinutes, &u.LogCount)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) RetrieveUserByEmail() error {
	err := db.DB.QueryRow(`SELECT * FROM practice_db.users WHERE email=?`, u.Email).Scan(
		&u.ID, &u.Name, &u.Email, &u.TotalHours, &u.TotalMinutes, &u.LogCount)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) RetrieveUserIDByEmail() error {
	log.Print("Hit RetrieveUserIDByEmail")

	err := db.DB.QueryRow("SELECT id FROM practice_db.users WHERE email=?", u.Email).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}
