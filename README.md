# GoReconX - Advanced OSINT & Reconnaissance Platform

<div align="center">

![GoReconX Logo](assets/logo.png)

**A Comprehensive Open-Source Intelligence & Network Reconnaissance Toolkit**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Security](https://img.shields.io/badge/Security-Focused-red?style=for-the-badge)](SECURITY.md)

</div>

## 🚀 Overview

GoReconX is a powerful, modular, and user-friendly Open-Source Intelligence (OSINT) and network reconnaissance toolkit designed for cybersecurity professionals, red teamers, and security enthusiasts. Built with Go for exceptional performance and featuring a modern web-based GUI, GoReconX provides a centralized platform for ethical information gathering.

### ✨ Key Features

- 🎯 **Modular Architecture** - Easily extensible with plug-and-play reconnaissance modules
- 🖥️ **Modern GUI** - Intuitive web-based interface with real-time updates
- 🔒 **Security First** - Encrypted data storage and secure API key management
- ⚡ **High Performance** - Built with Go for speed and concurrency
- 🌐 **Multi-Platform** - Cross-platform support (Windows, Linux, macOS)
- 📊 **Advanced Reporting** - Comprehensive reports in multiple formats
- 🔌 **API Integration** - Support for popular OSINT APIs and services

## 🏗️ Architecture

```
GoReconX/
├── cmd/                    # CLI entry points
├── internal/               # Private application code
│   ├── api/               # REST API handlers
│   ├── core/              # Core business logic
│   ├── database/          # Database layer
│   ├── modules/           # Reconnaissance modules
│   ├── gui/               # Web GUI components
│   └── utils/             # Utility functions
├── pkg/                   # Public library code
├── web/                   # Frontend assets
├── configs/               # Configuration files
├── docs/                  # Documentation
└── scripts/               # Build and deployment scripts
```

## 🔧 Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/GoReconX.git
cd GoReconX

# Build the application
go mod tidy
go build -o bin/gorconx cmd/main.go

# Run GoReconX
./bin/gorconx
```

### Using Docker (Recommended)

```bash
# Run with Docker
docker run -p 8080:8080 gorconx/gorconx:latest
```

## 🎯 Modules

### Passive OSINT Modules
- **Domain & Subdomain Enumeration** - WHOIS, DNS records, Certificate Transparency
- **Email & People Search** - Email harvesting, social media profiling
- **Website Analysis** - Technology detection, metadata extraction
- **IP Geolocation** - GeoIP lookup, ASN information
- **Code Repository Recon** - GitHub/GitLab intelligence gathering

### Active Reconnaissance Modules
- **Port Scanning** - TCP/UDP port discovery and service detection
- **Directory Enumeration** - Web directory and file discovery
- **Vulnerability Assessment** - Basic security checks

### Utility Modules
- **Results Management** - Advanced filtering and search
- **Report Generation** - HTML, PDF, JSON exports
- **API Management** - Secure credential storage
- **Session Management** - Save and restore scan sessions

## 🚦 Usage

### Web Interface
1. Start GoReconX: `./bin/gorconx`
2. Open your browser to `http://localhost:8080`
3. Configure your API keys in Settings
4. Select your reconnaissance modules
5. Start your intelligence gathering

### Command Line Interface
```bash
# Domain enumeration
gorconx domain -t example.com

# Port scanning
gorconx portscan -t 192.168.1.1 -p 1-1000

# Generate report
gorconx report -f html -o /path/to/report.html
```

## ⚖️ Legal Disclaimer

**IMPORTANT**: GoReconX is designed for ethical hacking, educational purposes, and legitimate security assessments only. Users must:

- ✅ Have explicit written permission before scanning any target
- ✅ Comply with local laws and regulations
- ✅ Use the tool responsibly and ethically
- ❌ Never use for malicious purposes
- ❌ Never scan systems without authorization

## 🔒 Security Features

- **Encrypted Storage** - All sensitive data encrypted at rest
- **API Key Protection** - Secure credential management
- **No Data Exfiltration** - All data stays local unless explicitly exported
- **Audit Logging** - Comprehensive activity logging
- **Rate Limiting** - Built-in API rate limiting

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup
```bash
# Clone and setup development environment
git clone https://github.com/yourusername/GoReconX.git
cd GoReconX
go mod tidy

# Run in development mode
go run cmd/main.go --dev
```

## 📚 Documentation

- [User Guide](docs/user-guide.md)
- [API Documentation](docs/api.md)
- [Module Development](docs/module-development.md)
- [Configuration Guide](docs/configuration.md)

## 🏆 Why GoReconX?

- **Performance**: Built with Go for exceptional speed and concurrency
- **User Experience**: Modern, intuitive interface designed for professionals
- **Extensibility**: Modular architecture allows easy customization
- **Security**: Security-first design with encrypted data storage
- **Community**: Open-source with active community support

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- The Go community for excellent libraries and tools
- OSINT practitioners and security researchers
- Contributors and beta testers

---

<div align="center">

**Made with ❤️ by the cybersecurity community**

[Website](https://gorconx.com) • [Documentation](docs/) • [Community](https://github.com/yourusername/GoReconX/discussions)

</div>
