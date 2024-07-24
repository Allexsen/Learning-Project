package ws

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
