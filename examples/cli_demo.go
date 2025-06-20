package main

import (
	"GoReconX/internal/ai"
	"GoReconX/internal/config"
	"GoReconX/internal/database"
	"GoReconX/internal/logging"
	"GoReconX/internal/modules"
	"GoReconX/internal/reports"
	"fmt"
	"log"
	"os"
	
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logging
	logger := logging.InitLogger()
	logger.Info("Starting GoReconX CLI Demo")

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize module manager
	moduleManager := modules.NewModuleManager(db, cfg, logger)
	defer moduleManager.Close()

	// Example target
	target := "example.com"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	fmt.Printf("🎯 GoReconX CLI Demo - Scanning target: %s\n\n", target)

	// Run subdomain enumeration
	fmt.Println("🔍 Running Subdomain Enumeration...")
	subdomainResult, err := runSubdomainEnum(moduleManager, target)
	if err != nil {
		logger.WithError(err).Error("Subdomain enumeration failed")
	} else {
		fmt.Printf("✅ Found %d subdomains\n", len(subdomainResult.Results))
	}

	// Run port scanning on discovered subdomains
	fmt.Println("\n🔌 Running Port Scanning...")
	portResults := runPortScanning(moduleManager, target, logger)
	fmt.Printf("✅ Completed port scans on %d targets\n", len(portResults))

	// Collect all results
	var allResults []*modules.ScanResult
	if subdomainResult != nil {
		allResults = append(allResults, subdomainResult)
	}
	allResults = append(allResults, portResults...)

	// Initialize AI client for analysis
	var aiClient *ai.GeminiClient
	if cfg.API.GeminiKey != "" {
		aiClient, err = ai.NewGeminiClient(cfg.API.GeminiKey, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize AI client")
		} else {
			fmt.Println("\n🤖 AI Analysis enabled")
		}
	} else {
		fmt.Println("\n⚠️  No Gemini API key configured - AI features disabled")
	}

	// Generate comprehensive report
	fmt.Println("\n📊 Generating Report...")
	reportGen := reports.NewReportGenerator(logger, aiClient, cfg.Output.OutputDir)
	report, err := reportGen.GenerateReport(target, allResults)
	if err != nil {
		logger.WithError(err).Error("Failed to generate report")
		return
	}

	// Export reports in multiple formats
	jsonFile, _ := reportGen.ExportJSON(report)
	htmlFile, _ := reportGen.ExportHTML(report)
	csvFile, _ := reportGen.ExportCSV(report)

	fmt.Printf("✅ Reports generated:\n")
	fmt.Printf("   📄 JSON: %s\n", jsonFile)
	fmt.Printf("   🌐 HTML: %s\n", htmlFile)
	fmt.Printf("   📊 CSV:  %s\n", csvFile)

	// Display summary
	fmt.Printf("\n📋 Summary:\n")
	fmt.Printf("   Target: %s\n", report.Target)
	fmt.Printf("   Total Scans: %d\n", len(report.Results))
	fmt.Printf("   Success Rate: %.1f%%\n", report.Statistics["success_rate"])
	fmt.Printf("   Total Findings: %d\n", report.Statistics["total_results"])

	if report.AIAnalysis != nil {
		fmt.Printf("\n🤖 AI Analysis:\n")
		fmt.Printf("   Threat Level: %s\n", report.AIAnalysis.ThreatLevel)
		fmt.Printf("   Confidence: %.1f%%\n", report.AIAnalysis.Confidence*100)
		fmt.Printf("   Summary: %s\n", truncateString(report.AIAnalysis.Summary, 200))
		
		if len(report.AIAnalysis.Recommendations) > 0 {
			fmt.Printf("\n💡 Key Recommendations:\n")
			for i, rec := range report.AIAnalysis.Recommendations {
				if i >= 3 { // Limit to top 3 recommendations
					break
				}
				fmt.Printf("   %d. %s\n", i+1, truncateString(rec, 100))
			}
		}
	}

	fmt.Println("\n🎉 GoReconX CLI Demo completed successfully!")
	fmt.Println("💡 Tip: Use the GUI for a more interactive experience with real-time updates")
}

func runSubdomainEnum(moduleManager *modules.ModuleManager, target string) (*modules.ScanResult, error) {
	options := map[string]interface{}{
		"threads":     20,
		"timeout":     3,
		"resolve_ips": true,
	}

	result, err := moduleManager.ExecuteModule("subdomain_enumeration", target, options)
	if err != nil {
		return nil, err
	}

	// Display some results
	if len(result.Results) > 0 {
		fmt.Printf("   Sample subdomains found:\n")
		for i, item := range result.Results {
			if i >= 5 { // Show only first 5
				break
			}
			if subdomain, ok := item.(map[string]interface{}); ok {
				if name, exists := subdomain["subdomain"]; exists {
					fmt.Printf("   - %s\n", name)
				}
			}
		}
		if len(result.Results) > 5 {
			fmt.Printf("   ... and %d more\n", len(result.Results)-5)
		}
	}

	return result, nil
}

