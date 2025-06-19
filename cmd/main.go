package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorconx/internal/api"
	"gorconx/internal/core"
	"gorconx/internal/database"
	"gorconx/internal/gui"
	"gorconx/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "1.0.0"
	cfgFile string
	port    int
	host    string
	devMode bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gorconx",
	Short: "GoReconX - Advanced OSINT & Reconnaissance Platform",
	Long: `
  ██████   ██████  ██████  ███████  ██████  ██████  ███    ██ ██   ██ 
 ██       ██    ██ ██   ██ ██      ██      ██    ██ ████   ██  ██ ██  
 ██   ███ ██    ██ ██████  █████   ██      ██    ██ ██ ██  ██   ███   
 ██    ██ ██    ██ ██   ██ ██      ██      ██    ██ ██  ██ ██  ██ ██  
  ██████   ██████  ██   ██ ███████  ██████  ██████  ██   ████ ██   ██ 

GoReconX - A Comprehensive OSINT & Reconnaissance Platform

Built for cybersecurity professionals, red teamers, and security enthusiasts.
Ethical use only - Ensure you have explicit permission before scanning any target.

Features:
  • Modular reconnaissance modules
  • Modern web-based GUI
  • Encrypted data storage
  • Real-time results
  • Comprehensive reporting
  • API integrations
`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func main() {
	Execute()
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/config.yaml)")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to run the web server on")
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "host to bind the web server to")
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "run in development mode")

	// Bind flags to viper
	viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("app.dev_mode", rootCmd.PersistentFlags().Lookup("dev"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// Create default config if it doesn't exist
		createDefaultConfig()
	}
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("app.name", "GoReconX")
	viper.SetDefault("app.version", version)
	viper.SetDefault("app.dev_mode", false)
	viper.SetDefault("database.path", "./data/gorconx.db")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.file", "./logs/gorconx.log")
}

func createDefaultConfig() {
	configDir := "./configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v", err)
		return
	}

	configContent := `# GoReconX Configuration File

server:
  host: localhost
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

app:
  name: GoReconX
  version: ` + version + `
  dev_mode: false

database:
  path: ./data/gorconx.db

logging:
  level: info
  file: ./logs/gorconx.log

# External API Configuration (encrypted storage)
apis:
  # Add your API keys here - they will be encrypted before storage
  # virustotal_api_key: ""
  # shodan_api_key: ""
  # hunter_io_api_key: ""

# Module Configuration
modules:
  domain_enum:
    enabled: true
    timeout: 30s
    wordlists:
      - ./wordlists/subdomains.txt
  
  port_scan:
    enabled: true
    timeout: 60s
    max_threads: 100
  
  web_enum:
    enabled: true
    timeout: 30s
    user_agent: "GoReconX/` + version + ` (Ethical Security Scanner)"

# Rate limiting for API calls
rate_limits:
  default: 10  # requests per second
  virustotal: 4
  shodan: 1
`

	configPath := "./configs/config.yaml"
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		log.Printf("Failed to create default config file: %v", err)
	} else {
		fmt.Printf("Created default config file: %s\n", configPath)
		viper.SetConfigFile(configPath)
		viper.ReadInConfig()
	}
}

func startServer() {
	// Initialize logger
	logger := utils.InitLogger()

	// Create necessary directories
	createDirectories()

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize core application
	app := core.NewApplication(db, logger)

	// Initialize API server
	server := api.NewServer(app)

	// Initialize GUI
	guiServer := gui.NewGUIServer(app)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start servers
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Failed to start API server: %v", err)
		}
	}()

	go func() {
		if err := guiServer.Start(); err != nil {
			logger.Fatalf("Failed to start GUI server: %v", err)
		}
	}()

	// Print startup information
	printStartupInfo()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Failed to shutdown API server: %v", err)
	}

	if err := guiServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Failed to shutdown GUI server: %v", err)
	}

	logger.Info("Servers stopped")
}

func createDirectories() {
	dirs := []string{
		"./data",
		"./logs",
		"./reports",
		"./wordlists",
		"./configs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
		}
	}
}

func printStartupInfo() {
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")

	fmt.Printf("\n")
	fmt.Printf("🚀 GoReconX v%s is starting...\n", version)
	fmt.Printf("🌐 Web Interface: http://%s:%d\n", host, port)
	fmt.Printf("🔧 API Endpoint: http://%s:%d/api\n", host, port)
	fmt.Printf("📊 Dashboard: http://%s:%d/dashboard\n", host, port)
	fmt.Printf("\n")
	fmt.Printf("⚖️  ETHICAL USE ONLY - Ensure you have explicit permission before scanning any target\n")
	fmt.Printf("\n")
	fmt.Printf("Press Ctrl+C to stop the server\n")
	fmt.Printf("\n")
}
