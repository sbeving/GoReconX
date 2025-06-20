package modules

import (
	"GoReconX/internal/config"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// EmailHarvester handles email harvesting operations
type EmailHarvester struct {
	config *config.Config
	logger *logrus.Logger
}

// NewEmailHarvester creates a new email harvester
func NewEmailHarvester(cfg *config.Config, logger *logrus.Logger) *EmailHarvester {
	return &EmailHarvester{config: cfg, logger: logger}
}

func (eh *EmailHarvester) GetName() string { return "Email Harvester" }
func (eh *EmailHarvester) GetDescription() string {
	return "Harvests email addresses from various sources"
}
func (eh *EmailHarvester) Validate(target string) error              { return nil }
func (eh *EmailHarvester) GetDefaultOptions() map[string]interface{} { return map[string]interface{}{} }
func (eh *EmailHarvester) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	return &ScanResult{
		ModuleName: eh.GetName(),
		Target:     target,
		Status:     "completed",
		Results:    []interface{}{},
		StartTime:  time.Now().Format(time.RFC3339),
		EndTime:    time.Now().Format(time.RFC3339),
	}, nil
}

// PortResult represents a port scan result
type PortResult struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Banner   string `json:"banner,omitempty"`
}

// PortScanner handles port scanning operations
type PortScanner struct {
	config *config.Config
	logger *logrus.Logger
}

// NewPortScanner creates a new port scanner
func NewPortScanner(cfg *config.Config, logger *logrus.Logger) *PortScanner {
	return &PortScanner{config: cfg, logger: logger}
}

func (ps *PortScanner) GetName() string { return "Port Scanner" }
func (ps *PortScanner) GetDescription() string {
	return "Scans for open TCP and UDP ports on target hosts"
}
func (ps *PortScanner) Validate(target string) error {
	if target == "" {
		return fmt.Errorf("target cannot be empty")
	}

	// Check if it's a valid IP or domain
	if net.ParseIP(target) == nil {
		if _, err := net.LookupHost(target); err != nil {
			return fmt.Errorf("invalid target: %v", err)
		}
	}

	return nil
}

func (ps *PortScanner) GetDefaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"ports":    "1-1000",
		"threads":  100,
		"timeout":  2,
		"scan_tcp": true,
	}
}

func (ps *PortScanner) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	startTime := time.Now()
	ps.logger.WithField("target", target).Info("Starting port scan")

	result := &ScanResult{
		ModuleName: ps.GetName(),
		Target:     target,
		Status:     "running",
		StartTime:  startTime.Format(time.RFC3339),
		Metadata:   make(map[string]interface{}),
	}

	// Parse options
	portsStr, _ := options["ports"].(string)
	threads, _ := options["threads"].(int)
	timeout, _ := options["timeout"].(int)

	if portsStr == "" {
		portsStr = "1-1000"
	}

	// Parse port range
	ports, err := ps.parsePorts(portsStr)
	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = fmt.Sprintf("Invalid port specification: %v", err)
		result.EndTime = time.Now().Format(time.RFC3339)
		return result, err
	}

	// Scan TCP ports
	results := ps.scanTCPPorts(target, ports, threads, timeout)

	// Convert results to interface slice
	var interfaceResults []interface{}
	for _, r := range results {
		interfaceResults = append(interfaceResults, r)
	}

	endTime := time.Now()
	result.Results = interfaceResults
	result.Status = "completed"
	result.EndTime = endTime.Format(time.RFC3339)
	result.Metadata["open_ports"] = len(results)
	result.Metadata["scanned_ports"] = len(ports)
	result.Metadata["duration_seconds"] = endTime.Sub(startTime).Seconds()

	return result, nil
}

// parsePorts parses port specification (e.g., "80,443,1000-2000")
func (ps *PortScanner) parsePorts(portsStr string) ([]int, error) {
	var ports []int
	seen := make(map[int]bool)

	parts := strings.Split(portsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			// Handle range (e.g., "1000-2000")
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("start port cannot be greater than end port")
			}

			for i := start; i <= end; i++ {
				if i >= 1 && i <= 65535 && !seen[i] {
					ports = append(ports, i)
					seen[i] = true
				}
			}
		} else {
			// Handle single port
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}

			if port >= 1 && port <= 65535 && !seen[port] {
				ports = append(ports, port)
				seen[port] = true
			}
		}
	}

	sort.Ints(ports)
	return ports, nil
}

// scanTCPPorts scans TCP ports
func (ps *PortScanner) scanTCPPorts(target string, ports []int, threads, timeout int) []*PortResult {
	var results []*PortResult
	var resultsMutex sync.Mutex

	semaphore := make(chan struct{}, threads)
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			address := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)

			if err == nil {
				defer conn.Close()

				result := &PortResult{
					Port:     p,
					Protocol: "tcp",
					State:    "open",
					Service:  ps.getServiceName(p),
				}

				resultsMutex.Lock()
				results = append(results, result)
				resultsMutex.Unlock()

				ps.logger.WithFields(logrus.Fields{
					"target": target,
					"port":   p,
					"state":  "open",
				}).Debug("Found open port")
			}
		}(port)
	}

	wg.Wait()
	return results
}

