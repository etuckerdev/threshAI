# syntax=docker/dockerfile:1.4
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with security flags
ARG VERSION="dev"
ARG COMMIT="unknown"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o thresh ./cmd/web

# Security scan stage
FROM golang:1.21-alpine AS security-check
RUN apk add --no-cache git
WORKDIR /scan
COPY --from=builder /build .
RUN go install golang.org/x/vuln/cmd/govulncheck@latest && \
    govulncheck ./...

# Final stage
FROM alpine:3.19 AS final

# Add non-root user
RUN adduser -D -H -s /bin/false appuser && \
    apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary and set permissions
COPY --from=builder /build/thresh .
RUN chown appuser:appuser /app/thresh && \
    chmod 500 /app/thresh

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV TZ=UTC \
    GO_ENV=production

# Run the application
ENTRYPOINT ["/app/thresh"]

# Labels for container metadata
LABEL org.opencontainers.image.title="ThreshAI" \
    org.opencontainers.image.description="ThreshAI Server" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.revision="${COMMIT}" \
    org.opencontainers.image.vendor="ThreshAI" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.created="${BUILD_TIME}"