# ─── Barca-Strategos Phoenix GUI ──────────────────────────────────────────────
# Standard online build with Rust 1.85 and pure Rust dependencies

# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM rust:1.85 AS builder

WORKDIR /app

# Copy manifests first for better caching
COPY Cargo.toml Cargo.lock ./

# Create dummy binaries to cache dependency compile layer
RUN mkdir -p src/bin && \
    echo "fn main() {}" > src/bin/phoenix-core.rs && \
    echo "fn main() {}" > src/bin/phoenix-agent.rs && \
    echo "fn main() {}" > src/bin/phoenix-deploy.rs

# Build dependencies (will use rustls, no native-tls)
RUN cargo build --release
RUN rm -rf src

# Build real source
COPY src ./src/
COPY static ./static/
RUN touch src/bin/phoenix-core.rs src/bin/phoenix-agent.rs src/bin/phoenix-deploy.rs
RUN cargo build --release

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM debian:bookworm-slim

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/target/release/phoenix-core   /usr/local/bin/phoenix-core
COPY --from=builder /app/target/release/phoenix-agent  /usr/local/bin/phoenix-agent
COPY --from=builder /app/target/release/phoenix-deploy /usr/local/bin/phoenix-deploy
COPY --from=builder /app/static /app/static

RUN useradd -r -s /bin/false phoenix && \
    mkdir -p /app/logs /app/data /app/config && \
    chown -R phoenix:phoenix /app /usr/local/bin/phoenix-*

USER phoenix
EXPOSE 8080

ENV RUST_LOG=info
ENV PHOENIX_HOST=0.0.0.0
ENV PHOENIX_PORT=8080
ENV PHOENIX_STATIC_PATH=/app/static
ENV PHOENIX_MAX_CONNECTIONS=1000

VOLUME ["/app/logs", "/app/data", "/app/config"]

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD ["/usr/local/bin/phoenix-core", "--health-check"]

LABEL maintainer="Barca-Strategos Team"
LABEL version="1.0.0"

CMD ["/usr/local/bin/phoenix-core"]