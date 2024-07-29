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

func (c *BaseChat) GetMembers() []int64 {
	return c.Members
}

func (c *BaseChat) GetManager() *ws.WsManager {
	return c.Manager
}

func (c *BaseChat) SendMessage(msg msg.Message) {
	c.Manager.Broadcast(msg)
}
