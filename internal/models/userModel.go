package models

import (
	"fmt"

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
		return -1, fmt.Errorf("couldn't add the user: %v", err)
	}

	return result.LastInsertId()
}

func (u *User) RetrieveUserbyID() error {
	err := db.DB.QueryRow(`SELECT * FROM practice_db.users WHERE id=?`, u.ID).Scan(
		&u.ID, &u.Name, &u.Email, &u.LogCount, &u.TotalHours, &u.TotalMinutes)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the user by id: %v", err)
	}

	return nil
}

func (u *User) RetrieveUserByEmail() error {
	err := db.DB.QueryRow(`SELECT * FROM practice_db.users WHERE email=?`, u.Email).Scan(
		&u.ID, &u.Name, &u.Email, &u.TotalHours, &u.TotalMinutes, &u.LogCount)
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

func RetrieveAllRecordsByUserID(uid int64) ([]Record, error) {
	rows, err := db.DB.Query(`SELECT id, hours, minutes FROM practice_db.records WHERE user_id=?`, uid)
	if err != nil {
		return []Record{}, fmt.Errorf("couldn't query the database: %v", err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		r := Record{UserID: uid}
		err = rows.Scan(&r.ID, &r.Hours, &r.Minutes)
		if err != nil {
			return []Record{}, fmt.Errorf("couldn't scan a row: %v", err)
		}

		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return []Record{}, fmt.Errorf("error during the iteration of the rows: %v", err)
	}

	return records, nil
}
