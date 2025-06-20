package reports

import (
	"GoReconX/internal/ai"
	"GoReconX/internal/modules"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ReportGenerator handles report generation in various formats
type ReportGenerator struct {
	logger    *logrus.Logger
	aiClient  *ai.GeminiClient
	outputDir string
}

// Report represents a comprehensive reconnaissance report
type Report struct {
	ID          string                    `json:"id"`
	Target      string                    `json:"target"`
	Title       string                    `json:"title"`
	GeneratedAt time.Time                 `json:"generated_at"`
	Summary     string                    `json:"summary"`
	Results     []*modules.ScanResult     `json:"results"`
	AIAnalysis  *ai.AnalysisResponse      `json:"ai_analysis,omitempty"`
	Statistics  map[string]interface{}    `json:"statistics"`
	Metadata    map[string]interface{}    `json:"metadata"`
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(logger *logrus.Logger, aiClient *ai.GeminiClient, outputDir string) *ReportGenerator {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logger.WithError(err).Warn("Failed to create output directory")
	}

	return &ReportGenerator{
		logger:    logger,
		aiClient:  aiClient,
		outputDir: outputDir,
	}
}

// GenerateReport creates a comprehensive report from scan results
func (rg *ReportGenerator) GenerateReport(target string, results []*modules.ScanResult) (*Report, error) {
	reportID := fmt.Sprintf("report_%s_%d", strings.ReplaceAll(target, ".", "_"), time.Now().Unix())
	
	report := &Report{
		ID:          reportID,
		Target:      target,
		Title:       fmt.Sprintf("GoReconX Security Assessment - %s", target),
		GeneratedAt: time.Now(),
		Results:     results,
		Statistics:  rg.calculateStatistics(results),
		Metadata:    make(map[string]interface{}),
	}

	// Generate AI analysis if available
	if rg.aiClient != nil && rg.aiClient.IsConfigured() {
		rg.logger.Info("Generating AI analysis for report")
		
		var allResults []interface{}
		for _, result := range results {
			allResults = append(allResults, result)
		}
		
		aiAnalysis, err := rg.aiClient.GenerateReport(allResults, target)
		if err != nil {
			rg.logger.WithError(err).Warn("Failed to generate AI analysis")
		} else {
			report.AIAnalysis = aiAnalysis
			report.Summary = aiAnalysis.Summary
		}
	}

	// Generate fallback summary if AI analysis failed
	if report.Summary == "" {
		report.Summary = rg.generateBasicSummary(results)
	}

	return report, nil
}

// ExportJSON exports the report in JSON format
func (rg *ReportGenerator) ExportJSON(report *Report) (string, error) {
	filename := filepath.Join(rg.outputDir, report.ID+".json")
	
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal report: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write JSON report: %v", err)
	}

	rg.logger.WithField("file", filename).Info("JSON report exported")
	return filename, nil
}

// ExportHTML exports the report in HTML format
func (rg *ReportGenerator) ExportHTML(report *Report) (string, error) {
	filename := filepath.Join(rg.outputDir, report.ID+".html")
	
	tmpl := rg.getHTMLTemplate()
	
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create HTML file: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, report); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	rg.logger.WithField("file", filename).Info("HTML report exported")
	return filename, nil
}

// ExportCSV exports scan results in CSV format
func (rg *ReportGenerator) ExportCSV(report *Report) (string, error) {
	filename := filepath.Join(rg.outputDir, report.ID+".csv")
	
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Write CSV header
	file.WriteString("Module,Target,Status,Start Time,End Time,Results Count\n")

	// Write scan results
	for _, result := range report.Results {
		file.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%d\n",
			result.ModuleName,
			result.Target,
			result.Status,
			result.StartTime,
			result.EndTime,
			len(result.Results),
		))
	}

	rg.logger.WithField("file", filename).Info("CSV report exported")
	return filename, nil
}

// calculateStatistics generates statistics from scan results
func (rg *ReportGenerator) calculateStatistics(results []*modules.ScanResult) map[string]interface{} {
	stats := make(map[string]interface{})
	
	totalScans := len(results)
	completedScans := 0
	failedScans := 0
	totalResults := 0
	
	moduleStats := make(map[string]int)
	
	for _, result := range results {
		switch result.Status {
		case "completed":
			completedScans++
		case "failed":
			failedScans++
		}
		
		totalResults += len(result.Results)
		moduleStats[result.ModuleName]++
	}
	
	stats["total_scans"] = totalScans
	stats["completed_scans"] = completedScans
	stats["failed_scans"] = failedScans
	stats["total_results"] = totalResults
	stats["module_usage"] = moduleStats
	
	if totalScans > 0 {
		stats["success_rate"] = float64(completedScans) / float64(totalScans) * 100
	}
	
	return stats
}

// generateBasicSummary creates a basic summary when AI analysis is not available
func (rg *ReportGenerator) generateBasicSummary(results []*modules.ScanResult) string {
	var summary strings.Builder
	
	summary.WriteString("GoReconX Security Assessment Summary\n\n")
	
	completedCount := 0
	totalFindings := 0
	
	for _, result := range results {
		if result.Status == "completed" {
			completedCount++
			totalFindings += len(result.Results)
		}
	}
	
	summary.WriteString(fmt.Sprintf("Completed %d reconnaissance modules with %d total findings.\n", 
		completedCount, totalFindings))
	
	// Module-specific summaries
	for _, result := range results {
		if result.Status == "completed" && len(result.Results) > 0 {
			summary.WriteString(fmt.Sprintf("- %s: Found %d items\n", 
				result.ModuleName, len(result.Results)))
		}
	}
	
	return summary.String()
}

