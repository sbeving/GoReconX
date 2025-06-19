package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// NetworkReconModule implements comprehensive network reconnaissance
type NetworkReconModule struct {
	*BaseModule
	client *http.Client
}

// NetworkReconResult represents network reconnaissance results
type NetworkReconResult struct {
	Target          string          `json:"target"`
	IPInfo          IPInfo          `json:"ip_info"`
	GeolocationInfo GeolocationInfo `json:"geolocation_info"`
	ASNInfo         ASNInfo         `json:"asn_info"`
	ReverseDNS      []string        `json:"reverse_dns"`
	PortScan        NetworkPortScan `json:"port_scan"`
	NetworkRange    string          `json:"network_range"`
	ThreatIntel     ThreatIntelInfo `json:"threat_intel"`
	ScanTime        string          `json:"scan_time"`
}

// IPInfo contains IP address information
type IPInfo struct {
	IP         string   `json:"ip"`
	Type       string   `json:"type"` // IPv4 or IPv6
	IsPublic   bool     `json:"is_public"`
	IsPrivate  bool     `json:"is_private"`
	Hostnames  []string `json:"hostnames"`
	PTRRecords []string `json:"ptr_records"`
}

// GeolocationInfo contains IP geolocation data
type GeolocationInfo struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
}

// ASNInfo contains Autonomous System Number information
type ASNInfo struct {
	ASN         int    `json:"asn"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Registry    string `json:"registry"`
	Allocated   string `json:"allocated"`
	Description string `json:"description"`
}

// NetworkPortScan contains network port scan results
type NetworkPortScan struct {
	TotalPorts  int        `json:"total_ports"`
	OpenPorts   []PortInfo `json:"open_ports"`
	CommonPorts []PortInfo `json:"common_ports"`
}

// ThreatIntelInfo contains threat intelligence information
type ThreatIntelInfo struct {
	IsMalicious bool           `json:"is_malicious"`
	ThreatTypes []string       `json:"threat_types"`
	Reputation  int            `json:"reputation"` // 0-100 scale
	LastSeen    string         `json:"last_seen"`
	Sources     []string       `json:"sources"`
	Reports     []ThreatReport `json:"reports"`
}

// ThreatReport contains individual threat reports
type ThreatReport struct {
	Source      string `json:"source"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Date        string `json:"date"`
}

// NewNetworkReconModule creates a new network reconnaissance module
func NewNetworkReconModule() *NetworkReconModule {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	info := ModuleInfo{
		Name:        "network_recon",
		Category:    "passive_osint",
		Description: "Comprehensive network reconnaissance including IP analysis, geolocation, and threat intelligence",
		Version:     "1.0.0",
		Author:      "GoReconX Team",
		Tags:        []string{"network", "ip", "geolocation", "asn", "threat", "intelligence"},
		Options: []ModuleOption{
			{
				Name:        "include_geolocation",
				Type:        "bool",
				Description: "Include IP geolocation lookup",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "include_asn",
				Type:        "bool",
				Description: "Include ASN (Autonomous System) lookup",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "include_threat_intel",
				Type:        "bool",
				Description: "Include threat intelligence lookup",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "include_port_scan",
				Type:        "bool",
				Description: "Include basic port scanning",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "use_virustotal",
				Type:        "bool",
				Description: "Use VirusTotal API for threat intelligence",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "use_shodan",
				Type:        "bool",
				Description: "Use Shodan API for additional information",
				Required:    false,
				Default:     false,
			},
			{
				Name:        "timeout",
				Type:        "int",
				Description: "Request timeout in seconds",
				Required:    false,
				Default:     15,
			},
		},
		Requirements: []string{"network"},
	}

	return &NetworkReconModule{
		BaseModule: NewBaseModule(info),
		client:     client,
	}
}

// Validate validates the module input
func (n *NetworkReconModule) Validate(input ModuleInput) error {
	if err := n.ValidateInput(input); err != nil {
		return err
	}

	// Validate target is IP or resolvable hostname
	if net.ParseIP(input.Target) == nil {
		if _, err := net.LookupHost(input.Target); err != nil {
			return NewModuleError("target must be a valid IP address or resolvable hostname", "INVALID_TARGET")
		}
	}

	return nil
}

