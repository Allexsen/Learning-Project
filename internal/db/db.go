package db

import (
	"database/sql"
	"log"
	"os"
	"time"
)

var DB *sql.DB

func InitDB() {
	DB, err := sql.Open("mysql", os.Getenv("MYSQL_AUTH_CREDS"))
	if err != nil {
		log.Fatal("[Error]: Coul dn't open a database connection: " + err.Error())
		panic(err)
	}

	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetConnMaxIdleTime(time.Minute * 1)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
}
