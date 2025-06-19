package modules

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WebEnumModule implements web directory and file enumeration
type WebEnumModule struct {
	*BaseModule
	client    *http.Client
	semaphore chan bool
}

// WebEnumResult represents web enumeration results
type WebEnumResult struct {
	Target          string     `json:"target"`
	BaseURL         string     `json:"base_url"`
	FoundPaths      []PathInfo `json:"found_paths"`
	TotalTested     int        `json:"total_tested"`
	ScanTime        string     `json:"scan_time"`
	TechStack       []string   `json:"tech_stack"`
	Vulnerabilities []VulnInfo `json:"vulnerabilities"`
}

// PathInfo contains information about a discovered path
type PathInfo struct {
	Path         string            `json:"path"`
	StatusCode   int               `json:"status_code"`
	Size         int64             `json:"size"`
	ContentType  string            `json:"content_type"`
	Title        string            `json:"title"`
	ResponseTime string            `json:"response_time"`
	Headers      map[string]string `json:"headers"`
	IsDirectory  bool              `json:"is_directory"`
	Technology   []string          `json:"technology"`
}

// VulnInfo contains vulnerability information
type VulnInfo struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

// NewWebEnumModule creates a new web enumeration module
func NewWebEnumModule() *WebEnumModule {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects to avoid infinite loops
			return http.ErrUseLastResponse
		},
	}

	info := ModuleInfo{
		Name:        "web_enum",
		Category:    "active_recon",
		Description: "Web directory and file enumeration with technology detection",
		Version:     "1.0.0",
		Author:      "GoReconX Team",
		Tags:        []string{"web", "directory", "enumeration", "http", "active"},
		Options: []ModuleOption{
			{
				Name:        "wordlist",
				Type:        "choice",
				Description: "Wordlist to use for enumeration",
				Required:    false,
				Default:     "common",
				Choices:     []string{"common", "extensive", "quick"},
			},
			{
				Name:        "extensions",
				Type:        "string",
				Description: "File extensions to test (comma-separated)",
				Required:    false,
				Default:     "php,html,js,txt,xml,json",
			},
			{
				Name:        "threads",
				Type:        "int",
				Description: "Number of concurrent threads",
				Required:    false,
				Default:     20,
			},
			{
				Name:        "timeout",
				Type:        "int",
				Description: "HTTP request timeout in seconds",
				Required:    false,
				Default:     10,
			},
			{
				Name:        "user_agent",
				Type:        "string",
				Description: "User agent string to use",
				Required:    false,
				Default:     "GoReconX/1.0 (Security Scanner)",
			},
			{
				Name:        "recursive",
				Type:        "bool",
				Description: "Perform recursive directory scanning",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "status_codes",
				Type:        "string",
				Description: "Status codes to consider as found (comma-separated)",
				Required:    false,
				Default:     "200,201,204,301,302,307,401,403",
			},
		},
		Requirements: []string{"network"},
	}

	return &WebEnumModule{
		BaseModule: NewBaseModule(info),
		client:     client,
		semaphore:  make(chan bool, 20),
	}
}

// Validate validates the module input
func (w *WebEnumModule) Validate(input ModuleInput) error {
	if err := w.ValidateInput(input); err != nil {
		return err
	}

	// Validate URL format
	if _, err := url.Parse(input.Target); err != nil {
		return NewModuleError("invalid URL format", "INVALID_URL")
	}

	return nil
}

