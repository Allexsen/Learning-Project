package controllers

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/chat"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// RoomCreate creates a new room
func RoomCreate(name string) (*chat.Room, error) {
	log.Printf("[CONTROLLER] Creating room %s", name)

	if err := utils.IsValidName(name); err != nil {
		return nil, err
	}

	room := chat.NewRoom(name)
	return room, nil
}

// RoomsGet retrieves all rooms
func RoomsGet() ([]*chat.Room, error) {
	log.Printf("[CONTROLLER] Getting all rooms")

	rooms, err := chat.GetRooms()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

// RoomGet retrieves a room by ID
func RoomGet(idStr string) (*chat.Room, error) {
	log.Printf("[CONTROLLER] Getting room %s", idStr)

	id, err := utils.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	room, err := chat.GetRoomByID(int64(id))
	if err != nil {
		return nil, err
	}

	return room, nil
}

// RoomAddUser adds a user to a room
func RoomAddUser(c *gin.Context, roomIDStr string, userDTO user.UserDTO) (*chat.Room, error) {
	log.Printf("[CONTROLLER] Adding user %+v to room %s", userDTO, roomIDStr)
	room, err := RoomGet(roomIDStr)
	if err != nil {
		return nil, err
	}

	err = room.AddClient(c, userDTO)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func RoomJoin(c *gin.Context, roomIDStr string, userDTO *user.UserDTO) error {
	log.Printf("[CONTROLLER] Joining user %+v to room %s", userDTO, roomIDStr)

	room, err := RoomGet(roomIDStr)
	if err != nil {
		return err
	}

	err = room.AddClient(c, *userDTO)
	if err != nil {
		return err
	}

	return nil
}

// RoomRemove removes a room and all associated data by room  ID.
// If the room doesn't exist, apperrors.ErrNotFound is returned
func RoomRemove(idStr string) error {
	log.Printf("[CONTROLLER] Removing room %s", idStr)

	room, err := RoomGet(idStr)
	if err != nil {
		return err
	}

	err = room.DeleteRoom()
	if err != nil {
		return err
	}

	return nil
}
