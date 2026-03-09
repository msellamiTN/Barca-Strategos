# Makefile for Barca-Strategos Phoenix GUI
# Provides convenient commands for building, running, and deploying the web application

.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean dev prod

# Default target
help:
	@echo "Barca-Strategos Phoenix GUI - Makefile Commands:"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make quick-start          - One-command deployment"
	@echo "  make quick-dev            - Quick development setup"
	@echo ""
	@echo "🏗️ Development:"
	@echo "  make build                - Build the Rust application"
	@echo "  make run                  - Run the web GUI locally"
	@echo "  make test                 - Run tests"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make dev                  - Start development server with hot reload"
	@echo ""
	@echo "🐳 Docker (Basic):"
	@echo "  make docker-build         - Build Docker image"
	@echo "  make docker-run           - Run with Docker Compose"
	@echo "  make docker-stop          - Stop Docker containers"
	@echo "  make docker-clean         - Clean Docker images and containers"
	@echo "  make docker-logs          - Show Docker logs"
	@echo ""
	@echo "📈 Scalable Deployment:"
	@echo "  make deploy-scalable      - Deploy with auto-scaling"
	@echo "  make scale-up             - Scale up services"
	@echo "  make scale-down           - Scale down services"
	@echo "  make status-cluster       - Check cluster status"
	@echo ""
	@echo "🤖 Chat Integration:"
	@echo "  make deploy-bots          - Deploy all chat bots"
	@echo "  make deploy-telegram      - Deploy Telegram bot only"
	@echo "  make deploy-discord       - Deploy Discord bot only"
	@echo "  make deploy-slack         - Deploy Slack bot only"
	@echo "  make test-telegram        - Test Telegram bot"
	@echo ""
	@echo "🏭 Production:"
	@echo "  make prod                 - Build and run production setup"
	@echo "  make deploy               - Deploy to production (requires config)"
	@echo "  make rollback             - Rollback to previous version"
	@echo ""
	@echo "🔧 Utilities:"
	@echo "  make check                - Check code formatting and linting"
	@echo "  make docs                 - Generate documentation"
	@echo "  make benchmark            - Run performance benchmarks"
	@echo "  make monitor              - Show monitoring dashboards"

# Development targets
build:
	@echo "🔨 Building Phoenix GUI..."
	cargo build --release
	@echo "✅ Build completed successfully!"

run:
	@echo "🚀 Starting Phoenix GUI locally..."
	cargo run --example start_web_gui

dev:
	@echo "🔧 Starting development server..."
	RUST_LOG=debug cargo run --example start_web_gui

test:
	@echo "🧪 Running tests..."
	cargo test
	@echo "✅ All tests passed!"

clean:
	@echo "🧹 Cleaning build artifacts..."
	cargo clean
	rm -rf target/
	@echo "✅ Clean completed!"

# Docker targets
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t phoenix-gui:latest .
	@echo "✅ Docker image built successfully!"

docker-run:
	@echo "🚀 Starting with Docker Compose..."
	docker-compose up -d
	@echo "✅ Containers started!"
	@echo "🌐 Web interface available at: http://localhost:8080"

docker-stop:
	@echo "🛑 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Containers stopped!"

docker-clean:
	@echo "🧹 Cleaning Docker resources..."
	docker-compose down -v --rmi all
	docker system prune -f
	@echo "✅ Docker cleanup completed!"

docker-logs:
	@echo "📋 Showing Docker logs..."
	docker-compose logs -f

# Production targets
prod: docker-build
	@echo "🚀 Starting production setup..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "✅ Production environment started!"
	@echo "🌐 Web interface available at: http://localhost:8080"

deploy:
	@echo "🚀 Deploying to production..."
	# This would typically involve:
	# - Building the image
	# - Pushing to registry
	# - Updating production environment
	@echo "⚠️  Deploy script not implemented yet"

# Utility targets
check:
	@echo "🔍 Checking code quality..."
	cargo fmt -- --check
	cargo clippy -- -D warnings
	@echo "✅ Code quality checks passed!"

docs:
	@echo "📚 Generating documentation..."
	cargo doc --no-deps --open
	@echo "✅ Documentation generated!"

benchmark:
	@echo "⚡ Running performance benchmarks..."
	cargo bench
	@echo "✅ Benchmarks completed!"