// Execute runs the web enumeration module
func (w *WebEnumModule) Execute(ctx context.Context, input ModuleInput, output chan<- ModuleResult) error {
	startTime := time.Now()
	w.SetStatus("running", 0.0, "Starting web enumeration")

	// Parse options
	threads, _ := input.Options["threads"].(int)
	if threads <= 0 {
		threads = 20
	}
	w.semaphore = make(chan bool, threads)

	timeout, _ := input.Options["timeout"].(int)
	if timeout <= 0 {
		timeout = 10
	}
	w.client.Timeout = time.Duration(timeout) * time.Second

	userAgent, _ := input.Options["user_agent"].(string)
	if userAgent == "" {
		userAgent = "GoReconX/1.0 (Security Scanner)"
	}

	wordlistType, _ := input.Options["wordlist"].(string)
	if wordlistType == "" {
		wordlistType = "common"
	}

	extensions, _ := input.Options["extensions"].(string)
	if extensions == "" {
		extensions = "php,html,js,txt,xml,json"
	}

	statusCodes, _ := input.Options["status_codes"].(string)
	if statusCodes == "" {
		statusCodes = "200,201,204,301,302,307,401,403"
	}

	// Parse target URL
	baseURL := input.Target
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "https://" + baseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return NewModuleError("invalid URL: "+err.Error(), "INVALID_URL")
	}

	result := &WebEnumResult{
		Target:          input.Target,
		BaseURL:         baseURL,
		FoundPaths:      []PathInfo{},
		TechStack:       []string{},
		Vulnerabilities: []VulnInfo{},
	}

	// Phase 1: Technology Detection
	w.SetStatus("running", 0.1, "Detecting web technologies")
	w.SendResult(output, "progress", "Detecting web technologies", nil, input.SessionID)

	techStack := w.detectTechnologies(baseURL, userAgent)
	result.TechStack = techStack

	w.SendResult(output, "data", map[string]interface{}{
		"type": "tech_stack",
		"data": techStack,
	}, nil, input.SessionID)

	// Phase 2: Generate wordlist
	w.SetStatus("running", 0.2, "Preparing wordlist")
	wordlist := w.getWordlist(wordlistType)
	extensionList := w.parseExtensions(extensions)

	// Generate full path list
	paths := w.generatePaths(wordlist, extensionList)
	result.TotalTested = len(paths)

	w.SendResult(output, "progress", fmt.Sprintf("Testing %d paths", len(paths)), nil, input.SessionID)

	// Phase 3: Enumerate paths
	w.SetStatus("running", 0.3, "Enumerating web paths")

	validStatusCodes := w.parseStatusCodes(statusCodes)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	testedCount := 0

	for _, path := range paths {
		if w.IsStopped() {
			break
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			defer func() {
				mutex.Lock()
				testedCount++
				progress := 0.3 + (0.6 * float64(testedCount) / float64(len(paths)))
				w.SetStatus("running", progress, fmt.Sprintf("Tested %d/%d paths", testedCount, len(paths)))
				mutex.Unlock()
			}()

			w.semaphore <- true
			defer func() { <-w.semaphore }()

			pathInfo := w.testPath(parsedURL, path, userAgent, validStatusCodes)
			if pathInfo != nil {
				mutex.Lock()
				result.FoundPaths = append(result.FoundPaths, *pathInfo)
				mutex.Unlock()

				w.SendResult(output, "data", map[string]interface{}{
					"type": "found_path",
					"path": *pathInfo,
				}, nil, input.SessionID)
			}
		}(path)
	}

	wg.Wait()

	// Phase 4: Vulnerability Analysis
	w.SetStatus("running", 0.9, "Analyzing for common vulnerabilities")
	vulns := w.analyzeVulnerabilities(result.FoundPaths, baseURL, userAgent)
	result.Vulnerabilities = vulns

	for _, vuln := range vulns {
		w.SendResult(output, "data", map[string]interface{}{
			"type": "vulnerability",
			"vuln": vuln,
		}, nil, input.SessionID)
	}

	// Sort results
	sort.Slice(result.FoundPaths, func(i, j int) bool {
		return result.FoundPaths[i].Path < result.FoundPaths[j].Path
	})

	result.ScanTime = time.Since(startTime).String()

	// Send final result
	w.SetStatus("completed", 1.0, fmt.Sprintf("Web enumeration completed: %d paths found", len(result.FoundPaths)))
	w.SendResult(output, "complete", result, map[string]interface{}{
		"found_paths":     len(result.FoundPaths),
		"vulnerabilities": len(result.Vulnerabilities),
		"scan_time":       result.ScanTime,
	}, input.SessionID)

	return nil
}

