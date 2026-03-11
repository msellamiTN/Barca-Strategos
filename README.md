# 🚀 Barca-Strategos Phoenix

> **Cognitive Collaboration Platform with Web GUI** - Where Human-AI Teamwork Meets Scalable Security

![Go](https://img.shields.io/badge/Go-1.22+-blue.svg)![Next.js](https://img.shields.io/badge/Next.js-15-black.svg)![Docker](https://img.shields.io/badge/docker-ready-blue.svg)![License](https://img.shields.io/badge/license-MIT-green.svg)

## 🌟 Overview

Barca-Strategos Phoenix is a modular Go platform delivering security operations, compliance automation, and collaborative dashboards through a Fiber/gRPC service layer and a Next.js web GUI. It features an agentic AI orchestration layer inspired by OpenClaw principles, enabling autonomous workflows with safety guardrails.

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- 2GB+ RAM recommended
- 2+ CPU cores for optimal performance

### One-Command Deployment

```bash
# Clone the repository
git clone https://github.com/barca-strategos/phoenix.git
cd phoenix

# Quick start with web GUI and all features
make build-docker && make up
```

### Access Phoenix

- **🌐 Web GUI**: <http://localhost:3000> (Next.js frontend)
- **🔧 API**: <http://localhost:8080> (Go backend)
- **📊 Health**: <http://localhost:8080/api/system/health>

## 🛠️ Development

### Build & Test

```bash
# Go backend
make build-go
make test-go
make lint-go

# Next.js frontend
make build-frontend
make lint-frontend

# Docker images
make build-docker
```

### Local Development

```bash
# Start full stack
make up

# Stop stack
make down

# Clean images
make clean
```

### Smoke Test

```bash
chmod +x scripts/smoke.sh
./scripts/smoke.sh
```

## 📁 Project Structure

```
barca-strategos/
├── cmd/                    # Go entrypoints
│   └── api/
│       └── main.go
├── internal/               # Go internal packages
│   ├── config/
│   ├── server/
│   ├── routes/
│   ├── monitoring/
│   ├── compliance/
│   ├── security/
│   └── agentic/
├── frontend/               # Next.js app
│   ├── app/
│   ├── components/
│   └── public/
├── docs/                   # Architecture & reports
├── scripts/                # Automation scripts
├── .github/workflows/      # CI/CD
├── Dockerfile
├── docker-compose.full.yml
├── Makefile
└── README.md
```

## 🔒 Security & Compliance

- Zero-Trust Architecture with mutual authentication
- Agentic AI guardrails (broker, shared memory, referee)
- Continuous compliance monitoring (SOC2, PCI, ISO)
- Comprehensive audit logging and evidence export

## 📈 Scalability

- Horizontal scaling via Docker Compose/Kubernetes
- Connection pooling (Postgres, Redis)
- Observability with OpenTelemetry
- Auto-scaling targets: 5,000 req/s, sub-80ms latency, 99.99% availability

## 🤖 Agentic AI

- Tool registry with safety tiers
- Shared scratchpad memory for agents
- Adaptive autonomy levels (manual/assisted/autonomous)
- Referee loop for policy enforcement

## 📚 Documentation

- [Architecture Review](docs/architecture-review-2026.md)
- [Language Migration Evaluation](docs/language-migration-evaluation.md)
- [Deep Review & Alignment](docs/deep-review-alignment-report.md)

## 🙏 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and documentation
5. Submit a pull request

---

*Barca-Strategos Phoenix — Where human-AI collaboration meets enterprise security*
