# Contributing to threshAI

We welcome contributions from the community! Here are the guidelines to help you get started.

## Getting Started

1. **Fork the Repository**: Start by forking the threshAI repository to your GitHub account.

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/[your-username]/threshAI.git
   cd threshAI
   ```

3. **Set Up the Environment**:
   - Ensure you have Go 1.21+ installed.
   - Install Node.js 16+ for frontend development.
   - Install Docker 20.10+ for containerized deployment.
   - Install golangci-lint for development.
   - Install Git for version control.

4. **Install Dependencies**:
   ```bash
   make deps
   ```

5. **Run Security Checks**:
   ```bash
   make security-check
   ```

6. **Build and Test**:
   ```bash
   make all
   ```

## Making Changes

1. **Create a Feature Branch**:
   ```bash
   git checkout -b feature/AmazingFeature
   ```

2. **Make Your Changes**: Implement your feature or fix.

3. **Run Tests and Security Checks**:
   ```bash
   make all security-check
   ```

4. **Commit Your Changes**:
   ```bash
   git commit -m 'Add some AmazingFeature'
   ```

5. **Push to the Branch**:
   ```bash
   git push origin feature/AmazingFeature
   ```

6. **Open a Pull Request**: Go to the threshAI repository on GitHub and open a pull request from your feature branch.

## Code Style

- Follow the existing code style and conventions.
- Write clear and concise comments.
- Ensure your code is well-documented.

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