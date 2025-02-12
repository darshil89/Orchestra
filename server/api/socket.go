package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader    = websocket.Upgrader{} // Upgrader to handle WebSocket requests
	clients     = make(map[*websocket.Conn]bool)
	clientsLock sync.Mutex
)

// WebSocket handler
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // Allow all origins
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Failed:", err)
		return
	}

	// Add client to the map
	clientsLock.Lock()
	clients[conn] = true
	clientsLock.Unlock()
	log.Println("✅ New WebSocket client connected!")

	// Listen for client disconnect
	// defer func() {
	// 	clientsLock.Lock()
	// 	delete(clients, conn)
	// 	clientsLock.Unlock()
	// 	conn.Close()
	// 	log.Println("WebSocket client disconnected")
	// }()
}

// Function to broadcast Redis messages to all WebSocket clients
func broadcastToClients(message []byte) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error sending message to WebSocket client:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