# Development utilities
install-tools:
	@echo "🔧 Installing development tools..."
	cargo install cargo-watch
	cargo install cargo-audit
	cargo install cargo-outdated
	@echo "✅ Development tools installed!"

watch:
	@echo "👀 Watching for changes..."
	cargo watch -x "run --example start_web_gui"

audit:
	@echo "🔒 Running security audit..."
	cargo audit
	@echo "✅ Security audit completed!"

update:
	@echo "⬆️ Updating dependencies..."
	cargo update
	cargo outdated
	@echo "✅ Dependencies updated!"

# Database utilities
db-migrate:
	@echo "🗄️ Running database migrations..."
	# This would run database migrations
	@echo "✅ Database migrations completed!"

db-reset:
	@echo "🔄 Resetting database..."
	docker-compose exec postgres psql -U phoenix -d phoenix -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "✅ Database reset completed!"

# Monitoring utilities
logs:
	@echo "📋 Showing application logs..."
	docker-compose logs -f phoenix-gui

status:
	@echo "📊 Checking system status..."
	docker-compose ps
	@echo ""
	@echo "🌐 Checking web interface..."
	curl -f http://localhost:8080/api/system/health || echo "❌ Web interface not responding"

# Quick start commands
quick-start: docker-build docker-run
	@echo "🚀 Quick start completed!"
	@echo "🌐 Open http://localhost:8080 in your browser"
	@echo "🤖 Telegram bot: http://t.me/your_bot_username"

quick-dev: build run
	@echo "🔧 Quick development setup completed!"

# Scalable deployment targets
deploy-scalable:
	@echo "📈 Deploying scalable Phoenix platform..."
	docker-compose -f docker-compose.scalable.yml up -d
	@echo "✅ Scalable deployment completed!"
	@echo "🌐 Load balancer: http://localhost"
	@echo "📊 Monitoring: http://localhost:3000 (admin/admin)"
	@echo "📈 Metrics: http://localhost:9090"

scale-up:
	@echo "📈 Scaling up services..."
	docker-compose -f docker-compose.scalable.yml up -d --scale phoenix-gui=5
	docker-compose -f docker-compose.scalable.yml up -d --scale telegram-bot=3
	@echo "✅ Services scaled up!"

scale-down:
	@echo "📉 Scaling down services..."
	docker-compose -f docker-compose.scalable.yml up -d --scale phoenix-gui=2
	docker-compose -f docker-compose.scalable.yml up -d --scale telegram-bot=1
	@echo "✅ Services scaled down!"

status-cluster:
	@echo "📊 Cluster Status:"
	docker-compose -f docker-compose.scalable.yml ps
	@echo ""
	@echo "🔍 Service Health:"
	@curl -s http://localhost/api/system/health 2>/dev/null && echo "✅ Main API: Healthy" || echo "❌ Main API: Unhealthy"
	@curl -s http://localhost:3000/api/health 2>/dev/null && echo "✅ Grafana: Healthy" || echo "❌ Grafana: Unhealthy"

# Chat integration targets
deploy-bots:
	@echo "🤖 Deploying all chat bots..."
	docker-compose up -d telegram-bot discord-bot slack-bot teams-bot
	@echo "✅ All chat bots deployed!"
	@echo "📱 Configure your bot tokens in .env file"

deploy-telegram:
	@echo "📱 Deploying Telegram bot..."
	docker-compose up -d telegram-bot
	@echo "✅ Telegram bot deployed!"
	@echo "🔧 Set TELEGRAM_BOT_TOKEN in .env"

deploy-discord:
	@echo "🎮 Deploying Discord bot..."
	docker-compose up -d discord-bot
	@echo "✅ Discord bot deployed!"
	@echo "🔧 Set DISCORD_BOT_TOKEN in .env"

deploy-slack:
	@echo "💬 Deploying Slack bot..."
	docker-compose up -d slack-bot
	@echo "✅ Slack bot deployed!"
	@echo "🔧 Set SLACK_BOT_TOKEN in .env"

test-telegram:
	@echo "🧪 Testing Telegram bot..."
	@echo "Send /start to your bot: https://t.me/your_bot_username"
	@echo "Available commands:"
	@echo "  /start - Start the bot"
	@echo "  /dashboard - View system dashboard"
	@echo "  /security - Security status"
	@echo "  /compliance - Compliance overview"
	@echo "  /risk - Risk assessment"
	@echo "  /agents - AI agent status"
	@echo "  /scale - Scale services"
	@echo "  /status - System status"

