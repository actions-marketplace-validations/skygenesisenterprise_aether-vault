# Build stage
FROM golang:1.25.5-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vault ./cmd/vault

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S vault && \
    adduser -u 1001 -S vault -G vault

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/vault .

# Copy configuration files
COPY --from=builder /app/config ./config

# Change ownership to vault user
RUN chown -R vault:vault /app

# Switch to non-root user
USER vault

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ./vault health || exit 1

# Run the application
CMD ["./vault"]