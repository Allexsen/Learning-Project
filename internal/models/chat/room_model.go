package chat

import "github.com/Allexsen/Learning-Project/internal/models/ws"

// Room  represents a group chat
type Room struct {
	BaseChat
	Name string `json:"name"`
}

// NewRoom creates a new room
func NewRoom(name string, manager *ws.WsManager) *Room {
	return &Room{
		BaseChat: *NewBaseChat(manager),
		Name:     name,
	}
}

// GetRooms returns all rooms
func GetRooms() ([]*Room, error) {
	// TODO: Add Rooms table to the database
	return []*Room{}, nil
}

// GetRoomByID gets a room by its ID
func (room *Room) GetRoomByID() error {
	// TODO: Add room fetching logic
	return nil
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
