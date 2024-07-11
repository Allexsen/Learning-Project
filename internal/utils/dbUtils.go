package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/db"
	"github.com/Allexsen/Learning-Project/internal/models"
	"github.com/gin-gonic/gin"
)

func RetrieveAllRecordsByUserID(uid int64) ([]models.Record, error) {
	rows, err := db.DB.Query(`SELECT id, hours, minutes FROM practice_db.records WHERE user_id=?`, uid)
	if err != nil {
		return []models.Record{}, fmt.Errorf("couldn't query the database: %v", err)
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		r := models.Record{UserID: uid}
		err = rows.Scan(&r.ID, &r.Hours, &r.Minutes)
		if err != nil {
			return []models.Record{}, fmt.Errorf("couldn't scan a row: %v", err)
		}

		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return []models.Record{}, fmt.Errorf("error during the iteration of the rows: %v", err)
	}

	return records, nil
}

func GetPasswordHashByUsername(username string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE username=?`
	var pswdHash string
	err := db.DB.QueryRow(q, username).Scan(&pswdHash)
	if err != nil {
		return "", fmt.Errorf("couldn't scan a row: %v", err)
	}

	return pswdHash, nil
}

func GetPasswordHashByEmail(email string) (string, error) {
	q := `SELECT password FROM practice_db.users WHERE email=?`
	var pswdHash string
	err := db.DB.QueryRow(q, email).Scan(&pswdHash)
	if err != nil {
		return "", fmt.Errorf("couldn't scan a row: %v", err)
	}

	return pswdHash, nil
}

func IsExistingEmail(email string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE email=?)`
	err := db.DB.QueryRow(q, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("couldn't scan a row: %v", err)
	}

	return exists, nil
}

func IsExistingUsername(username string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS (SELECT 1 FROM practice_db.users WHERE username=?)`
	err := db.DB.QueryRow(q, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("couldn't scan a row: %v", err)
	}

	return exists, nil
}

func IsExistingCreds(email, username string, c *gin.Context) bool {
	if email != "" {
		exists, err := IsExistingEmail(email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking email availability"})
			log.Printf("couldn't check email availability: %v", err)
			return true
		}

		if exists {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			log.Printf("couldn't register a user, existing email: %v", email)
			return true
		}
	}

	if username != "" {
		exists, err := IsExistingUsername(username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking username availability"})
			log.Printf("couldn't check username availability: %v", err)
			return true
		}

		if exists {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
			log.Printf("couldn't register a user, existing username: %v", username)
			return true
		}
	}

	return false
}
