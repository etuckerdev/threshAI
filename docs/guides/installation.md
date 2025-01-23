# Installation Guide

## Overview
This guide provides step-by-step instructions for installing and setting up threshAI on your system.

## Prerequisites
Ensure you have the following installed:
- Go 1.21+
- Node.js 16+ (for frontend development)
- Docker 20.10+ (if you plan to use Docker)

## Installation Steps

### Step 1: Clone the Repository
```bash
git clone https://github.com/[org]/threshAI.git
cd threshAI
```

### Step 2: Verify Prerequisites
Run the following command to verify that all prerequisites are installed:
```bash
make verify-prereqs
```

### Step 3: Install Dependencies and Run Security Check
```bash
make deps security-check
```

### Step 4: Build, Test, and Lint
```bash
make all
```

### Step 5: Verify Installation
```bash
make verify-install
```

## Using Docker

### Step 1: Build with Security Scanning
```bash
docker build --target security-check -t threshai-security .
docker build -t threshai .
```

### Step 2: Verify Container
```bash
docker run --rm threshai --version
docker inspect threshai | jq '.[].Config.Healthcheck'
```

### Step 3: Run the Application
```bash
docker-compose up --build
```

## Environment Setup

### Step 1: Create `.env` File
```bash
cp .env.example .env
```

### Step 2: Configure Environment Variables
Edit the `.env` file to configure environment variables:
```bash
# API Keys and Secrets
DEEPSEEK_API_KEY=your_key_here
```

## Security Checks

### Step 1: Run Vulnerability Scanning
```bash
make security-check
```

### Step 2: Run Security Audit
```bash
go list -json -m all | docker run -i sonatypescan/nancy:latest sleuth
```

## Testing

### Step 1: Run All Tests with Coverage Report
```bash
make test
```

### Step 2: View Coverage Report in Browser
```bash
open coverage.html
```

### Step 3: Run Specific Test Suite
```bash
go test ./pkg/...
```

## Conclusion
Follow these steps to successfully install and set up threshAI on your system. Ensure you have all prerequisites installed and configure environment variables as needed.