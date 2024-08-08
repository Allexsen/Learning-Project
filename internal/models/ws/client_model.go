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
	conn    *websocket.Conn  // Websocket connection
	send    chan msg.Message // Channel on which messages are sent to the client
	manager *WsManager       // ClientManager associated with the client
	userDTO *user.UserDTO    // User associated with the client
}

// readLoop spins off an infinite for loop reading incoming messages from a client.
// In case of error, breaks the loop, closes connection and unregisters the client.
func (client *Client) readLoop(manager *WsManager) {
	defer func() {
		manager.unregister <- client
		client.conn.Close()
		log.Printf("[WS] Client %+v has been unregistered", client)
	}()

	log.Printf("[WS] Spinning off read loop for client %+v", client)
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
		log.Printf("[WS] Client %+v has been disconnected", client)
	}()

	log.Printf("[WS] Spinning off write loop for client %+v", client)
	for msg := range client.send {
		message, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[WS] Error marshalling message: %v", err)
			continue
		}

		log.Printf("[WS] Sending message to client %+v: %s", client, message)
		err = client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("[WS] Error writing message: %v", err)
			log.Printf("[WS] Breaking write loop for client %+v", client)
			break
		}
	}
}
