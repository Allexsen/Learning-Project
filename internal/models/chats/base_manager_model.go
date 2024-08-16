package chat

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type BaseWsManager struct {
	Clients    map[*Client]bool
	Broadcast  chan msg.BaseMessage
	Register   chan *Client
	Unregister chan *Client
	Stop       chan struct{}
	sync.RWMutex
}

type WsManager interface {
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message msg.BaseMessage)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewBaseWsManager creates a new BaseWsManager
func NewBaseWsManager() *BaseWsManager {
	return &BaseWsManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan msg.BaseMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Stop:       make(chan struct{}),
	}
}

// Run starts the BaseWsManager
func (manager *BaseWsManager) Run() {
	log.Printf("[CHAT-manager] Starting manager")
	for {
		select {
		case client := <-manager.Register:
			manager.RegisterClient(client)
		case client := <-manager.Unregister:
			manager.UnregisterClient(client)
		case message := <-manager.Broadcast:
			manager.BroadcastMessage(message)
		case <-manager.Stop:
			return
		}
	}
}

// RegisterClient registers a client
func (manager *BaseWsManager) RegisterClient(client *Client) {
	log.Printf("[CHAT-manager] Registering client: %+v", client.UserDTO)
	manager.Lock()
	manager.Clients[client] = true
	manager.Unlock()
	manager.BroadcastMessage(msg.BaseMessage{
		ID:        int64(uuid.New().ID()),
		Type:      "chatMessage",
		SenderID:  0, // System message
		Timestamp: time.Now().Unix(),
		Content:   fmt.Sprintf("%s has joined the chat", client.UserDTO.Username),
		Status:    "received",
	})
	log.Printf("[CHAT-manager] Client registered: %+v", client.UserDTO)
}

// UnregisterClient unregisters a client
func (manager *BaseWsManager) UnregisterClient(client *Client) {
	manager.Lock()
	if _, ok := manager.Clients[client]; ok {
		log.Printf("[CHAT-manager] Unregistering client: %+v", client.UserDTO)
		client.Conn.Close()
		delete(manager.Clients, client)
		close(client.Send)

		manager.Unlock()
		manager.BroadcastMessage(msg.BaseMessage{
			ID:        int64(uuid.New().ID()),
			Type:      "chatMessage",
			SenderID:  0, // System message
			Timestamp: time.Now().Unix(),
			Content:   fmt.Sprintf("%s has left the chat", client.UserDTO.Username),
			Status:    "received",
		})
		manager.Lock()
	}

	client = nil
	manager.Unlock()
}

// BroadcastMessage sends a message to all clients associated with the manager
func (manager *BaseWsManager) BroadcastMessage(message msg.BaseMessage) {
	manager.Lock()
	for client := range manager.Clients {
		select {
		case client.Send <- message:
		default:
			manager.Unlock()
			manager.UnregisterClient(client)
			manager.Lock()
		}
	}
	manager.Unlock()
}

// Close closes the BaseWsManager
func (manager *BaseWsManager) Close() {
	log.Printf("[CHAT-manager] Closing manager %+v", manager)
	manager.Lock()

	manager.Stop <- struct{}{}
	close(manager.Broadcast)
	close(manager.Register)
	close(manager.Unregister)
	close(manager.Stop)

	manager.Unlock()
	for client := range manager.Clients {
		manager.UnregisterClient(client)
	}

	manager.Clients = nil
	manager = nil
}
