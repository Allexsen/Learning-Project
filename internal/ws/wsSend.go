package ws

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
