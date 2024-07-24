package ws

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte // Message to send to client
}

// WsManager manages WebSocket connections and messages
type WsManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	sync.Mutex
}

// NewWsManager creates a new wsManager
func NewWsManager() *WsManager {
	return &WsManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// upgrader sets up the Upgrader.Upgrade() method
// for http to websocket connection upgrade
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

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	manager.register <- client

	go client.writeLoop()
	go client.readLoop(manager)
}
