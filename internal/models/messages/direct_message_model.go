package messages

import (
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/ws"
)

// DirectMessage represents a 1-1 chat message
type DirectMessage struct {
	ID        int64      `json:"id"`
	ChatID    int64      `json:"chatID"`
	Content   string     `json:"content"`
	Sender    *ws.Client `json:"sender"`
	Timestamp time.Time  `json:"timestamp"`
	Status    string     `json:"status"`
}
