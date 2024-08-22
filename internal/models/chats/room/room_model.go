// Package room provides room model, room ws, and its methods for database interaction.
package room

import (
	"fmt"
	"log"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	chat "github.com/Allexsen/Learning-Project/internal/models/chats"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/gin-gonic/gin"
)

// Room  represents a group chat
type Room struct {
	chat.BaseChat
	Name    string     `json:"name,omitempty"`
	Manager *wsManager `json:"-"` // Manager for the room
}

var roomsList map[int64]*Room

func init() {
	roomsList = make(map[int64]*Room)
}

// NewRoom creates a new room
func NewRoom(name string) *Room {
	room := &Room{
		BaseChat: *chat.NewBaseChat(),
		Name:     name,
	}

	manager := newWsManager(room)
	go manager.Run()
	room.Manager = manager

	roomsList[room.ID] = room
	return room
}

// GetRooms returns all rooms
func GetRooms() ([]*Room, error) {
	rooms := make([]*Room, 0, len(roomsList))
	for _, room := range roomsList {
		newRoom := Room{
			BaseChat: chat.BaseChat{
				ID:        room.ID,
				CreatedAt: room.CreatedAt,
				UpdatedAt: room.UpdatedAt,
				Messages:  room.Messages,
			},
			Name:    room.Name,
			Manager: room.Manager,
		}

		rooms = append(rooms, &newRoom)
	}

	return rooms, nil
}

func GetRoomByID(id int64) (*Room, error) {
	room, exists := roomsList[id]
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

// GetParticipants returns all participants in the room
func (room *Room) GetParticipants() []*user.UserDTO {
	members := make([]*user.UserDTO, 0, len(room.Manager.Clients))
	for client := range room.Manager.Clients {
		members = append(members, client.UserDTO)
	}

	return members
}

// DeleteRoom deletes a room
func (room *Room) DeleteRoom() error {
	_, exists := roomsList[room.ID]
	if !exists {
		return apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("Room with ID %d not found", room.ID),
			apperrors.ErrNotFound,
			nil,
		)
	}

	room.Manager.close()
	delete(roomsList, room.ID)
	room = nil
	return nil
}

// AddClient adds a client to the room
func (room *Room) AddClient(c *gin.Context, userDTO *user.UserDTO) error {
	log.Printf("[ROOM] Adding user %+v to room %d", userDTO, room.ID)
	err := chat.NewClient(c, userDTO, &room.Manager.BaseWsManager)
	if err != nil {
		return err
	}

	return nil
}
