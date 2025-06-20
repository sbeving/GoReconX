# GoReconX - Comprehensive OSINT & Reconnaissance Platform

<p align="center">
  <img src="icon.png" alt="GoReconX Logo" width="200"/>
</p>

<p align="center">
  <strong>A powerful, user-friendly, and highly modular Open-Source Intelligence (OSINT) and network reconnaissance toolkit written entirely in Golang.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version"/>
  <img src="https://img.shields.io/badge/License-MIT-green" alt="License"/>
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey" alt="Platform"/>
  <img src="https://img.shields.io/badge/AI%20Powered-Google%20Gemini-blue" alt="AI Powered"/>
</p>

## 🎯 Overview

GoReconX provides cybersecurity professionals, red teamers, and enthusiasts with a centralized, visually appealing, and intuitive tool that supports both passive and active reconnaissance. The platform is enhanced with AI-powered analysis via the Google Gemini API, offering smart insights and professional reporting capabilities.

## ✨ Features

### 🔍 Core Reconnaissance Modules

#### Passive OSINT
- **Subdomain Enumeration**: Advanced DNS-based subdomain discovery with wordlist support
- **Email Harvesting**: Collect email addresses from various public sources
- **Website Analysis**: Analyze web technologies, headers, and content
- **IP Geolocation**: Determine geographical location and ASN information
- **GitHub Reconnaissance**: Search for sensitive information in public repositories

#### Active Reconnaissance
- **Port Scanning**: Fast TCP/UDP port scanning with service detection
- **Directory Enumeration**: Discover hidden directories and files on web servers
- **Service Detection**: Identify running services and their versions

### 🤖 AI-Powered Analysis
- **Smart Summarization**: AI-generated executive summaries of findings
- **Threat Assessment**: Automated threat level classification
- **Security Recommendations**: Actionable security advice based on results
- **Natural Language Insights**: Easy-to-understand analysis of technical data

### 📊 Professional Reporting
- **Multiple Formats**: Export reports in JSON, HTML, PDF, and CSV formats
- **Executive Summaries**: AI-enhanced summaries for management presentations
- **Detailed Technical Reports**: Comprehensive findings for technical teams
- **Custom Branding**: Professional report templates with your organization's branding

### 🛡️ Security & Ethics
- **Ethical Usage Disclaimer**: Prominent warnings and usage agreements
- **Encrypted Storage**: Secure storage of API keys and sensitive data
- **Audit Logging**: Comprehensive logging of all activities
- **Rate Limiting**: Built-in protections against API abuse

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Git
- X11 development libraries (for GUI on Linux)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-org/goreconx.git
   cd goreconx
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Install system dependencies (Linux):**
   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install libgl1-mesa-dev xorg-dev

   # CentOS/RHEL/Fedora
   sudo dnf install mesa-libGL-devel libX11-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel
   ```

4. **Build the application:**
   ```bash
   go build -o goreconx ./cmd/main.go
   ```

5. **Run GoReconX:**
   ```bash
   ./goreconx
   ```

### Configuration

1. **Accept the ethical usage disclaimer** when prompted
2. **Configure API keys** in the Settings tab:
   - Google Gemini API key (for AI features)
   - VirusTotal API key (optional)
   - Shodan API key (optional)
   - Hunter.io API key (optional)
3. **Set up wordlists** or use the built-in defaults
4. **Configure output preferences** and proxy settings if needed

## 📖 Usage Guide

### Basic Workflow

1. **Create a New Project**: Organize your reconnaissance activities
2. **Select Target**: Enter domain, IP address, or organization name
3. **Choose Modules**: Select appropriate passive or active reconnaissance modules
4. **Configure Options**: Adjust threads, timeouts, and module-specific settings
5. **Run Scans**: Execute reconnaissance modules individually or in batch
6. **Analyze Results**: Review findings in the structured results viewer
7. **Generate Reports**: Create professional reports with AI-enhanced insights

### Module Details

#### Subdomain Enumeration
```
Target: example.com
Options:
  - Wordlist: /path/to/subdomains.txt
  - Threads: 50
  - Timeout: 5 seconds
  - Resolve IPs: Yes
```

#### Port Scanning
```
Target: 192.168.1.1
Options:
  - Ports: 1-1000,3389,5432,8080-8090
  - Threads: 100
  - Timeout: 2 seconds
  - TCP/UDP: Both
