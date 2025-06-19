package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// EmailEnumModule implements email enumeration and people search
type EmailEnumModule struct {
	*BaseModule
	client *http.Client
}

// EmailEnumResult represents email enumeration results
type EmailEnumResult struct {
	Domain      string       `json:"domain"`
	Emails      []EmailInfo  `json:"emails"`
	SocialMedia []SocialInfo `json:"social_media"`
	People      []PersonInfo `json:"people"`
	Sources     []string     `json:"sources"`
	TotalFound  int          `json:"total_found"`
	ScanTime    string       `json:"scan_time"`
}

// EmailInfo contains email information
type EmailInfo struct {
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	Position   string   `json:"position"`
	Department string   `json:"department"`
	Sources    []string `json:"sources"`
	Confidence int      `json:"confidence"`
	LastSeen   string   `json:"last_seen"`
}

// SocialInfo contains social media profile information
type SocialInfo struct {
	Platform  string `json:"platform"`
	Username  string `json:"username"`
	URL       string `json:"url"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Followers int    `json:"followers"`
	Verified  bool   `json:"verified"`
}

// PersonInfo contains person information
type PersonInfo struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Position string   `json:"position"`
	Company  string   `json:"company"`
	Location string   `json:"location"`
	LinkedIn string   `json:"linkedin"`
	Twitter  string   `json:"twitter"`
	Sources  []string `json:"sources"`
}

// NewEmailEnumModule creates a new email enumeration module
func NewEmailEnumModule() *EmailEnumModule {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	info := ModuleInfo{
		Name:        "email_enum",
		Category:    "passive_osint",
		Description: "Email enumeration and people search using various OSINT techniques",
		Version:     "1.0.0",
		Author:      "GoReconX Team",
		Tags:        []string{"email", "people", "osint", "social", "passive"},
		Options: []ModuleOption{
			{
				Name:        "search_engines",
				Type:        "bool",
				Description: "Use search engines for email discovery",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "social_media",
				Type:        "bool",
				Description: "Search social media platforms",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "use_hunter_io",
				Type:        "bool",
				Description: "Use Hunter.io API (requires API key)",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "use_clearbit",
				Type:        "bool",
				Description: "Use Clearbit API (requires API key)",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "deep_search",
				Type:        "bool",
				Description: "Perform deep search (slower but more thorough)",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "max_results",
				Type:        "int",
				Description: "Maximum number of results to return",
				Required:    false,
				Default:     100,
			},
		},
		Requirements: []string{"network"},
	}

	return &EmailEnumModule{
		BaseModule: NewBaseModule(info),
		client:     client,
	}
}

// Validate validates the module input
func (e *EmailEnumModule) Validate(input ModuleInput) error {
	if err := e.ValidateInput(input); err != nil {
		return err
	}

	// Validate domain format
	domain := strings.ToLower(strings.TrimSpace(input.Target))
	if !isValidDomain(domain) {
		return NewModuleError("invalid domain format", "INVALID_DOMAIN")
	}

	return nil
}

// Execute runs the email enumeration module
func (e *EmailEnumModule) Execute(ctx context.Context, input ModuleInput, output chan<- ModuleResult) error {
	startTime := time.Now()
	e.SetStatus("running", 0.0, "Starting email enumeration")

	domain := strings.ToLower(strings.TrimSpace(input.Target))
	// Parse options
	useSearchEngines, _ := input.Options["search_engines"].(bool)
	useSocialMedia, _ := input.Options["social_media"].(bool)
	useHunterIO, _ := input.Options["use_hunter_io"].(bool)
	_ = input.Options["use_clearbit"].(bool) // useClearbit for future use
	deepSearch, _ := input.Options["deep_search"].(bool)
	maxResults, _ := input.Options["max_results"].(int)
	if maxResults <= 0 {
		maxResults = 100
	}

	result := &EmailEnumResult{
		Domain:      domain,
		Emails:      []EmailInfo{},
		SocialMedia: []SocialInfo{},
		People:      []PersonInfo{},
		Sources:     []string{},
	}

	emailMap := make(map[string]*EmailInfo)
	peopleMap := make(map[string]*PersonInfo)

	// Phase 1: Search engines
	if useSearchEngines {
		e.SetStatus("running", 0.1, "Searching with search engines")
		e.SendResult(output, "progress", "Searching with search engines", nil, input.SessionID)

		searchEmails := e.searchEngineEmails(domain)
		for _, email := range searchEmails {
			if existing, exists := emailMap[email.Email]; exists {
				existing.Sources = append(existing.Sources, email.Sources...)
				existing.Confidence = max(existing.Confidence, email.Confidence)
			} else {
				emailMap[email.Email] = &email
			}
		}

		if len(searchEmails) > 0 {
			result.Sources = append(result.Sources, "Search Engines")
		}
	}

	if e.IsStopped() {
		return nil
	}

	// Phase 2: Hunter.io API
	if useHunterIO {
		e.SetStatus("running", 0.3, "Querying Hunter.io API")
		e.SendResult(output, "progress", "Querying Hunter.io API", nil, input.SessionID)

		hunterEmails := e.hunterIOSearch(domain, input.Options)
		for _, email := range hunterEmails {
			if existing, exists := emailMap[email.Email]; exists {
				existing.Sources = append(existing.Sources, email.Sources...)
				existing.Confidence = max(existing.Confidence, email.Confidence)
			} else {
				emailMap[email.Email] = &email
			}
		}

		if len(hunterEmails) > 0 {
			result.Sources = append(result.Sources, "Hunter.io")
		}
	}

	if e.IsStopped() {
		return nil
	}

	// Phase 3: Social media search
	if useSocialMedia {
		e.SetStatus("running", 0.5, "Searching social media platforms")
		e.SendResult(output, "progress", "Searching social media platforms", nil, input.SessionID)

		socialProfiles := e.searchSocialMedia(domain)
		result.SocialMedia = socialProfiles

		if len(socialProfiles) > 0 {
			result.Sources = append(result.Sources, "Social Media")
		}
	}

	if e.IsStopped() {
		return nil
	}

	// Phase 4: Website crawling for emails
	e.SetStatus("running", 0.7, "Crawling website for emails")
	e.SendResult(output, "progress", "Crawling website for emails", nil, input.SessionID)

	websiteEmails := e.crawlWebsiteEmails(domain)
	for _, email := range websiteEmails {
		if existing, exists := emailMap[email.Email]; exists {
			existing.Sources = append(existing.Sources, email.Sources...)
			existing.Confidence = max(existing.Confidence, email.Confidence)
		} else {
			emailMap[email.Email] = &email
		}
	}

	if len(websiteEmails) > 0 {
		result.Sources = append(result.Sources, "Website Crawling")
	}

	// Phase 5: Deep search (if enabled)
	if deepSearch {
		e.SetStatus("running", 0.8, "Performing deep search")
		e.SendResult(output, "progress", "Performing deep search", nil, input.SessionID)

		deepEmails := e.deepEmailSearch(domain)
		for _, email := range deepEmails {
			if existing, exists := emailMap[email.Email]; exists {
				existing.Sources = append(existing.Sources, email.Sources...)
				existing.Confidence = max(existing.Confidence, email.Confidence)
			} else {
				emailMap[email.Email] = &email
			}
		}

		if len(deepEmails) > 0 {
			result.Sources = append(result.Sources, "Deep Search")
		}
	}

	// Convert maps to slices and send individual results
	for _, email := range emailMap {
		if len(result.Emails) < maxResults {
			result.Emails = append(result.Emails, *email)
			e.SendResult(output, "data", map[string]interface{}{
				"type": "email",
				"data": *email,
			}, nil, input.SessionID)
		}
	}

	for _, person := range peopleMap {
		result.People = append(result.People, *person)
		e.SendResult(output, "data", map[string]interface{}{
			"type": "person",
			"data": *person,
		}, nil, input.SessionID)
	}

	result.TotalFound = len(result.Emails)
	result.ScanTime = time.Since(startTime).String()

	// Send final result
	e.SetStatus("completed", 1.0, fmt.Sprintf("Email enumeration completed: %d emails found", len(result.Emails)))
	e.SendResult(output, "complete", result, map[string]interface{}{
		"total_emails": len(result.Emails),
		"total_people": len(result.People),
		"sources":      len(result.Sources),
	}, input.SessionID)

	return nil
}

// searchEngineEmails searches for emails using search engines
func (e *EmailEnumModule) searchEngineEmails(domain string) []EmailInfo {
	var emails []EmailInfo

	// Google dork for emails
	queries := []string{
		fmt.Sprintf("site:%s \"@%s\"", domain, domain),
		fmt.Sprintf("\"%s\" email contact", domain),
		fmt.Sprintf("\"@%s\" -site:%s", domain, domain),
	}

	for _, query := range queries {
		if e.IsStopped() {
			break
		}

		searchResults := e.performGoogleSearch(query)
		foundEmails := e.extractEmailsFromText(searchResults, domain)

		for _, email := range foundEmails {
			emails = append(emails, EmailInfo{
				Email:      email,
				Sources:    []string{"Google Search"},
				Confidence: 70,
				LastSeen:   time.Now().Format("2006-01-02"),
			})
		}

		// Rate limiting
		time.Sleep(2 * time.Second)
	}

	return emails
}

// hunterIOSearch searches using Hunter.io API
func (e *EmailEnumModule) hunterIOSearch(domain string, options map[string]interface{}) []EmailInfo {
	var emails []EmailInfo

	// This would require a real Hunter.io API key
	// For demonstration, we'll return a placeholder

	apiKey, exists := options["hunter_io_api_key"].(string)
	if !exists || apiKey == "" {
		return emails
	}

	url := fmt.Sprintf("https://api.hunter.io/v2/domain-search?domain=%s&api_key=%s", domain, apiKey)

	resp, err := e.client.Get(url)
	if err != nil {
		return emails
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Emails []struct {
				Value      string `json:"value"`
				Type       string `json:"type"`
				Confidence int    `json:"confidence"`
				Sources    []struct {
					Domain string `json:"domain"`
					URI    string `json:"uri"`
				} `json:"sources"`
				FirstName  string `json:"first_name"`
				LastName   string `json:"last_name"`
				Position   string `json:"position"`
				Department string `json:"department"`
			} `json:"emails"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return emails
	}

	for _, hunterEmail := range result.Data.Emails {
		var sources []string
		for _, source := range hunterEmail.Sources {
			sources = append(sources, source.Domain)
		}

		name := strings.TrimSpace(hunterEmail.FirstName + " " + hunterEmail.LastName)

		emails = append(emails, EmailInfo{
			Email:      hunterEmail.Value,
			Name:       name,
			Position:   hunterEmail.Position,
			Department: hunterEmail.Department,
			Sources:    append(sources, "Hunter.io"),
			Confidence: hunterEmail.Confidence,
			LastSeen:   time.Now().Format("2006-01-02"),
		})
	}

	return emails
}

