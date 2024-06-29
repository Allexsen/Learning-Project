package models

import (
	"fmt"
	"log"

	"github.com/Allexsen/Learning-Project/internal/db"
)

type User struct {
	ID           int64    `db:"id" json:"id"`
	Name         string   `db:"name" json:"name"`
	Email        string   `db:"email" json:"email"`
	LogCount     int      `db:"log_count" json:"log_count"`
	TotalHours   int      `db:"total_hours" json:"total_hours"`
	TotalMinutes int      `db:"total_minutes" json:"total_minutes"`
	Records      []Record `db:"-" json:"worklog"`
}

func (u *User) AddUser() (int64, error) {
	result, err := db.DB.Exec(`INSERT INTO practice_db.users (name, email, log_count) VALUES (?, ?, ?)`, u.Name, u.Email, u.LogCount)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (u *User) RetrieveUserbyID() error {
	err := db.DB.QueryRow(`SELECT * FROM practice_db.users WHERE id=?`, u.ID).Scan(
		&u.ID, &u.Name, &u.Email, &u.LogCount, &u.TotalHours, &u.TotalMinutes)
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
	err := db.DB.QueryRow(`SELECT id FROM practice_db.users WHERE email=?`, u.Email).Scan(&u.ID)
	if err != nil {
		log.Print("hit model retrieval")
		return err
	}

	return nil
}

func (u User) UpdateUserWorklogInfoByID() error {
	q := `UPDATE practice_db.users
		SET total_hours=?, total_minutes=?, log_count=?
		WHERE id=?`

	res, err := db.DB.Exec(q, u.TotalHours, u.TotalMinutes, u.LogCount, u.ID)
	if err != nil {
		return fmt.Errorf("couldn't update the user: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user find with ID %d", u.ID)
	}

	return nil
}
