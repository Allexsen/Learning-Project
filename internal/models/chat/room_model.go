package chat

import (
	"fmt"
	"log"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
)

// Room  represents a group chat
type Room struct {
	BaseChat
	Name string `json:"name"`
}

type RoomsManager struct {
	Rooms map[int64]*Room
}

var (
	roomsManager = &RoomsManager{
		Rooms: make(map[int64]*Room),
	}
)

// NewRoom creates a new room
func NewRoom(name string) *Room {
	manager := ws.NewManager()
	go manager.Run()
	room := &Room{
		BaseChat: *NewBaseChat(manager),
		Name:     name,
	}

	log.Println(room.ID) // Temporary log to join the room
	roomsManager.Rooms[room.ID] = room
	return room
}

// GetRooms returns all rooms
func GetRooms() ([]*Room, error) {
	// TODO: Add Rooms table to the storage and/or database
	rooms := make([]*Room, 0, len(roomsManager.Rooms))
	for _, room := range roomsManager.Rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func GetRoomByID(id int64) (*Room, error) {
	room, exists := roomsManager.Rooms[id]
	if !exists {
		return nil, apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("Room with ID %d not found", id),
			apperrors.ErrNotFound,
			nil,
		)
	}

	return room, nil
}

// AddUser adds a user to the room
func (room *Room) AddUser(userID int64) error {
	// TODO: Add database logic
	room.Members = append(room.Members, userID)
	room.Manager.AddClient(userID)
	return nil
}

// DeleteRoom deletes a room
func (room *Room) DeleteRoom() error {
	// TODO: Implement removing the room from the roomsManager
	room.Manager.Close()
	err := removeRoomFromDB(room.ID)
	if err != nil {
		return err
	}

	room = nil
	return nil
}

// RemoveRoomFromDB removes a room from the database
func removeRoomFromDB(_ int64) error {
	// TODO: Add database logic
	return nil
}