// searchSocialMedia searches social media platforms
func (e *EmailEnumModule) searchSocialMedia(domain string) []SocialInfo {
	var profiles []SocialInfo

	// Search for company profiles on major platforms
	platforms := []struct {
		name string
		urls []string
	}{
		{
			name: "LinkedIn",
			urls: []string{
				fmt.Sprintf("https://www.linkedin.com/company/%s", domain),
				fmt.Sprintf("https://www.linkedin.com/search/results/companies/?keywords=%s", domain),
			},
		},
		{
			name: "Twitter",
			urls: []string{
				fmt.Sprintf("https://twitter.com/%s", domain),
				fmt.Sprintf("https://twitter.com/search?q=%s", domain),
			},
		},
		{
			name: "Facebook",
			urls: []string{
				fmt.Sprintf("https://www.facebook.com/%s", domain),
			},
		},
	}

	for _, platform := range platforms {
		if e.IsStopped() {
			break
		}

		for _, url := range platform.urls {
			if profile := e.checkSocialProfile(platform.name, url); profile != nil {
				profiles = append(profiles, *profile)
			}
		}
	}

	return profiles
}

// crawlWebsiteEmails crawls the target website for emails
func (e *EmailEnumModule) crawlWebsiteEmails(domain string) []EmailInfo {
	var emails []EmailInfo

	// Common pages that might contain emails
	pages := []string{
		"/",
		"/contact",
		"/about",
		"/team",
		"/staff",
		"/people",
		"/directory",
		"/support",
		"/help",
	}

	baseURLs := []string{
		fmt.Sprintf("https://%s", domain),
		fmt.Sprintf("http://%s", domain),
	}

	for _, baseURL := range baseURLs {
		if e.IsStopped() {
			break
		}

		for _, page := range pages {
			if e.IsStopped() {
				break
			}

			fullURL := baseURL + page
			if content := e.fetchWebContent(fullURL); content != "" {
				foundEmails := e.extractEmailsFromText(content, domain)

				for _, email := range foundEmails {
					emails = append(emails, EmailInfo{
						Email:      email,
						Sources:    []string{fmt.Sprintf("Website: %s", page)},
						Confidence: 85,
						LastSeen:   time.Now().Format("2006-01-02"),
					})
				}
			}
		}

		// If HTTPS worked, skip HTTP
		if len(emails) > 0 {
			break
		}
	}

	return emails
}

