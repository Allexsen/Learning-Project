package chat

import (
	"fmt"
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/gin-gonic/gin"
)

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
	manager := ws.NewManager()
	go manager.Run()
	room := &Room{
		BaseChat: *NewBaseChat(manager),
		Name:     name,
	}

	roomsManager.Rooms[room.ID] = room
	return room
}

// GetRooms returns all rooms
func GetRooms() ([]*Room, error) {
	rooms := make([]*Room, 0, len(roomsManager.Rooms))
	for _, room := range roomsManager.Rooms {
		newRoom := Room{
			BaseChat: BaseChat{
				ID: room.ID,
			},
			Name: room.Name,
		}

		rooms = append(rooms, &newRoom)
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

// AddClient adds a user to the room
func (room *Room) AddClient(c *gin.Context, userDTO user.UserDTO) error {
	_, exists := roomsManager.Rooms[room.ID]
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

	room.Members = append(room.Members, cl)
	return nil
}

// DeleteRoom deletes a room
func (room *Room) DeleteRoom() error {
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
	delete(roomsManager.Rooms, room.ID)
	room = nil
	return nil
}