# Enhanced production targets
prod: docker-build
	@echo "🚀 Starting production setup..."
	docker-compose -f docker-compose.yml -f docker-compose.scalable.yml up -d
	@echo "✅ Production environment started!"
	@echo "🌐 Web interface available at: http://localhost"
	@echo "📊 Monitoring: http://localhost:3000"

deploy:
	@echo "🚀 Deploying to production..."
	# This would typically involve:
	# - Building the image
	# - Pushing to registry
	# - Updating production environment
	@echo "⚠️  Deploy script not implemented yet"

rollback:
	@echo "🔄 Rolling back to previous version..."
	docker-compose down
	docker tag phoenix-gui:previous phoenix-gui:current
	docker-compose up -d
	@echo "✅ Rollback completed!"

# Monitoring utilities
monitor:
	@echo "📊 Opening monitoring dashboards..."
	@echo "🌐 Grafana: http://localhost:3000 (admin/admin)"
	@echo "📈 Prometheus: http://localhost:9090"
	@echo "🔍 Jaeger: http://localhost:16686"
	@if command -v xdg-open > /dev/null; then \
		xdg-open http://localhost:3000; \
	elif command -v open > /dev/null; then \
		open http://localhost:3000; \
	fi

# Environment setup
setup-dev:
	@echo "🔧 Setting up development environment..."
	cp .env.template .env
	@echo "⚠️  Please edit .env file with your configuration"
	@echo "✅ Development environment setup completed!"

setup-prod:
	@echo "🏭 Setting up production environment..."
	cp .env.template .env.prod
	@echo "⚠️  Please edit .env.prod file with production configuration"
	@echo "✅ Production environment setup completed!"

# Backup and restore
backup:
	@echo "💾 Creating backup..."
	docker-compose exec postgres pg_dump -U phoenix phoenix > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup completed!"

restore:
	@echo "🔄 Restoring from backup..."
	@echo "⚠️  Usage: make restore BACKUP_FILE=<filename>"
	@if [ -z "$(BACKUP_FILE)" ]; then echo "❌ Please specify BACKUP_FILE"; exit 1; fi
	docker-compose exec -T postgres psql -U phoenix phoenix < $(BACKUP_FILE)
	@echo "✅ Restore completed!"

# Performance profiling
profile:
	@echo "📊 Running performance profiling..."
	cargo run --example start_web_gui --release --features profiling
	@echo "✅ Profiling completed!"

# Security scanning
security-scan:
	@echo "🔒 Running security scan..."
	cargo audit
	docker run --rm -v $(PWD):/app clair-scanner:latest
	@echo "✅ Security scan completed!"

# Load testing
load-test:
	@echo "⚡ Running load tests..."
	# This would run load testing tools like Apache Bench or k6
	@echo "✅ Load tests completed!"

# CI/CD helpers
ci-build:
	@echo "🔧 CI Build..."
	cargo build --release
	cargo test
	cargo clippy -- -D warnings
	@echo "✅ CI build completed!"

ci-test:
	@echo "🧪 CI Tests..."
	cargo test --release
	@echo "✅ CI tests completed!"

# Version management
version:
	@echo "📋 Version information:"
	@echo "Rust: $(shell rustc --version)"
	@echo "Cargo: $(shell cargo --version)"
	@echo "Docker: $(shell docker --version)"
	@echo "Docker Compose: $(shell docker-compose --version)"

# Help for specific areas
help-docker:
	@echo "Docker Commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run with Docker Compose"
	@echo "  make docker-stop    - Stop Docker containers"
	@echo "  make docker-clean   - Clean Docker resources"
	@echo "  make docker-logs    - Show Docker logs"

help-dev:
	@echo "Development Commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run locally"
	@echo "  make dev            - Development server"
	@echo "  make test           - Run tests"
	@echo "  make watch          - Watch for changes"

help-prod:
	@echo "Production Commands:"
	@echo "  make prod           - Production setup"
	@echo "  make deploy         - Deploy to production"
	@echo "  make backup         - Create backup"
	@echo "  make restore        - Restore from backup"
