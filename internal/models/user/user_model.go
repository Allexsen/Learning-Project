// Package user provides user model and its methods for database interaction.
package user

import (
	"database/sql"
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/common"
)

// UserDTO represents a user data transfer object.
// It is used to transfer semi-public user data between the client and the server.
type UserDTO struct {
	ID        int64  `db:"id" json:"id,omitempty"`               // Unique user id
	Firstname string `db:"firstname" json:"firstName,omitempty"` // Firstname
	Lastname  string `db:"lastname" json:"lastName,omitempty"`   // Lastname
	Email     string `db:"email" json:"email,omitempty"`         // Email
	Username  string `db:"username" json:"username,omitempty"`   // Unique username
}

// User represents a user model.
// It is used to store user data in the database.
type User struct {
	UserDTO
	Password string `db:"password" json:"-"` // Password hash
}

// AddUser adds a new user to the database.
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

	return id, nil
}

// RetrieveUserbyID retrieves user by user id.
func (u *User) RetrieveUserbyID(db *sql.DB) error {
	log.Printf("[USER] Retrieving user by id %d from the database", u.ID)

	q := `SELECT firstname, lastname, email, username,
		FROM practice_db.users
		WHERE id=?`
	err := db.QueryRow(q, u.ID).Scan(
		&u.Firstname, &u.Lastname, &u.Email, &u.Username)

	err = common.GetQueryError(q, "Couldn't retrieve user by id", u, err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserByEmail retrieves user by email.
func (u *User) RetrieveUserByEmail(db *sql.DB) error {
	log.Printf("[USER] Retrieving user by email %s from the database", u.Email)

	q := `SELECT id, firstname, lastname, username,
		FROM practice_db.users
		WHERE email=?`
	err := db.QueryRow(q, u.Email).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username)

	err = common.GetQueryError(q, "Couldn't retrieve user by email", u, err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserByUsername retrieves user by username.
func (u *User) RetrieveUserByUsername(db *sql.DB) error {
	log.Printf("[USER] Retrieving user by username %s from the database", u.Username)

	q := `SELECT id, firstname, lastname, email
		FROM practice_db.users
		WHERE username=?`
	err := db.QueryRow(q, u.Username).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email)

	err = common.GetQueryError(q, "Couldn't retrieve user by username", u, err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserIDByEmail retrieves user id by email.
func (u *User) RetrieveUserIDByEmail(db *sql.DB) error {
	log.Printf("[USER] Retrieving user id by email %s from the database", u.Email)

	u.ID = -1
	q := `SELECT id FROM practice_db.users WHERE email=?`
	err := db.QueryRow(q, u.Email).Scan(&u.ID)

	err = common.GetQueryError(q, "Couldn't retrieve user id by email", u, err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserIDByUsername retrieves user id by username.
func (u *User) RetrieveUserIDByUsername(db *sql.DB) error { // TODO: Implement Unit Tests
	log.Printf("[USER] Retrieving user id by username %s from the database", u.Username)
	q := `SELECT id FROM practice_db.users WHERE username=?`
	err := db.QueryRow(q, u.Username).Scan(&u.ID)

	err = common.GetQueryError(q, "Couldn't retrieve user id by username", u, err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserDTOByID retrieves userDTO by user id.
func (u *UserDTO) RetrieveUserDTOByID(db *sql.DB) error { // TODO: Implement Unit Tests
	log.Printf("[USER] Retrieving userDTO by id %d from the database", u.ID)

	q := `SELECT firstname, lastname, email, username FROM practice_db.users WHERE id=?`
	err := db.QueryRow(q, u.ID).Scan(&u.Firstname, &u.Lastname, &u.Email, &u.Username)

	err = common.GetQueryError(q, "Couldn't retrieve userDTO by id", u, err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserDTOByCred retrieves userDTO by email or username.
func (u *UserDTO) RetrieveUserDTOByCred(db *sql.DB) error { // TODO: Implement Unit Tests
	cred := u.Email
	if cred == "" {
		cred = u.Username
	}

	log.Printf("[USER] Retrieving userDTO by credential %s from the database", cred)

	q := `SELECT id, firstname, lastname, email, username FROM practice_db.users WHERE email=? OR username=?`
	err := db.QueryRow(q, u.Email, u.Username).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email, &u.Username)

	err = common.GetQueryError(q, "Couldn't retrieve userDTO by email or username", u, err)
	if err != nil {
		return err
	}

	return nil
}
