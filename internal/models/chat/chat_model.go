package chat

import (
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
	"github.com/google/uuid"
)

type BaseChat struct {
	ID        int64         `json:"id,omitempty"`
	CreatedAt int64         `json:"created_at,omitempty"`
	UpdatedAt int64         `json:"updated_at,omitempty"`
	Manager   *ws.WsManager `json:"manager,omitempty"`
	Members   []*ws.Client  `json:"members,omitempty"`
	Messages  []msg.Message `json:"messages,omitempty"`
}

type Chat interface {
	GetManager() *ws.WsManager
}

// NewBaseChat creates a new chat
func NewBaseChat(manager *ws.WsManager, members ...*ws.Client) *BaseChat {
	return &BaseChat{
		ID:        int64(uuid.New().ID()),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Manager:   manager,
		Members:   members,
	}
}

// GetManager implements the Chat interface
func (chat BaseChat) GetManager() *ws.WsManager {
	return chat.Manager
}

// SendMessage sends a message to the chat
func (chat *BaseChat) SendMessage(message msg.Message) {
	chat.Manager.Broadcast(message)
}

// AddMessage adds a message to the chat
func (chat *BaseChat) AddMessage(message msg.Message) {
	chat.Messages = append(chat.Messages, message.(msg.GroupMessage))
}
