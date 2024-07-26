package ws

func (manager *WsManager) send(message string) {
	for client := range manager.clients {
		client.send <- []byte(message)
	}
}
