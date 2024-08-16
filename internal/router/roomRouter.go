package router

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/handlers"
	"github.com/Allexsen/Learning-Project/internal/middlewares"
)

// initRoomRouter sets up routes associated with rooms
func initRoomRouter() {
	log.Println("Setting up room routes...")
	roomRouter := r.Group("/rooms")
	roomRouter.Use(middlewares.CheckJWT())
	{
		roomRouter.GET("", handlers.GetRooms)
		roomRouter.GET("/room/:id", handlers.GetRoom)
		roomRouter.POST("/new", handlers.CreateRoom)
		roomRouter.GET("/:id/ws", handlers.JoinRoom)
		roomRouter.POST("/join/:id", handlers.GetRoom)
		// TODO: Implement: roomRouter.POST("/leave/:id", handlers.LeaveRoom)
		roomRouter.DELETE("/remove/:id", handlers.DeleteRoom)
	}
}
