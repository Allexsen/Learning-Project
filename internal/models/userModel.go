package models

import (
	"fmt"

	"github.com/Allexsen/Learning-Project/internal/db"
)

type User struct {
	ID           int64    `db:"id" json:"id"`
	Firstname    string   `db:"firstname" json:"firstName"`
	Lastname     string   `db:"lastname" json:"lastName"`
	Email        string   `db:"email" json:"email"`
	Username     string   `db:"username" json:"username"`
	Password     string   `db:"password" json:"-"`
	LogCount     int      `db:"log_count" json:"log_count"`
	TotalHours   int      `db:"total_hours" json:"total_hours"`
	TotalMinutes int      `db:"total_minutes" json:"total_minutes"`
	Records      []Record `db:"-" json:"worklog"`
}

func (u User) AddUser() (int64, error) {
	q := `INSERT INTO practice_db.users (firstname, lastname, email, log_count) VALUES(?, ?, ?, ?)`
	result, err := db.DB.Exec(q, u.Firstname, u.Lastname, u.Email, u.LogCount)
	if err != nil {
		return -1, fmt.Errorf("couldn't add the user: %v", err)
	}

	return result.LastInsertId()
}

func (u User) Register() (int64, error) {
	q := `INSERT INTO practice_db.users (firstname, lastname, email, username, password) VALUES(?, ?, ?, ?, ?)`
	result, err := db.DB.Exec(q, u.Firstname, u.Lastname, u.Email, u.Username, u.Password)
	if err != nil {
		return -1, fmt.Errorf("couldn't register a new user: %v", err)
	}

	return result.LastInsertId()
}

func (u *User) RetrieveUserbyID() error {
	q := `SELECT (firstname, lastname, email, username, log_count, total_hours, total_minutes)
	FROM practice_db.users
	WHERE id=?`
	err := db.DB.QueryRow(q, u.ID).Scan(
		&u.Firstname, u.Lastname, &u.Email, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the user by id: %v", err)
	}

	return nil
}

func (u *User) RetrieveUserByEmail() error {
	q := `SELECT (id, firstname, lastname, username, log_count, total_hours, total_minutes)
		FROM practice_db.users
		WHERE email=?`
	err := db.DB.QueryRow(q, u.Email).Scan(&u.ID, &u.Firstname, u.Lastname, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the user by email: %v", err)
	}

	return nil
}

func (u *User) RetrieveUserIDByEmail() error {
	err := db.DB.QueryRow(`SELECT id FROM practice_db.users WHERE email=?`, u.Email).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the user id by email: %v", err)
	}

	return nil
}

func (u User) UpdateUserWorklogInfoByID() error {
	q := `UPDATE practice_db.users
		SET log_count=?, total_hours=?, total_minutes=?
		WHERE id=?`

	res, err := db.DB.Exec(q, u.LogCount, u.TotalHours, u.TotalMinutes, u.ID)
	if err != nil {
		return fmt.Errorf("couldn't update the user info by id: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %d", u.ID)
	}

	return nil
}
