package modules

import (
	"GoReconX/internal/config"
	"GoReconX/internal/database"
	"GoReconX/internal/ai"
	"fmt"

	"github.com/sirupsen/logrus"
)

// ModuleManager manages all reconnaissance modules
type ModuleManager struct {
	DB     *database.DB
	Config *config.Config
	Logger *logrus.Logger
	
	// AI client
	AIClient       *ai.GeminiClient
	
	// Module instances
	SubdomainEnum    *SubdomainEnumerator
	EmailHarvester   *EmailHarvester
	PortScanner      *PortScanner
	DirEnumerator    *DirectoryEnumerator
	WebAnalyzer      *WebAnalyzer
	IPGeolocation    *IPGeolocator
	GitHubRecon      *GitHubRecon
}

// NewModuleManager creates a new module manager instance
func NewModuleManager(db *database.DB, cfg *config.Config, logger *logrus.Logger) *ModuleManager {
	mm := &ModuleManager{
		DB:     db,
		Config: cfg,
		Logger: logger,
		
		// Initialize modules
		SubdomainEnum:    NewSubdomainEnumerator(cfg, logger),
		EmailHarvester:   NewEmailHarvester(cfg, logger),
		PortScanner:      NewPortScanner(cfg, logger),
		DirEnumerator:    NewDirectoryEnumerator(cfg, logger),
		WebAnalyzer:      NewWebAnalyzer(cfg, logger),
		IPGeolocation:    NewIPGeolocator(cfg, logger),
		GitHubRecon:      NewGitHubRecon(cfg, logger),
	}
	
	// Initialize AI client if API key is available
	if cfg.API.GeminiKey != "" {
		aiClient, err := ai.NewGeminiClient(cfg.API.GeminiKey, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize AI client")
		} else {
			mm.AIClient = aiClient
			logger.Info("AI client initialized successfully")
		}
	}
	
	return mm
}

// GetAvailableModules returns a list of all available modules
func (mm *ModuleManager) GetAvailableModules() map[string]ModuleInterface {
	return map[string]ModuleInterface{
		"subdomain_enumeration": mm.SubdomainEnum,
		"email_harvesting":      mm.EmailHarvester,
		"port_scanning":         mm.PortScanner,
		"directory_enumeration": mm.DirEnumerator,
		"web_analysis":          mm.WebAnalyzer,
		"ip_geolocation":        mm.IPGeolocation,
		"github_reconnaissance": mm.GitHubRecon,
	}
}

// ExecuteModule executes a specific module
func (mm *ModuleManager) ExecuteModule(moduleName, target string, options map[string]interface{}) (*ScanResult, error) {
	modules := mm.GetAvailableModules()
	
	module, exists := modules[moduleName]
	if !exists {
		return nil, fmt.Errorf("module not found: %s", moduleName)
	}
	
	// Validate target
	if err := module.Validate(target); err != nil {
		return nil, fmt.Errorf("target validation failed: %v", err)
	}
	
	// Execute module
	return module.Execute(target, options)
}

// Close closes any open connections
func (mm *ModuleManager) Close() error {
	if mm.AIClient != nil {
		return mm.AIClient.Close()
	}
	return nil
}

// ScanResult represents the result of a scan operation
type ScanResult struct {
	ModuleName   string                 `json:"module_name"`
	Target       string                 `json:"target"`
	Status       string                 `json:"status"`
	Results      []interface{}          `json:"results"`
	Metadata     map[string]interface{} `json:"metadata"`
	StartTime    string                 `json:"start_time"`
	EndTime      string                 `json:"end_time"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}

// ModuleInterface defines the interface that all modules must implement
type ModuleInterface interface {
	GetName() string
	GetDescription() string
	Validate(target string) error
	Execute(target string, options map[string]interface{}) (*ScanResult, error)
	GetDefaultOptions() map[string]interface{}
}
