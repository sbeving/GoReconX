package core

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"sync"
	"time"

	"gorconx/internal/modules"

	"github.com/sirupsen/logrus"
)

// Module interface for compatibility
type Module interface {
	GetName() string
	Execute(target string) (interface{}, error)
}

// ModuleAdapter adapts the new module interface to the old one
type ModuleAdapter struct {
	module modules.Module
}

func (m *ModuleAdapter) GetName() string {
	return m.module.GetInfo().Name
}

func (m *ModuleAdapter) Execute(target string) (interface{}, error) {
	// Create a simple adapter that converts the old interface to new
	input := modules.ModuleInput{
		Target:    target,
		Options:   make(map[string]interface{}),
		SessionID: "default",
		Timeout:   30 * time.Second,
	}

	output := make(chan modules.ModuleResult, 100)
	defer close(output)

	ctx := context.Background()
	err := m.module.Execute(ctx, input, output)
	if err != nil {
		return nil, err
	}

	// Collect results
	var results []modules.ModuleResult
	for {
		select {
		case result, ok := <-output:
			if !ok {
				goto done
			}
			results = append(results, result)
		case <-time.After(100 * time.Millisecond):
			goto done
		}
	}
done:

	return results, nil
}

// Application represents the core application structure
type Application struct {
	db      *sql.DB
	logger  *logrus.Logger
	modules map[string]Module
	config  *Config
	mutex   sync.RWMutex

	// Session management
	sessions map[string]*Session

	// Real-time communication
	subscribers map[string]chan *Event

	// API key management
	apiKeyMgr *APIKeyManager

	// Scan management
	scanMgr *ScanManager
}

// Config holds application configuration
type Config struct {
	DatabasePath string                  `json:"database_path"`
	LogLevel     string                  `json:"log_level"`
	APIKeys      map[string]string       `json:"api_keys"`
	RateLimits   map[string]int          `json:"rate_limits"`
	Modules      map[string]ModuleConfig `json:"modules"`
}

// ModuleConfig holds configuration for individual modules
type ModuleConfig struct {
	Enabled bool                   `json:"enabled"`
	Timeout string                 `json:"timeout"`
	Options map[string]interface{} `json:"options"`
}

// Session represents a reconnaissance session
type Session struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	CreatedAt    int64                  `json:"created_at"`
	UpdatedAt    int64                  `json:"updated_at"`
	Status       string                 `json:"status"`
	Target       string                 `json:"target"`
	Results      map[string]interface{} `json:"results"`
	ModuleStates map[string]interface{} `json:"module_states"`
}

// Event represents a real-time event
type Event struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id"`
	Module    string      `json:"module"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// NewApplication creates a new application instance
func NewApplication(db *sql.DB, logger *logrus.Logger) *Application {
	app := &Application{
		db:          db,
		logger:      logger,
		modules:     make(map[string]Module),
		sessions:    make(map[string]*Session),
		subscribers: make(map[string]chan *Event),
	}

	// Load configuration
	app.loadConfig()
	// Initialize API key manager with a default master password
	// In production, this should be user-configurable
	app.apiKeyMgr = NewAPIKeyManager(app, "gorconx-master-key-2024")

	// Initialize scan manager
	app.scanMgr = NewScanManager(app)

	// Initialize modules
	app.initializeModules()

	return app
}

// loadConfig loads application configuration
func (a *Application) loadConfig() {
	// Implementation will load from database or config file
	a.config = &Config{
		DatabasePath: "./data/gorconx.db",
		LogLevel:     "info",
		APIKeys:      make(map[string]string),
		RateLimits:   make(map[string]int),
		Modules:      make(map[string]ModuleConfig),
	}
}

// initializeModules initializes all available reconnaissance modules
func (a *Application) initializeModules() {
	a.logger.Info("Initializing reconnaissance modules...")

	// Get all modules from the global registry
	allModules := modules.GlobalRegistry.GetAll()

	a.mutex.Lock()
	for name, module := range allModules {
		adapter := &ModuleAdapter{module: module}
		a.modules[name] = adapter
		a.logger.Infof("Registered module: %s", name)
	}
	a.mutex.Unlock()

	a.logger.Infof("Loaded %d reconnaissance modules", len(allModules))
}

// RegisterModule registers a new reconnaissance module
func (a *Application) RegisterModule(name string, module Module) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.modules[name] = module
	a.logger.Infof("Registered module: %s", name)
}

// GetModule returns a module by name
func (a *Application) GetModule(name string) (Module, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	module, exists := a.modules[name]
	return module, exists
}

// GetModules returns all registered modules
func (a *Application) GetModules() map[string]Module {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	result := make(map[string]Module)
	for name, module := range a.modules {
		result[name] = module
	}
	return result
}

// CreateSession creates a new reconnaissance session
func (a *Application) CreateSession(name, target string) *Session {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	session := &Session{
		ID:           generateSessionID(),
		Name:         name,
		CreatedAt:    getCurrentTimestamp(),
		UpdatedAt:    getCurrentTimestamp(),
		Status:       "created",
		Target:       target,
		Results:      make(map[string]interface{}),
		ModuleStates: make(map[string]interface{}),
	}

	a.sessions[session.ID] = session
	a.logger.Infof("Created session: %s (%s)", session.Name, session.ID)

	return session
}

// GetSession returns a session by ID
func (a *Application) GetSession(id string) (*Session, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	session, exists := a.sessions[id]
	return session, exists
}

// GetSessions returns all sessions
func (a *Application) GetSessions() []*Session {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	sessions := make([]*Session, 0, len(a.sessions))
	for _, session := range a.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// UpdateSession updates a session
func (a *Application) UpdateSession(session *Session) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	session.UpdatedAt = getCurrentTimestamp()
	a.sessions[session.ID] = session
}

// DeleteSession deletes a session
func (a *Application) DeleteSession(id string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	delete(a.sessions, id)
	a.logger.Infof("Deleted session: %s", id)
}

// Subscribe subscribes to real-time events
func (a *Application) Subscribe(clientID string) chan *Event {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ch := make(chan *Event, 100)
	a.subscribers[clientID] = ch
	return ch
}

// Unsubscribe unsubscribes from real-time events
func (a *Application) Unsubscribe(clientID string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if ch, exists := a.subscribers[clientID]; exists {
		close(ch)
		delete(a.subscribers, clientID)
	}
}

// Publish publishes an event to all subscribers
func (a *Application) Publish(event *Event) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, ch := range a.subscribers {
		select {
		case ch <- event:
		default:
			// Channel is full, skip this subscriber
		}
	}
}

// GetDatabase returns the database connection
func (a *Application) GetDatabase() *sql.DB {
	return a.db
}

// GetLogger returns the logger
func (a *Application) GetLogger() *logrus.Logger {
	return a.logger
}

// GetAPIKeyManager returns the API key manager
func (a *Application) GetAPIKeyManager() *APIKeyManager {
	return a.apiKeyMgr
}

// GetConfig returns the application configuration
func (a *Application) GetConfig() *Config {
	return a.config
}

// GetScanManager returns the scan manager
func (a *Application) GetScanManager() *ScanManager {
	return a.scanMgr
}

// Utility functions
func generateSessionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "session_" + hex.EncodeToString(bytes)
}

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
