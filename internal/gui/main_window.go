package gui

import (
	"GoReconX/internal/config"
	"GoReconX/internal/database"
	"GoReconX/internal/modules"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

// MainWindow represents the main application window
type MainWindow struct {
	Window  fyne.Window
	App     fyne.App
	DB      *database.DB
	Config  *config.Config
	Logger  *logrus.Logger
	Modules *modules.ModuleManager

	// GUI components
	content   *container.AppTabs
	dashboard *DashboardTab
	passive   *PassiveOSINTTab
	active    *ActiveReconTab
	utilities *UtilitiesTab
	settings  *SettingsTab
	results   *ResultsTab
}

// NewMainWindow creates a new main window
func NewMainWindow(app fyne.App, db *database.DB, cfg *config.Config, logger *logrus.Logger) *MainWindow {
	window := app.NewWindow("GoReconX - Comprehensive OSINT & Reconnaissance Platform")
	window.Resize(fyne.NewSize(1200, 800))
	window.CenterOnScreen()

	// Initialize module manager
	moduleManager := modules.NewModuleManager(db, cfg, logger)

	mainWindow := &MainWindow{
		Window:  window,
		App:     app,
		DB:      db,
		Config:  cfg,
		Logger:  logger,
		Modules: moduleManager,
	}

	mainWindow.setupUI()
	return mainWindow
}

// setupUI initializes the user interface
func (mw *MainWindow) setupUI() {
	// Create tabs
	mw.dashboard = NewDashboardTab(mw.DB, mw.Config, mw.Logger)
	mw.passive = NewPassiveOSINTTab(mw.Modules, mw.Logger)
	mw.active = NewActiveReconTab(mw.Modules, mw.Logger)
	mw.utilities = NewUtilitiesTab(mw.DB, mw.Config, mw.Logger)
	mw.settings = NewSettingsTab(mw.DB, mw.Config, mw.Logger)
	mw.results = NewResultsTab(mw.DB, mw.Logger)

	// Create main tab container
	mw.content = container.NewAppTabs(
		container.NewTabItemWithIcon("Dashboard", theme.HomeIcon(), mw.dashboard.Content()),
		container.NewTabItemWithIcon("Passive OSINT", theme.SearchIcon(), mw.passive.Content()),
		container.NewTabItemWithIcon("Active Recon", theme.ComputerIcon(), mw.active.Content()),
		container.NewTabItemWithIcon("Results", theme.DocumentIcon(), mw.results.Content()),
		container.NewTabItemWithIcon("Utilities", theme.FolderIcon(), mw.utilities.Content()),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), mw.settings.Content()),
	)

	// Set main content
	mw.Window.SetContent(mw.content)

	// Set window icon and other properties
	mw.Window.SetMaster()
}

// Show displays the main window
func (mw *MainWindow) Show() {
	mw.Window.Show()
}

// DashboardTab represents the dashboard tab
type DashboardTab struct {
	db      *database.DB
	config  *config.Config
	logger  *logrus.Logger
	content fyne.CanvasObject
}

// NewDashboardTab creates a new dashboard tab
func NewDashboardTab(db *database.DB, cfg *config.Config, logger *logrus.Logger) *DashboardTab {
	tab := &DashboardTab{
		db:     db,
		config: cfg,
		logger: logger,
	}
	tab.setupContent()
	return tab
}

// setupContent initializes the dashboard content
func (dt *DashboardTab) setupContent() {
	// Welcome message
	welcome := widget.NewCard("Welcome to GoReconX", "Comprehensive OSINT & Reconnaissance Platform",
		widget.NewRichTextFromMarkdown(`
## Getting Started

GoReconX is your all-in-one solution for ethical reconnaissance and OSINT activities.

### Quick Actions:
- **Create a new project** to organize your reconnaissance activities
- **Run passive OSINT** to gather information without direct interaction
- **Perform active reconnaissance** for detailed target analysis
- **Generate reports** with AI-powered insights

### Features:
- 🔍 **Passive OSINT**: Domain enumeration, email harvesting, public code search
- 🎯 **Active Reconnaissance**: Port scanning, directory enumeration
- 🤖 **AI-Powered Analysis**: Google Gemini integration for smart insights
- 📊 **Professional Reports**: HTML, PDF, and JSON export options
- 🛡️ **Ethical Focus**: Built with security best practices

Remember: Always use GoReconX ethically and with proper authorization!
		`))

	// Recent projects section
	projectsCard := widget.NewCard("Recent Projects", "", widget.NewLabel("Loading projects..."))

	// Quick stats
	statsCard := widget.NewCard("Statistics", "",
		container.NewVBox(
			widget.NewLabel("Total Projects: 0"),
			widget.NewLabel("Completed Scans: 0"),
			widget.NewLabel("Active Sessions: 0"),
		))

	// Quick actions
	actionsCard := widget.NewCard("Quick Actions", "",
		container.NewVBox(
			widget.NewButton("Create New Project", func() {
				// TODO: Implement project creation dialog
			}),
			widget.NewButton("Import Wordlists", func() {
				// TODO: Implement wordlist import
			}),
			widget.NewButton("Export Results", func() {
				// TODO: Implement results export
			}),
		))

	// Layout the dashboard
	dt.content = container.NewVBox(
		welcome,
		container.NewGridWithColumns(3,
			projectsCard,
			statsCard,
			actionsCard,
		),
	)
}

// Content returns the tab content
func (dt *DashboardTab) Content() fyne.CanvasObject {
	return dt.content
}
