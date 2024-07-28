package messages

import (
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/ws"
)

// GroupMessage represents a group chat message
type GroupMessage struct {
	ID        int64      `json:"id"`
	RoomID    int64      `json:"roomID"`
	Content   string     `json:"content"`
	Sender    *ws.Client `json:"sender"`
	Timestamp time.Time  `json:"timestamp"`
	Status    string     `json:"status"`
}
