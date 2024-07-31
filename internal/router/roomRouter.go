package router

import (
	"github.com/Allexsen/Learning-Project/internal/handlers"
)

func initRoomRouter() {
	roomRouter := r.Group("/room")
	{
		roomRouter.GET("/", handlers.GetRooms)
		roomRouter.POST("/new", handlers.CreateRoom)
		roomRouter.POST("/join/:id", handlers.JoinRoom)
		roomRouter.DELETE("/remove/:id", handlers.DeleteRoom)
	}
}