func runPortScanning(moduleManager *modules.ModuleManager, target string, logger *log.Logger) []*modules.ScanResult {
	var results []*modules.ScanResult

	// Scan common ports on the main target
	targets := []string{target}

	for _, scanTarget := range targets {
		options := map[string]interface{}{
			"ports":    "22,80,443,8080,8443",
			"threads":  50,
			"timeout":  2,
			"scan_tcp": true,
		}

		result, err := moduleManager.ExecuteModule("port_scanning", scanTarget, options)
		if err != nil {
			continue // Skip failed scans
		}

		results = append(results, result)

		// Display results
		if len(result.Results) > 0 {
			fmt.Printf("   %s: %d open ports\n", scanTarget, len(result.Results))
		}
	}

	return results
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// demonstrateModules shows available modules and their capabilities
func demonstrateModules() {
	fmt.Println("🧩 Available GoReconX Modules:")
	fmt.Println("")

	modules := []struct {
		name        string
		description string
		category    string
		features    []string
	}{
		{
			name:        "Subdomain Enumeration",
			description: "Advanced DNS-based subdomain discovery",
			category:    "Passive OSINT",
			features:    []string{"Wordlist-based discovery", "DNS resolution", "Concurrent scanning", "IP resolution"},
		},
		{
			name:        "Port Scanner",
			description: "Fast TCP/UDP port scanning with service detection",
			category:    "Active Reconnaissance",
			features:    []string{"TCP/UDP scanning", "Service detection", "Banner grabbing", "Custom port ranges"},
		},
		{
			name:        "Email Harvester",
			description: "Collect email addresses from various sources",
			category:    "Passive OSINT",
			features:    []string{"Search engine scraping", "Social media lookup", "WHOIS extraction", "DNS enumeration"},
		},
		{
			name:        "Web Analyzer",
			description: "Analyze web technologies and content",
			category:    "Passive OSINT",
			features:    []string{"Technology detection", "Header analysis", "Content extraction", "Vulnerability hints"},
		},
		{
			name:        "Directory Enumerator",
			description: "Discover hidden directories and files",
			category:    "Active Reconnaissance",
			features:    []string{"Wordlist-based discovery", "Response analysis", "Recursive scanning", "Custom headers"},
		},
		{
			name:        "IP Geolocator",
			description: "Determine geographical location and ASN info",
			category:    "Passive OSINT",
			features:    []string{"GeoIP lookup", "ASN information", "ISP details", "Country/region data"},
		},
		{
			name:        "GitHub Reconnaissance",
			description: "Search for sensitive information in repositories",
			category:    "Passive OSINT",
			features:    []string{"Repository scanning", "Sensitive file detection", "Commit history analysis", "API integration"},
		},
	}

	for _, module := range modules {
		fmt.Printf("📦 %s (%s)\n", module.name, module.category)
		fmt.Printf("   %s\n", module.description)
		fmt.Printf("   Features: %s\n", fmt.Sprintf("%v", module.features))
		fmt.Println()
	}
}

// ExampleUsage demonstrates how to use GoReconX programmatically
func ExampleUsage() {
	// This function shows how developers can integrate GoReconX into their own tools

	fmt.Println("🔧 Programmatic Usage Example:")
	fmt.Println("")

	example := `
package main

import (
    "GoReconX/internal/modules"
    "GoReconX/internal/config"
    "GoReconX/internal/logging"
    "GoReconX/internal/database"
)

func main() {
    // Initialize components
    logger := logging.InitLogger()
    cfg, _ := config.LoadConfig()
    db, _ := database.InitDB()
    defer db.Close()
    
    // Create module manager
    moduleManager := modules.NewModuleManager(db, cfg, logger)
    defer moduleManager.Close()
    
    // Configure scan options
    options := map[string]interface{}{
        "threads": 50,
        "timeout": 5,
        "resolve_ips": true,
    }
    
    // Execute subdomain enumeration
    result, err := moduleManager.ExecuteModule(
        "subdomain_enumeration", 
        "example.com", 
        options,
    )
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    fmt.Printf("Found %d subdomains\n", len(result.Results))
    
    // Generate report with AI analysis
    reportGen := reports.NewReportGenerator(logger, aiClient, "output/")
    report, _ := reportGen.GenerateReport("example.com", []*modules.ScanResult{result})
    
    // Export in multiple formats
    reportGen.ExportJSON(report)
    reportGen.ExportHTML(report)
    reportGen.ExportCSV(report)
}
`

	fmt.Println(example)
}
