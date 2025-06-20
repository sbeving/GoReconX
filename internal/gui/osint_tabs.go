package gui

import (
	"GoReconX/internal/modules"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

// PassiveOSINTTab represents the passive OSINT tab
type PassiveOSINTTab struct {
	modules *modules.ModuleManager
	logger  *logrus.Logger
	content fyne.CanvasObject
}

// NewPassiveOSINTTab creates a new passive OSINT tab
func NewPassiveOSINTTab(moduleManager *modules.ModuleManager, logger *logrus.Logger) *PassiveOSINTTab {
	tab := &PassiveOSINTTab{
		modules: moduleManager,
		logger:  logger,
	}
	tab.setupContent()
	return tab
}

// setupContent initializes the passive OSINT content
func (pot *PassiveOSINTTab) setupContent() {
	// Target input
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Enter target domain (e.g., example.com)")

	// Module selection
	moduleSelect := widget.NewSelect([]string{
		"Subdomain Enumeration",
		"Email Harvesting",
		"Web Analysis",
		"IP Geolocation",
		"GitHub Reconnaissance",
	}, nil)
	moduleSelect.SetSelected("Subdomain Enumeration")

	// Options
	optionsCard := widget.NewCard("Options", "",
		container.NewVBox(
			widget.NewLabel("Threads:"),
			widget.NewSlider(1, 100),
			widget.NewLabel("Timeout (seconds):"),
			widget.NewSlider(1, 30),
			widget.NewCheck("Resolve IPs", nil),
		))

	// Control buttons
	runButton := widget.NewButton("Run Scan", func() {
		target := targetEntry.Text
		if target == "" {
			pot.logger.Warn("No target specified")
			return
		}

		pot.logger.WithField("target", target).Info("Starting passive OSINT scan")
		// TODO: Implement scan execution
	})

	clearButton := widget.NewButton("Clear", func() {
		targetEntry.SetText("")
	})

	// Output console
	outputText := widget.NewRichTextFromMarkdown("Ready to start passive reconnaissance...")
	outputText.Resize(fyne.NewSize(600, 300))
	outputScroll := container.NewScroll(outputText)
	outputScroll.SetMinSize(fyne.NewSize(600, 300))

	// Layout
	inputSection := widget.NewCard("Target & Module", "",
		container.NewVBox(
			widget.NewLabel("Target:"),
			targetEntry,
			widget.NewLabel("Module:"),
			moduleSelect,
			container.NewHBox(runButton, clearButton),
		))

	outputSection := widget.NewCard("Output", "", outputScroll)

	pot.content = container.NewHSplit(
		container.NewVBox(inputSection, optionsCard),
		outputSection,
	)
}

// Content returns the tab content
func (pot *PassiveOSINTTab) Content() fyne.CanvasObject {
	return pot.content
}

// ActiveReconTab represents the active reconnaissance tab
type ActiveReconTab struct {
	modules *modules.ModuleManager
	logger  *logrus.Logger
	content fyne.CanvasObject
}

// NewActiveReconTab creates a new active reconnaissance tab
func NewActiveReconTab(moduleManager *modules.ModuleManager, logger *logrus.Logger) *ActiveReconTab {
	tab := &ActiveReconTab{
		modules: moduleManager,
		logger:  logger,
	}
	tab.setupContent()
	return tab
}

// setupContent initializes the active recon content
func (art *ActiveReconTab) setupContent() {
	// Warning message
	warningCard := widget.NewCard("⚠️ Warning", "",
		widget.NewRichTextFromMarkdown(`
**Important**: Active reconnaissance involves direct interaction with target systems.
Only use these modules on systems you own or have explicit permission to test.
Unauthorized scanning may be illegal in your jurisdiction.
		`))

	// Target input
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Enter target IP or domain")

	// Module selection
	moduleSelect := widget.NewSelect([]string{
		"Port Scanner",
		"Directory Enumeration",
		"Service Detection",
	}, nil)
	moduleSelect.SetSelected("Port Scanner")

	// Control buttons
	runButton := widget.NewButton("Run Scan", func() {
		target := targetEntry.Text
		if target == "" {
			art.logger.Warn("No target specified")
			return
		}

		art.logger.WithField("target", target).Info("Starting active reconnaissance")
		// TODO: Implement scan execution
	})

	// Output console
	outputText := widget.NewRichTextFromMarkdown("Ready for active reconnaissance...")
	outputScroll := container.NewScroll(outputText)
	outputScroll.SetMinSize(fyne.NewSize(600, 400))

	// Layout
	inputSection := widget.NewCard("Target & Module", "",
		container.NewVBox(
			widget.NewLabel("Target:"),
			targetEntry,
			widget.NewLabel("Module:"),
			moduleSelect,
			runButton,
		))

	outputSection := widget.NewCard("Output", "", outputScroll)

	art.content = container.NewVBox(
		warningCard,
		container.NewHSplit(inputSection, outputSection),
	)
}

// Content returns the tab content
func (art *ActiveReconTab) Content() fyne.CanvasObject {
	return art.content
}
