package chat

import (
	"fmt"
	"log"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
)

// TODO: Design and add database table for rooms, or completely remove db logic
// TODO: Change to UserDTO and implement proper retrievals

// Room  represents a group chat
type Room struct {
	BaseChat
	Name string `json:"name,omitempty"`
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
	log.Printf("[CHAT] Creating room %s", name)

	manager := ws.NewManager()
	go manager.Run()
	room := &Room{
		BaseChat: *NewBaseChat(manager),
		Name:     name,
	}

	log.Printf("[CHAT] Room has been successfully created: %+v", room)
	roomsManager.Rooms[room.ID] = room
	return room
}

// GetRooms returns all rooms
func GetRooms() ([]*Room, error) {
	// TODO: Add Rooms table to the storage and/or database
	log.Printf("[CHAT] Getting all rooms")

	rooms := make([]*Room, 0, len(roomsManager.Rooms))
	for _, room := range roomsManager.Rooms {
		rooms = append(rooms, room)
	}

	log.Printf("[CHAT] %d rooms have been successfully retrieved", len(rooms))
	return rooms, nil
}

func GetRoomByID(id int64) (*Room, error) {
	log.Printf("[CHAT] Getting room %d", id)

	room, exists := roomsManager.Rooms[id]
	if !exists {
		return nil, apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("Room with ID %d not found", id),
			apperrors.ErrNotFound,
			nil,
		)
	}

	log.Printf("[CHAT] Room %d has been successfully retrieved", id)
	return room, nil
}

// AddUser adds a user to the room
func (room *Room) AddUser(userID int64) error {
	// TODO: Add database logic
	log.Printf("[CHAT] Adding user %d to room %d", userID, room.ID)

	room.Members = append(room.Members, userID)
	room.Manager.AddClient(userID)

	log.Printf("[CHAT] User %d has been successfully added to room %d", userID, room.ID)
	return nil
}

// DeleteRoom deletes a room
func (room *Room) DeleteRoom() error {
	log.Printf("[CHAT] Removing room %d", room.ID)

	_, exists := roomsManager.Rooms[room.ID]
	if !exists {
		return apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("Room with ID %d not found", room.ID),
			apperrors.ErrNotFound,
			nil,
		)
	}

	room.Manager.Close()
	err := removeRoomFromDB(room.ID)
	if err != nil {
		return err
	}

	delete(roomsManager.Rooms, room.ID)
	log.Printf("[CHAT] Room %d has been successfully removed", room.ID)

	room = nil
	return nil
}

// RemoveRoomFromDB removes a room from the database
func removeRoomFromDB(_ int64) error {
	// TODO: Add database logic
	log.Printf("[CHAT] Removing room from the database")

	log.Printf("[CHAT] Room has been successfully removed from the database")
	return nil
}
