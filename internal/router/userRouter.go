package router

import "github.com/Allexsen/Learning-Project/internal/handlers"

func initUserRouter() {
	r.POST("/user/retrieve", handlers.UserGet())
}
