package gui

import (
	"GoReconX/internal/config"
	"GoReconX/internal/database"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

// ResultsTab represents the results viewer tab
type ResultsTab struct {
	db      *database.DB
	logger  *logrus.Logger
	content fyne.CanvasObject
}

// NewResultsTab creates a new results tab
func NewResultsTab(db *database.DB, logger *logrus.Logger) *ResultsTab {
	tab := &ResultsTab{
		db:     db,
		logger: logger,
	}
	tab.setupContent()
	return tab
}

// setupContent initializes the results content
func (rt *ResultsTab) setupContent() {
	// Results table
	table := widget.NewTable(
		func() (int, int) { return 0, 4 }, // rows, cols
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, template fyne.CanvasObject) {
			template.(*widget.Label).SetText("No data")
		},
	)

	// Filter options
	filterCard := widget.NewCard("Filters", "",
		container.NewVBox(
			widget.NewLabel("Project:"),
			widget.NewSelect([]string{"All Projects"}, nil),
			widget.NewLabel("Scan Type:"),
			widget.NewSelect([]string{"All Types"}, nil),
			widget.NewLabel("Status:"),
			widget.NewSelect([]string{"All", "Completed", "Failed"}, nil),
			widget.NewButton("Apply Filters", func() {
				rt.logger.Info("Applying result filters")
			}),
		))

	// Export options
	exportCard := widget.NewCard("Export", "",
		container.NewVBox(
			widget.NewButton("Export to JSON", func() {
				rt.logger.Info("Exporting results to JSON")
			}),
			widget.NewButton("Export to HTML", func() {
				rt.logger.Info("Exporting results to HTML")
			}),
			widget.NewButton("Export to PDF", func() {
				rt.logger.Info("Exporting results to PDF")
			}),
		))

	// Layout
	sidebar := container.NewVBox(filterCard, exportCard)
	rt.content = container.NewHSplit(sidebar, table)
}

// Content returns the tab content
func (rt *ResultsTab) Content() fyne.CanvasObject {
	return rt.content
}

// UtilitiesTab represents the utilities tab
type UtilitiesTab struct {
	db      *database.DB
	config  *config.Config
	logger  *logrus.Logger
	content fyne.CanvasObject
}

// NewUtilitiesTab creates a new utilities tab
func NewUtilitiesTab(db *database.DB, cfg *config.Config, logger *logrus.Logger) *UtilitiesTab {
	tab := &UtilitiesTab{
		db:     db,
		config: cfg,
		logger: logger,
	}
	tab.setupContent()
	return tab
}

// setupContent initializes the utilities content
func (ut *UtilitiesTab) setupContent() {
	// Project management
	projectCard := widget.NewCard("Project Management", "",
		container.NewVBox(
			widget.NewButton("Create New Project", func() {
				ut.logger.Info("Creating new project")
			}),
			widget.NewButton("Import Project", func() {
				ut.logger.Info("Importing project")
			}),
			widget.NewButton("Delete Project", func() {
				ut.logger.Info("Deleting project")
			}),
		))

	// Wordlist management
	wordlistCard := widget.NewCard("Wordlist Management", "",
		container.NewVBox(
			widget.NewButton("Import Wordlists", func() {
				ut.logger.Info("Importing wordlists")
			}),
			widget.NewButton("Create Custom Wordlist", func() {
				ut.logger.Info("Creating custom wordlist")
			}),
			widget.NewButton("Download SecLists", func() {
				ut.logger.Info("Downloading SecLists")
			}),
		))

	// Session management
	sessionCard := widget.NewCard("Session Management", "",
		container.NewVBox(
			widget.NewButton("Save Current Session", func() {
				ut.logger.Info("Saving current session")
			}),
			widget.NewButton("Load Session", func() {
				ut.logger.Info("Loading session")
			}),
			widget.NewButton("Export Session Data", func() {
				ut.logger.Info("Exporting session data")
			}),
		))

	// AI Analysis
	aiCard := widget.NewCard("AI-Powered Analysis", "",
		container.NewVBox(
			widget.NewButton("Analyze Recent Results", func() {
				ut.logger.Info("Analyzing results with AI")
			}),
			widget.NewButton("Generate Executive Summary", func() {
				ut.logger.Info("Generating AI summary")
			}),
			widget.NewButton("Get Recommendations", func() {
				ut.logger.Info("Getting AI recommendations")
			}),
		))

	// Layout
	ut.content = container.NewGridWithColumns(2,
		projectCard, wordlistCard,
		sessionCard, aiCard,
	)
}

