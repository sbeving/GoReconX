package appinstance

import (
	"GoReconX/internal/config"
	"GoReconX/internal/database"
	"GoReconX/internal/gui"
	"GoReconX/internal/reports"
	"GoReconX/internal/ai"

	"fyne.io/fyne/v2"
	"github.com/sirupsen/logrus"
)

// App represents the main application instance
type App struct {
	FyneApp fyne.App
	DB      *database.DB
	Config  *config.Config
	Logger  *logrus.Logger
	GUI     *gui.MainWindow
	
	// AI and Reporting
	AIClient  *ai.GeminiClient
	ReportGen *reports.ReportGenerator
}

// NewApp creates a new application instance
func NewApp(fyneApp fyne.App, db *database.DB, cfg *config.Config, logger *logrus.Logger) *App {
	app := &App{
		FyneApp: fyneApp,
		DB:      db,
		Config:  cfg,
		Logger:  logger,
	}
	
	// Initialize AI client if API key is available
	if cfg.API.GeminiKey != "" {
		aiClient, err := ai.NewGeminiClient(cfg.API.GeminiKey, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize AI client")
		} else {
			app.AIClient = aiClient
			logger.Info("AI-powered analysis enabled")
		}
	} else {
		logger.Info("No Gemini API key configured - AI features disabled")
	}
	
	// Initialize report generator
	app.ReportGen = reports.NewReportGenerator(logger, app.AIClient, cfg.Output.OutputDir)

	// Initialize the GUI
	app.GUI = gui.NewMainWindow(fyneApp, db, cfg, logger)

	return app
}

// Run starts the application
func (app *App) Run() {
	app.Logger.Info("Starting GoReconX...")
	app.GUI.Show()
	app.FyneApp.Run()
	
	// Cleanup
	if app.AIClient != nil {
		app.AIClient.Close()
	}
}
