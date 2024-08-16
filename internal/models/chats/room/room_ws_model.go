package room

import (
	"log"

	chat "github.com/Allexsen/Learning-Project/internal/models/chats"
)

type wsManager struct {
	chat.BaseChat
	chat.BaseWsManager
	room *Room
}

// newWsManager creates a new RoomWSM
func newWsManager(room *Room) *wsManager {
	return &wsManager{
		BaseChat:      room.BaseChat,
		BaseWsManager: *chat.NewBaseWsManager(),
		room:          room,
	}
}

// Close closes the RoomWsManager
func (manager *wsManager) close() {
	log.Printf("[ROOM-manager] Closing manager for room %s", manager.room.Name)
	manager.room = nil
	manager.Close()
	manager = nil
}
