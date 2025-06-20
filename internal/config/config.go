package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	
	API struct {
		GeminiKey    string            `yaml:"gemini_key"`
		VirusTotal   string            `yaml:"virustotal_key"`
		Shodan       string            `yaml:"shodan_key"`
		Hunter       string            `yaml:"hunter_key"`
		GitHub       string            `yaml:"github_key"`
		CustomAPIs   map[string]string `yaml:"custom_apis"`
	} `yaml:"api"`
	
	Network struct {
		Timeout    int    `yaml:"timeout"`
		Retries    int    `yaml:"retries"`
		ProxyURL   string `yaml:"proxy_url"`
		UserAgent  string `yaml:"user_agent"`
	} `yaml:"network"`
	
	Wordlists struct {
		Subdomains   string `yaml:"subdomains"`
		Directories  string `yaml:"directories"`
		Files        string `yaml:"files"`
		Ports        string `yaml:"ports"`
	} `yaml:"wordlists"`
	
	Output struct {
		DefaultFormat string `yaml:"default_format"`
		OutputDir     string `yaml:"output_dir"`
	} `yaml:"output"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Database: struct {
			Path string `yaml:"path"`
		}{
			Path: "data/goreconx.db",
		},
		Network: struct {
			Timeout    int    `yaml:"timeout"`
			Retries    int    `yaml:"retries"`
			ProxyURL   string `yaml:"proxy_url"`
			UserAgent  string `yaml:"user_agent"`
		}{
			Timeout:   30,
			Retries:   3,
			UserAgent: "GoReconX/1.0 (OSINT Tool)",
		},
		Wordlists: struct {
			Subdomains   string `yaml:"subdomains"`
			Directories  string `yaml:"directories"`
			Files        string `yaml:"files"`
			Ports        string `yaml:"ports"`
		}{
			Subdomains:  "wordlists/subdomains.txt",
			Directories: "wordlists/directories.txt",
			Files:       "wordlists/files.txt",
			Ports:       "wordlists/ports.txt",
		},
		Output: struct {
			DefaultFormat string `yaml:"default_format"`
			OutputDir     string `yaml:"output_dir"`
		}{
			DefaultFormat: "json",
			OutputDir:     "output",
		},
	}
}

// LoadConfig loads configuration from file or creates default config
func LoadConfig() (*Config, error) {
	configPath := "config/config.yaml"
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}
	
	// If config file doesn't exist, create it with default values
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := SaveConfig(cfg, configPath); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config, configPath string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	
	return os.WriteFile(configPath, data, 0644)
}
