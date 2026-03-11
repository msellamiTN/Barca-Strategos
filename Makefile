.PHONY: build-go build-frontend build-docker test-go lint-frontend up down clean

# Go backend
build-go:
	go mod tidy
	go build -o bin/phoenix-api ./cmd/api

test-go:
	go test ./...

lint-go:
	go vet ./...
	staticcheck ./...

# Frontend
build-frontend:
	cd frontend && npm ci && npm run build

lint-frontend:
	cd frontend && npm run lint

# Docker
build-docker:
	docker build -t phoenix-api .
	docker build -t phoenix-frontend frontend/

up:
	docker-compose -f docker-compose.full.yml up -d

down:
	docker-compose -f docker-compose.full.yml down

clean:
	docker-compose -f docker-compose.full.yml down -r all
	docker rmi phoenix-api phoenix-frontend || true
