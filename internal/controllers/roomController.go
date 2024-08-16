package controllers

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/chats/room"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
)

// RoomCreate creates a new room
func RoomCreate(name string) (*room.Room, error) {
	log.Printf("[CONTROLLER] Creating room %s", name)

	if err := utils.IsValidName(name); err != nil {
		return nil, err
	}

	room := room.NewRoom(name)
	return room, nil
}

// RoomsGet retrieves all rooms
func RoomsGet() ([]*room.Room, error) {
	log.Printf("[CONTROLLER] Getting all rooms")

	rooms, err := room.GetRooms()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

// RoomGet retrieves a room by ID
func RoomGet(idStr string) (*room.Room, error) {
	log.Printf("[CONTROLLER] Getting room %s", idStr)

	id, err := utils.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	room, err := room.GetRoomByID(int64(id))
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

	err = room.AddClient(c, userDTO)
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
