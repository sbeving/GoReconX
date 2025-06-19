package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/likexian/whois"
	"github.com/miekg/dns"
)

// DomainEnumModule implements domain and subdomain enumeration
type DomainEnumModule struct {
	*BaseModule
	client *resty.Client
}

// DomainResult represents domain enumeration results
type DomainResult struct {
	Domain       string              `json:"domain"`
	Subdomains   []SubdomainInfo     `json:"subdomains"`
	DNSRecords   map[string][]string `json:"dns_records"`
	WhoisInfo    WhoisInfo           `json:"whois_info"`
	TechStack    []string            `json:"tech_stack"`
	Certificates []CertInfo          `json:"certificates"`
}

// SubdomainInfo contains subdomain information
type SubdomainInfo struct {
	Subdomain  string   `json:"subdomain"`
	IPs        []string `json:"ips"`
	Status     string   `json:"status"`
	Technology []string `json:"technology"`
	Title      string   `json:"title"`
	StatusCode int      `json:"status_code"`
}

// WhoisInfo contains WHOIS data
type WhoisInfo struct {
	Registrar   string   `json:"registrar"`
	CreatedDate string   `json:"created_date"`
	ExpiryDate  string   `json:"expiry_date"`
	NameServers []string `json:"name_servers"`
	Registrant  string   `json:"registrant"`
	AdminEmail  string   `json:"admin_email"`
}

// CertInfo contains certificate information
type CertInfo struct {
	CommonName string   `json:"common_name"`
	SANs       []string `json:"sans"`
	Issuer     string   `json:"issuer"`
	ValidFrom  string   `json:"valid_from"`
	ValidTo    string   `json:"valid_to"`
}

// NewDomainEnumModule creates a new domain enumeration module
func NewDomainEnumModule() *DomainEnumModule {
	info := ModuleInfo{
		Name:        "domain_enum",
		Category:    "passive_osint",
		Description: "Comprehensive domain and subdomain enumeration with DNS analysis",
		Version:     "1.0.0",
		Author:      "GoReconX Team",
		Tags:        []string{"domain", "subdomain", "dns", "whois", "passive"},
		Options: []ModuleOption{
			{
				Name:        "use_wordlist",
				Type:        "bool",
				Description: "Enable wordlist-based subdomain brute forcing",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "use_apis",
				Type:        "bool",
				Description: "Use external APIs for subdomain discovery",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "use_crt_sh",
				Type:        "bool",
				Description: "Query Certificate Transparency logs",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "deep_scan",
				Type:        "bool",
				Description: "Perform deep enumeration (slower but more thorough)",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "dns_timeout",
				Type:        "int",
				Description: "DNS query timeout in seconds",
				Required:    false,
				Default:     5,
			},
		},
		Requirements: []string{"network"},
	}

	module := &DomainEnumModule{
		BaseModule: NewBaseModule(info),
		client:     resty.New().SetTimeout(10 * time.Second),
	}

	return module
}

// Validate validates the module input
func (d *DomainEnumModule) Validate(input ModuleInput) error {
	if err := d.ValidateInput(input); err != nil {
		return err
	}

	// Validate domain format
	if !isValidDomain(input.Target) {
		return NewModuleError("invalid domain format", "INVALID_DOMAIN")
	}

	return nil
}

