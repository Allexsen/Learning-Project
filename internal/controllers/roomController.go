package controllers

import (
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/chat"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/Allexsen/Learning-Project/internal/utils"
)

func RoomCreate(name string, manager *ws.WsManager) *chat.Room {
	// TODO: Check back on this later on, seems like some part of logic is missing
	room := chat.NewRoom(name, manager)
	log.Printf("Room: %v", room)
	return room
}

func RoomsGet() ([]*chat.Room, error) {
	// TODO: Check back on this later on, seems like some part of logic is missingo
	rooms, err := chat.GetRooms()
	return rooms, err
}

func RoomGet(idStr string) (*chat.Room, error) {
	id, err := utils.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	room, err := chat.GetRoomByID(int64(id))
	log.Printf("Room: %v", room)
	if err != nil {
		return nil, err
	}

	log.Printf("Room: %v", room)

	return room, nil
}

func RoomAddUser(roomIDStr, userIDStr string) (*chat.Room, error) {
	userID, err := utils.Atoi(userIDStr)
	if err != nil {
		return nil, err
	}

	room, err := RoomGet(roomIDStr)
	if err != nil {
		return nil, err
	}

	log.Printf("Manager: %v", room.GetManager())

	// TODO: Implement retrieving UserDTO by UserID, and passing it to AddUser
	err = room.AddUser(int64(userID))
	if err != nil {
		return nil, err
	}

	return room, nil
}

func RoomRemove(idStr string) error {
	room, err := RoomGet(idStr)
	if err != nil {
		return err
	}

	return room.DeleteRoom()
}