// detectTechnologies detects web technologies
func (w *WebEnumModule) detectTechnologies(baseURL, userAgent string) []string {
	var technologies []string

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return technologies
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := w.client.Do(req)
	if err != nil {
		return technologies
	}
	defer resp.Body.Close()

	// Server header
	if server := resp.Header.Get("Server"); server != "" {
		technologies = append(technologies, server)
	}

	// X-Powered-By header
	if powered := resp.Header.Get("X-Powered-By"); powered != "" {
		technologies = append(technologies, powered)
	}

	// Read response body for analysis
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		bodyStr := string(body)
		bodyLower := strings.ToLower(bodyStr)

		// Common CMS/Framework detection
		if strings.Contains(bodyLower, "wordpress") || strings.Contains(bodyLower, "wp-content") {
			technologies = append(technologies, "WordPress")
		}
		if strings.Contains(bodyLower, "drupal") {
			technologies = append(technologies, "Drupal")
		}
		if strings.Contains(bodyLower, "joomla") {
			technologies = append(technologies, "Joomla")
		}
		if strings.Contains(bodyLower, "laravel") {
			technologies = append(technologies, "Laravel")
		}
		if strings.Contains(bodyLower, "django") {
			technologies = append(technologies, "Django")
		}
		if strings.Contains(bodyLower, "react") {
			technologies = append(technologies, "React")
		}
		if strings.Contains(bodyLower, "angular") {
			technologies = append(technologies, "Angular")
		}
		if strings.Contains(bodyLower, "vue") {
			technologies = append(technologies, "Vue.js")
		}
	}

	return technologies
}

// testPath tests a single path
func (w *WebEnumModule) testPath(baseURL *url.URL, path, userAgent string, validStatusCodes map[int]bool) *PathInfo {
	fullURL := baseURL.ResolveReference(&url.URL{Path: path})

	start := time.Now()
	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Check if status code is valid
	if !validStatusCodes[resp.StatusCode] {
		return nil
	}

	responseTime := time.Since(start)

	pathInfo := &PathInfo{
		Path:         path,
		StatusCode:   resp.StatusCode,
		ContentType:  resp.Header.Get("Content-Type"),
		ResponseTime: responseTime.String(),
		Headers:      make(map[string]string),
		IsDirectory:  strings.HasSuffix(path, "/"),
	}

	// Store important headers
	importantHeaders := []string{"Server", "X-Powered-By", "Content-Type", "Content-Length"}
	for _, header := range importantHeaders {
		if value := resp.Header.Get(header); value != "" {
			pathInfo.Headers[header] = value
		}
	}

	// Read body for analysis
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		pathInfo.Size = int64(len(body))

		// Extract title for HTML pages
		if strings.Contains(pathInfo.ContentType, "text/html") {
			if title := extractHTMLTitle(string(body)); title != "" {
				pathInfo.Title = title
			}
		}

		// Technology detection
		pathInfo.Technology = detectPathTechnology(string(body), resp.Header)
	}

	return pathInfo
}

// getWordlist returns a wordlist based on type
func (w *WebEnumModule) getWordlist(wordlistType string) []string {
	switch wordlistType {
	case "quick":
		return []string{
			"admin", "login", "dashboard", "panel", "config", "backup",
			"test", "dev", "staging", "api", "docs", "help",
		}
	case "extensive":
		return w.getExtensiveWordlist()
	default: // common
		return []string{
			"admin", "administrator", "login", "signin", "dashboard", "panel",
			"config", "configuration", "settings", "setup", "install",
			"backup", "backups", "db", "database", "sql", "test", "testing",
			"dev", "development", "staging", "production", "api", "v1", "v2",
			"docs", "documentation", "help", "support", "about", "contact",
			"uploads", "files", "images", "img", "css", "js", "scripts",
			"assets", "static", "media", "tmp", "temp", "cache", "logs",
			"error", "errors", "404", "403", "500", "robots.txt", "sitemap.xml",
			"favicon.ico", "crossdomain.xml", "clientaccesspolicy.xml",
			"web.config", ".htaccess", ".env", "readme.txt", "changelog.txt",
		}
	}
}

// getExtensiveWordlist returns an extensive wordlist
func (w *WebEnumModule) getExtensiveWordlist() []string {
	common := w.getWordlist("common")
	extensive := []string{
		"phpmyadmin", "phpMyAdmin", "mysql", "mssql", "oracle", "postgres",
		"wp-admin", "wp-content", "wp-includes", "wp-config.php",
		"xmlrpc.php", "readme.html", "license.txt",
		"cpanel", "cPanel", "plesk", "webmin", "virtualmin",
		"manager", "tomcat", "jboss", "weblogic", "websphere",
		"jenkins", "bamboo", "teamcity", "gitlab", "github",
		"svn", "git", ".git", ".svn", "CVS", ".hg",
		"private", "secret", "hidden", "internal", "confidential",
		"old", "new", "beta", "alpha", "demo", "example", "sample",
		"vendor", "node_modules", "bower_components", "packages",
		"include", "includes", "lib", "libs", "library", "libraries",
		"src", "source", "resources", "public", "www", "html",
		"cgi-bin", "bin", "sbin", "usr", "var", "etc", "opt",
		"mail", "email", "ftp", "sftp", "ssh", "telnet", "rdp",
		"proxy", "load-balancer", "firewall", "router", "switch",
	}

	return append(common, extensive...)
}

