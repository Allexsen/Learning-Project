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
		roomRouter.GET("/:id", handlers.GetRoom)
		roomRouter.GET("/:id/participants", handlers.GetRoomParticipants)
		roomRouter.POST("/new", handlers.CreateRoom)
		roomRouter.GET("/:id/ws", handlers.JoinRoom)
		roomRouter.POST("/join/:id", handlers.GetRoom)
		roomRouter.DELETE("/remove/:id", handlers.DeleteRoom)
	}
}
