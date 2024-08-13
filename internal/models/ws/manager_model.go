package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/Allexsen/Learning-Project/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WsManager manages Clients and WebSocket connections
type WsManager struct {
	clients    map[*Client]bool
	broadcast  chan msg.BaseMessage
	register   chan *Client
	unregister chan *Client
	stop       chan struct{}
	sync.RWMutex
}

// upgrader sets up the Upgrader.Upgrade() method to
// be used for http to websocket connection upgrade
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default
	},
}

// NewManager creates a new ClientManager
func NewManager() *WsManager {
	return &WsManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan msg.BaseMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan struct{}),
	}
}

// WsHandler handles WebSocket requests from the peer
func WsHandler(manager *WsManager, c *gin.Context) {
	log.Printf("[WS] Upgrading connection: %v", c.Request.RemoteAddr)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	log.Printf("[WS] Connection established: %v", c.Request.RemoteAddr)
	var userDTO *user.UserDTO
	if !utils.ShouldBindJSON(c, userDTO) {
		return
	}

	client := &Client{
		conn:    conn,
		send:    make(chan msg.Message, 256),
		userDTO: userDTO,
	}

	manager.register <- client

	go client.writeLoop()
	go client.readLoop(manager)
}

// Run starts the WsManager to handle connections and messages
func (manager *WsManager) Run() {
	log.Println("[WS] Starting a WebSocket Manager")
	for {
		select {
		case client := <-manager.register:
			manager.Lock() // Must be unlocked
			log.Printf("[WS] Client registered: %v", client.userDTO)
			manager.clients[client] = true
			client.manager = manager
			manager.Broadcast(msg.BaseMessage{
				ID:        int64(uuid.New().ID()),
				Type:      "chatMessage",
				SenderID:  0, // System message
				Timestamp: time.Now().Unix(),
				Content:   fmt.Sprintf("%s has joined the chat", client.userDTO.Username),
				Status:    "received",
			})
			manager.Unlock()
		case client := <-manager.unregister:
			manager.Lock() // Must be unlocked
			if _, ok := manager.clients[client]; ok {
				log.Printf("[WS] Client unregistered: %v", client.userDTO)
				delete(manager.clients, client)
				close(client.send)
				manager.Broadcast(msg.BaseMessage{
					ID:        int64(uuid.New().ID()),
					Type:      "chatMessage",
					SenderID:  0, // System message
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
		case <-manager.stop:
			return
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

// Close closes the manager and all clients
func (manager *WsManager) Close() {
	manager.Lock()
	defer manager.Unlock()
	close(manager.stop)
	for client := range manager.clients {
		close(client.send)
		delete(manager.clients, client)
	}
}
