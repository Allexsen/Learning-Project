package models

import (
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

func (u User) AddUser() (int64, error) { // here, should it be "error" or "*customErrors.Err" ??
	q := `INSERT INTO practice_db.users (firstname, lastname, email, log_count) VALUES(?, ?, ?, ?)`
	result, err := db.DB.Exec(q, u.Firstname, u.Lastname, u.Email, u.LogCount)
	if err != nil {
		return -1, getQueryError(q, "Couldn't add a new user", u, err)
	}

	return getLastInsertId(result, q, u)
}

func (u User) Register() (int64, error) {
	q := `INSERT INTO practice_db.users (firstname, lastname, email, username, password) VALUES(?, ?, ?, ?, ?)`
	result, err := db.DB.Exec(q, u.Firstname, u.Lastname, u.Email, u.Username, u.Password)
	if err != nil {
		return -1, getQueryError(q, "Couldn't register a new user", u, err)
	}

	return getLastInsertId(result, q, u)
}

func (u *User) RetrieveUserbyID() error {
	q := `SELECT firstname, lastname, email, username, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE id=?`
	err := db.DB.QueryRow(q, u.ID).Scan(
		&u.Firstname, &u.Lastname, &u.Email, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	return getQueryError(q, "Couldn't retrieve user by id", u, err)
}

func (u *User) RetrieveUserByEmail() error {
	q := `SELECT id, firstname, lastname, username, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE email=?`
	err := db.DB.QueryRow(q, u.Email).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	return getQueryError(q, "Couldn't retrieve user by email", u, err)
}

func (u *User) RetrieveUserIDByEmail() error {
	q := `SELECT id FROM practice_db.users WHERE email=?`
	err := db.DB.QueryRow(q, u.Email).Scan(&u.ID)

	return getQueryError(q, "Couldn't retrieve user id by email", u, err)
}

func (u User) UpdateUserWorklogInfoByID() error {
	q := `UPDATE practice_db.users
		SET log_count=?, total_hours=?, total_minutes=?
		WHERE id=?`
	result, err := db.DB.Exec(q, u.LogCount, u.TotalHours, u.TotalMinutes, u.ID)

	return handleUpdateQuery(result, err, q, u)
}

func (u *User) RetrieveAllRecordsByUserID() error {
	q := `SELECT id, hours, minutes FROM practice_db.records WHERE user_id=?`
	rows, err := db.DB.Query(q, u.ID)
	if err != nil {
		return getQueryError(q, "Couldn't query the database", u, err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		r := Record{UserID: u.ID}
		if err = rows.Scan(&r.ID, &r.Hours, &r.Minutes); err != nil {
			return getQueryError(q, "Couldn't scan a row", u, err)
		}
		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return getQueryError(q, "Error during the iteration of the rows", u, err)
	}

	u.Records = records
	return nil
}
