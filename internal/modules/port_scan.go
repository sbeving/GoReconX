package modules

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PortScanModule implements network port scanning
type PortScanModule struct {
	*BaseModule
	semaphore chan bool
}

// PortScanResult represents port scan results
type PortScanResult struct {
	Target      string     `json:"target"`
	OpenPorts   []PortInfo `json:"open_ports"`
	ClosedPorts []int      `json:"closed_ports"`
	TotalPorts  int        `json:"total_ports"`
	ScanTime    string     `json:"scan_time"`
	ScanType    string     `json:"scan_type"`
}

// PortInfo contains information about an open port
type PortInfo struct {
	Port         int    `json:"port"`
	Protocol     string `json:"protocol"`
	State        string `json:"state"`
	Service      string `json:"service"`
	Version      string `json:"version"`
	Banner       string `json:"banner"`
	ResponseTime string `json:"response_time"`
}

// NewPortScanModule creates a new port scanning module
func NewPortScanModule() *PortScanModule {
	info := ModuleInfo{
		Name:        "port_scan",
		Category:    "active_recon",
		Description: "Network port scanner with service detection",
		Version:     "1.0.0",
		Author:      "GoReconX Team",
		Tags:        []string{"port", "scan", "network", "service", "active"},
		Options: []ModuleOption{
			{
				Name:        "ports",
				Type:        "string",
				Description: "Port range to scan (e.g., 1-1000, 80,443,22 or 'common')",
				Required:    false,
				Default:     "common",
			},
			{
				Name:        "threads",
				Type:        "int",
				Description: "Number of concurrent threads",
				Required:    false,
				Default:     100,
			},
			{
				Name:        "timeout",
				Type:        "int",
				Description: "Connection timeout in seconds",
				Required:    false,
				Default:     2,
			},
			{
				Name:        "service_detection",
				Type:        "bool",
				Description: "Enable service version detection",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "banner_grab",
				Type:        "bool",
				Description: "Enable banner grabbing",
				Required:    false,
				Default:     true,
			},
			{
				Name:        "scan_type",
				Type:        "choice",
				Description: "Type of scan to perform",
				Required:    false,
				Default:     "tcp_connect",
				Choices:     []string{"tcp_connect", "syn_scan", "udp_scan"},
			},
		},
		Requirements: []string{"network"},
	}

	return &PortScanModule{
		BaseModule: NewBaseModule(info),
		semaphore:  make(chan bool, 100), // Default thread limit
	}
}

// Validate validates the module input
func (p *PortScanModule) Validate(input ModuleInput) error {
	if err := p.ValidateInput(input); err != nil {
		return err
	}

	// Validate target is a valid IP or hostname
	if net.ParseIP(input.Target) == nil {
		if _, err := net.LookupHost(input.Target); err != nil {
			return NewModuleError("invalid target: must be valid IP or hostname", "INVALID_TARGET")
		}
	}

	return nil
}

// Execute runs the port scanning module
func (p *PortScanModule) Execute(ctx context.Context, input ModuleInput, output chan<- ModuleResult) error {
	startTime := time.Now()
	p.SetStatus("running", 0.0, "Starting port scan")

	// Parse options
	threads, _ := input.Options["threads"].(int)
	if threads <= 0 {
		threads = 100
	}
	p.semaphore = make(chan bool, threads)

	timeout, _ := input.Options["timeout"].(int)
	if timeout <= 0 {
		timeout = 2
	}

	serviceDetection, _ := input.Options["service_detection"].(bool)
	bannerGrab, _ := input.Options["banner_grab"].(bool)
	scanType, _ := input.Options["scan_type"].(string)
	if scanType == "" {
		scanType = "tcp_connect"
	}

	// Parse ports
	portsOption, _ := input.Options["ports"].(string)
	if portsOption == "" {
		portsOption = "common"
	}

	ports := p.parsePorts(portsOption)
	if len(ports) == 0 {
		return NewModuleError("no valid ports to scan", "NO_PORTS")
	}

	p.SendResult(output, "progress", fmt.Sprintf("Scanning %d ports on %s", len(ports), input.Target), nil, input.SessionID)

	result := &PortScanResult{
		Target:      input.Target,
		OpenPorts:   []PortInfo{},
		ClosedPorts: []int{},
		TotalPorts:  len(ports),
		ScanType:    scanType,
	}

	// Scan ports concurrently
	var wg sync.WaitGroup
	var mutex sync.Mutex
	scannedCount := 0

	for _, port := range ports {
		if p.IsStopped() {
			break
		}

		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			defer func() {
				mutex.Lock()
				scannedCount++
				progress := float64(scannedCount) / float64(len(ports))
				p.SetStatus("running", progress, fmt.Sprintf("Scanned %d/%d ports", scannedCount, len(ports)))
				mutex.Unlock()
			}()

			p.semaphore <- true              // Acquire semaphore
			defer func() { <-p.semaphore }() // Release semaphore

			if p.scanPort(input.Target, port, time.Duration(timeout)*time.Second) {
				portInfo := PortInfo{
					Port:     port,
					Protocol: "tcp",
					State:    "open",
					Service:  getServiceName(port),
				}

				if serviceDetection {
					portInfo.Version = p.detectService(input.Target, port)
				}

				if bannerGrab {
					portInfo.Banner = p.grabBanner(input.Target, port)
				}

				mutex.Lock()
				result.OpenPorts = append(result.OpenPorts, portInfo)
				mutex.Unlock()

				// Send individual port result
				p.SendResult(output, "data", map[string]interface{}{
					"type": "open_port",
					"port": portInfo,
				}, nil, input.SessionID)
			} else {
				mutex.Lock()
				result.ClosedPorts = append(result.ClosedPorts, port)
				mutex.Unlock()
			}
		}(port)
	}

	wg.Wait()

	// Sort results
	sort.Slice(result.OpenPorts, func(i, j int) bool {
		return result.OpenPorts[i].Port < result.OpenPorts[j].Port
	})
	sort.Ints(result.ClosedPorts)

	result.ScanTime = time.Since(startTime).String()

	// Send final result
	p.SetStatus("completed", 1.0, fmt.Sprintf("Scan completed: %d open ports found", len(result.OpenPorts)))
	p.SendResult(output, "complete", result, map[string]interface{}{
		"open_ports": len(result.OpenPorts),
		"scan_time":  result.ScanTime,
	}, input.SessionID)

	return nil
}

