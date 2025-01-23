# threshAI

## Overview
threshAI is an innovative AI platform designed for building, running, and monitoring AI prompts. It provides tools for advanced prompt engineering, real-time performance monitoring, and an extensible plugin system to customize and enhance prompt processing workflows.

## Features
- **Advanced Prompt Engineering:** Tools and capabilities for designing, refining, and optimizing AI prompts for various tasks.
- **Real-time Monitoring:** Comprehensive monitoring of prompt execution, performance metrics, and system health.
- **Extensible Plugin System:** A flexible plugin architecture allowing users to extend and customize prompt processing with custom functionalities.
- **CLI and Web Interface:** User-friendly command-line interface (CLI) and web-based interface for interacting with the platform.
- **Robust Security Features:** Built-in security measures to ensure the safety and integrity of prompt processing and system operations.

## Installation Matrix

| Platform | Architecture | Requirements | Status |
|----------|-------------|--------------|--------|
| Linux    | amd64       | Go 1.21+     | ✅     |
| macOS    | amd64       | Go 1.21+     | ✅     |
| Windows  | amd64       | Go 1.21+     | ✅     |

### Prerequisites
- Go 1.21+ (required)
- Node.js 16+ (for frontend)
- Docker 20.10+ (for containerized deployment)
- golangci-lint (for development)
- Git (for version information)

### Quick Start

To get started with threshAI, follow these steps:

#### Prerequisites
Ensure you have the following installed:
- Go 1.21+
- Node.js 16+ (for frontend development)
- Docker 20.10+ (if you plan to use Docker)

#### Installation Steps
```bash
# Clone the repository
git clone https://github.com/[org]/threshAI.git
cd threshAI

# Verify prerequisites
make verify-prereqs

# Install dependencies and run security check
make deps security-check

# Build, test, and lint
make all

# Verify installation
make verify-install
```

## Build System

### Using Makefile
```bash
make help           # Show all available commands
make build         # Build the binary
make release       # Build cross-platform binaries
make run          # Run the web server
make run-cli      # Run the CLI
make test         # Run tests with coverage report
make lint         # Lint the code
make clean        # Clean up artifacts
make all          # Run full verification and build suite
```

### Build Artifact Verification
After building (especially for releases), verify the artifacts:

```bash
# For release builds
make release
cd dist/
sha256sum -c checksums.txt

# Verify binary health
./thresh --version
make verify-install
```

### Using Docker
```bash
# Build with security scanning
docker build --target security-check -t threshai-security .
docker build -t threshai .

# Verify container
docker run --rm threshai --version
docker inspect threshai | jq '.[].Config.Healthcheck'

# Run the application
docker-compose up --build
```

## Development

### Environment Setup
1. Create `.env` file:
```bash
cp .env.example .env
```

2. Configure environment variables:
```bash
# API Keys and Secrets
DEEPSEEK_API_KEY=your_key_here
```

### Security Checks
```bash
# Run vulnerability scanning
make security-check

# Run security audit
go list -json -m all | docker run -i sonatypescan/nancy:latest sleuth
```

### Testing
```bash
# Run all tests with coverage report
make test

# View coverage report in browser
open coverage.html

# Run specific test suite
go test ./pkg/...
```

## Troubleshooting Guide

### Common Issues

1. **Build Failures**
   - Ensure Go 1.21+ is installed: `go version`
   - Clear build cache: `go clean -cache`
   - Verify dependencies: `make deps`

2. **Docker Issues**
   - Verify Docker daemon: `docker info`
   - Clean Docker system: `make clean`
   - Check Docker logs: `docker logs <container-id>`

3. **Test Failures**
   - Update dependencies: `go mod tidy`
   - Clear test cache: `go clean -testcache`
   - Run with verbose output: `go test -v ./...`

4. **Permission Issues**
   - Docker socket access: `sudo usermod -aG docker $USER`
   - Binary permissions: `chmod +x bin/thresh`
   - Web port access: Use `sudo setcap CAP_NET_BIND_SERVICE=+eip bin/thresh`

### Getting Help

1. Check logs:
```bash
docker-compose logs
journalctl -u thresh.service
```

2. Debug mode:
```bash
make run-cli -- --debug
GO_ENV=development make run
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Run tests and security checks (`make all security-check`)
4. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
5. Push to the branch (`git push origin feature/AmazingFeature`)
6. Open a Pull Request

Please read our [Contribution Guidelines](CONTRIBUTING.md) for more details.

## Reporting Issues

Found a bug? Please [open an issue](https://github.com/[org]/threshAI/issues) and include:
- Detailed description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Screenshots (if applicable)
- Version information (`thresh --version`)
- Environment details (OS, Go version, Docker version)

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## License
This project is licensed under the [MIT License](LICENSE).
