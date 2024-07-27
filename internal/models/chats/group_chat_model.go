package chats

import (
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/messages"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
)

// Room  represents a group chat
type Room struct {
	ID        int64                   `json:"id"`
	Name      string                  `json:"name"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	WsManager *ws.WsManager           `json:"-"` // Websocket manager associated with the chat
	Messages  []messages.GroupMessage `json:"messages"`
}
