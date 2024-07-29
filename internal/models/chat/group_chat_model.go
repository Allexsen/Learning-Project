package chat

import (
	"github.com/Allexsen/Learning-Project/internal/models/msg"
)

// GroupChat  represents a group chat
type GroupChat struct {
	BaseChat
	Name     string             `json:"name"`
	Messages []msg.GroupMessage `json:"messages"`
}
