package ws

import (
	"os/user"

	"github.com/gorilla/websocket"
)

// Client represents a websocket client
type Client struct {
	conn    *websocket.Conn `json:"-"` // Websocket connection
	send    chan []byte     `json:"-"` // Channel on which messages are sent to the client
	manager *WsManager      `json:"-"` // ClientManager associated with the client
	user    *user.User      `json:"-"` // User associated with the client
}
