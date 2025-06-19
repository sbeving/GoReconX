package modules

import (
	"context"
	"time"
)

// Module defines the interface that all reconnaissance modules must implement
type Module interface {
	// GetInfo returns basic information about the module
	GetInfo() ModuleInfo

	// Validate validates the input parameters for the module
	Validate(input ModuleInput) error

	// Execute runs the module with the given input and returns results
	Execute(ctx context.Context, input ModuleInput, output chan<- ModuleResult) error

	// Stop stops the module execution
	Stop() error

	// GetStatus returns the current status of the module
	GetStatus() ModuleStatus
}

// ModuleInfo contains metadata about a module
type ModuleInfo struct {
	Name         string         `json:"name"`
	Category     string         `json:"category"`
	Description  string         `json:"description"`
	Version      string         `json:"version"`
	Author       string         `json:"author"`
	Tags         []string       `json:"tags"`
	Options      []ModuleOption `json:"options"`
	Requirements []string       `json:"requirements"`
}

// ModuleOption defines a configurable option for a module
type ModuleOption struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // string, int, bool, choice
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
	Choices     []string    `json:"choices,omitempty"`
	Validation  string      `json:"validation,omitempty"`
}

// ModuleInput contains input parameters for module execution
type ModuleInput struct {
	Target    string                 `json:"target"`
	Options   map[string]interface{} `json:"options"`
	SessionID string                 `json:"session_id"`
	Timeout   time.Duration          `json:"timeout"`
}

// ModuleResult represents a single result from module execution
type ModuleResult struct {
	Type      string                 `json:"type"` // progress, data, error, complete
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	SessionID string                 `json:"session_id"`
	Module    string                 `json:"module"`
}

// ModuleStatus represents the current status of a module
type ModuleStatus struct {
	IsRunning   bool          `json:"is_running"`
	Progress    float64       `json:"progress"` // 0.0 to 1.0
	StartTime   time.Time     `json:"start_time"`
	ElapsedTime time.Duration `json:"elapsed_time"`
	Status      string        `json:"status"` // idle, running, completed, error, stopped
	Message     string        `json:"message"`
}

// BaseModule provides a basic implementation that modules can embed
type BaseModule struct {
	info   ModuleInfo
	status ModuleStatus
	stopCh chan bool
}

// NewBaseModule creates a new base module
func NewBaseModule(info ModuleInfo) *BaseModule {
	return &BaseModule{
		info: info,
		status: ModuleStatus{
			IsRunning: false,
			Progress:  0.0,
			Status:    "idle",
		},
		stopCh: make(chan bool, 1),
	}
}

// GetInfo returns the module information
func (b *BaseModule) GetInfo() ModuleInfo {
	return b.info
}

// GetStatus returns the current module status
func (b *BaseModule) GetStatus() ModuleStatus {
	return b.status
}

// Stop stops the module execution
func (b *BaseModule) Stop() error {
	select {
	case b.stopCh <- true:
	default:
	}
	b.status.IsRunning = false
	b.status.Status = "stopped"
	return nil
}

// SetStatus updates the module status
func (b *BaseModule) SetStatus(status string, progress float64, message string) {
	b.status.Status = status
	b.status.Progress = progress
	b.status.Message = message
	if status == "running" && !b.status.IsRunning {
		b.status.IsRunning = true
		b.status.StartTime = time.Now()
	} else if status == "completed" || status == "error" || status == "stopped" {
		b.status.IsRunning = false
		b.status.ElapsedTime = time.Since(b.status.StartTime)
	}
}

// IsStopped checks if the module should stop
func (b *BaseModule) IsStopped() bool {
	select {
	case <-b.stopCh:
		return true
	default:
		return false
	}
}

// SendResult sends a result through the output channel
func (b *BaseModule) SendResult(output chan<- ModuleResult, resultType string, data interface{}, metadata map[string]interface{}, sessionID string) {
	result := ModuleResult{
		Type:      resultType,
		Data:      data,
		Metadata:  metadata,
		Timestamp: time.Now(),
		SessionID: sessionID,
		Module:    b.info.Name,
	}

	select {
	case output <- result:
	default:
		// Channel is full or closed
	}
}

// ValidateInput provides basic input validation
func (b *BaseModule) ValidateInput(input ModuleInput) error {
	// Basic validation - submodules can override
	if input.Target == "" {
		return NewModuleError("target is required", "INVALID_INPUT")
	}
	return nil
}

// ModuleError represents an error from module execution
type ModuleError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

func (e *ModuleError) Error() string {
	return e.Message
}

// NewModuleError creates a new module error
func NewModuleError(message, code string) *ModuleError {
	return &ModuleError{
		Message: message,
		Code:    code,
	}
}

// Registry manages all available modules
type Registry struct {
	modules map[string]Module
}

// NewRegistry creates a new module registry
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register registers a module in the registry
func (r *Registry) Register(module Module) {
	info := module.GetInfo()
	r.modules[info.Name] = module
}

// Get returns a module by name
func (r *Registry) Get(name string) (Module, bool) {
	module, exists := r.modules[name]
	return module, exists
}

// List returns all registered modules
func (r *Registry) List() map[string]Module {
	result := make(map[string]Module)
	for name, module := range r.modules {
		result[name] = module
	}
	return result
}

// ListByCategory returns modules filtered by category
func (r *Registry) ListByCategory(category string) map[string]Module {
	result := make(map[string]Module)
	for name, module := range r.modules {
		if module.GetInfo().Category == category {
			result[name] = module
		}
	}
	return result
}
