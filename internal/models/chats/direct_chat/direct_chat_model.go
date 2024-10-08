package direct_chat

import chat "github.com/Allexsen/Learning-Project/internal/models/chats"

// DirectChat represents a 1-1 chat
type DirectChat struct {
	chat.BaseChat
	User1ID int64 `json:"user1_id,omitempty"`
	User2ID int64 `json:"user2_id,omitempty"`
}

// NewDirectChat creates a new DirectChat
func NewDirectChat(user1ID, user2ID int64) *DirectChat {
	return &DirectChat{
		User1ID: user1ID,
		User2ID: user2ID,
	}
}
