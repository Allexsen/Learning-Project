package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WsManager manages Clients and WebSocket connections
type WsManager struct {
	clients    map[*Client]bool
	broadcast  chan msg.Message
	register   chan *Client
	unregister chan *Client
	sync.RWMutex
}

// NewManager creates a new ClientManager
func NewManager() *WsManager {
	return &WsManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan msg.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// upgrader sets up the Upgrader.Upgrade() method to
// be used for http to websocket connection upgrade
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// TODO: Perhaps, implement a proper origin checking, and other security measurements
	},
}

// WsHandler handles WebSocket requests from the peer
func WsHandler(manager *WsManager, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	userDTO, exists := c.Get("userDTO") // TODO: Implement proper JWT validation
	if !exists {
		http.NotFound(c.Writer, c.Request) // TODO: Cetralize error handling
		return
	}

	client := &Client{
		conn:    conn,
		send:    make(chan msg.Message, 256),
		userDTO: userDTO.(*user.UserDTO),
	}

	manager.register <- client

	go client.writeLoop()
	go client.readLoop(manager)
}

// Run starts the WsManager to handle connections and messages
func (manager *WsManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.Lock()                                      // Must be unlocked
			log.Printf("Client registered: %v", client.userDTO) // Temporary log
			manager.clients[client] = true
			client.manager = manager
			manager.Broadcast(msg.Message{ // Placeholder message
				ID:        0, // Placeholder
				SenderID:  0, // Placeholder
				Timestamp: time.Now().Unix(),
				Content:   fmt.Sprintf("%s has joined the chat", client.userDTO.Username),
				Status:    "received",
			})
			manager.Unlock()
		case client := <-manager.unregister:
			manager.Lock() // Must be unlocked
			if _, ok := manager.clients[client]; ok {
				log.Printf("Client unregistered: %v", client.userDTO) // Temporary log
				delete(manager.clients, client)
				close(client.send)
				manager.Broadcast(msg.Message{ // Placeholder message
					ID:        0, // Placeholder
					SenderID:  0, // Placeholder
					Timestamp: time.Now().Unix(),
					Content:   fmt.Sprintf("%s has left the chat", client.userDTO.Username),
					Status:    "received",
				})
			}
			manager.Unlock()
		case msg := <-manager.broadcast:
			manager.Lock() // Must be unlocked
			manager.Broadcast(msg)
			manager.Unlock()
		}
	}
}

// Broadcast sends a message to all clients associated with the manager
func (manager *WsManager) Broadcast(msg msg.Message) {
	for client := range manager.clients {
		select {
		case client.send <- msg:
		default:
			delete(manager.clients, client)
			close(client.send)
		}
	}
}
