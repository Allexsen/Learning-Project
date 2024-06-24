package main

import (
	config "github.com/Allexsen/Learning-Project/cmd/config"
	db "github.com/Allexsen/Learning-Project/internal/db"
	router "github.com/Allexsen/Learning-Project/internal/router"
)

func main() {
	config.LoadConfig()

	db.InitDB()
	defer db.DB.Close()

	router.InitRouter()
}