// scanPort scans a single port
func (p *PortScanModule) scanPort(target string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// detectService attempts to detect the service running on a port
func (p *PortScanModule) detectService(target string, port int) string {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	// Send HTTP request for web services
	if port == 80 || port == 443 || port == 8080 || port == 8443 {
		conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			if strings.Contains(response, "Server:") {
				lines := strings.Split(response, "\n")
				for _, line := range lines {
					if strings.HasPrefix(strings.ToLower(line), "server:") {
						return strings.TrimSpace(line[7:])
					}
				}
			}
			return "HTTP"
		}
	}

	// Try SSH detection
	if port == 22 {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			if strings.Contains(response, "SSH") {
				return strings.TrimSpace(response)
			}
		}
	}

	// Try FTP detection
	if port == 21 {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			if strings.Contains(response, "FTP") {
				return strings.TrimSpace(response)
			}
		}
	}

	return getServiceName(port)
}

// grabBanner attempts to grab service banner
func (p *PortScanModule) grabBanner(target string, port int) string {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return ""
	}
	defer conn.Close()

	// Read any immediate response
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buffer)
	if err == nil && n > 0 {
		return strings.TrimSpace(string(buffer[:n]))
	}

	// For HTTP ports, send a request
	if port == 80 || port == 443 || port == 8080 || port == 8443 {
		conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response := string(buffer[:n])
			lines := strings.Split(response, "\n")
			if len(lines) > 0 {
				return strings.TrimSpace(lines[0])
			}
		}
	}

	return ""
}

// parsePorts parses port specification into a list of ports
func (p *PortScanModule) parsePorts(portsSpec string) []int {
	var ports []int

	switch strings.ToLower(portsSpec) {
	case "common":
		// Most common ports
		commonPorts := []int{
			21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 993, 995,
			1723, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 27017,
		}
		return commonPorts
	case "all":
		// All ports 1-65535 (be careful!)
		for i := 1; i <= 65535; i++ {
			ports = append(ports, i)
		}
		return ports
	case "well-known":
		// Well-known ports 1-1023
		for i := 1; i <= 1023; i++ {
			ports = append(ports, i)
		}
		return ports
	}

	// Parse custom port specification
	parts := strings.Split(portsSpec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			// Range specification
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))

				if err1 == nil && err2 == nil && start <= end && start > 0 && end <= 65535 {
					for i := start; i <= end; i++ {
						ports = append(ports, i)
					}
				}
			}
		} else {
			// Single port
			if port, err := strconv.Atoi(part); err == nil && port > 0 && port <= 65535 {
				ports = append(ports, port)
			}
		}
	}

	return ports
}

// getServiceName returns the common service name for a port
func getServiceName(port int) string {
	serviceMap := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		111:   "RPC",
		135:   "MS-RPC",
		139:   "NetBIOS-SSN",
		143:   "IMAP",
		443:   "HTTPS",
		445:   "SMB",
		993:   "IMAPS",
		995:   "POP3S",
		1433:  "MSSQL",
		1521:  "Oracle",
		1723:  "PPTP",
		3306:  "MySQL",
		3389:  "RDP",
		5432:  "PostgreSQL",
		5900:  "VNC",
		6379:  "Redis",
		8080:  "HTTP-Proxy",
		8443:  "HTTPS-Alt",
		27017: "MongoDB",
	}

	if service, exists := serviceMap[port]; exists {
		return service
	}
	return "unknown"
}
