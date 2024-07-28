package ws

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WsManager manages Clients and WebSocket connections
type WsManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	sync.RWMutex
}

// NewWsManager creates a new ClientManager
func NewWsManager() *WsManager {
	return &WsManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
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
		// TODO: Perhaps, implement a proper origin checking, and other security measurements
		return true
	},
}

// ServeWs handles WebSocket requests from the peer
func ServeWs(manager *WsManager, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	username := c.Query("username")
	if username == "" {
		username = "anonymous"
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		user: &user.User{Username: username}, // Placeholder
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
			manager.Lock()                                            // Must be unlocked
			log.Printf("Client registered: %s", client.user.Username) // Temporary log
			manager.clients[client] = true
			client.manager = manager
			manager.send(fmt.Sprintf("User %s has joined the chat", client.user.Username)) // Placeholder message
			manager.Unlock()
		case client := <-manager.unregister:
			manager.Lock() // Must be unlocked
			if _, ok := manager.clients[client]; ok {
				log.Printf("Client unregistered: %s", client.user.Username) // Temporary log
				delete(manager.clients, client)
				close(client.send)
				manager.send(fmt.Sprintf("User %s has left the chat", client.user.Username)) // Placeholder message
			}
			manager.Unlock()
		case message := <-manager.broadcast:
			manager.Lock() // Must be unlocked
			manager.send(string(message))
			manager.Unlock()
		}
	}
}

// send sends a message to all clients associated with the manager
func (manager *WsManager) send(message string) {
	for client := range manager.clients {
		select {
		case client.send <- []byte(message):
		default:
			delete(manager.clients, client)
			close(client.send)
		}
	}
}

// writeLoop spins off an infinite for loop iterating over a send channel.
// If there is a new message in the channel, sends it to the client.
// If the send channel gets closed, writeLoop closes the client connection.
func (client *Client) writeLoop() {
	defer func() {
		client.conn.Close()
		client.manager = nil
	}()

	for msg := range client.send {
		client.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// readLoop spins off an infinite for loop reading incoming messages from a client.
// In case of error, breaks the loop, closes connection and unregisters the client.
func (client *Client) readLoop(manager *WsManager) {
	defer func() {
		manager.unregister <- client
		client.conn.Close()
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			break
		}

		manager.broadcast <- msg
	}
}
