// Package router sets up API endpoints to serve requests
package router

import "github.com/Allexsen/Learning-Project/internal/handlers"

// initUserRouter sets up routes associated with users
func initUserRouter() {
	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", handlers.UserRegister)
		userRouter.POST("/login", handlers.UserLogin)
		userRouter.POST("/retrieve", handlers.UserGet)
		userRouter.POST("/check-email", handlers.IsAvailableEmail)
		userRouter.POST("/check-username", handlers.IsAvailableUsername)
	}
}
