package chat

import (
	"github.com/Allexsen/Learning-Project/internal/models/msg"
)

// DirectChat represents a 1-1 chat
type DirectChat struct {
	BaseChat
	User1ID  int64               `json:"user1_id"`
	User2ID  int64               `json:"user2_id"`
	Messages []msg.DirectMessage `json:"messages"`
}
