package ws

import "github.com/gorilla/websocket"

// writeLoop spins off an infinite for loop iterating over a send channel.
// If there is a new message in the channel, sends it to the client.
// If the send channel gets closed, writeLoop closes the client connection.
func (client *Client) writeLoop() {
	defer client.conn.Close()

	for msg := range client.send {
		client.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
