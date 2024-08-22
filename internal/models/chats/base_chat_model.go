// Package chat provides base chat & ws model, managing clients and broadcasting messages.
// It also provides database interaction methods for chat messages.
package chat

import (
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/google/uuid"
)

type BaseChat struct {
	ID        int64         `json:"id,omitempty"`
	CreatedAt int64         `json:"created_at,omitempty"`
	UpdatedAt int64         `json:"updated_at,omitempty"`
	Messages  []msg.Message `json:"messages,omitempty"`
}

// NewBaseChat creates a new chat
func NewBaseChat() *BaseChat {
	return &BaseChat{
		ID:        int64(uuid.New().ID()),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Messages:  make([]msg.Message, 0),
	}
}
