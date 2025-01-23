# Troubleshooting Guide

## Overview
This guide provides solutions to common issues and troubleshooting steps for threshAI.

## Common Issues

### Build Failures
1. **Ensure Go 1.21+ is installed**:
   ```bash
   go version
   ```

2. **Clear build cache**:
   ```bash
   go clean -cache
   ```

3. **Verify dependencies**:
   ```bash
   make deps
   ```

### Docker Issues
1. **Verify Docker daemon**:
   ```bash
   docker info
   ```

2. **Clean Docker system**:
   ```bash
   make clean
   ```

3. **Check Docker logs**:
   ```bash
   docker logs <container-id>
   ```

### Test Failures
1. **Update dependencies**:
   ```bash
   go mod tidy
   ```

2. **Clear test cache**:
   ```bash
   go clean -testcache
   ```

3. **Run with verbose output**:
   ```bash
   go test -v ./...
   ```

### Permission Issues
1. **Docker socket access**:
   ```bash
   sudo usermod -aG docker $USER
   ```

2. **Binary permissions**:
   ```bash
   chmod +x bin/thresh
   ```

3. **Web port access**:
   ```bash
   sudo setcap CAP_NET_BIND_SERVICE=+eip bin/thresh
   ```

## Getting Help

### Check Logs
1. **Docker Compose Logs**:
   ```bash
   docker-compose logs
   ```

2. **Systemd Journal**:
   ```bash
   journalctl -u thresh.service
   ```

### Debug Mode
1. **Run CLI in Debug Mode**:
   ```bash
   make run-cli -- --debug
   ```

2. **Run Web Server in Development Mode**:
   ```bash
   GO_ENV=development make run
   ```

## Conclusion
Follow these troubleshooting steps to resolve common issues with threshAI. If you encounter problems not covered here, please [open an issue](https://github.com/[org]/threshAI/issues) and include detailed information about the problem.