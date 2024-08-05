package controllers

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/chat"
	"github.com/Allexsen/Learning-Project/internal/utils"
)

// RoomCreate creates a new room
func RoomCreate(name string) *chat.Room {
	log.Printf("[CONTROLLER] Creating room %s", name)

	// TODO: Validate room name
	room := chat.NewRoom(name)

	log.Printf("[CONTROLLER] Room %s has been successfully created", name)
	return room
}

// RoomsGet retrieves all rooms
func RoomsGet() ([]*chat.Room, error) {
	log.Printf("[CONTROLLER] Getting all rooms")

	rooms, err := chat.GetRooms()

	log.Printf("[CONTROLLER] %d rooms have been successfully retrieved", len(rooms))
	return rooms, err
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

	log.Printf("[CONTROLLER] Room %s has been successfully retrieved", idStr)
	return room, nil
}

// RoomAddUser adds a user to a room
func RoomAddUser(roomIDStr, userIDStr string) (*chat.Room, error) {
	log.Printf("[CONTROLLER] Adding user %s to room %s", userIDStr, roomIDStr)
	userID, err := utils.Atoi(userIDStr)
	if err != nil {
		return nil, err
	}

	room, err := RoomGet(roomIDStr)
	if err != nil {
		return nil, err
	}

	// TODO: Implement retrieving UserDTO by UserID, and passing it to AddUser
	err = room.AddUser(int64(userID))
	if err != nil {
		return nil, err
	}

	log.Printf("[CONTROLLER] User %s has been successfully added to room %s", userIDStr, roomIDStr)
	return room, nil
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

	log.Printf("[CONTROLLER] Room %s has been successfully removed", idStr)
	return nil
}
