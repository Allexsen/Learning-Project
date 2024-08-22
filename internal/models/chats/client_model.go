// Package chat provides base chat & ws model, managing clients and broadcasting messages.
// It also provides database interaction methods for chat messages.
package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Allexsen/Learning-Project/internal/models/msg"
	"github.com/Allexsen/Learning-Project/internal/models/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Client represents a websocket client
type Client struct {
	Conn    *websocket.Conn      // Websocket connection
	Send    chan msg.BaseMessage // Channel to send messages to the client
	UserDTO *user.UserDTO        // User associated with the client
}

// NewClient creates a new Client
func NewClient(c *gin.Context, userDTO *user.UserDTO, manager *BaseWsManager) error {
	log.Printf("[CHAT-client] Upgrading connection: %s", c.Request.RemoteAddr)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return err
	}
	log.Printf("[CHAT-client] Connection established: %s", c.Request.RemoteAddr)

	client := &Client{
		Conn:    conn,
		Send:    make(chan msg.BaseMessage),
		UserDTO: userDTO,
	}

	// THIS IS IMPORTANT. DO NOT REMOVE
	// This ensures write loop spins off before registering a client.
	// Unlocks in the writeLoop function.
	manager.Lock()
	go client.writeLoop(manager)

	// THIS IS IMPORTANT. DO NOT REMOVE
	// This ensures read loop spins off before registering a client.
	// Unlocks in the readLoop function.
	manager.Lock()
	go client.readLoop(manager)

	manager.Register <- client

	return nil
}

// readLoop spins off an infinite for loop reading incoming messages of a client.
// In case of error, breaks the loop, closes connection and unregisters the client.
func (client *Client) readLoop(manager *BaseWsManager) {
	log.Printf("[CHAT-client] Spinning off read loop for client %+v", client)

	defer func() {
		manager.UnregisterClient(client)
	}()

	// THIS IS IMPORTANT. DO NOT REMOVE
	// This unlock is paired with the lock before spinning off the loop from NewClient
	manager.Unlock()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		log.Printf("[CHAT-client] Received message from client %+v: %s", client, message)
		var msg msg.BaseMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("[CHAT-client] Error unmarshalling message: %+v", err)
			continue
		}

		msg.SenderID = client.UserDTO.ID
		msg.Timestamp = time.Now().Unix()
		msg.Status = "received"

		manager.BroadcastMessage(msg)
	}
}

// writeLoop spins off an infinite for loop writing outgoing messages of a client.
// If there is a new message in the channel, sends it to the client.
// If the send channel gets closed, writeLoop closes the client connection.
func (client *Client) writeLoop(manager *BaseWsManager) {
	log.Printf("[CHAT-client] Spinning off write loop for client %+v", client)

	defer func() {
		manager.UnregisterClient(client)
	}()

	// THIS IS IMPORTANT. DO NOT REMOVE
	// This unlock is paired with the lock before spinning off the loop from NewClient
	manager.Unlock()

	for msg := range client.Send {
		message, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[CHAT-client] Error marshalling message: %+v", err)
			continue
		}

		err = client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("[CHAT-client] Error writing message: %+v", err)
			log.Printf("[CHAT-client] Breaking write loop for client %+v", client)
			break
		}
	}
}
