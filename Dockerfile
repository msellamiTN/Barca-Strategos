# Multi-stage Dockerfile for Barca-Strategos Phoenix GUI
# Optimized for production deployment with web interface

# Stage 1: Build the Rust application
FROM rust:1.75-slim as builder

# Install system dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    libpq-dev \
    build-essential \
    curl-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy Cargo files
COPY Cargo.toml Cargo.lock ./

# Create dummy main.rs to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies (this layer is cached if Cargo.toml doesn't change)
RUN cargo build --release && rm -rf src

# Copy the actual source code
COPY src ./src/
COPY static ./static/

# Build the application
RUN cargo build --release

# Stage 2: Runtime image
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    libpq5 \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user for security
RUN useradd -r -s /bin/false phoenix

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/target/release/phoenix-core ./bin/phoenix-core
COPY --from=builder /app/target/release/start_web_gui ./bin/start_web_gui

# Copy static files for web interface
COPY --from=builder /app/static ./static

# Create necessary directories
RUN mkdir -p /app/logs /app/data /app/config && \
    chown -R phoenix:phoenix /app

# Switch to non-root user
USER phoenix

# Expose ports
EXPOSE 8080 8443

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/system/health || exit 1

# Environment variables
ENV RUST_LOG=info
ENV PHOENIX_HOST=0.0.0.0
ENV PHOENIX_PORT=8080
ENV PHOENIX_STATIC_PATH=/app/static
ENV PHOENIX_MAX_CONNECTIONS=1000

# Volume mounts for persistent data
VOLUME ["/app/logs", "/app/data", "/app/config"]

# Labels for metadata
LABEL maintainer="Barca-Strategos Team"
LABEL version="1.0.0"
LABEL description="Barca-Strategos Phoenix - Cognitive Collaboration Platform with Web GUI"
LABEL org.opencontainers.image.title="Phoenix GUI"
LABEL org.opencontainers.image.description="Web-based cognitive collaboration platform"
LABEL org.opencontainers.image.vendor="Barca-Strategos"
LABEL org.opencontainers.image.licenses="MIT"

# Start the web GUI application
CMD ["./bin/start_web_gui"]
