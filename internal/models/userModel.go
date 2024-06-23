package models

import "github.com/Allexsen/Learning-Project/internal/db"

type User struct {
	ID           string
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
