# Multi-stage build for orchestrating web and server components
FROM alpine:latest AS orchestrator

# Install tools for container orchestration
RUN apk add --no-cache docker-cli curl wget

# Create non-root user
RUN addgroup -g 1001 -S aether && \
    adduser -u 1001 -S aether -G aether

# Set working directory
WORKDIR /app

# Copy docker-compose files from subdirectories
COPY web/docker-compose.yml ./web-compose.yml
COPY server/docker-compose.yml ./server-compose.yml

# Create orchestration script
COPY <<'EOF' ./orchestrate.sh
#!/bin/sh

set -e

echo "🚀 Starting Aether Vault CLI Orchestration..."

# Build and start server components
echo "📦 Building server components..."
docker-compose -f server-compose.yml build

echo "🔧 Starting server infrastructure..."
docker-compose -f server-compose.yml up -d postgres redis

echo "⏳ Waiting for database to be ready..."
sleep 10

echo "🚀 Starting server application..."
docker-compose -f server-compose.yml up -d aether-cli-server

# Build and start web components
echo "📦 Building web components..."
docker-compose -f web-compose.yml build

echo "🚀 Starting web application..."
docker-compose -f web-compose.yml up -d

echo "✅ All services started successfully!"
echo ""
echo "🌐 Web Interface: http://localhost:3000"
echo "🔧 API Server: http://localhost:8080"
echo "🐘 PostgreSQL: localhost:5432"
echo "🔴 Redis: localhost:6379"
echo ""
echo "📋 To view logs:"
echo "  Server: docker-compose -f server-compose.yml logs -f"
echo "  Web: docker-compose -f web-compose.yml logs -f"
echo ""
echo "🛑 To stop all services:"
echo "  docker-compose -f server-compose.yml down"
echo "  docker-compose -f web-compose.yml down"

# Keep container running for monitoring
tail -f /dev/null
EOF

# Make script executable
RUN chmod +x orchestrate.sh

# Switch to non-root user
USER aether

# Expose management ports
EXPOSE 3000 8080 5432 6379

# Health check for orchestration
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/system/health && \
       curl -f http://localhost:3000 || exit 1

# Start orchestration
CMD ["./orchestrate.sh"]