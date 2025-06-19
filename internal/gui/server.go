package gui

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"gorconx/internal/core"
)

// GUIServer represents the web GUI server
type GUIServer struct {
	app       *core.Application
	server    *http.Server
	wsManager *WebSocketManager
}

// NewGUIServer creates a new GUI server
func NewGUIServer(app *core.Application) *GUIServer {
	return &GUIServer{
		app:       app,
		wsManager: NewWebSocketManager(app),
	}
}

// Start starts the GUI server
func (g *GUIServer) Start() error {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// WebSocket endpoint
	mux.HandleFunc("/ws", g.wsManager.HandleWebSocket)	// Main routes
	mux.HandleFunc("/", g.handleIndex)
	mux.HandleFunc("/dashboard", g.handleDashboard)
	mux.HandleFunc("/modules", g.handleModules)
	mux.HandleFunc("/sessions/", g.handleSessionDetail) // Handle individual session pages FIRST
	mux.HandleFunc("/sessions", g.handleSessions)       // Handle sessions list page AFTER
	mux.HandleFunc("/settings", g.handleSettings)
	mux.HandleFunc("/reports", g.handleReports)

	// API proxy
	mux.HandleFunc("/api/", g.handleAPIProxy)

	g.server = &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	g.app.GetLogger().Infof("GUI server starting on port 8080")
	return g.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (g *GUIServer) Shutdown(ctx context.Context) error {
	if g.server != nil {
		return g.server.Shutdown(ctx)
	}
	return nil
}

func (g *GUIServer) handleAPIProxy(w http.ResponseWriter, r *http.Request) {
	// Forward to API server on port 8081
	apiURL := "http://localhost:8081" + r.URL.Path
	if r.URL.RawQuery != "" {
		apiURL += "?" + r.URL.RawQuery
	}

	g.app.GetLogger().Infof("Proxying API request: %s %s -> %s", r.Method, r.URL.Path, apiURL)

	// Create a new request to the API server
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(r.Method, apiURL, r.Body)
	if err != nil {
		g.app.GetLogger().Errorf("Failed to create API proxy request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
	// Make the request to API server
	resp, err := client.Do(req)
	if err != nil {
		g.app.GetLogger().Errorf("Failed to reach API server: %v", err)
		http.Error(w, "Failed to reach API server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	g.app.GetLogger().Infof("API proxy response: %d", resp.StatusCode)

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

// Handle index page
func (g *GUIServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := getIndexHTML()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Handle dashboard
func (g *GUIServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	html := getDashboardHTML()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Handle other pages
func (g *GUIServer) handleModules(w http.ResponseWriter, r *http.Request) {
	html := getModulesHTML()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (g *GUIServer) handleSessions(w http.ResponseWriter, r *http.Request) {
	html := getSessionsHTML()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (g *GUIServer) handleSettings(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Settings Page - Coming Soon</h1><a href='/dashboard'>Back to Dashboard</a>"))
}

func (g *GUIServer) handleReports(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Reports Page - Coming Soon</h1><a href='/dashboard'>Back to Dashboard</a>"))
}

func (g *GUIServer) handleSessionDetail(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL path
	sessionID := extractSessionIDFromPath(r.URL.Path)
	if sessionID == "" {
		g.app.GetLogger().Warnf("No session ID found in path: %s", r.URL.Path)
		http.Redirect(w, r, "/sessions", http.StatusFound)
		return
	}

	g.app.GetLogger().Infof("Attempting to load session: %s", sessionID)

	// Check if session exists
	session, exists := g.app.GetSession(sessionID)
	if !exists {
		g.app.GetLogger().Warnf("Session not found: %s", sessionID)
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	g.app.GetLogger().Infof("Session found: %s - %s", sessionID, session.Name)

	// Serve session detail page
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(getSessionDetailHTML(session)))
}

// extractSessionIDFromPath extracts session ID from URL path like /sessions/session_xxx
func extractSessionIDFromPath(path string) string {
	// Remove trailing slash if present
	path = strings.TrimSuffix(path, "/")
	
	parts := strings.Split(path, "/")
	if len(parts) >= 3 && parts[1] == "sessions" && parts[2] != "" {
		return parts[2]
	}
	return ""
}

// GetWebSocketManager returns the WebSocket manager
func (g *GUIServer) GetWebSocketManager() *WebSocketManager {
	return g.wsManager
}
