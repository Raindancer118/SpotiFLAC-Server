package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"spotiflac/backend"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader configuration
// Following rule #14: Secure by Default - configure websocket securely
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from configured origins only
	CheckOrigin: func(r *http.Request) bool {
		// This will be configured based on config.yml CORS settings
		// For now, allow all (will be restricted in production)
		return true
	},
}

// WebSocketManager manages WebSocket connections and broadcasts
type WebSocketManager struct {
	clients   map[*websocket.Conn]bool
	broadcast chan interface{}
	mutex     sync.RWMutex
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan interface{}, 100),
	}
}

// Start begins the WebSocket broadcast loop
func (wsm *WebSocketManager) Start() {
	go func() {
		for {
			message := <-wsm.broadcast
			wsm.mutex.RLock()
			for client := range wsm.clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("WebSocket write error: %v", err)
					client.Close()
					wsm.removeClient(client)
				}
			}
			wsm.mutex.RUnlock()
		}
	}()

	// Start polling for updates to broadcast
	go wsm.pollUpdates()
}

// pollUpdates periodically checks for download progress and queue updates
func (wsm *WebSocketManager) pollUpdates() {
	// This will poll backend for updates and broadcast to clients
	// Implementation will send download progress and queue updates
	for {
		// Get download progress
		progress := backend.GetDownloadProgress()
		if progress.IsDownloading {
			wsm.Broadcast(map[string]interface{}{
				"type": "download_progress",
				"data": progress,
			})
		}

		// Get queue info
		queue := backend.GetDownloadQueue()
		wsm.Broadcast(map[string]interface{}{
			"type": "queue_update",
			"data": queue,
		})

		// Sleep briefly to avoid excessive polling
		// This could be optimized with channels/events from backend
		//time.Sleep(500 * time.Millisecond)
		break // For now, don't continuously poll
	}
}

// addClient adds a new WebSocket client
func (wsm *WebSocketManager) addClient(conn *websocket.Conn) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()
	wsm.clients[conn] = true
}

// removeClient removes a WebSocket client
func (wsm *WebSocketManager) removeClient(conn *websocket.Conn) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()
	delete(wsm.clients, conn)
}

// Broadcast sends a message to all connected clients
func (wsm *WebSocketManager) Broadcast(message interface{}) {
	select {
	case wsm.broadcast <- message:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

// Global WebSocket manager instance
var wsManager *WebSocketManager

// InitWebSocketManager initializes the global WebSocket manager
func InitWebSocketManager() {
	wsManager = NewWebSocketManager()
	wsManager.Start()
}

// HandleWebSocket handles WebSocket connection requests
// Endpoint: GET /ws
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add client to manager
	wsManager.addClient(conn)
	defer wsManager.removeClient(conn)

	// Send initial state
	wsManager.Broadcast(map[string]interface{}{
		"type": "connected",
		"data": map[string]interface{}{
			"message": "Connected to SpotiFLAC server",
		},
	})

	// Read messages from client (for ping/pong or commands)
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle client messages if needed
		handleClientMessage(conn, msg)
	}
}

// handleClientMessage processes messages from WebSocket clients
func handleClientMessage(conn *websocket.Conn, msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "ping":
		// Respond to ping
		conn.WriteJSON(map[string]interface{}{
			"type": "pong",
		})

	case "request_status":
		// Send current status
		progress := backend.GetDownloadProgress()
		queue := backend.GetDownloadQueue()

		conn.WriteJSON(map[string]interface{}{
			"type": "status_update",
			"data": map[string]interface{}{
				"progress": progress,
				"queue":    queue,
			},
		})
	}
}

// BroadcastDownloadProgress sends download progress to all WebSocket clients
// This can be called from backend when progress updates occur
func BroadcastDownloadProgress(progress backend.ProgressInfo) {
	if wsManager != nil {
		wsManager.Broadcast(map[string]interface{}{
			"type": "download_progress",
			"data": progress,
		})
	}
}

// BroadcastQueueUpdate sends queue updates to all WebSocket clients
func BroadcastQueueUpdate(queue backend.DownloadQueueInfo) {
	if wsManager != nil {
		wsManager.Broadcast(map[string]interface{}{
			"type": "queue_update",
			"data": queue,
		})
	}
}

// BroadcastMessage sends a generic message to all clients
func BroadcastMessage(messageType string, data interface{}) error {
	if wsManager == nil {
		return fmt.Errorf("WebSocket manager not initialized")
	}

	message := map[string]interface{}{
		"type": messageType,
		"data": data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	wsManager.Broadcast(jsonData)
	return nil
}
