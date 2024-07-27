package chats

import (
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/messages"
)

// Chat represents a 1-1 chat
type Chat struct {
	ID        int64                    `json:"id"`
	User1ID   int64                    `json:"user1ID"`
	User2ID   int64                    `json:"user2ID"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	Messages  []messages.DirectMessage `json:"messages"`
}
