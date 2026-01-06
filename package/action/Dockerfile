# Multi-stage build for Aether Vault Action
# Stage 1: Build stage
FROM golang:1.25.5-alpine AS builder

# Set the working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary for multiple platforms
ARG TARGETPLATFORM
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Build the binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bin/aether-vault-${TARGETOS}-${TARGETARCH} \
    ./cmd

# Stage 2: Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    jq \
    && rm -rf /var/cache/apk/*

# Create app user
RUN addgroup -g 1001 -S aether && \
    adduser -u 1001 -S aether -G aether

# Set the working directory
WORKDIR /app

# Copy binary from builder stage
ARG TARGETOS=linux
ARG TARGETARCH=amd64
COPY --from=builder /app/bin/aether-vault-${TARGETOS}-${TARGETARCH} /usr/local/bin/aether-vault

# Make binary executable
RUN chmod +x /usr/local/bin/aether-vault

# Create necessary directories
RUN mkdir -p /github/workspace && \
    chown -R aether:aether /github

# Switch to non-root user
USER aether

# Set environment variables
ENV LOG_LEVEL=info
ENV VAULT_URL=""
ENV AUTH_METHOD=github-oidc
ENV ROLE=""
ENV POLICY_MODE=enforce
ENV AUDIENCE=aether-vault
ENV ALLOW_TOKEN_OUTPUT=false

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD aether-vault --help || exit 1

# Expose port (if needed for future web interface)
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["aether-vault"]

# Default command
CMD ["--help"]