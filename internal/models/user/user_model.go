// Package user provides user model and its methods for database interaction.
package user

import (
	"database/sql"

	"github.com/Allexsen/Learning-Project/internal/models/common"
	"github.com/Allexsen/Learning-Project/internal/models/record"
)

// User represents an internal object.
// It stores user data (:D)
type User struct {
	ID           int64           `db:"id" json:"id"`                       // Unique user id
	Firstname    string          `db:"firstname" json:"firstName"`         // Firstname
	Lastname     string          `db:"lastname" json:"lastName"`           // Lastname
	Email        string          `db:"email" json:"email"`                 // Email
	Username     string          `db:"username" json:"username"`           // Unique username
	Password     string          `db:"password" json:"-"`                  // Password hash
	LogCount     int             `db:"log_count" json:"log_count"`         // Total number of logs
	TotalHours   int             `db:"total_hours" json:"total_hours"`     // Total hours worked
	TotalMinutes int             `db:"total_minutes" json:"total_minutes"` // Total minutes worked
	Records      []record.Record `db:"-" json:"worklog"`                   // List of records
}

func (u User) AddUser(db *sql.DB) (int64, error) {
	q := `INSERT INTO practice_db.users (firstname, lastname, email, username, password) VALUES(?, ?, ?, ?, ?)`
	result, err := db.Exec(q, u.Firstname, u.Lastname, u.Email, u.Username, u.Password)
	if err != nil {
		return -1, common.GetQueryError(q, "Couldn't register a new user", u, err)
	}

	return common.GetLastInsertId(result, q, u)
}

func (u *User) RetrieveUserbyID(db *sql.DB) error {
	q := `SELECT firstname, lastname, email, username, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE id=?`
	err := db.QueryRow(q, u.ID).Scan(
		&u.Firstname, &u.Lastname, &u.Email, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	return common.GetQueryError(q, "Couldn't retrieve user by id", u, err)
}

func (u *User) RetrieveUserByEmail(db *sql.DB) error {
	q := `SELECT id, firstname, lastname, username, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE email=?`
	err := db.QueryRow(q, u.Email).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	return common.GetQueryError(q, "Couldn't retrieve user by email", u, err)
}

func (u *User) RetrieveUserIDByEmail(db *sql.DB) error {
	u.ID = -1

	q := `SELECT id FROM practice_db.users WHERE email=?`
	err := db.QueryRow(q, u.Email).Scan(&u.ID)

	return common.GetQueryError(q, "Couldn't retrieve user id by email", u, err)
}

// UpdateUserWorklogInfoByID changes the information about the user's worklog by user id,
// precisely: log count, hours and minutes worked.
func (u User) UpdateUserWorklogInfoByID(tx *sql.Tx) error {
	q := `UPDATE practice_db.users
		SET log_count=?, total_hours=?, total_minutes=?
		WHERE id=?`
	result, err := tx.Exec(q, u.LogCount, u.TotalHours, u.TotalMinutes, u.ID)

	return common.HandleUpdateQuery(result, err, q, u)
}

// RetrieveAllRecordsByUserID scans records table,
// looking for every record associated with the user.
func (u *User) RetrieveAllRecordsByUserID(db *sql.DB) error {
	q := `SELECT id, hours, minutes FROM practice_db.records WHERE user_id=?`
	rows, err := db.Query(q, u.ID)
	if err != nil {
		return common.GetQueryError(q, "Couldn't query the database", u, err)
	}
	defer rows.Close()

	var records []record.Record
	for rows.Next() {
		r := record.Record{UserID: u.ID}
		if err = rows.Scan(&r.ID, &r.Hours, &r.Minutes); err != nil {
			return common.GetQueryError(q, "Couldn't scan a row", u, err)
		}
		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return common.GetQueryError(q, "Error during the iteration of the rows", u, err)
	}

	u.Records = records
	return nil
}
