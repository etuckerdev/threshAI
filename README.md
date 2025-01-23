# threshAI

## Overview
threshAI is a cutting-edge AI platform for [brief description of what the project does].

## Installation

### Prerequisites
- Go 1.20+
- Node.js 16+ (for frontend)
- Docker (optional)

### Quick Start
```bash
# Clone the repository
git clone https://github.com/[org]/threshAI.git
cd threshAI

# Install dependencies
make deps

# Build, test, and lint
make all

# Run the CLI
make run-cli
```

### Build & Run

#### Using Makefile
```bash
make build   # Build the binary
make run     # Run the web server
make run-cli # Run the CLI
make test    # Run tests with coverage
make lint    # Lint the code
make clean   # Clean up artifacts
make all     # Run deps, build, test, and lint
```

#### Using Docker
```bash
# Build and run the CLI
docker build -t threshai-cli -f Dockerfile.cli .
docker run threshai-cli -provider deepseek -prompt "Hello world"

# Build and run the web server
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

### Running Tests
```bash
# Run all tests
make test

# Run specific test suite
go test ./pkg/...
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please read our [Contribution Guidelines](CONTRIBUTING.md) for more details.

## Reporting Issues

Found a bug? Please [open an issue](https://github.com/[org]/threshAI/issues) and include:
- Detailed description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Screenshots (if applicable)

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## License
This project is licensed under the [MIT License](LICENSE).