// parseExtensions parses extension string
func (w *WebEnumModule) parseExtensions(extensions string) []string {
	if extensions == "" {
		return []string{}
	}

	var result []string
	parts := strings.Split(extensions, ",")
	for _, part := range parts {
		ext := strings.TrimSpace(part)
		if ext != "" {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			result = append(result, ext)
		}
	}
	return result
}

// generatePaths generates all paths to test
func (w *WebEnumModule) generatePaths(wordlist, extensions []string) []string {
	var paths []string

	// Add directories (with trailing slash)
	for _, word := range wordlist {
		paths = append(paths, "/"+word+"/")
	}

	// Add files without extension
	for _, word := range wordlist {
		paths = append(paths, "/"+word)
	}

	// Add files with extensions
	for _, word := range wordlist {
		for _, ext := range extensions {
			paths = append(paths, "/"+word+ext)
		}
	}

	return paths
}

// parseStatusCodes parses status codes string
func (w *WebEnumModule) parseStatusCodes(codes string) map[int]bool {
	result := make(map[int]bool)

	parts := strings.Split(codes, ",")
	for _, part := range parts {
		if code := strings.TrimSpace(part); code != "" {
			if statusCode, err := strconv.Atoi(code); err == nil {
				result[statusCode] = true
			}
		}
	}

	return result
}

// analyzeVulnerabilities analyzes found paths for vulnerabilities
func (w *WebEnumModule) analyzeVulnerabilities(paths []PathInfo, baseURL, userAgent string) []VulnInfo {
	var vulns []VulnInfo

	for _, path := range paths {
		// Check for common vulnerabilities
		if strings.Contains(path.Path, ".env") {
			vulns = append(vulns, VulnInfo{
				Type:        "Information Disclosure",
				Path:        path.Path,
				Severity:    "High",
				Description: "Environment file exposed - may contain sensitive information",
			})
		}

		if strings.Contains(path.Path, "config") && path.StatusCode == 200 {
			vulns = append(vulns, VulnInfo{
				Type:        "Information Disclosure",
				Path:        path.Path,
				Severity:    "Medium",
				Description: "Configuration file accessible",
			})
		}

		if strings.Contains(path.Path, "backup") && path.StatusCode == 200 {
			vulns = append(vulns, VulnInfo{
				Type:        "Information Disclosure",
				Path:        path.Path,
				Severity:    "Medium",
				Description: "Backup file accessible",
			})
		}

		if strings.Contains(path.Path, "admin") && path.StatusCode == 200 {
			vulns = append(vulns, VulnInfo{
				Type:        "Administrative Interface",
				Path:        path.Path,
				Severity:    "Medium",
				Description: "Administrative interface accessible",
			})
		}

		if path.StatusCode == 403 && strings.Contains(path.Path, "/") {
			vulns = append(vulns, VulnInfo{
				Type:        "Directory Listing",
				Path:        path.Path,
				Severity:    "Low",
				Description: "Directory listing may be enabled",
			})
		}
	}

	return vulns
}

// Helper functions
func extractHTMLTitle(html string) string {
	start := strings.Index(strings.ToLower(html), "<title>")
	if start == -1 {
		return ""
	}
	start += 7

	end := strings.Index(strings.ToLower(html[start:]), "</title>")
	if end == -1 {
		return ""
	}

	title := strings.TrimSpace(html[start : start+end])
	if len(title) > 100 {
		title = title[:100] + "..."
	}
	return title
}

func detectPathTechnology(body string, headers http.Header) []string {
	var tech []string
	bodyLower := strings.ToLower(body)

	// Framework detection
	if strings.Contains(bodyLower, "laravel") {
		tech = append(tech, "Laravel")
	}
	if strings.Contains(bodyLower, "symfony") {
		tech = append(tech, "Symfony")
	}
	if strings.Contains(bodyLower, "codeigniter") {
		tech = append(tech, "CodeIgniter")
	}

	// JavaScript frameworks
	if strings.Contains(bodyLower, "jquery") {
		tech = append(tech, "jQuery")
	}
	if strings.Contains(bodyLower, "bootstrap") {
		tech = append(tech, "Bootstrap")
	}

	return tech
}
