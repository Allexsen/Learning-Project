package router

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/handlers"
)

func initRoomRouter() {
	log.Println("Setting up room routes...")
	roomRouter := r.Group("/rooms")
	{
		roomRouter.GET("", handlers.GetRooms)
		roomRouter.POST("/new", handlers.CreateRoom)
		roomRouter.POST("/join/:id", handlers.JoinRoom)
		roomRouter.DELETE("/remove/:id", handlers.DeleteRoom)
	}
}
