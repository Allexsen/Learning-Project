package handlers

import (
	"log"
	"net/http"

	"github.com/Allexsen/Learning-Project/internal/controllers"
	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/chat/room"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateRoom handles the request to create a new room
func CreateRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room creation request for %s", c.ClientIP())

	var reqData room.Room
	if !utils.ShouldBindJSON(c, &reqData) {
		return
	}

	room, err := controllers.RoomCreate(reqData.Name)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Room %v has been successfully created", room)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"room":    &room,
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

	log.Printf("[HANDLER] %d rooms have been successfully retrieved", len(rooms))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"rooms":   &rooms,
	})
}

// GetRoom handles the request to join a room
func GetRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room get request for %s", c.ClientIP())

	roomID := c.Param("id")
	room, err := controllers.RoomGet(roomID)
	if err != nil {
		handleError(c, apperrors.ErrNotFound)
		return
	}

	log.Printf("[HANDLER] Room %s has been successfully retrieved", roomID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"room":    &room,
	})
}

// JoinRoom handles the WebSocket connection request
func JoinRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room join request for %s", c.ClientIP())

	roomID := c.Param("id")
	userDTO, exists := c.Get("userDTO")
	if !exists {
		handleError(c, apperrors.ErrInternalServerError)
		return
	}

	err := controllers.RoomJoin(c, roomID, userDTO.(*user.UserDTO))
	if err != nil {
		handleError(c, err)
		return
	}

	log.Printf("[HANDLER] Room %s has been successfully joined", roomID)
}

// DeleteRoom handles the request to delete a room
func DeleteRoom(c *gin.Context) {
	log.Printf("[HANDLER] Handling room deletion request for %s", c.ClientIP())

	roomID := c.Param("id")

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
