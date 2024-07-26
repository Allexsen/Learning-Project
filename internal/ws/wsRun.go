package ws

import (
	"fmt"
	"log"
)

// Run starts the WsManager to handle connections and messages
func (manager *WsManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.Lock()                                       // Must be unlocked
			log.Printf("Client registered: %s", client.username) // Temporary log
			manager.clients[client] = true
			client.manager = manager
			manager.send(fmt.Sprintf("User %s has joined the chat", client.username)) // Placeholder message
			manager.Unlock()
		case client := <-manager.unregister:
			manager.Lock() // Must be unlocked
			if _, ok := manager.clients[client]; ok {
				log.Printf("Client unregistered: %s", client.username) // Temporary log
				delete(manager.clients, client)
				close(client.send)
				manager.send(fmt.Sprintf("User %s has left the chat", client.username)) // Placeholder message
			}
			manager.Unlock()
		case message := <-manager.broadcast:
			manager.Lock() // Must be unlocked
			manager.send(string(message))
			manager.Unlock()
		}
	}
}
