# Multi-stage build for Phoenix Core - Ultra-efficient AI Security Framework
# Where Tactical Excellence Meets Technological Revolution

# Build stage
FROM rust:1.75-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    musl-dev \
    pkgconfig \
    openssl-dev \
    curl-dev

# Set working directory
WORKDIR /app

# Copy Cargo files
COPY Cargo.toml Cargo.lock ./

# Create dummy main.rs to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies
RUN cargo build --release && rm -rf src

# Copy source code
COPY src ./src

# Build the application with optimizations
RUN cargo build --release

# Runtime stage - Ultra-minimal
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    curl \
    ca-certificates \
    openssl \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 phoenix && \
    adduser -D -s /bin/sh -u 1000 -G phoenix phoenix

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/target/release/barca-strategos /usr/local/bin/

# Copy configuration files
COPY config ./config
COPY nginx ./nginx

# Create data directories
RUN mkdir -p /app/data /app/logs /app/security && \
    chown -R phoenix:phoenix /app

# Switch to non-root user
USER phoenix

# Expose ports
EXPOSE 8080 8443 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables
ENV RUST_LOG=info
ENV RUST_BACKTRACE=1

# Run the application
CMD ["barca-strategos", "--config", "/app/config", "--data-dir", "/app/data"]