// Execute runs the domain enumeration module
func (d *DomainEnumModule) Execute(ctx context.Context, input ModuleInput, output chan<- ModuleResult) error {
	d.SetStatus("running", 0.0, "Starting domain enumeration")

	domain := strings.ToLower(strings.TrimSpace(input.Target))
	result := &DomainResult{
		Domain:       domain,
		Subdomains:   []SubdomainInfo{},
		DNSRecords:   make(map[string][]string),
		TechStack:    []string{},
		Certificates: []CertInfo{},
	}

	// Phase 1: WHOIS Lookup
	d.SetStatus("running", 0.1, "Performing WHOIS lookup")
	d.SendResult(output, "progress", "Performing WHOIS lookup", nil, input.SessionID)

	if whoisInfo, err := d.performWhoisLookup(domain); err == nil {
		result.WhoisInfo = whoisInfo
		d.SendResult(output, "data", map[string]interface{}{
			"type": "whois",
			"data": whoisInfo,
		}, nil, input.SessionID)
	}

	if d.IsStopped() {
		return nil
	}

	// Phase 2: DNS Record Enumeration
	d.SetStatus("running", 0.2, "Enumerating DNS records")
	d.SendResult(output, "progress", "Enumerating DNS records", nil, input.SessionID)

	dnsRecords := d.enumerateDNSRecords(domain)
	result.DNSRecords = dnsRecords
	d.SendResult(output, "data", map[string]interface{}{
		"type": "dns_records",
		"data": dnsRecords,
	}, nil, input.SessionID)

	if d.IsStopped() {
		return nil
	}

	// Phase 3: Certificate Transparency
	useCtSh, _ := input.Options["use_crt_sh"].(bool)
	if useCtSh {
		d.SetStatus("running", 0.4, "Querying Certificate Transparency logs")
		d.SendResult(output, "progress", "Querying Certificate Transparency logs", nil, input.SessionID)

		if certs, err := d.queryCertificateTransparency(domain); err == nil {
			result.Certificates = certs
			d.SendResult(output, "data", map[string]interface{}{
				"type": "certificates",
				"data": certs,
			}, nil, input.SessionID)
		}
	}

	if d.IsStopped() {
		return nil
	}

	// Phase 4: Subdomain Discovery
	d.SetStatus("running", 0.6, "Discovering subdomains")
	d.SendResult(output, "progress", "Discovering subdomains", nil, input.SessionID)

	subdomains := d.discoverSubdomains(ctx, domain, input.Options)

	// Phase 5: Subdomain Analysis
	d.SetStatus("running", 0.8, "Analyzing discovered subdomains")
	d.SendResult(output, "progress", "Analyzing discovered subdomains", nil, input.SessionID)

	for i, subdomain := range subdomains {
		if d.IsStopped() {
			break
		}

		subInfo := d.analyzeSubdomain(subdomain)
		result.Subdomains = append(result.Subdomains, subInfo)

		d.SendResult(output, "data", map[string]interface{}{
			"type": "subdomain",
			"data": subInfo,
		}, nil, input.SessionID)

		progress := 0.8 + (0.2 * float64(i+1) / float64(len(subdomains)))
		d.SetStatus("running", progress, fmt.Sprintf("Analyzed %d/%d subdomains", i+1, len(subdomains)))
	}

	// Send final result
	d.SetStatus("completed", 1.0, "Domain enumeration completed")
	d.SendResult(output, "complete", result, map[string]interface{}{
		"total_subdomains": len(result.Subdomains),
		"dns_records":      len(result.DNSRecords),
		"certificates":     len(result.Certificates),
	}, input.SessionID)

	return nil
}

