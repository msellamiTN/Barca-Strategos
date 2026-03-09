# ─── Barca-Strategos Phoenix GUI ──────────────────────────────────────────────
# Zero network calls during build.
# Requires: Cargo.toml with rustls (no native-tls, no libpq, no libssl-dev).

# ─── Stage 1: Build ───────────────────────────────────────────────────────────
# rust:1.75 (Debian Bookworm) ships with everything needed to compile
# a pure-Rust dependency tree: gcc, pkg-config, libssl-dev, make, build-essential.
FROM rust:1.75 AS builder

WORKDIR /app

# Cache dependency compilation separately from source changes
COPY Cargo.toml Cargo.lock ./
RUN mkdir -p src/bin && \
    echo "fn main() {}" > src/bin/phoenix-core.rs && \
    echo "fn main() {}" > src/bin/phoenix-agent.rs && \
    echo "fn main() {}" > src/bin/phoenix-deploy.rs
RUN cargo build --release
RUN rm -rf src

# Build real source
COPY src ./src/
COPY static ./static/
RUN touch src/bin/phoenix-core.rs src/bin/phoenix-agent.rs src/bin/phoenix-deploy.rs
RUN cargo build --release

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
# debian:bookworm-slim is compact and already has libc/libgcc — no apt needed.
FROM debian:bookworm-slim

# Pull SSL certs from builder so the app can make outbound HTTPS calls at runtime
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Copy all three binaries
COPY --from=builder /app/target/release/phoenix-core   /usr/local/bin/phoenix-core
COPY --from=builder /app/target/release/phoenix-agent  /usr/local/bin/phoenix-agent
COPY --from=builder /app/target/release/phoenix-deploy /usr/local/bin/phoenix-deploy

# Static assets
COPY --from=builder /app/static /app/static

# Non-root user + directory layout
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
LABEL org.opencontainers.image.title="Phoenix GUI"
LABEL org.opencontainers.image.licenses="MIT"

CMD ["/usr/local/bin/phoenix-core"]