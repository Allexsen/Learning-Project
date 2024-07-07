package router

import "github.com/Allexsen/Learning-Project/internal/handlers"

func initUserRouter() {
	r.POST("/user/register", handlers.UserRegister)
	r.POST("/user/login", handlers.UserLogin)
	r.POST("/user/retrieve", handlers.UserGet)
}