// performWhoisLookup performs WHOIS lookup for domain
func (d *DomainEnumModule) performWhoisLookup(domain string) (WhoisInfo, error) {
	whoisData, err := whois.Whois(domain)
	if err != nil {
		return WhoisInfo{}, err
	}

	// Parse WHOIS data (simplified)
	info := WhoisInfo{}
	lines := strings.Split(whoisData, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "registrar:") {
			info.Registrar = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(strings.ToLower(line), "creation date:") || strings.Contains(strings.ToLower(line), "created:") {
			info.CreatedDate = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(strings.ToLower(line), "expiry date:") || strings.Contains(strings.ToLower(line), "expires:") {
			info.ExpiryDate = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return info, nil
}

// enumerateDNSRecords enumerates various DNS record types
func (d *DomainEnumModule) enumerateDNSRecords(domain string) map[string][]string {
	records := make(map[string][]string)

	recordTypes := []uint16{
		dns.TypeA, dns.TypeAAAA, dns.TypeMX, dns.TypeNS,
		dns.TypeTXT, dns.TypeCNAME, dns.TypeSOA, dns.TypeSRV,
	}

	c := dns.Client{Timeout: 5 * time.Second}

	for _, recordType := range recordTypes {
		typeName := dns.TypeToString[recordType]

		m := &dns.Msg{}
		m.SetQuestion(dns.Fqdn(domain), recordType)

		if r, _, err := c.Exchange(m, "8.8.8.8:53"); err == nil {
			var values []string
			for _, ans := range r.Answer {
				values = append(values, strings.TrimSpace(strings.Fields(ans.String())[4:][0]))
			}
			if len(values) > 0 {
				records[typeName] = values
			}
		}
	}

	return records
}

// queryCertificateTransparency queries CT logs for certificates
func (d *DomainEnumModule) queryCertificateTransparency(domain string) ([]CertInfo, error) {
	var certs []CertInfo

	// Query crt.sh
	url := fmt.Sprintf("https://crt.sh/?q=%%.%s&output=json", domain)

	resp, err := d.client.R().Get(url)
	if err != nil {
		return certs, err
	}

	var ctResults []map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &ctResults); err != nil {
		return certs, err
	}

	for _, ct := range ctResults {
		cert := CertInfo{
			CommonName: getString(ct, "common_name"),
			Issuer:     getString(ct, "issuer_name"),
			ValidFrom:  getString(ct, "not_before"),
			ValidTo:    getString(ct, "not_after"),
		}

		// Parse SANs if available
		if sans, ok := ct["name_value"].(string); ok {
			cert.SANs = strings.Split(sans, "\n")
		}

		certs = append(certs, cert)
	}

	return certs, nil
}

// discoverSubdomains discovers subdomains using various techniques
func (d *DomainEnumModule) discoverSubdomains(ctx context.Context, domain string, options map[string]interface{}) []string {
	subdomains := make(map[string]bool)

	// Add main domain
	subdomains[domain] = true

	// Wordlist-based discovery
	useWordlist, _ := options["use_wordlist"].(bool)
	if useWordlist {
		wordlistSubs := d.wordlistSubdomains(domain)
		for _, sub := range wordlistSubs {
			subdomains[sub] = true
		}
	}

	// Certificate Transparency parsing for subdomains
	useCrtSh, _ := options["use_crt_sh"].(bool)
	if useCrtSh {
		if certs, err := d.queryCertificateTransparency(domain); err == nil {
			for _, cert := range certs {
				if strings.HasSuffix(cert.CommonName, domain) {
					subdomains[cert.CommonName] = true
				}
				for _, san := range cert.SANs {
					if strings.HasSuffix(san, domain) && san != domain {
						subdomains[san] = true
					}
				}
			}
		}
	}

	// Search engine enumeration
	useAPIs, _ := options["use_apis"].(bool)
	if useAPIs {
		searchSubs := d.searchEngineSubdomains(domain)
		for _, sub := range searchSubs {
			subdomains[sub] = true
		}
	}

	// Convert map to slice
	var result []string
	for sub := range subdomains {
		if sub != "" && isValidDomain(sub) {
			result = append(result, sub)
		}
	}

	return result
}

// wordlistSubdomains performs wordlist-based subdomain discovery
func (d *DomainEnumModule) wordlistSubdomains(domain string) []string {
	var subdomains []string

	// Common subdomain wordlist
	wordlist := []string{
		"www", "mail", "ftp", "admin", "test", "dev", "staging", "api",
		"blog", "forum", "shop", "store", "m", "mobile", "app", "portal",
		"secure", "vpn", "remote", "cdn", "static", "assets", "img", "images",
		"video", "media", "docs", "support", "help", "status", "monitor",
		"demo", "beta", "alpha", "labs", "research", "news", "careers",
		"jobs", "legal", "privacy", "terms", "about", "contact", "login",
		"register", "dashboard", "panel", "control", "cpanel", "webmail",
		"email", "smtp", "pop", "imap", "ns1", "ns2", "dns", "mx",
	}

	c := dns.Client{Timeout: 2 * time.Second}

	for _, prefix := range wordlist {
		if d.IsStopped() {
			break
		}

		subdomain := fmt.Sprintf("%s.%s", prefix, domain)

		m := &dns.Msg{}
		m.SetQuestion(dns.Fqdn(subdomain), dns.TypeA)

		if _, _, err := c.Exchange(m, "8.8.8.8:53"); err == nil {
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains
}

// searchEngineSubdomains uses search engines for subdomain discovery
func (d *DomainEnumModule) searchEngineSubdomains(domain string) []string {
	var subdomains []string

	// Google dork for subdomains
	query := fmt.Sprintf("site:%s", domain)

	// This is a simplified implementation
	// In a real scenario, you'd use proper search APIs
	resp, err := d.client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(fmt.Sprintf("https://www.google.com/search?q=%s", query))

	if err == nil {
		// Parse response for subdomains (simplified)
		body := string(resp.Body())
		// This would need proper HTML parsing and subdomain extraction
		_ = body // Placeholder
	}

	return subdomains
}

// analyzeSubdomain analyzes a subdomain for additional information
func (d *DomainEnumModule) analyzeSubdomain(subdomain string) SubdomainInfo {
	info := SubdomainInfo{
		Subdomain:  subdomain,
		IPs:        []string{},
		Status:     "unknown",
		Technology: []string{},
		Title:      "",
		StatusCode: 0,
	}

	// Resolve IP addresses
	if ips, err := net.LookupIP(subdomain); err == nil {
		for _, ip := range ips {
			info.IPs = append(info.IPs, ip.String())
		}
		info.Status = "active"
	} else {
		info.Status = "inactive"
		return info
	}

	// HTTP analysis
	urls := []string{
		fmt.Sprintf("https://%s", subdomain),
		fmt.Sprintf("http://%s", subdomain),
	}
	for _, url := range urls {
		resp, err := d.client.R().
			SetHeader("User-Agent", "GoReconX/1.0").
			Get(url)

		if err == nil {
			info.StatusCode = resp.StatusCode()
			if resp.StatusCode() < 400 {
				// Extract title and technology detection would go here
				body := string(resp.Body())
				if title := extractTitle(body); title != "" {
					info.Title = title
				}

				// Basic technology detection
				info.Technology = detectTechnology(resp.Header(), body)
			}
			break
		}
	}

	return info
}

// Helper functions
func isValidDomain(domain string) bool {
	if domain == "" || len(domain) > 253 {
		return false
	}

	// Basic domain validation
	return strings.Contains(domain, ".") && !strings.HasPrefix(domain, ".") && !strings.HasSuffix(domain, ".")
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func extractTitle(html string) string {
	// Simplified title extraction
	start := strings.Index(strings.ToLower(html), "<title>")
	if start == -1 {
		return ""
	}
	start += 7

	end := strings.Index(strings.ToLower(html[start:]), "</title>")
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(html[start : start+end])
}

func detectTechnology(headers map[string][]string, body string) []string {
	var tech []string

	// Server header
	if server := headers["Server"]; len(server) > 0 {
		tech = append(tech, server[0])
	}

	// X-Powered-By header
	if powered := headers["X-Powered-By"]; len(powered) > 0 {
		tech = append(tech, powered[0])
	}

	// Body analysis for frameworks/CMS
	bodyLower := strings.ToLower(body)
	if strings.Contains(bodyLower, "wordpress") {
		tech = append(tech, "WordPress")
	}
	if strings.Contains(bodyLower, "drupal") {
		tech = append(tech, "Drupal")
	}
	if strings.Contains(bodyLower, "joomla") {
		tech = append(tech, "Joomla")
	}
	if strings.Contains(bodyLower, "react") {
		tech = append(tech, "React")
	}
	if strings.Contains(bodyLower, "angular") {
		tech = append(tech, "Angular")
	}
	if strings.Contains(bodyLower, "vue") {
		tech = append(tech, "Vue.js")
	}

	return tech
}
