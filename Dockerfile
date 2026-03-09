# ─── Barca-Strategos Phoenix GUI ──────────────────────────────────────────────
# Completely offline build - no network calls required

# ─── Stage 1: Runtime only ─────────────────────────────────────────────────────
FROM debian:bookworm-slim

# Install minimal runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -s /bin/false phoenix

# Set working directory
WORKDIR /app

# Create directory structure
RUN mkdir -p /app/logs /app/data /app/config /app/static && \
    chown -R phoenix:phoenix /app

# Copy pre-built binaries (will be built locally first)
COPY target/release/phoenix-core   /usr/local/bin/phoenix-core
COPY target/release/phoenix-agent  /usr/local/bin/phoenix-agent
COPY target/release/phoenix-deploy /usr/local/bin/phoenix-deploy

# Copy static files
COPY static /app/static

# Switch to non-root user
USER phoenix

# Expose port
EXPOSE 8080

# Environment variables
ENV RUST_LOG=info
ENV PHOENIX_HOST=0.0.0.0
ENV PHOENIX_PORT=8080
ENV PHOENIX_STATIC_PATH=/app/static
ENV PHOENIX_MAX_CONNECTIONS=1000

# Volume mounts
VOLUME ["/app/logs", "/app/data", "/app/config"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD ["/usr/local/bin/phoenix-core", "--health-check"]

# Labels
LABEL maintainer="Barca-Strategos Team"
LABEL version="1.0.0"

# Start the application
CMD ["/usr/local/bin/phoenix-core"]