# ─── Barca-Strategos Phoenix GUI ──────────────────────────────────────────────
# Fully offline build — all crates vendored locally, zero network calls.
# Prerequisites on HOST before building:
#   cargo vendor && mkdir -p .cargo && cargo vendor >> .cargo/config.toml

# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM rust:1.85 AS builder

WORKDIR /app

# Copy cargo config (tells cargo to use vendor/ instead of crates.io)
COPY .cargo .cargo/

# Copy manifests + vendored crates
COPY Cargo.toml Cargo.lock ./
COPY vendor ./vendor/

# Dummy binaries to cache the dependency compile layer
RUN mkdir -p src/bin && \
    echo "fn main() {}" > src/bin/phoenix-core.rs && \
    echo "fn main() {}" > src/bin/phoenix-agent.rs && \
    echo "fn main() {}" > src/bin/phoenix-deploy.rs

# Build deps from vendor — no network needed
RUN cargo build --release --offline
RUN rm -rf src

# Build real source
COPY src ./src/
COPY static ./static/
RUN touch src/bin/phoenix-core.rs src/bin/phoenix-agent.rs src/bin/phoenix-deploy.rs
RUN cargo build --release --offline

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