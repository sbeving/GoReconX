package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gorconx/internal/core"
)

// Server represents the API server
type Server struct {
	app    *core.Application
	server *http.Server
}

// NewServer creates a new API server
func NewServer(app *core.Application) *Server {
	return &Server{
		app: app,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register routes
	s.registerRoutes(mux)

	// Create server
	s.server = &http.Server{
		Addr:         ":8081", // API on different port
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	s.app.GetLogger().Infof("API server starting on port 8081")
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// registerRoutes registers all API routes
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Sessions
	mux.HandleFunc("/api/sessions", s.handleSessions)
	mux.HandleFunc("/api/sessions/", s.handleSession)

	// Modules
	mux.HandleFunc("/api/modules", s.handleModules)
	mux.HandleFunc("/api/modules/", s.handleModule)
	// Scans
	mux.HandleFunc("/api/scans", s.handleScansEnhanced)
	mux.HandleFunc("/api/scans/", s.handleScanEnhanced)

	// API Keys
	mux.HandleFunc("/api/apikeys", s.handleAPIKeys)
	mux.HandleFunc("/api/apikeys/", s.handleAPIKey)

	// Configuration
	mux.HandleFunc("/api/config", s.handleConfig)

	// Health check
	mux.HandleFunc("/api/health", s.handleHealth)
}

// Handle sessions
func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sessions := s.app.GetSessions()
		s.writeJSON(w, sessions)
	case "POST":
		var req struct {
			Name   string `json:"name"`
			Target string `json:"target"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		session := s.app.CreateSession(req.Name, req.Target)
		s.writeJSON(w, session)
	default:
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle individual session
func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	sessionID := extractIDFromPath(r.URL.Path, "/api/sessions/")

	switch r.Method {
	case "GET":
		session, exists := s.app.GetSession(sessionID)
		if !exists {
			s.writeError(w, "Session not found", http.StatusNotFound)
			return
		}
		s.writeJSON(w, session)
	case "DELETE":
		s.app.DeleteSession(sessionID)
		s.writeJSON(w, map[string]string{"status": "deleted"})
	default:
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle modules
func (s *Server) handleModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	modules := s.app.GetModules()
	moduleInfo := make(map[string]interface{})

	for name, module := range modules {
		moduleInfo[name] = map[string]string{
			"name": module.GetName(),
		}
	}

	s.writeJSON(w, moduleInfo)
}

// Handle individual module
func (s *Server) handleModule(w http.ResponseWriter, r *http.Request) {
	moduleName := extractIDFromPath(r.URL.Path, "/api/modules/")

	module, exists := s.app.GetModule(moduleName)
	if !exists {
		s.writeError(w, "Module not found", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		// Execute module
		var req struct {
			Target    string `json:"target"`
			SessionID string `json:"session_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := module.Execute(req.Target)
		if err != nil {
			s.writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, map[string]interface{}{
			"result": result,
			"status": "completed",
		})
	} else {
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle scans
// Handle API keys
func (s *Server) handleAPIKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiKeys, err := s.app.GetAPIKeyManager().ListAPIKeys()
		if err != nil {
			s.writeError(w, "Failed to list API keys", http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, apiKeys)
	case "POST":
		var req struct {
			Service string `json:"service"`
			APIKey  string `json:"api_key"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Service == "" || req.APIKey == "" {
			s.writeError(w, "Service and API key are required", http.StatusBadRequest)
			return
		}

		err := s.app.GetAPIKeyManager().StoreAPIKey(req.Service, req.APIKey)
		if err != nil {
			s.writeError(w, "Failed to store API key", http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, map[string]string{"status": "success"})
	default:
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle individual API key
func (s *Server) handleAPIKey(w http.ResponseWriter, r *http.Request) {
	service := extractIDFromPath(r.URL.Path, "/api/apikeys/")

	switch r.Method {
	case "GET":
		// Don't return the actual key for security
		s.writeJSON(w, map[string]string{"service": service, "status": "exists"})
	case "DELETE":
		err := s.app.GetAPIKeyManager().DeleteAPIKey(service)
		if err != nil {
			s.writeError(w, "Failed to delete API key", http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, map[string]string{"status": "deleted"})
	default:
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle scans (enhanced)
func (s *Server) handleScansEnhanced(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get query parameters
		sessionID := r.URL.Query().Get("session_id")
		if sessionID != "" {
			scans := s.app.GetScanManager().GetSessionScans(sessionID)
			s.writeJSON(w, scans)
		} else {
			// Return all scans (you might want to paginate this)
			s.writeJSON(w, map[string]string{"error": "session_id parameter required"})
		}
	case "POST":
		var req struct {
			SessionID  string                 `json:"session_id"`
			ModuleName string                 `json:"module_name"`
			Target     string                 `json:"target"`
			Options    map[string]interface{} `json:"options"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.SessionID == "" || req.ModuleName == "" || req.Target == "" {
			s.writeError(w, "Session ID, module name, and target are required", http.StatusBadRequest)
			return
		}

		scan, err := s.app.GetScanManager().StartScan(req.SessionID, req.ModuleName, req.Target, req.Options)
		if err != nil {
			s.writeError(w, "Failed to start scan: "+err.Error(), http.StatusInternalServerError)
			return
		}

		s.writeJSON(w, scan)
	default:
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle individual scan
func (s *Server) handleScanEnhanced(w http.ResponseWriter, r *http.Request) {
	scanID := extractIDFromPath(r.URL.Path, "/api/scans/")

	switch r.Method {
	case "GET":
		scan, exists := s.app.GetScanManager().GetScan(scanID)
		if !exists {
			s.writeError(w, "Scan not found", http.StatusNotFound)
			return
		}
		s.writeJSON(w, scan)
	case "DELETE":
		err := s.app.GetScanManager().CancelScan(scanID)
		if err != nil {
			s.writeError(w, "Failed to cancel scan: "+err.Error(), http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, map[string]string{"status": "cancelled"})
	default:
		s.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle configuration
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	config := s.app.GetConfig()
	s.writeJSON(w, config)
}

// Handle health check
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Utility functions
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.app.GetLogger().Printf("Error encoding JSON: %v", err)
	}
}

func (s *Server) writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code":  http.StatusText(code),
	})
}

func extractIDFromPath(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	return path[len(prefix):]
}
