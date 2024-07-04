package config

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[Error]: Couldn't load .env file: " + err.Error())
		panic(err)
	}
}
