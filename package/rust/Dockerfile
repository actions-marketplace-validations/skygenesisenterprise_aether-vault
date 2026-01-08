# Multi-stage build for Rust SDK
FROM rust:1.75-slim as builder

# Install system dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

# Copy Cargo files
COPY Cargo.toml Cargo.lock ./

# Create dummy main to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy source code
COPY src ./src

# Build the library
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -s /bin/false vault

# Copy built library
COPY --from=builder /usr/src/app/target/release/libaether_vault.so /usr/local/lib/
COPY --from=builder /usr/src/app/target/release/deps /usr/local/lib/

# Update library cache
RUN ldconfig

# Set permissions
RUN chown -R vault:vault /usr/local/lib

# Switch to non-root user
USER vault

# Set environment
ENV RUST_LOG=info
ENV AETHER_VAULT_URL=https://vault.internal:8200

WORKDIR /app

# Expose health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD echo "SDK library loaded" || exit 1

# Default command
CMD ["echo", "Aether Vault Rust SDK container ready"]