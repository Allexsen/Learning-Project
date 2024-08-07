// Package db sets up a connection to a database
package database

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	_ "github.com/go-sql-driver/mysql"
)

// DB variable provides a way for the project to interact with the database
var DB *sql.DB

// InitDB connects to the db, and sets a connection parameters.
// Invokes a critical error in case of failure.
func InitDB() {
	log.Println("Connecting to the database...")

	var err error
	DB, err = sql.Open("mysql", os.Getenv("MYSQL_AUTH_CREDS"))
	if err != nil {
		apperrors.HandleCriticalError(apperrors.New(
			http.StatusInternalServerError,
			"Could not open a database connection",
			apperrors.ErrDBConnection,
			map[string]interface{}{"details": err.Error()},
		))
	}

	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetConnMaxIdleTime(time.Minute * 1)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	log.Println("Connected to the database")
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) {
	log.Println("Closing the database connection...")

	err := db.Close()
	if err != nil {
		apperrors.HandleCriticalError(apperrors.New(
			http.StatusInternalServerError,
			"Could not close a database connection",
			apperrors.ErrDBConnection,
			map[string]interface{}{"details": err.Error()},
		))
	}

	log.Println("Database connection closed")
}
