package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/gin-gonic/gin"
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
	log.Printf("[WS] Upgrading connection: %v", c.Request.RemoteAddr) // Temporary log
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	log.Printf("[WS] Connection established: %v", c.Request.RemoteAddr) // Temporary log
	userDTO, exists := c.Get("userDTO")
	if !exists {
		apperrors.HandleError(c, apperrors.New(
			http.StatusInternalServerError,
			"UserDTO not set",
			apperrors.ErrMissingRequiredField,
			map[string]interface{}{"details": "userDTO not set in the context"},
		))
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
	log.Println("Starting WsManager")
	for {
		select {
		case client := <-manager.register:
			manager.Lock()                                      // Must be unlocked
			log.Printf("Client registered: %v", client.userDTO) // Temporary log
			manager.clients[client] = true
			client.manager = manager
			manager.Broadcast(msg.BaseMessage{ // Placeholder message
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
				manager.Broadcast(msg.BaseMessage{ // Placeholder message
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

// AddClient adds a client to the manager
func (manager *WsManager) AddClient(userID int64) {
	manager.register <- &Client{
		userDTO: &user.UserDTO{ // TODO: Swap userID with UserDTO once the logic is implemented
			ID:       userID,
			Email:    "Place@holder.com", // Placeholder
			Username: "Placeholder",      // Placeholder
		},
		send: make(chan msg.Message),
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