// Content returns the tab content
func (ut *UtilitiesTab) Content() fyne.CanvasObject {
	return ut.content
}

// SettingsTab represents the settings tab
type SettingsTab struct {
	db      *database.DB
	config  *config.Config
	logger  *logrus.Logger
	content fyne.CanvasObject
}

// NewSettingsTab creates a new settings tab
func NewSettingsTab(db *database.DB, cfg *config.Config, logger *logrus.Logger) *SettingsTab {
	tab := &SettingsTab{
		db:     db,
		config: cfg,
		logger: logger,
	}
	tab.setupContent()
	return tab
}

// setupContent initializes the settings content
func (st *SettingsTab) setupContent() {
	// API Keys section
	apiCard := widget.NewCard("API Keys", "",
		container.NewVBox(
			widget.NewLabel("Google Gemini API Key:"),
			widget.NewPasswordEntry(),
			widget.NewLabel("VirusTotal API Key:"),
			widget.NewPasswordEntry(),
			widget.NewLabel("Shodan API Key:"),
			widget.NewPasswordEntry(),
			widget.NewLabel("Hunter.io API Key:"),
			widget.NewPasswordEntry(),
			widget.NewButton("Save API Keys", func() {
				st.logger.Info("Saving API keys")
			}),
		))

	// Network settings
	networkCard := widget.NewCard("Network Settings", "",
		container.NewVBox(
			widget.NewLabel("HTTP Timeout (seconds):"),
			widget.NewSlider(5, 60),
			widget.NewLabel("Retry Attempts:"),
			widget.NewSlider(1, 10),
			widget.NewLabel("User Agent:"),
			widget.NewEntry(),
			widget.NewLabel("Proxy URL (optional):"),
			widget.NewEntry(),
			widget.NewButton("Save Network Settings", func() {
				st.logger.Info("Saving network settings")
			}),
		))

	// Application settings
	appCard := widget.NewCard("Application Settings", "",
		container.NewVBox(
			widget.NewLabel("Default Output Format:"),
			widget.NewSelect([]string{"JSON", "HTML", "PDF"}, nil),
			widget.NewLabel("Output Directory:"),
			widget.NewEntry(),
			widget.NewLabel("Log Level:"),
			widget.NewSelect([]string{"DEBUG", "INFO", "WARN", "ERROR"}, nil),
			widget.NewCheck("Enable AI Features", nil),
			widget.NewButton("Save App Settings", func() {
				st.logger.Info("Saving application settings")
			}),
		))

	// About section
	aboutCard := widget.NewCard("About GoReconX", "",
		widget.NewRichTextFromMarkdown(`
**GoReconX v1.0**

A comprehensive OSINT & reconnaissance platform built with Go.

**Features:**
- Passive OSINT modules
- Active reconnaissance tools  
- AI-powered analysis with Google Gemini
- Professional reporting
- Ethical security focus

**Developed for cybersecurity professionals and researchers.**

Remember to always use this tool ethically and with proper authorization.
		`))

	// Layout
	st.content = container.NewVBox(
		container.NewGridWithColumns(2, apiCard, networkCard),
		container.NewGridWithColumns(2, appCard, aboutCard),
	)
}

// Content returns the tab content
func (st *SettingsTab) Content() fyne.CanvasObject {
	return st.content
}
