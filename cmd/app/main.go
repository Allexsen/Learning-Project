// This is an ever-changing practice project. There is no end-goal but to keep adding
// New features along the way as I study new stuff. For now, the following is applicable
// However, most likely will change into a social media mock project, as I start learning
// Kafka, real time updates and similar topics.
//
// This is a comprehensive tool designed to manage user registrations, logins, and
// other user-related functionalities. It provides a RESTful API for interacting with
// the user data stored in a relational database.
//
// Key functionalities include:
// - Registering new users
// - User authentication and login
// - Retrieving user information
// - Handling user worklog records
//
// This application is built using the Go programming language and follows best practices
// for security and error handling(I'm trying my best :D)
package main

import (
	config "github.com/Allexsen/Learning-Project/cmd/config"
	database "github.com/Allexsen/Learning-Project/internal/db"
	router "github.com/Allexsen/Learning-Project/internal/router"
)

func main() {
	config.LoadEnv()

	database.InitDB()
	defer database.CloseDB(database.DB)

	router.InitRouter()
}
