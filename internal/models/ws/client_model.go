package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/gorilla/websocket"
)

// Client represents a websocket client
type Client struct {
	conn    *websocket.Conn  `json:"-"` // Websocket connection
	send    chan msg.Message `json:"-"` // Channel on which messages are sent to the client
	manager *WsManager       `json:"-"` // ClientManager associated with the client
	userDTO *user.UserDTO    `json:"-"` // User associated with the client
}

// readLoop spins off an infinite for loop reading incoming messages from a client.
// In case of error, breaks the loop, closes connection and unregisters the client.
func (client *Client) readLoop(manager *WsManager) {
	defer func() {
		manager.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg msg.BaseMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		msg.SenderID = client.userDTO.ID
		msg.Timestamp = time.Now().Unix()
		msg.Status = "received"

		manager.broadcast <- msg
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
		message, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %v", err)
			continue
		}
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}
