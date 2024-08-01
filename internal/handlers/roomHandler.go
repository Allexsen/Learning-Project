package handlers

import (
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/Allexsen/Learning-Project/internal/models/chat"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	var reqData chat.Room
	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	manager := ws.NewManager()
	go manager.Run()
	room := controllers.RoomCreate(reqData.Name, manager)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"room":    room,
	})
}

func GetRooms(c *gin.Context) {
	rooms, err := controllers.RoomsGet()
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"rooms":   rooms,
	})
}

// TODO: Rewise error handling
func JoinRoom(c *gin.Context) {
	roomID := c.Param("id")
	userID := "1" // TODO: Retrieve UserDTO from gin.Context
	room, err := controllers.RoomAddUser(roomID, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"room":    room,
	})
}

func DeleteRoom(c *gin.Context) {
	roomID := c.Param("id")
	err := controllers.RoomRemove(roomID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
