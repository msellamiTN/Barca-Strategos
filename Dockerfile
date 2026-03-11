# Minimal Dockerfile for restricted network environments
# Uses busybox as base - no package installation required
FROM busybox:latest

# Create necessary directories
RUN mkdir -p /app/reports /app/logs /app/scripts

# Copy static binary
COPY phoenix-api-static /app/phoenix-api
RUN chmod +x /app/phoenix-api

# Copy scripts
COPY scripts/*.sh /app/scripts/
RUN chmod +x /app/scripts/*.sh 2>/dev/null || true

# Create .env file
RUN echo 'PHOENIX_ENV=development' > /app/.env && \
    echo 'PHOENIX_PORT=8080' >> /app/.env

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 phoenix && \
    adduser -D -u 1000 -G phoenix phoenix && \
    chown -R phoenix:phoenix /app

USER phoenix

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD /app/phoenix-api --health-check || exit 1

# Labels
LABEL maintainer="Barca Strategos Team"
LABEL version="1.0.0"
LABEL description="Phoenix API with SOC2 and PCI DSS Compliance Modules"
LABEL org.opencontainers.image.source="https://github.com/barca-strategos/phoenix"

CMD ["/app/phoenix-api"]