// Execute runs the network reconnaissance module
func (n *NetworkReconModule) Execute(ctx context.Context, input ModuleInput, output chan<- ModuleResult) error {
	startTime := time.Now()
	n.SetStatus("running", 0.0, "Starting network reconnaissance")

	// Parse options
	includeGeo, _ := input.Options["include_geolocation"].(bool)
	includeASN, _ := input.Options["include_asn"].(bool)
	includeThreat, _ := input.Options["include_threat_intel"].(bool)
	includePortScan, _ := input.Options["include_port_scan"].(bool)
	useVirusTotal, _ := input.Options["use_virustotal"].(bool)
	useShodan, _ := input.Options["use_shodan"].(bool)
	timeout, _ := input.Options["timeout"].(int)

	if timeout > 0 {
		n.client.Timeout = time.Duration(timeout) * time.Second
	}

	// Resolve target to IP if it's a hostname
	target := input.Target
	var targetIP string

	if ip := net.ParseIP(target); ip != nil {
		targetIP = target
	} else {
		n.SetStatus("running", 0.05, "Resolving hostname to IP")
		ips, err := net.LookupIP(target)
		if err != nil {
			return NewModuleError("failed to resolve hostname: "+err.Error(), "RESOLUTION_FAILED")
		}
		if len(ips) == 0 {
			return NewModuleError("no IP addresses found for hostname", "NO_IPS_FOUND")
		}
		targetIP = ips[0].String()
	}

	result := &NetworkReconResult{
		Target: target,
	}

	// Phase 1: Basic IP Information
	n.SetStatus("running", 0.1, "Gathering basic IP information")
	n.SendResult(output, "progress", "Gathering basic IP information", nil, input.SessionID)

	ipInfo := n.gatherIPInfo(targetIP)
	result.IPInfo = ipInfo

	n.SendResult(output, "data", map[string]interface{}{
		"type": "ip_info",
		"data": ipInfo,
	}, nil, input.SessionID)

	if n.IsStopped() {
		return nil
	}

	// Phase 2: Geolocation Lookup
	if includeGeo {
		n.SetStatus("running", 0.3, "Performing geolocation lookup")
		n.SendResult(output, "progress", "Performing geolocation lookup", nil, input.SessionID)

		geoInfo := n.performGeolocationLookup(targetIP)
		result.GeolocationInfo = geoInfo

		n.SendResult(output, "data", map[string]interface{}{
			"type": "geolocation",
			"data": geoInfo,
		}, nil, input.SessionID)
	}

	if n.IsStopped() {
		return nil
	}

	// Phase 3: ASN Lookup
	if includeASN {
		n.SetStatus("running", 0.5, "Performing ASN lookup")
		n.SendResult(output, "progress", "Performing ASN lookup", nil, input.SessionID)

		asnInfo := n.performASNLookup(targetIP)
		result.ASNInfo = asnInfo

		n.SendResult(output, "data", map[string]interface{}{
			"type": "asn_info",
			"data": asnInfo,
		}, nil, input.SessionID)
	}

	if n.IsStopped() {
		return nil
	}

	// Phase 4: Reverse DNS
	n.SetStatus("running", 0.6, "Performing reverse DNS lookup")
	n.SendResult(output, "progress", "Performing reverse DNS lookup", nil, input.SessionID)

	reverseDNS := n.performReverseDNS(targetIP)
	result.ReverseDNS = reverseDNS

	n.SendResult(output, "data", map[string]interface{}{
		"type": "reverse_dns",
		"data": reverseDNS,
	}, nil, input.SessionID)

	if n.IsStopped() {
		return nil
	}

	// Phase 5: Port Scanning (if enabled)
	if includePortScan {
		n.SetStatus("running", 0.7, "Performing basic port scan")
		n.SendResult(output, "progress", "Performing basic port scan", nil, input.SessionID)

		portScan := n.performBasicPortScan(targetIP)
		result.PortScan = portScan

		n.SendResult(output, "data", map[string]interface{}{
			"type": "port_scan",
			"data": portScan,
		}, nil, input.SessionID)
	}

	if n.IsStopped() {
		return nil
	}

	// Phase 6: Threat Intelligence
	if includeThreat {
		n.SetStatus("running", 0.8, "Gathering threat intelligence")
		n.SendResult(output, "progress", "Gathering threat intelligence", nil, input.SessionID)

		threatInfo := n.gatherThreatIntelligence(targetIP, useVirusTotal, useShodan, input.Options)
		result.ThreatIntel = threatInfo

		n.SendResult(output, "data", map[string]interface{}{
			"type": "threat_intel",
			"data": threatInfo,
		}, nil, input.SessionID)
	}

	if n.IsStopped() {
		return nil
	}

	// Phase 7: Shodan Integration (if enabled)
	if useShodan {
		n.SetStatus("running", 0.9, "Querying Shodan database")
		n.SendResult(output, "progress", "Querying Shodan database", nil, input.SessionID)

		shodanInfo := n.queryShodan(targetIP, input.Options)
		if shodanInfo != nil {
			n.SendResult(output, "data", map[string]interface{}{
				"type": "shodan_info",
				"data": shodanInfo,
			}, nil, input.SessionID)
		}
	}

	result.ScanTime = time.Since(startTime).String()

	// Send final result
	n.SetStatus("completed", 1.0, "Network reconnaissance completed")
	n.SendResult(output, "complete", result, map[string]interface{}{
		"target_ip":    targetIP,
		"has_geo_info": includeGeo && result.GeolocationInfo.Country != "",
		"has_asn_info": includeASN && result.ASNInfo.ASN != 0,
		"open_ports":   len(result.PortScan.OpenPorts),
		"scan_time":    result.ScanTime,
	}, input.SessionID)

	return nil
}