```

### AI-Powered Analysis

When configured with a Google Gemini API key, GoReconX provides:

- **Intelligent Summaries**: Automatically generated overviews of findings
- **Risk Assessment**: AI-evaluated threat levels and confidence scores
- **Actionable Recommendations**: Specific steps to improve security posture
- **Natural Language Queries**: Ask questions about your results

### Report Generation

Generate professional reports with:

```go
// Example: Generate HTML report
reportGen := reports.NewReportGenerator(logger, aiClient, "output/")
report, _ := reportGen.GenerateReport("example.com", scanResults)
htmlFile, _ := reportGen.ExportHTML(report)
```

## 🏗️ Architecture

### Project Structure

```
GoReconX/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── ai/
│   │   └── gemini.go          # Google Gemini AI integration
│   ├── appinstance/
│   │   └── app.go             # Main application instance
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── database/
│   │   └── database.go        # SQLite database operations
│   ├── gui/
│   │   ├── main_window.go     # Main GUI window
│   │   ├── osint_tabs.go      # OSINT module tabs
│   │   └── utility_tabs.go    # Utility and settings tabs
│   ├── logging/
│   │   └── logger.go          # Logging configuration
│   ├── modules/
│   │   ├── manager.go         # Module management
│   │   ├── subdomain.go       # Subdomain enumeration
│   │   └── placeholder_modules.go # Other reconnaissance modules
│   └── reports/
│       └── generator.go       # Report generation
├── config/
│   └── config.yaml            # Application configuration
├── data/
│   └── goreconx.db           # SQLite database
├── logs/
│   └── goreconx.log          # Application logs
├── output/
│   └── reports/              # Generated reports
├── wordlists/
│   ├── subdomains.txt        # Subdomain wordlist
│   ├── directories.txt       # Directory wordlist
│   └── ports.txt             # Port list
├── go.mod
├── go.sum
└── README.md
```

### Module Interface

All reconnaissance modules implement the `ModuleInterface`:

```go
type ModuleInterface interface {
    GetName() string
    GetDescription() string
    Validate(target string) error
    Execute(target string, options map[string]interface{}) (*ScanResult, error)
    GetDefaultOptions() map[string]interface{}
}
```

### Adding Custom Modules

1. **Implement the ModuleInterface**
2. **Add to ModuleManager**
3. **Update GUI components**
4. **Add configuration options**

## 🔧 Configuration

### config.yaml

```yaml
database:
  path: "data/goreconx.db"

api:
  gemini_key: "your-gemini-api-key"
  virustotal_key: "your-virustotal-api-key"
  shodan_key: "your-shodan-api-key"
  hunter_key: "your-hunter-api-key"

network:
  timeout: 30
  retries: 3
  user_agent: "GoReconX/1.0 (OSINT Tool)"
  proxy_url: ""

wordlists:
  subdomains: "wordlists/subdomains.txt"
  directories: "wordlists/directories.txt"
  files: "wordlists/files.txt"
  ports: "wordlists/ports.txt"

output:
  default_format: "json"
  output_dir: "output"
```

### Environment Variables

```bash
export GORECONX_GEMINI_KEY="your-api-key"
export GORECONX_DB_PATH="/custom/path/to/db"
export GORECONX_OUTPUT_DIR="/custom/output/path"
```

## 🔒 Security Considerations

### Ethical Usage

GoReconX is designed for:
- ✅ Authorized penetration testing
- ✅ Bug bounty programs
- ✅ Security research
- ✅ Educational purposes
- ✅ Personal network testing

**Never use GoReconX for:**
- ❌ Unauthorized scanning
- ❌ Malicious activities
- ❌ Privacy violations
- ❌ Illegal reconnaissance

### Data Protection

- **API Key Encryption**: All API keys are encrypted before storage
- **Local Data Only**: No telemetry or data exfiltration
- **Secure Defaults**: Conservative timeout and rate limiting settings
- **Audit Trail**: Comprehensive logging of all activities

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Add tests**: Ensure your code is well-tested
5. **Commit changes**: `git commit -m 'Add amazing feature'`
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Code Style

- Follow Go conventions
- Use `gofmt` for formatting
- Add comprehensive comments
- Include unit tests for new features

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Documentation**: [docs.goreconx.io](https://docs.goreconx.io)
- **Issue Tracker**: [GitHub Issues](https://github.com/your-org/goreconx/issues)
- **Discord Community**: [Join our Discord](https://discord.gg/goreconx)
- **Twitter**: [@GoReconX](https://twitter.com/goreconx)

## 🙏 Acknowledgments

- **Fyne Framework**: For the excellent Go GUI framework
- **Google Gemini**: For AI-powered analysis capabilities
- **Go Community**: For the robust standard library and ecosystem
- **Security Community**: For wordlists, techniques, and best practices

## 📧 Support

- **Email**: support@goreconx.io
- **Documentation**: [docs.goreconx.io](https://docs.goreconx.io)
- **Community Forum**: [community.goreconx.io](https://community.goreconx.io)

---

<p align="center">
  <strong>Remember: Use GoReconX responsibly and ethically. Always obtain proper authorization before conducting any reconnaissance activities.</strong>
</p>

<p align="center">
  Made with ❤️ by the GoReconX Team
</p>
