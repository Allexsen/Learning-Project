package chat

import (
	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/ws"
)

type BaseChat struct {
	ID        int64         `json:"id"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
	Manager   *ws.WsManager `json:"manager"`
	Members   []int64       `json:"members"`
}

type Chat interface {
	SendMessage(msg msg.Message)
	GetMessages() []msg.Message
	GetMembers() []int64
	GetManager() *ws.WsManager
}

func (chat *BaseChat) GetMembers() []int64 {
	return chat.Members
}

func (chat *BaseChat) GetManager() *ws.WsManager {
	return chat.Manager
}

func (chat *BaseChat) SendMessage(msg msg.Message) {
	chat.Manager.Broadcast(msg)
}