// gatherIPInfo gathers basic IP information
func (n *NetworkReconModule) gatherIPInfo(ip string) IPInfo {
	info := IPInfo{
		IP:         ip,
		Hostnames:  []string{},
		PTRRecords: []string{},
	}

	// Determine IP type
	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		if parsedIP.To4() != nil {
			info.Type = "IPv4"
		} else {
			info.Type = "IPv6"
		}

		info.IsPrivate = parsedIP.IsPrivate()
		info.IsPublic = !parsedIP.IsPrivate()
	}

	// Reverse DNS lookup
	names, err := net.LookupAddr(ip)
	if err == nil {
		info.Hostnames = names
		info.PTRRecords = names
	}

	return info
}

// performGeolocationLookup performs IP geolocation lookup
func (n *NetworkReconModule) performGeolocationLookup(ip string) GeolocationInfo {
	info := GeolocationInfo{}

	// Using ip-api.com (free service)
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := n.client.Get(url)
	if err != nil {
		return info
	}
	defer resp.Body.Close()

	var apiResult struct {
		Status      string  `json:"status"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"regionName"`
		City        string  `json:"city"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		Timezone    string  `json:"timezone"`
		ISP         string  `json:"isp"`
		Org         string  `json:"org"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return info
	}

	if apiResult.Status == "success" {
		info = GeolocationInfo{
			Country:     apiResult.Country,
			CountryCode: apiResult.CountryCode,
			Region:      apiResult.Region,
			City:        apiResult.City,
			Latitude:    apiResult.Lat,
			Longitude:   apiResult.Lon,
			Timezone:    apiResult.Timezone,
			ISP:         apiResult.ISP,
			Org:         apiResult.Org,
		}
	}

	return info
}

// performASNLookup performs ASN lookup
func (n *NetworkReconModule) performASNLookup(ip string) ASNInfo {
	info := ASNInfo{}

	// Using ipinfo.io ASN API (simplified)
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)

	resp, err := n.client.Get(url)
	if err != nil {
		return info
	}
	defer resp.Body.Close()

	var apiResult struct {
		IP       string `json:"ip"`
		Hostname string `json:"hostname"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Org      string `json:"org"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return info
	}

	// Parse ASN from org field (format: "AS#### Organization Name")
	if strings.HasPrefix(apiResult.Org, "AS") {
		parts := strings.SplitN(apiResult.Org, " ", 2)
		if len(parts) >= 2 {
			asnStr := strings.TrimPrefix(parts[0], "AS")
			if asn, err := parseASN(asnStr); err == nil {
				info.ASN = asn
				info.Name = parts[1]
				info.Country = apiResult.Country
				info.Description = apiResult.Org
			}
		}
	}

	return info
}

// performReverseDNS performs reverse DNS lookup
func (n *NetworkReconModule) performReverseDNS(ip string) []string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return []string{}
	}
	return names
}

// performBasicPortScan performs a basic port scan on common ports
func (n *NetworkReconModule) performBasicPortScan(ip string) NetworkPortScan {
	scan := NetworkPortScan{
		OpenPorts:   []PortInfo{},
		CommonPorts: []PortInfo{},
	}

	// Common ports to check
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 993, 995, 3389, 5432, 3306}

	scan.TotalPorts = len(commonPorts)

	for _, port := range commonPorts {
		if n.IsStopped() {
			break
		}

		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err == nil {
			conn.Close()

			portInfo := PortInfo{
				Port:     port,
				Protocol: "tcp",
				State:    "open",
				Service:  getServiceName(port),
			}

			scan.OpenPorts = append(scan.OpenPorts, portInfo)
		}
	}

	scan.CommonPorts = scan.OpenPorts
	return scan
}

