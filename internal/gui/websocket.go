package gui

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"

	"gorconx/internal/core"

	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections for real-time updates
type WebSocketManager struct {
	connections map[string]*websocket.Conn
	mutex       sync.RWMutex
	app         *core.Application
	upgrader    websocket.Upgrader
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(app *core.Application) *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]*websocket.Conn),
		app:         app,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, implement proper origin checking
				return true
			},
		},
	}
}

// HandleWebSocket handles WebSocket connections
func (wsm *WebSocketManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wsm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Generate client ID
	clientID := generateClientID()

	// Register connection
	wsm.mutex.Lock()
	wsm.connections[clientID] = conn
	wsm.mutex.Unlock()

	// Remove connection when done
	defer func() {
		wsm.mutex.Lock()
		delete(wsm.connections, clientID)
		wsm.mutex.Unlock()
	}()

	// Send welcome message
	welcome := map[string]interface{}{
		"type": "welcome",
		"data": map[string]string{
			"client_id": clientID,
			"message":   "Connected to GoReconX WebSocket",
		},
	}
	wsm.sendToClient(conn, welcome)
	// Listen for messages
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			// Check if it's a normal close
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket unexpected close error: %v", err)
			} else {
				log.Printf("WebSocket client disconnected: %s", clientID)
			}
			break
		}

		// Handle different message types
		wsm.handleMessage(clientID, msg)
	}
}

// handleMessage processes incoming WebSocket messages
func (wsm *WebSocketManager) handleMessage(clientID string, msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		wsm.sendError(clientID, "Invalid message format")
		return
	}

	switch msgType {
	case "subscribe_session":
		sessionID, ok := msg["session_id"].(string)
		if !ok {
			wsm.sendError(clientID, "Session ID required")
			return
		}
		wsm.subscribeToSession(clientID, sessionID)

	case "unsubscribe_session":
		sessionID, ok := msg["session_id"].(string)
		if !ok {
			wsm.sendError(clientID, "Session ID required")
			return
		}
		wsm.unsubscribeFromSession(clientID, sessionID)

	case "get_session_status":
		sessionID, ok := msg["session_id"].(string)
		if !ok {
			wsm.sendError(clientID, "Session ID required")
			return
		}
		wsm.sendSessionStatus(clientID, sessionID)

	default:
		wsm.sendError(clientID, "Unknown message type: "+msgType)
	}
}

// subscribeToSession subscribes a client to session updates
func (wsm *WebSocketManager) subscribeToSession(clientID, sessionID string) {
	// Send confirmation
	response := map[string]interface{}{
		"type": "subscription_confirmed",
		"data": map[string]string{
			"session_id": sessionID,
		},
	}
	wsm.sendToClientByID(clientID, response)
}

// unsubscribeFromSession unsubscribes a client from session updates
func (wsm *WebSocketManager) unsubscribeFromSession(clientID, sessionID string) {
	response := map[string]interface{}{
		"type": "unsubscription_confirmed",
		"data": map[string]string{
			"session_id": sessionID,
		},
	}
	wsm.sendToClientByID(clientID, response)
}

// sendSessionStatus sends current session status to client
func (wsm *WebSocketManager) sendSessionStatus(clientID, sessionID string) {
	session, exists := wsm.app.GetSession(sessionID)
	if !exists {
		wsm.sendError(clientID, "Session not found")
		return
	}

	response := map[string]interface{}{
		"type": "session_status",
		"data": session,
	}
	wsm.sendToClientByID(clientID, response)
}

// BroadcastSessionUpdate broadcasts session updates to all connected clients
func (wsm *WebSocketManager) BroadcastSessionUpdate(sessionID string, update interface{}) {
	message := map[string]interface{}{
		"type": "session_update",
		"data": map[string]interface{}{
			"session_id": sessionID,
			"update":     update,
		},
	}

	wsm.broadcastMessage(message)
}

// BroadcastModuleProgress broadcasts module execution progress
func (wsm *WebSocketManager) BroadcastModuleProgress(sessionID, moduleName string, progress float64, status string) {
	message := map[string]interface{}{
		"type": "module_progress",
		"data": map[string]interface{}{
			"session_id":  sessionID,
			"module_name": moduleName,
			"progress":    progress,
			"status":      status,
		},
	}

	wsm.broadcastMessage(message)
}

// BroadcastModuleResult broadcasts module execution results
func (wsm *WebSocketManager) BroadcastModuleResult(sessionID, moduleName string, result interface{}) {
	message := map[string]interface{}{
		"type": "module_result",
		"data": map[string]interface{}{
			"session_id":  sessionID,
			"module_name": moduleName,
			"result":      result,
		},
	}

	wsm.broadcastMessage(message)
}

// sendError sends an error message to a specific client
func (wsm *WebSocketManager) sendError(clientID, errorMsg string) {
	response := map[string]interface{}{
		"type": "error",
		"data": map[string]string{
			"message": errorMsg,
		},
	}
	wsm.sendToClientByID(clientID, response)
}

// sendToClient sends a message to a specific WebSocket connection
func (wsm *WebSocketManager) sendToClient(conn *websocket.Conn, message interface{}) {
	if err := conn.WriteJSON(message); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// sendToClientByID sends a message to a client by ID
func (wsm *WebSocketManager) sendToClientByID(clientID string, message interface{}) {
	wsm.mutex.RLock()
	conn, exists := wsm.connections[clientID]
	wsm.mutex.RUnlock()

	if exists {
		wsm.sendToClient(conn, message)
	}
}

// broadcastMessage sends a message to all connected clients
func (wsm *WebSocketManager) broadcastMessage(message interface{}) {
	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()

	for _, conn := range wsm.connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("WebSocket broadcast error: %v", err)
		}
	}
}

// GetConnectionCount returns the number of active WebSocket connections
func (wsm *WebSocketManager) GetConnectionCount() int {
	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()
	return len(wsm.connections)
}

// generateClientID generates a unique client ID
func generateClientID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "client_" + hex.EncodeToString(bytes)
}
