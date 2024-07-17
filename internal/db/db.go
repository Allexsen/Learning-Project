// Package db sets up a connection to a database
package database

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	customErrors "github.com/Allexsen/Learning-Project/internal/errors"
)

// DB variable provides a way for the project to interact with the database
var DB *sql.DB

// InitDB connects to the db, and sets a connection parameters.
// Invokes a critical error in case of failure.
func InitDB() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("MYSQL_AUTH_CREDS"))
	if err != nil {
		customErrors.HandleCriticalError(customErrors.New(
			http.StatusInternalServerError,
			"Could not open a database connection",
			customErrors.ErrDBConnection,
			map[string]interface{}{"details": err.Error()},
		))
	}

	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetConnMaxIdleTime(time.Minute * 1)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
}