// gatherThreatIntelligence gathers threat intelligence information
func (n *NetworkReconModule) gatherThreatIntelligence(ip string, useVirusTotal, useShodan bool, options map[string]interface{}) ThreatIntelInfo {
	info := ThreatIntelInfo{
		IsMalicious: false,
		ThreatTypes: []string{},
		Reputation:  50, // Neutral
		Sources:     []string{},
		Reports:     []ThreatReport{},
	}

	// Check VirusTotal if enabled and API key is available
	if useVirusTotal {
		if apiKey, exists := options["virustotal_api_key"].(string); exists && apiKey != "" {
			vtInfo := n.queryVirusTotal(ip, apiKey)
			if vtInfo != nil {
				info.Sources = append(info.Sources, "VirusTotal")
				// Merge VirusTotal data
				if vtInfo.IsMalicious {
					info.IsMalicious = true
					info.Reputation = 10 // Low reputation for malicious IPs
				}
			}
		}
	}

	// Basic reputation check using public blacklists
	blacklistInfo := n.checkPublicBlacklists(ip)
	if blacklistInfo.IsMalicious {
		info.IsMalicious = true
		info.ThreatTypes = append(info.ThreatTypes, blacklistInfo.ThreatTypes...)
		info.Sources = append(info.Sources, blacklistInfo.Sources...)
		info.Reports = append(info.Reports, blacklistInfo.Reports...)
		info.Reputation = 20 // Low reputation
	}

	return info
}

// queryShodan queries Shodan API for additional information
func (n *NetworkReconModule) queryShodan(ip string, options map[string]interface{}) interface{} {
	apiKey, exists := options["shodan_api_key"].(string)
	if !exists || apiKey == "" {
		return nil
	}

	url := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ip, apiKey)

	resp, err := n.client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}

	return result
}

// queryVirusTotal queries VirusTotal API
func (n *NetworkReconModule) queryVirusTotal(ip string, apiKey string) *ThreatIntelInfo {
	url := fmt.Sprintf("https://www.virustotal.com/vtapi/v2/ip-address/report?apikey=%s&ip=%s", apiKey, ip)

	resp, err := n.client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var result struct {
		ResponseCode int `json:"response_code"`
		Positives    int `json:"positives"`
		Total        int `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}

	info := &ThreatIntelInfo{
		IsMalicious: result.Positives > 0,
		Sources:     []string{"VirusTotal"},
	}

	if result.Positives > 0 {
		info.ThreatTypes = append(info.ThreatTypes, "Malicious")
		info.Reputation = max(0, 100-(result.Positives*10))
	}

	return info
}

// checkPublicBlacklists checks public blacklists
func (n *NetworkReconModule) checkPublicBlacklists(ip string) ThreatIntelInfo {
	info := ThreatIntelInfo{
		IsMalicious: false,
		ThreatTypes: []string{},
		Sources:     []string{},
		Reports:     []ThreatReport{},
	}

	// This is a simplified implementation
	// In reality, you'd check multiple reputation services

	// Example: Check if IP is in known bad ranges (simplified)
	if n.isKnownBadIP(ip) {
		info.IsMalicious = true
		info.ThreatTypes = append(info.ThreatTypes, "Known Malicious")
		info.Sources = append(info.Sources, "Public Blacklists")
		info.Reports = append(info.Reports, ThreatReport{
			Source:      "Public Blacklist",
			Type:        "Malicious IP",
			Description: "IP found in public blacklist",
			Severity:    "High",
			Date:        time.Now().Format("2006-01-02"),
		})
	}

	return info
}

// Helper functions
func (n *NetworkReconModule) isKnownBadIP(ip string) bool {
	// Simplified check - in reality, you'd check against real threat feeds
	knownBadRanges := []string{
		"127.0.0.", // Localhost (for demo)
	}

	for _, badRange := range knownBadRanges {
		if strings.HasPrefix(ip, badRange) {
			return false // Don't flag localhost as malicious
		}
	}

	return false
}

func parseASN(asnStr string) (int, error) {
	// Simple ASN parsing - would be more robust in production
	var asn int
	_, err := fmt.Sscanf(asnStr, "%d", &asn)
	return asn, err
}
