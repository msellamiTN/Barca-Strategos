# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /phoenix-api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl
WORKDIR /root/

# Create directories for compliance reports and logs
RUN mkdir -p /app/reports /app/logs

# Copy binary and configuration
COPY --from=builder /phoenix-api .
COPY --from=builder /app/.env .env

# Copy compliance-specific scripts
COPY scripts/* /app/scripts/
RUN chmod +x /app/scripts/*.sh

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/system/health || exit 1

# Create non-root user for security
RUN addgroup -g 1000 -S phoenix && \
    adduser -u 1000 -S phoenix -G phoenix
USER phoenix

EXPOSE 8080

# Add labels for monitoring
LABEL maintainer="Barca Strategos Team"
LABEL version="1.0.0"
LABEL description="Phoenix API with SOC2 and PCI DSS Compliance Modules"
LABEL org.opencontainers.image.source="https://github.com/barca-strategos/phoenix"

CMD ["./phoenix-api"]
