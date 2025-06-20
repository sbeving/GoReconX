package modules

import (
	"GoReconX/internal/config"
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SubdomainEnumerator handles subdomain enumeration
type SubdomainEnumerator struct {
	config *config.Config
	logger *logrus.Logger
}

// SubdomainResult represents a discovered subdomain
type SubdomainResult struct {
	Subdomain string   `json:"subdomain"`
	IPs       []string `json:"ips"`
	Resolved  bool     `json:"resolved"`
}

// NewSubdomainEnumerator creates a new subdomain enumerator
func NewSubdomainEnumerator(cfg *config.Config, logger *logrus.Logger) *SubdomainEnumerator {
	return &SubdomainEnumerator{
		config: cfg,
		logger: logger,
	}
}

// GetName returns the module name
func (se *SubdomainEnumerator) GetName() string {
	return "Subdomain Enumerator"
}

// GetDescription returns the module description
func (se *SubdomainEnumerator) GetDescription() string {
	return "Enumerates subdomains using wordlist-based DNS resolution"
}

// Validate validates the target domain
func (se *SubdomainEnumerator) Validate(target string) error {
	if target == "" {
		return fmt.Errorf("target domain cannot be empty")
	}

	// Basic domain validation
	if !strings.Contains(target, ".") {
		return fmt.Errorf("invalid domain format")
	}

	return nil
}

// GetDefaultOptions returns default options for the module
func (se *SubdomainEnumerator) GetDefaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"wordlist":    se.config.Wordlists.Subdomains,
		"threads":     50,
		"timeout":     5,
		"resolve_ips": true,
	}
}

// Execute performs subdomain enumeration
func (se *SubdomainEnumerator) Execute(target string, options map[string]interface{}) (*ScanResult, error) {
	startTime := time.Now()
	se.logger.WithField("target", target).Info("Starting subdomain enumeration")

	result := &ScanResult{
		ModuleName: se.GetName(),
		Target:     target,
		Status:     "running",
		StartTime:  startTime.Format(time.RFC3339),
		Metadata:   make(map[string]interface{}),
	}

	// Get options
	wordlistPath, _ := options["wordlist"].(string)
	threads, _ := options["threads"].(int)
	timeout, _ := options["timeout"].(int)
	resolveIPs, _ := options["resolve_ips"].(bool)

	if wordlistPath == "" {
		wordlistPath = se.config.Wordlists.Subdomains
	}

	// Load wordlist
	subdomains, err := se.loadWordlist(wordlistPath)
	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = fmt.Sprintf("Failed to load wordlist: %v", err)
		result.EndTime = time.Now().Format(time.RFC3339)
		return result, err
	}

	se.logger.WithField("wordlist_size", len(subdomains)).Info("Loaded subdomain wordlist")

	// Perform enumeration
	results := se.enumerateSubdomains(target, subdomains, threads, timeout, resolveIPs)

	// Convert results to interface slice
	var interfaceResults []interface{}
	for _, r := range results {
		interfaceResults = append(interfaceResults, r)
	}

	endTime := time.Now()
	result.Results = interfaceResults
	result.Status = "completed"
	result.EndTime = endTime.Format(time.RFC3339)
	result.Metadata["found_subdomains"] = len(results)
	result.Metadata["duration_seconds"] = endTime.Sub(startTime).Seconds()

	se.logger.WithFields(logrus.Fields{
		"target":   target,
		"found":    len(results),
		"duration": endTime.Sub(startTime),
	}).Info("Subdomain enumeration completed")

	return result, nil
}

// loadWordlist loads subdomains from a wordlist file
func (se *SubdomainEnumerator) loadWordlist(filename string) ([]string, error) {
	// Create default wordlist if it doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		se.logger.Warn("Wordlist not found, creating default wordlist")
		if err := se.createDefaultWordlist(filename); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var subdomains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains, scanner.Err()
}

// createDefaultWordlist creates a basic subdomain wordlist
func (se *SubdomainEnumerator) createDefaultWordlist(filename string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll("wordlists", 0755); err != nil {
		return err
	}

	defaultSubdomains := []string{
		"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop", "ns1", "ns2",
		"webdisk", "ns", "test", "blog", "pop3", "dev", "www2", "admin", "forum",
		"news", "vpn", "ns3", "mail2", "new", "mysql", "old", "www1", "beta",
		"exchange", "mx", "linux", "ftp2", "test2", "ns4", "www3", "dns1", "api",
		"dns2", "web", "email", "git", "mobile", "demo", "secure", "vpn2", "server",
		"staging", "app", "cdn", "images", "static", "media", "docs", "help",
		"support", "portal", "shop", "store", "payment", "checkout", "cart", "my",
		"account", "profile", "user", "backup", "archive", "data", "files",
		"assets", "resources", "analytics", "stats", "reports", "logs", "api2",
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, subdomain := range defaultSubdomains {
		if _, err := file.WriteString(subdomain + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// enumerateSubdomains performs concurrent subdomain enumeration
func (se *SubdomainEnumerator) enumerateSubdomains(domain string, subdomains []string, threads, timeout int, resolveIPs bool) []*SubdomainResult {
	var results []*SubdomainResult
	var resultsMutex sync.Mutex

	// Create worker pool
	semaphore := make(chan struct{}, threads)
	var wg sync.WaitGroup

	for _, subdomain := range subdomains {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			fullDomain := fmt.Sprintf("%s.%s", sub, domain)

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			defer cancel()

			// Resolve domain
			resolver := &net.Resolver{}
			ips, err := resolver.LookupIPAddr(ctx, fullDomain)

			if err == nil && len(ips) > 0 {
				var ipStrings []string
				if resolveIPs {
					for _, ip := range ips {
						ipStrings = append(ipStrings, ip.IP.String())
					}
				}

				resultsMutex.Lock()
				results = append(results, &SubdomainResult{
					Subdomain: fullDomain,
					IPs:       ipStrings,
					Resolved:  true,
				})
				resultsMutex.Unlock()

				se.logger.WithFields(logrus.Fields{
					"subdomain": fullDomain,
					"ips":       ipStrings,
				}).Debug("Found subdomain")
			}
		}(subdomain)
	}

	wg.Wait()
	return results
}
