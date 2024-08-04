// Package user provides user model and its methods for database interaction.
package user

import (
	"database/sql"
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/common"
	"github.com/Allexsen/Learning-Project/internal/models/record"
)

type UserDTO struct {
	ID        int64  `db:"id" json:"id"`               // Unique user id
	Firstname string `db:"firstname" json:"firstName"` // Firstname
	Lastname  string `db:"lastname" json:"lastName"`   // Lastname
	Email     string `db:"email" json:"email"`         // Email
	Username  string `db:"username" json:"username"`   // Unique username
}

// User represents an internal object.
// It stores user data (:D)
type User struct {
	UserDTO
	Password     string          `db:"password" json:"-"`                  // Password hash
	LogCount     int             `db:"log_count" json:"log_count"`         // Total number of logs
	TotalHours   int             `db:"total_hours" json:"total_hours"`     // Total hours worked
	TotalMinutes int             `db:"total_minutes" json:"total_minutes"` // Total minutes worked
	Records      []record.Record `db:"-" json:"worklog"`                   // List of records
}

func (u User) AddUser(db *sql.DB) (int64, error) {
	log.Printf("[USER] Adding user %s to the database", u.Username)

	q := `INSERT INTO practice_db.users (firstname, lastname, email, username, password) VALUES(?, ?, ?, ?, ?)`
	result, err := db.Exec(q, u.Firstname, u.Lastname, u.Email, u.Username, u.Password)
	if err != nil {
		return -1, common.GetQueryError(q, "Couldn't register a new user", u, err)
	}

	id, err := common.GetLastInsertId(result, q, u)
	if err != nil {
		return -1, err
	}

	log.Printf("[USER] User %s has been successfully added to the database", u.Username)
	return id, nil
}

func (u *User) RetrieveUserbyID(db *sql.DB) error {
	log.Printf("[USER] Retrieving user by id %d", u.ID)

	q := `SELECT firstname, lastname, email, username, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE id=?`
	err := db.QueryRow(q, u.ID).Scan(
		&u.Firstname, &u.Lastname, &u.Email, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	err = common.GetQueryError(q, "Couldn't retrieve user by id", u, err)
	if err != nil {
		return err
	}

	log.Printf("[USER] User %s has been successfully retrieved", u.Username)
	return nil
}

func (u *User) RetrieveUserByEmail(db *sql.DB) error {
	log.Printf("[USER] Retrieving user by email %s", u.Email)

	q := `SELECT id, firstname, lastname, username, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE email=?`
	err := db.QueryRow(q, u.Email).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	err = common.GetQueryError(q, "Couldn't retrieve user by email", u, err)
	if err != nil {
		return err
	}

	log.Printf("[USER] User %s has been successfully retrieved", u.Username)
	return nil
}

func (u *User) RetrieveUserByUsername(db *sql.DB) error {
	log.Printf("[USER] Retrieving user by username %s", u.Username)

	q := `SELECT id, firstname, lastname, email, log_count, total_hours, total_minutes
		FROM practice_db.users
		WHERE username=?`
	err := db.QueryRow(q, u.Username).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.LogCount, &u.TotalHours, &u.TotalMinutes)

	err = common.GetQueryError(q, "Couldn't retrieve user by username", u, err)
	if err != nil {
		return err
	}

	log.Printf("[USER] User %s has been successfully retrieved", u.Username)
	return nil
}

func (u *User) RetrieveUserIDByEmail(db *sql.DB) error {
	log.Printf("[USER] Retrieving user id by email %s", u.Email)

	u.ID = -1
	q := `SELECT id FROM practice_db.users WHERE email=?`
	err := db.QueryRow(q, u.Email).Scan(&u.ID)

	err = common.GetQueryError(q, "Couldn't retrieve user id by email", u, err)
	if err != nil {
		return err
	}

	log.Printf("[USER] User id %d has been successfully retrieved", u.ID)
	return nil
}

// New, not tested
func (u *User) RetrieveUserIDByUsername(db *sql.DB) error {
	log.Printf("[USER] Retrieving user id by username %s", u.Username)
	q := `SELECT id FROM practice_db.users WHERE username=?`
	err := db.QueryRow(q, u.Username).Scan(&u.ID)

	err = common.GetQueryError(q, "Couldn't retrieve user id by username", u, err)
	if err != nil {
		return err
	}

	log.Printf("[USER] User id %d has been successfully retrieved", u.ID)
	return nil
}

// New, not tested
// RetrieveUserDTOByCred retrieves userDTO by email or username.
func (u *UserDTO) RetrieveUserDTOByCred(db *sql.DB) error {
	log.Printf("[USER] Retrieving userDTO by email: %s, or by username: %s", u.Email, u.Username)

	q := `SELECT id, firstname, lastname, email, username FROM practice_db.users WHERE email=? OR username=?`
	err := db.QueryRow(q, u.Email, u.Username).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Username)

	err = common.GetQueryError(q, "Couldn't retrieve userDTO by email or username", u, err)
	if err != nil {
		return err
	}

	log.Printf("[USER] UserDTO %s has been successfully retrieved", u.Username)
	return nil
}

// UpdateUserWorklogInfoByID changes the information about the user's worklog by user id.
// Precisely: log count, and hours & minutes worked.
func (u User) UpdateUserWorklogInfoByID(tx *sql.Tx) error {
	log.Printf("[USER] Updating user %s worklog info", u.Username)

	q := `UPDATE practice_db.users
		SET log_count=?, total_hours=?, total_minutes=?
		WHERE id=?`
	result, err := tx.Exec(q, u.LogCount, u.TotalHours, u.TotalMinutes, u.ID)

	err = common.HandleUpdateQuery(result, err, q, u)
	if err != nil {
		return err
	}

	log.Printf("[USER] User %s worklog info has been successfully updated", u.Username)
	return nil
}

// RetrieveAllRecordsByUserID scans records table,
// looking for every record associated with the user.
func (u *User) RetrieveAllRecordsByUserID(db *sql.DB) error {
	log.Printf("[USER] Retrieving all records by user id %d", u.ID)

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
	log.Printf("[USER] All records by user id %d have been successfully retrieved", u.ID)
	return nil
}
