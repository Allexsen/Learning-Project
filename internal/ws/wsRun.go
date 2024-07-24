package ws

// Run starts the WsManager to handle connections and messages
func (manager *WsManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.Lock() // Must be unlocked
			manager.clients[client] = true
			manager.Unlock()

		case client := <-manager.unregister:
			manager.Lock() // Must be unlocked
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
			}
			manager.Unlock()

		case message := <-manager.broadcast:
			manager.Lock() // Must be unlocked
			for client := range manager.clients {
				select {
				case client.send <- message:
				default:
					delete(manager.clients, client)
					close(client.send)
				}
			}
			manager.Unlock()
		}
	}
}
