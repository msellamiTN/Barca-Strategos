# Multi-stage Dockerfile for Barca-Strategos Phoenix GUI
# Using Alpine Linux to avoid package installation issues

# Stage 1: Build the Rust application
FROM rust:1.75-alpine as builder

# Install system dependencies (Alpine packages are more reliable)
RUN apk add --no-cache \
    pkgconfig \
    openssl-dev \
    postgresql-dev \
    musl-dev \
    curl \
    build-base

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
FROM alpine:latest

# Install runtime dependencies only
RUN apk add --no-cache \
    ca-certificates \
    openssl \
    postgresql \
    curl

# Create non-root user for security
RUN adduser -D -s /bin/sh phoenix

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/target/release/phoenix-core ./bin/phoenix-core

# Copy static files for web interface
COPY --from=builder /app/static ./static

# Create necessary directories
RUN mkdir -p /app/logs /app/data /app/config && \
    chown -R phoenix:phoenix /app

# Switch to non-root user
USER phoenix

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD curl -f http://localhost:8080/api/system/health || exit 1

# Expose port
EXPOSE 8080

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
LABEL description="Barca-Strategos Phoenix GUI - Cognitive Collaboration Platform"
LABEL version="1.0.0"
LABEL org.opencontainers.image.title="Phoenix GUI"
LABEL org.opencontainers.image.description="Web-based cognitive collaboration platform"
LABEL org.opencontainers.image.vendor="Barca-Strategos"
LABEL org.opencontainers.image.licenses="MIT"

# Start the application
CMD ["./bin/phoenix-core"]
