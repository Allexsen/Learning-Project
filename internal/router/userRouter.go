// Package router sets up API endpoints to serve requests
package router

import "github.com/Allexsen/Learning-Project/internal/handlers"

// initUserRouter sets up routes associated with users
func initUserRouter() {
	r.POST("/user/register", handlers.UserRegister)
	r.POST("/user/login", handlers.UserLogin)
	r.POST("/user/retrieve", handlers.UserGet)
	r.POST("/user/check-email", handlers.IsAvailableEmail)
	r.POST("/user/check-username", handlers.IsAvailableUsername)
}
