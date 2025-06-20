package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// GeminiClient represents the Google Gemini AI client
type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
	logger     *logrus.Logger
	baseURL    string
}

// GeminiRequest represents a request to the Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents content in the request
type Content struct {
	Parts []Part `json:"parts"`
}

// Part represents a part of the content
type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

// Candidate represents a candidate response
type Candidate struct {
	Content struct {
		Parts []Part `json:"parts"`
	} `json:"content"`
}

// AnalysisRequest represents a request for AI analysis
type AnalysisRequest struct {
	Type        string                 `json:"type"`        // "summary", "recommendations", "threat_analysis", "report"
	Data        interface{}            `json:"data"`        // The data to analyze
	Context     string                 `json:"context"`     // Additional context
	Target      string                 `json:"target"`      // Target being analyzed
	Metadata    map[string]interface{} `json:"metadata"`    // Additional metadata
}

// AnalysisResponse represents the response from AI analysis
type AnalysisResponse struct {
	Type            string                 `json:"type"`
	Summary         string                 `json:"summary"`
	Insights        []string               `json:"insights"`
	Recommendations []string               `json:"recommendations"`
	ThreatLevel     string                 `json:"threat_level"`
	Confidence      float64                `json:"confidence"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NewGeminiClient creates a new Gemini AI client
func NewGeminiClient(apiKey string, logger *logrus.Logger) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("google Gemini API key is required")
	}

	return &GeminiClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:  logger,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent",
	}, nil
}

// AnalyzeResults performs AI analysis on reconnaissance results
func (gc *GeminiClient) AnalyzeResults(req *AnalysisRequest) (*AnalysisResponse, error) {
	ctx := context.Background()
	
	prompt := gc.buildPrompt(req)
	
	gc.logger.WithFields(logrus.Fields{
		"type":   req.Type,
		"target": req.Target,
	}).Info("Performing AI analysis")

	geminiReq := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", gc.baseURL+"?key="+gc.apiKey, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := gc.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned from Gemini")
	}

	content := ""
	for _, part := range geminiResp.Candidates[0].Content.Parts {
		content += part.Text
	}

	return gc.parseResponse(content, req.Type), nil
}

// buildPrompt constructs the prompt for the AI model
func (gc *GeminiClient) buildPrompt(req *AnalysisRequest) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are a cybersecurity expert analyzing reconnaissance data. ")
	promptBuilder.WriteString("Provide professional, actionable insights based on the following data.\n\n")

	switch req.Type {
	case "summary":
		promptBuilder.WriteString("TASK: Provide a comprehensive summary of the reconnaissance findings.\n")
	case "recommendations":
		promptBuilder.WriteString("TASK: Provide security recommendations based on the findings.\n")
	case "threat_analysis":
		promptBuilder.WriteString("TASK: Analyze potential security threats and vulnerabilities.\n")
	case "report":
		promptBuilder.WriteString("TASK: Generate an executive summary for a security report.\n")
	default:
		promptBuilder.WriteString("TASK: Analyze the reconnaissance data and provide insights.\n")
	}

	promptBuilder.WriteString(fmt.Sprintf("TARGET: %s\n", req.Target))
	
	if req.Context != "" {
		promptBuilder.WriteString(fmt.Sprintf("CONTEXT: %s\n", req.Context))
	}

	promptBuilder.WriteString("DATA:\n")
	
	// Convert data to JSON string for analysis
	dataJSON, err := json.MarshalIndent(req.Data, "", "  ")
	if err != nil {
		promptBuilder.WriteString(fmt.Sprintf("%v", req.Data))
	} else {
		promptBuilder.WriteString(string(dataJSON))
	}

	promptBuilder.WriteString("\n\nPlease provide your analysis in the following format:\n")
	promptBuilder.WriteString("SUMMARY: [Brief overview]\n")
	promptBuilder.WriteString("KEY INSIGHTS: [Bullet points of key findings]\n")
	promptBuilder.WriteString("RECOMMENDATIONS: [Security recommendations]\n")
	promptBuilder.WriteString("THREAT LEVEL: [LOW/MEDIUM/HIGH/CRITICAL]\n")
	promptBuilder.WriteString("CONFIDENCE: [0.0-1.0]\n")

	return promptBuilder.String()
}

// parseResponse parses the AI response into structured format
func (gc *GeminiClient) parseResponse(content, analysisType string) *AnalysisResponse {
	response := &AnalysisResponse{
		Type:            analysisType,
		Insights:        []string{},
		Recommendations: []string{},
		ThreatLevel:     "UNKNOWN",
		Confidence:      0.5,
		Metadata:        make(map[string]interface{}),
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(strings.ToUpper(line), "SUMMARY:") {
			currentSection = "summary"
			response.Summary = strings.TrimSpace(strings.TrimPrefix(line, "SUMMARY:"))
			continue
		} else if strings.HasPrefix(strings.ToUpper(line), "KEY INSIGHTS:") {
			currentSection = "insights"
			continue
		} else if strings.HasPrefix(strings.ToUpper(line), "RECOMMENDATIONS:") {
			currentSection = "recommendations"
			continue
		} else if strings.HasPrefix(strings.ToUpper(line), "THREAT LEVEL:") {
			response.ThreatLevel = strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(line), "THREAT LEVEL:"))
			currentSection = ""
			continue
		} else if strings.HasPrefix(strings.ToUpper(line), "CONFIDENCE:") {
			confidenceStr := strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(line), "CONFIDENCE:"))
			if conf := parseConfidence(confidenceStr); conf > 0 {
				response.Confidence = conf
			}
			currentSection = ""
			continue
		}

		// Process content based on current section
		switch currentSection {
		case "summary":
			if response.Summary == "" {
				response.Summary = line
			} else {
				response.Summary += " " + line
			}
		case "insights":
			if strings.HasPrefix(line, "•") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				response.Insights = append(response.Insights, strings.TrimSpace(line[1:]))
			} else if line != "" {
				response.Insights = append(response.Insights, line)
			}
		case "recommendations":
			if strings.HasPrefix(line, "•") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				response.Recommendations = append(response.Recommendations, strings.TrimSpace(line[1:]))
			} else if line != "" {
				response.Recommendations = append(response.Recommendations, line)
			}
		}
	}

	// If we couldn't parse specific sections, use the entire content as summary
	if response.Summary == "" && len(response.Insights) == 0 && len(response.Recommendations) == 0 {
		response.Summary = content
	}

	return response
}

// parseConfidence attempts to parse confidence value from string
func parseConfidence(confidenceStr string) float64 {
	// Remove any non-numeric characters except decimal point
	cleaned := ""
	for _, char := range confidenceStr {
		if char >= '0' && char <= '9' || char == '.' {
			cleaned += string(char)
		}
	}
	
	if cleaned == "" {
		return 0.5 // Default confidence
	}
	
	var conf float64
	if _, err := fmt.Sscanf(cleaned, "%f", &conf); err == nil {
		if conf > 1.0 {
			conf = conf / 100.0 // Convert percentage to decimal
		}
		if conf >= 0.0 && conf <= 1.0 {
			return conf
		}
	}
	
	return 0.5 // Default confidence if parsing fails
}

// GenerateReport generates a comprehensive security report
func (gc *GeminiClient) GenerateReport(results []interface{}, target string) (*AnalysisResponse, error) {
	req := &AnalysisRequest{
		Type:    "report",
		Data:    results,
		Target:  target,
		Context: "Generate a comprehensive security assessment report",
	}
	
	return gc.AnalyzeResults(req)
}

// IsConfigured checks if the client is properly configured
func (gc *GeminiClient) IsConfigured() bool {
	return gc.apiKey != ""
}

// Close closes the Gemini client (placeholder for interface compatibility)
func (gc *GeminiClient) Close() error {
	return nil
}