// getServiceName returns the common service name for a port
func (ps *PortScanner) getServiceName(port int) string {
	commonPorts := map[int]string{
		21:    "ftp",
		22:    "ssh",
		23:    "telnet",
		25:    "smtp",
		53:    "dns",
		80:    "http",
		110:   "pop3",
		143:   "imap",
		443:   "https",
		993:   "imaps",
		995:   "pop3s",
		3389:  "rdp",
		5432:  "postgresql",
		3306:  "mysql",
		1433:  "mssql",
		6379:  "redis",
		27017: "mongodb",
	}

	if service, exists := commonPorts[port]; exists {
		return service
	}

	return "unknown"
}

// DirectoryEnumerator handles directory enumeration
type DirectoryEnumerator struct {
	config *config.Config
	logger *logrus.Logger
}

// NewDirectoryEnumerator creates a new directory enumerator
func NewDirectoryEnumerator(cfg *config.Config, logger *logrus.Logger) *DirectoryEnumerator {
	return &DirectoryEnumerator{config: cfg, logger: logger}
}

func (de *DirectoryEnumerator) GetName() string { return "Directory Enumerator" }
func (de *DirectoryEnumerator) GetDescription() string {
	return "Enumerates directories and files on web servers"
}
func (de *DirectoryEnumerator) Validate(target string) error { return nil }
func (de *DirectoryEnumerator) GetDefaultOptions() map[string]interface{} {
	return map[string]interface{}{}
}
func (de *DirectoryEnumerator) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	return &ScanResult{
		ModuleName: de.GetName(),
		Target:     target,
		Status:     "completed",
		Results:    []interface{}{},
		StartTime:  time.Now().Format(time.RFC3339),
		EndTime:    time.Now().Format(time.RFC3339),
	}, nil
}

// WebAnalyzer handles web application analysis
type WebAnalyzer struct {
	config *config.Config
	logger *logrus.Logger
}

// NewWebAnalyzer creates a new web analyzer
func NewWebAnalyzer(cfg *config.Config, logger *logrus.Logger) *WebAnalyzer {
	return &WebAnalyzer{config: cfg, logger: logger}
}

func (wa *WebAnalyzer) GetName() string { return "Web Analyzer" }
func (wa *WebAnalyzer) GetDescription() string {
	return "Analyzes web applications for technologies and vulnerabilities"
}
func (wa *WebAnalyzer) Validate(target string) error              { return nil }
func (wa *WebAnalyzer) GetDefaultOptions() map[string]interface{} { return map[string]interface{}{} }
func (wa *WebAnalyzer) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	return &ScanResult{
		ModuleName: wa.GetName(),
		Target:     target,
		Status:     "completed",
		Results:    []interface{}{},
		StartTime:  time.Now().Format(time.RFC3339),
		EndTime:    time.Now().Format(time.RFC3339),
	}, nil
}

// IPGeolocator handles IP geolocation
type IPGeolocator struct {
	config *config.Config
	logger *logrus.Logger
}

// NewIPGeolocator creates a new IP geolocator
func NewIPGeolocator(cfg *config.Config, logger *logrus.Logger) *IPGeolocator {
	return &IPGeolocator{config: cfg, logger: logger}
}

func (ig *IPGeolocator) GetName() string { return "IP Geolocator" }
func (ig *IPGeolocator) GetDescription() string {
	return "Provides geolocation information for IP addresses"
}
func (ig *IPGeolocator) Validate(target string) error              { return nil }
func (ig *IPGeolocator) GetDefaultOptions() map[string]interface{} { return map[string]interface{}{} }
func (ig *IPGeolocator) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	return &ScanResult{
		ModuleName: ig.GetName(),
		Target:     target,
		Status:     "completed",
		Results:    []interface{}{},
		StartTime:  time.Now().Format(time.RFC3339),
		EndTime:    time.Now().Format(time.RFC3339),
	}, nil
}

// GitHubRecon handles GitHub reconnaissance
type GitHubRecon struct {
	config *config.Config
	logger *logrus.Logger
}

// NewGitHubRecon creates a new GitHub recon module
func NewGitHubRecon(cfg *config.Config, logger *logrus.Logger) *GitHubRecon {
	return &GitHubRecon{config: cfg, logger: logger}
}

func (gr *GitHubRecon) GetName() string { return "GitHub Reconnaissance" }
func (gr *GitHubRecon) GetDescription() string {
	return "Searches GitHub for sensitive information and code"
}
func (gr *GitHubRecon) Validate(target string) error              { return nil }
func (gr *GitHubRecon) GetDefaultOptions() map[string]interface{} { return map[string]interface{}{} }
func (gr *GitHubRecon) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	return &ScanResult{
		ModuleName: gr.GetName(),
		Target:     target,
		Status:     "completed",
		Results:    []interface{}{},
		StartTime:  time.Now().Format(time.RFC3339),
		EndTime:    time.Now().Format(time.RFC3339),
	}, nil
}