// getHTMLTemplate returns the HTML template for reports
func (rg *ReportGenerator) getHTMLTemplate() *template.Template {
	tmplContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px; background: #f4f4f4; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
        .header { text-align: center; border-bottom: 2px solid #333; padding-bottom: 20px; margin-bottom: 30px; }
        .header h1 { color: #333; margin-bottom: 10px; }
        .meta-info { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
        .summary { background: #e9ecef; padding: 20px; border-radius: 5px; margin-bottom: 30px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin-bottom: 30px; }
        .stat-card { background: #007bff; color: white; padding: 15px; border-radius: 5px; text-align: center; }
        .stat-card h3 { margin: 0 0 10px 0; font-size: 2em; }
        .stat-card p { margin: 0; opacity: 0.9; }
        .results { margin-bottom: 30px; }
        .result-card { background: #f8f9fa; border: 1px solid #dee2e6; border-radius: 5px; margin-bottom: 15px; overflow: hidden; }
        .result-header { background: #343a40; color: white; padding: 10px 15px; font-weight: bold; }
        .result-body { padding: 15px; }
        .status-completed { color: #28a745; }
        .status-failed { color: #dc3545; }
        .ai-analysis { background: #d4edda; border: 1px solid #c3e6cb; border-radius: 5px; padding: 20px; margin-bottom: 30px; }
        .recommendations { background: #fff3cd; border: 1px solid #ffeaa7; border-radius: 5px; padding: 20px; }
        .footer { text-align: center; color: #6c757d; margin-top: 30px; padding-top: 20px; border-top: 1px solid #dee2e6; }
        pre { background: #f1f3f4; padding: 10px; border-radius: 3px; overflow-x: auto; }
        .threat-level { padding: 4px 8px; border-radius: 3px; font-weight: bold; }
        .threat-low { background: #d4edda; color: #155724; }
        .threat-medium { background: #fff3cd; color: #856404; }
        .threat-high { background: #f8d7da; color: #721c24; }
        .threat-critical { background: #f5c6cb; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p>Target: <strong>{{.Target}}</strong></p>
            <p>Generated: {{.GeneratedAt.Format "January 2, 2006 15:04:05 MST"}}</p>
        </div>

        <div class="meta-info">
            <h2>Report Information</h2>
            <p><strong>Report ID:</strong> {{.ID}}</p>
            <p><strong>Target:</strong> {{.Target}}</p>
            <p><strong>Generated:</strong> {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</p>
        </div>

        {{if .Summary}}
        <div class="summary">
            <h2>Executive Summary</h2>
            <pre>{{.Summary}}</pre>
        </div>
        {{end}}

        <div class="stats">
            <div class="stat-card">
                <h3>{{index .Statistics "total_scans"}}</h3>
                <p>Total Scans</p>
            </div>
            <div class="stat-card">
                <h3>{{index .Statistics "completed_scans"}}</h3>
                <p>Completed</p>
            </div>
            <div class="stat-card">
                <h3>{{index .Statistics "total_results"}}</h3>
                <p>Total Findings</p>
            </div>
            <div class="stat-card">
                <h3>{{printf "%.1f%%" (index .Statistics "success_rate")}}</h3>
                <p>Success Rate</p>
            </div>
        </div>

        {{if .AIAnalysis}}
        <div class="ai-analysis">
            <h2>AI-Powered Analysis</h2>
            {{if .AIAnalysis.ThreatLevel}}
            <p><strong>Threat Level:</strong> 
                <span class="threat-level threat-{{.AIAnalysis.ThreatLevel | lower}}">{{.AIAnalysis.ThreatLevel}}</span>
            </p>
            {{end}}
            {{if .AIAnalysis.Confidence}}
            <p><strong>Confidence:</strong> {{printf "%.1f%%" (.AIAnalysis.Confidence | multiply 100)}}</p>
            {{end}}
            
            {{if .AIAnalysis.Insights}}
            <h3>Key Insights</h3>
            <ul>
                {{range .AIAnalysis.Insights}}
                <li>{{.}}</li>
                {{end}}
            </ul>
            {{end}}
            
            {{if .AIAnalysis.Recommendations}}
            <div class="recommendations">
                <h3>Recommendations</h3>
                <ul>
                    {{range .AIAnalysis.Recommendations}}
                    <li>{{.}}</li>
                    {{end}}
                </ul>
            </div>
            {{end}}
        </div>
        {{end}}

        <div class="results">
            <h2>Detailed Results</h2>
            {{range .Results}}
            <div class="result-card">
                <div class="result-header">
                    {{.ModuleName}} - <span class="status-{{.Status}}">{{.Status | title}}</span>
                </div>
                <div class="result-body">
                    <p><strong>Target:</strong> {{.Target}}</p>
                    <p><strong>Start Time:</strong> {{.StartTime}}</p>
                    <p><strong>End Time:</strong> {{.EndTime}}</p>
                    <p><strong>Results Count:</strong> {{len .Results}}</p>
                    {{if .ErrorMessage}}
                    <p><strong>Error:</strong> <span style="color: red;">{{.ErrorMessage}}</span></p>
                    {{end}}
                    {{if .Results}}
                    <details>
                        <summary>View Results ({{len .Results}} items)</summary>
                        <pre>{{.Results | marshal}}</pre>
                    </details>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>

        <div class="footer">
            <p>Generated by GoReconX - Comprehensive OSINT & Reconnaissance Platform</p>
            <p>Remember to use this tool ethically and with proper authorization</p>
        </div>
    </div>
</body>
</html>
`

	tmpl := template.New("report")
	tmpl = tmpl.Funcs(template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"marshal": func(v interface{}) string {
			data, _ := json.MarshalIndent(v, "", "  ")
			return string(data)
		},
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	})
	
	template.Must(tmpl.Parse(tmplContent))
	return tmpl
}