// deepEmailSearch performs deep email search using various techniques
func (e *EmailEnumModule) deepEmailSearch(domain string) []EmailInfo {
	var emails []EmailInfo

	// Search in code repositories (simplified)
	repoEmails := e.searchCodeRepositories(domain)
	emails = append(emails, repoEmails...)

	// Search in data breaches (passive, public sources only)
	breachEmails := e.searchDataBreaches(domain)
	emails = append(emails, breachEmails...)

	// Certificate transparency logs for email addresses
	certEmails := e.searchCertificateLogs(domain)
	emails = append(emails, certEmails...)

	return emails
}

// Helper functions

func (e *EmailEnumModule) performGoogleSearch(query string) string {
	// Simplified Google search - in real implementation, you'd use proper APIs
	url := fmt.Sprintf("https://www.google.com/search?q=%s", strings.ReplaceAll(query, " ", "+"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := e.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Read response - this would need proper HTML parsing in real implementation
	body := make([]byte, 10000)
	n, _ := resp.Body.Read(body)
	return string(body[:n])
}

func (e *EmailEnumModule) extractEmailsFromText(text, domain string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(text, -1)

	var emails []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if strings.HasSuffix(strings.ToLower(match), strings.ToLower(domain)) && !seen[match] {
			emails = append(emails, strings.ToLower(match))
			seen[match] = true
		}
	}

	return emails
}

func (e *EmailEnumModule) checkSocialProfile(platform, url string) *SocialInfo {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return &SocialInfo{
			Platform: platform,
			URL:      url,
			Name:     "", // Would extract from page content
		}
	}

	return nil
}

func (e *EmailEnumModule) fetchWebContent(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "GoReconX/1.0 (OSINT Scanner)")

	resp, err := e.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	body := make([]byte, 50000)
	n, _ := resp.Body.Read(body)
	return string(body[:n])
}

func (e *EmailEnumModule) searchCodeRepositories(domain string) []EmailInfo {
	// Simplified implementation - would search GitHub, GitLab, etc.
	return []EmailInfo{}
}

func (e *EmailEnumModule) searchDataBreaches(domain string) []EmailInfo {
	// Simplified implementation - would query HaveIBeenPwned API, etc.
	return []EmailInfo{}
}

func (e *EmailEnumModule) searchCertificateLogs(domain string) []EmailInfo {
	// Simplified implementation - would search CT logs for email addresses
	return []EmailInfo{}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
