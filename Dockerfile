# Multi-stage Dockerfile for Barca-Strategos Phoenix GUI
# Fully static binary → zero runtime package installation required

# ─── Stage 1: Build ───────────────────────────────────────────────────────────
# rust:1.75 (Debian-based) already ships with gcc, pkg-config, libssl-dev,
# make, etc. — no apk/apt needed for the core toolchain.
FROM rust:1.75 AS builder

# Add the musl target for fully static compilation
RUN rustup target add x86_64-unknown-linux-musl

# Install only the two things Debian doesn't ship in rust:1.75 by default.
# If apt is also blocked, see the NOTE below about pure-Rust alternatives.
RUN apt-get update && apt-get install -y --no-install-recommends \
    musl-tools \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Cache dependency layer
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN RUSTFLAGS="-C target-feature=+crt-static" \
    cargo build --release --target x86_64-unknown-linux-musl \
    && rm -rf src

# Build real binary
COPY src ./src/
COPY static ./static/
RUN RUSTFLAGS="-C target-feature=+crt-static" \
    cargo build --release --target x86_64-unknown-linux-musl

# ─── Stage 2: Runtime (scratch = zero OS, zero packages to install) ───────────
FROM scratch

# SSL certificates baked in from builder — no ca-certificates package needed
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Statically linked binary — no libc, libssl, or libpq dependencies at runtime
COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/phoenix-core /phoenix-core

COPY --from=builder /app/static /static

EXPOSE 8080

ENV RUST_LOG=info
ENV PHOENIX_HOST=0.0.0.0
ENV PHOENIX_PORT=8080
ENV PHOENIX_STATIC_PATH=/static
ENV PHOENIX_MAX_CONNECTIONS=1000

# Volume mounts for persistent data
VOLUME ["/logs", "/data", "/config"]

# Labels
LABEL maintainer="Barca-Strategos Team"
LABEL description="Barca-Strategos Phoenix GUI - Cognitive Collaboration Platform"
LABEL version="1.0.0"
LABEL org.opencontainers.image.title="Phoenix GUI"
LABEL org.opencontainers.image.licenses="MIT"

ENTRYPOINT ["/phoenix-core"]