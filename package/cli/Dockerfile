# Build stage
FROM golang:1.25.5-alpine AS builder

# Install git (required for go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aether-vault-cli .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S aether && \
    adduser -u 1001 -S aether -G aether

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/aether-vault-cli .

# Change ownership to non-root user
RUN chown aether:aether /app/aether-vault-cli

# Switch to non-root user
USER aether

# Expose port (if CLI needs to serve anything)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./aether-vault-cli"]

# Default command
CMD ["--help"]