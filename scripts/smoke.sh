#!/bin/bash
set -e

echo "=== Barca-Strategos Smoke Test ==="

# Build
echo "Building Go backend..."
go build -o bin/phoenix-api ./cmd/api

echo "Building frontend..."
cd frontend && npm ci && npm run build && cd ..

# Docker build
echo "Building Docker images..."
docker build -t phoenix-api .
docker build -t phoenix-frontend frontend/

# Start services
echo "Starting services..."
docker-compose -f docker-compose.full.yml up -d

# Wait for services
sleep 10

# Health checks
echo "Checking API health..."
curl -f http://localhost:8080/api/system/health || { echo "API health failed"; exit 1; }

echo "Checking frontend..."
curl -f http://localhost:3000 | grep -q "Barca-Strategos" || { echo "Frontend failed"; exit 1; }

echo "=== Smoke test passed ==="

# Cleanup
docker-compose -f docker-compose.full.yml down
