package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	DB *sql.DB
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[Error]: Couldn't load .env file: " + err.Error())
		panic(err)
	}
}
