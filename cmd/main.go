package main

import (
	"GoReconX/internal/appinstance"
	"GoReconX/internal/config"
	"GoReconX/internal/database"
	"GoReconX/internal/logging"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Initialize logging
	logger := logging.InitLogger()
	logger.Info("Starting GoReconX - Comprehensive OSINT & Reconnaissance Platform")

	// Show ethical usage disclaimer
	if !showDisclaimerDialog() {
		logger.Info("User declined disclaimer. Exiting.")
		return
	}

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}
	defer db.Close()

	// Create and run the main application
	myApp := app.NewWithID("com.goreconx.app")
	// TODO: Set icon once resource is generated
	
	mainApp := appinstance.NewApp(myApp, db, cfg, logger)
	mainApp.Run()
}

func showDisclaimerDialog() bool {
	disclaimerApp := app.New()
	disclaimerWindow := disclaimerApp.NewWindow("GoReconX - Ethical Usage Agreement")
	disclaimerWindow.Resize(fyne.NewSize(600, 400))

	disclaimerText := `ETHICAL USAGE DISCLAIMER

GoReconX is designed for legitimate cybersecurity professionals, ethical hackers, 
and security researchers for lawful reconnaissance and OSINT activities.

By using this tool, you agree to:
• Only use it on systems you own or have explicit permission to test
• Comply with all applicable laws and regulations
• Not use it for malicious purposes or unauthorized access
• Take responsibility for your actions and their consequences

The developers are not responsible for any misuse of this tool.

Do you agree to use GoReconX ethically and legally?`

	agreed := false
	
	content := container.NewVBox(
		widget.NewLabel("GoReconX - Ethical Usage Agreement"),
		widget.NewSeparator(),
		widget.NewRichTextFromMarkdown(disclaimerText),
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewButton("I Agree", func() {
				agreed = true
				disclaimerWindow.Close()
			}),
			widget.NewButton("Cancel", func() {
				agreed = false
				disclaimerWindow.Close()
			}),
		),
	)

	disclaimerWindow.SetContent(content)
	disclaimerWindow.ShowAndRun()
	
	return agreed
}