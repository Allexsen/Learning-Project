package room

import (
	"fmt"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/chat"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/gin-gonic/gin"
)

// Room  represents a group chat
type Room struct {
	chat.BaseChat
	Name string `json:"name,omitempty"`
}

var roomsList map[int64]*Room

func init() {
	roomsList = make(map[int64]*Room)
}

// NewRoom creates a new room
func NewRoom(name string) *Room {
	manager := ws.NewManager()
	go manager.Run()
	room := &Room{
		BaseChat: *chat.NewBaseChat(manager, nil),
		Name:     name,
	}

	roomsList[room.ID] = room
	return room
}

// GetRooms returns all rooms
func GetRooms() ([]*Room, error) {
	rooms := make([]*Room, 0, len(roomsList))
	for _, room := range roomsList {
		newRoom := Room{
			BaseChat: chat.BaseChat{
				ID: room.ID,
			},
			Name: room.Name,
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

// AddClient adds a user to the room
func (room *Room) AddClient(c *gin.Context, userDTO user.UserDTO) error {
	_, exists := roomsList[room.ID]
	if !exists {
		return apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("Room with ID %d not found", room.ID),
			apperrors.ErrNotFound,
			nil,
		)
	}

	cl, err := ws.NewClient(c, &userDTO, room.Manager)
	if err != nil {
		return err
	}

	room.Members[cl] = true
	return nil
}

// RemoveClient removes a user from the room
func (room *Room) RemoveClient(client *ws.Client) error {
	_, exists := roomsList[room.ID]
	if !exists {
		return apperrors.New(
			http.StatusNotFound,
			fmt.Sprintf("Room with ID %d not found", room.ID),
			apperrors.ErrNotFound,
			nil,
		)
	}

	delete(room.Members, client)
	return nil
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

	room.Manager.Close()
	delete(roomsList, room.ID)
	room = nil
	return nil
}
