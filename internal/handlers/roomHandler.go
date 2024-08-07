package handlers

import (
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	"github.com/Allexsen/Learning-Project/internal/models/chat"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateRoom handles the request to create a new room
func CreateRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room creation request for %s", c.ClientIP())

	var reqData chat.Room
	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	log.Printf("[HANDLER] Request Data: %+v", reqData)

	room, err := controllers.RoomCreate(reqData.Name)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Room %s has been successfully created", room.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"room":    room,
	})
}

// GetRooms handles the request to retrieve all rooms
func GetRooms(c *gin.Context) {
	log.Printf("[HANDLER] Handling room retrieval request for %s", c.ClientIP())

	rooms, err := controllers.RoomsGet()
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Rooms have been successfully retrieved")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"rooms":   rooms,
	})
}

// JoinRoom handles the request to join a room
func JoinRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room join request for %s", c.ClientIP())

	roomID := c.Param("id")
	var userDTO user.UserDTO
	if !utils.ShouldBindJSON(c, &userDTO) {
		return
	}

	log.Printf("[HANDLER] Request Data: RoomID: %s, User: %+v", roomID, userDTO)
	room, err := controllers.RoomAddUser(roomID, userDTO)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] User %+v has successfully joined room %s", userDTO, roomID) // Change to UserDTO
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"room":    room,
	})
}

// DeleteRoom handles the request to delete a room
func DeleteRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room deletion request for %s", c.ClientIP())

	roomID := c.Param("id")

	log.Printf("[HANDLER] Request Data: RoomID: %s", roomID)

	err := controllers.RoomRemove(roomID)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Room %s has been successfully deleted", roomID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
