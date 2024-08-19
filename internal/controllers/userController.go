// Package controllers is responsible for handler-model interaction, and business logic
package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	database "github.com/Allexsen/Learning-Project/internal/db"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserRegister takes user info, generates bcrypt
// password hash, and adds the user to the database
func UserRegister(firstname, lastname, username, email, pswd string) (*user.User, error) {
	log.Printf("[CONTROLLER] Registering user %s", username)

	log.Printf("[CONTROLLER] Generating bcrypt hash for user %s", username)
	pswdHash, err := bcrypt.GenerateFromPassword([]byte(pswd), 10)
	if err != nil {
		return nil, err
	}

	db := database.DB
	u := user.User{
		UserDTO: user.UserDTO{
			Firstname: firstname,
			Lastname:  lastname,
			Username:  username,
			Email:     email,
		},
		Password: string(pswdHash),
	}

	// Query adding to the db
	if u.ID, err = u.AddUser(db); err != nil {
		return nil, err
	}

	return &u, nil
}

// UserLogin retrieves bcrypt password hash, authenticates
// user credentials and, if successful, logs in the user
func UserLogin(credential, password string) (*user.UserDTO, error) {
	log.Printf("[CONTROLLER] Authenticating user %s", credential)
	var pswdHash string
	var err error
	db := database.DB
	userDTO := &user.UserDTO{}

	// Check if the provided credentail is email or username, and query accordingly
	if strings.Contains(credential, "@") {
		pswdHash, err = utils.GetPasswordHashByEmail(db, credential)
		userDTO.Email = credential
	} else {
		pswdHash, err = utils.GetPasswordHashByUsername(db, credential)
		userDTO.Username = credential
	}

	if err != nil {
		return nil, err
	}

	log.Printf("[CONTROLLER] Comparing password hash for user %s", credential)
	err = bcrypt.CompareHashAndPassword([]byte(pswdHash), []byte(password))
	if err != nil {
		return nil, apperrors.New(
			http.StatusUnauthorized,
			"Invalid credentials",
			apperrors.ErrUnauthorized,
			map[string]interface{}{"details": err.Error()},
		)
	}

	// Retrieve user info by email or username
	err = userDTO.RetrieveUserDTOByCred(db)
	if err != nil {
		return nil, err
	}

	return userDTO, nil
}

// UserGetByEmail retrieves user from the database by user email
func UserGetByEmail(email string) (*user.User, error) {
	log.Printf("[CONTROLLER] Retrieving user by email %s", email)
	db := database.DB
	u := user.User{
		UserDTO: user.UserDTO{
			Email: email,
		},
	}

	if err := u.RetrieveUserByEmail(db); err != nil {
		return nil, err
	}

	return &u, nil
}

// UserGetByUsername retrieves user from the database by username
func UserGetByUsername(username string) (*user.User, error) {
	log.Printf("[CONTROLLER] Retrieving user by username %s", username)
	db := database.DB
	u := user.User{
		UserDTO: user.UserDTO{
			ID:       -1,
			Username: username,
		},
	}

	if err := u.RetrieveUserByUsername(db); err != nil {
		return nil, err
	}

	return &u, nil
}

// UserGetIDByEmail retrieves user ID by user email.
// return -1 and error in case of failure
func UserGetIDByEmail(db *sql.DB, email string) (int64, error) {
	log.Printf("[CONTROLLER] Retrieving user ID by email %s", email)
	u := user.User{
		UserDTO: user.UserDTO{
			ID:    -1,
			Email: email,
		},
	}

	err := u.RetrieveUserIDByEmail(db)
	return u.ID, err
}
