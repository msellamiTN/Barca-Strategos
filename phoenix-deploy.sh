#!/bin/bash
# Phoenix Deployment Script - Ultra-efficient AI Security Framework
# Where Tactical Excellence Meets Technological Revolution

set -e

PHOENIX_VERSION="1.0.0"
PHOENIX_IMAGE="barca-strategos/phoenix:${PHOENIX_VERSION}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Print functions
print_status() { echo -e "${GREEN}[PHOENIX]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[PHOENIX]${NC} $1"; }
print_error() { echo -e "${RED}[PHOENIX]${NC} $1"; }
print_info() { echo -e "${BLUE}[PHOENIX]${NC} $1"; }
print_phoenix() { echo -e "${PURPLE}[🔥 PHOENIX]${NC} $1"; }

# Banner
print_banner() {
    echo -e "${PURPLE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    BARCA-STRATEGOS PHOENIX                   ║"
    echo "║              Cognitive Collaboration Platform               ║"
    echo "║                     � Human-AI Teamwork �              ║"
    echo "║                                                              ║"
    echo "║  • Web GUI Interface  • Multi-Agent System  • Chat Bots     ║"
    echo "║  • Auto-Scaling  • Enterprise Security  • Real-time       ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# System requirements check
check_requirements() {
    print_info "🔍 Checking system requirements..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "❌ Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check Docker permissions
    if ! docker info &> /dev/null; then
        print_warning "⚠️  Docker requires sudo permissions. Using sudo for Docker commands."
        DOCKER_SUDO="sudo"
    else
        DOCKER_SUDO=""
    fi
    
    # Check Docker Compose (v2 syntax: "docker compose", v1: "docker-compose")
    if $DOCKER_SUDO docker compose version &> /dev/null; then
        DOCKER_COMPOSE="$DOCKER_SUDO docker compose"
    elif command -v docker-compose &> /dev/null; then
        DOCKER_COMPOSE="$DOCKER_SUDO docker-compose"
    else
        print_error "❌ Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    # Check memory
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        AVAILABLE_MEMORY=$(free -m | awk 'NR==2{printf "%.0f", $7}')
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        AVAILABLE_MEMORY=$(vm_stat | grep "Pages free" | awk '{print $3}' | sed 's/\.//' | awk '{printf "%.0f", $1 * 4 / 1024}')
    else
        AVAILABLE_MEMORY=2048  # Default assumption
    fi
    
    if [ "$AVAILABLE_MEMORY" -lt 512 ]; then
        print_error "❌ Minimum 512MB RAM required. Available: ${AVAILABLE_MEMORY}MB"
        exit 1
    fi
    
    # Check disk space
    AVAILABLE_DISK=$(df . | tail -1 | awk '{print $4}')
    if [ "$AVAILABLE_DISK" -lt 1048576 ]; then  # 1GB in KB
        print_error "❌ Minimum 1GB disk space required."
        exit 1
    fi
    
    print_status "✅ System requirements met (${AVAILABLE_MEMORY}MB RAM, ${AVAILABLE_DISK}KB disk)"
    print_status "✅ Docker Compose command: $DOCKER_COMPOSE"
}

# Create directory structure
create_directories() {
    print_info "📁 Creating Phoenix directory structure..."
    
    directories=(
        "data/phoenix"
        "data/ai"
        "data/runtime"
        "data/collaboration"
        "config/security"
        "config/platforms"
        "agents"
        "logs"
        "nginx/ssl"
        "init-scripts"
    )
    
    for dir in "${directories[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            print_info "📂 Created: $dir"
        fi
    done
    
    # Set permissions
    chmod 755 data config agents logs
    chmod 700 config/security
    
    print_status "✅ Directory structure created"
}

# Generate secure configuration
generate_secure_config() {
    print_info "🔐 Generating secure configuration..."
    
    if [ ! -f .env ]; then
        print_info "📝 Creating .env file from template..."
        cp .env.template .env
        
        # Generate secure secrets
        JWT_SECRET=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32)
        AGENT_SECRET=$(openssl rand -base64 24 | tr -d "=+/" | cut -c1-24)
        DB_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
        REDIS_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
        
        # Update .env file
        sed -i.bak "s/your_jwt_secret_minimum_32_characters_here/$JWT_SECRET/" .env
        sed -i.bak "s/your_agent_secret_minimum_24_characters_here/$AGENT_SECRET/" .env
        sed -i.bak "s/your_secure_db_password_here/$DB_PASSWORD/" .env
        sed -i.bak "s/your_secure_redis_password_here/$REDIS_PASSWORD/" .env
        
        # Remove backup
        rm .env.bak
        
        print_status "✅ Secure configuration generated"
        print_warning "⚠️  Please edit .env file and add your platform tokens:"
        print_warning "   • TELEGRAM_BOT_TOKEN"
        print_warning "   • SLACK_BOT_TOKEN"
        print_warning "   • API_KEY (for AI)"
        print_warning "   • WEBHOOK_URL"
    else
        print_status "✅ .env file already exists"
    fi
}

# Pull Docker images
pull_images() {
    print_info "📦 Pulling Phoenix Docker images..."
    
    # Check if images already exist locally
    if $DOCKER_SUDO docker images | grep -q "barca-strategos"; then
        print_info "�️  Some Phoenix images already exist locally"
    fi
    
    # Pull images using docker-compose
    if $DOCKER_COMPOSE pull; then
        print_status "✅ Docker images pulled successfully"
    else
        print_warning "⚠️  Failed to pull some images, will build locally"
        build_images
    fi
}

# Build images locally if pull fails
build_images() {
    print_info "🔨 Building Phoenix images locally..."
    
    # Build core image
    if [ -f "Dockerfile" ]; then
        print_info "🏗️  Building phoenix-gui image..."
        $DOCKER_SUDO docker build -t barca-strategos/phoenix-gui:latest .
    fi
    
    # Build other images (if Dockerfiles exist)
    for service in agent-runtime ai-assistant collab-hub telegram-bot web-ui; do
        if [ -f "${service}/Dockerfile" ]; then
            print_info "🏗️  Building ${service} image..."
            $DOCKER_SUDO docker build -t "barca-strategos/${service}:latest" "./${service}"
        fi
    done
    
    print_status "✅ Images built successfully"
}

# Start Phoenix services
start_services() {
    print_phoenix "🚀 Starting Phoenix services..."
    
    # Start with ultra-fast boot optimization
    if $DOCKER_COMPOSE up -d --remove-orphans; then
        print_status "✅ Phoenix services started"
    else
        print_error "❌ Failed to start Phoenix services"
        $DOCKER_COMPOSE logs
        exit 1
    fi
}

# Wait for services to be ready
wait_for_services() {
    print_info "⏳ Waiting for Phoenix services to be ready..."
    
    # Wait for core service
    local max_wait=60
    local wait_time=0
    
    while [ $wait_time -lt $max_wait ]; do
        if curl -s http://localhost:8080/health >/dev/null 2>&1; then
            print_status "✅ Phoenix Core is ready"
            break
        fi
        
        if [ $wait_time -eq 0 ]; then
            echo -n "⏳ Waiting for Phoenix Core"
        else
            echo -n "."
        fi
        
        sleep 2
        wait_time=$((wait_time + 2))
    done
    
    if [ $wait_time -ge $max_wait ]; then
        print_error "❌ Phoenix Core failed to start within ${max_wait}s"
        $DOCKER_COMPOSE logs phoenix-gui
        exit 1
    fi
    
    # Check other services
    print_info "🔍 Checking service status..."
    sleep 5
    
    # Get service status
    if $DOCKER_COMPOSE ps | grep -q "Up"; then
        print_status "✅ All services are running"
    else
        print_warning "⚠️  Some services may not be ready yet"
        $DOCKER_COMPOSE ps
    fi
}

# Initialize database
init_database() {
    print_info "🗄️  Initializing Phoenix database..."
    
    # Wait for PostgreSQL to be ready
    local max_wait=30
    local wait_time=0
    
    while [ $wait_time -lt $max_wait ]; do
        if $DOCKER_COMPOSE exec -T postgres pg_isready -U phoenix -d phoenix >/dev/null 2>&1; then
            break
        fi
        sleep 2
        wait_time=$((wait_time + 2))
    done
    
    # Run migrations (if migration script exists)
    if $DOCKER_COMPOSE exec -T phoenix-gui /app/phoenix migrate >/dev/null 2>&1; then
        print_status "✅ Database initialized"
    else
        print_warning "⚠️  Database initialization may be needed manually"
    fi
}

# Show deployment information
show_deployment_info() {
    echo
    print_phoenix "🎉 Phoenix deployed successfully!"
    echo
    echo -e "${CYAN}� Barca-Strategos Phoenix Services:${NC}"
    echo "  • Web GUI:       http://localhost:8080"
    echo "  • Monitoring:    http://localhost:3000 (Grafana: admin/admin)"
    echo "  • Metrics:       http://localhost:9090 (Prometheus)"
    echo "  • Tracing:       http://localhost:16686 (Jaeger)"
    echo "  • API:           http://localhost:8080/api"
    echo "  • Telegram Bot:  Configure token in .env"
    echo
    echo -e "${CYAN}📊 Resource Usage:${NC}"
    $DOCKER_SUDO docker stats --no-stream --format "table {{.Container}}\t{{.MemUsage}}\t{{.CPUPerc}}"
    echo
    echo -e "${CYAN}🔧 Management Commands:${NC}"
    echo "  • Status:        ./phoenix-deploy.sh status"
    echo "  • Logs:          ./phoenix-deploy.sh logs [service]"
    echo "  • Stop:          ./phoenix-deploy.sh stop"
    echo "  • Restart:       ./phoenix-deploy.sh restart"
    echo "  • Scale:         ./phoenix-deploy.sh scale [count]"
    echo "  • Update:        ./phoenix-deploy.sh update"
    echo "  • Backup:        ./phoenix-deploy.sh backup"
    echo
    echo -e "${CYAN}🦐 Phoenix Features:${NC}"
    echo "  • Ultra-lightweight agents (<10MB each)"
    echo "  • 1-second boot time"
    echo "  • AI-powered security analysis"
    echo "  • Multi-platform collaboration"
    echo "  • Zero-trust security architecture"
    echo
    print_warning "⚠️  Remember to configure your platform tokens in .env file"
}

# Show service status
show_status() {
    print_info "📊 Phoenix Service Status:"
    docker-compose ps
    echo
    print_info "📈 Resource Usage:"
    docker stats --no-stream --format "table {{.Container}}\t{{.MemUsage}}\t{{.CPUPerc}}"
}

# Show logs
show_logs() {
    local service=${1:-""}
    if [ -n "$service" ]; then
        print_info "📋 Showing logs for: $service"
        docker-compose logs -f "$service"
    else
        print_info "📋 Showing logs for all services"
        docker-compose logs -f
    fi
}

# Stop services
stop_services() {
    print_info "🛑 Stopping Phoenix services..."
    docker-compose down
    print_status "✅ Phoenix services stopped"
}

# Restart services
restart_services() {
    print_info "🔄 Restarting Phoenix services..."
    docker-compose restart
    print_status "✅ Phoenix services restarted"
}

# Scale services
scale_services() {
    local count=${1:-3}
    print_info "📈 Scaling Phoenix agents to: $count"
    docker-compose up -d --scale phoenix-core="$count"
    print_status "✅ Services scaled"
}

# Update services
update_services() {
    print_info "🔄 Updating Phoenix..."
    docker-compose pull
    docker-compose up -d
    print_status "✅ Phoenix updated"
}

# Backup data
backup_data() {
    print_info "💾 Creating Phoenix backup..."
    
    local backup_dir="./backups/phoenix-$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"
    
    # Backup data directories
    cp -r data "$backup_dir/"
    cp -r config "$backup_dir/"
    cp .env "$backup_dir/"
    cp docker-compose.yml "$backup_dir/"
    
    # Compress backup
    tar -czf "${backup_dir}.tar.gz" -C "./backups" "$(basename "$backup_dir")"
    rm -rf "$backup_dir"
    
    print_status "✅ Backup created: ${backup_dir}.tar.gz"
}

# Health check
health_check() {
    print_info "🏥 Phoenix Health Check:"
    
    # Check core service
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        print_status "✅ Phoenix Core: Healthy"
    else
        print_error "❌ Phoenix Core: Unhealthy"
    fi
    
    # Check database
    if docker-compose exec -T postgres pg_isready -U barca -d barca_db >/dev/null 2>&1; then
        print_status "✅ Database: Healthy"
    else
        print_error "❌ Database: Unhealthy"
    fi
    
    # Check Redis
    if docker-compose exec -T redis redis-cli ping >/dev/null 2>&1; then
        print_status "✅ Redis: Healthy"
    else
        print_error "❌ Redis: Unhealthy"
    fi
    
    # Show overall status
    show_status
}

# Main deployment function
main() {
    print_banner
    
    case "${1:-}" in
        "stop")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "${2:-}"
            ;;
        "scale")
            scale_services "${2:-3}"
            ;;
        "update")
            update_services
            ;;
        "backup")
            backup_data
            ;;
        "health")
            health_check
            ;;
        "help"|"-h"|"--help")
            echo "Phoenix Deployment Script"
            echo
            echo "Usage: $0 [COMMAND] [OPTIONS]"
            echo
            echo "Commands:"
            echo "  (no args)    Deploy Phoenix (default)"
            echo "  stop          Stop all services"
            echo "  restart       Restart all services"
            echo "  status        Show service status"
            echo "  logs [svc]    Show logs (all services or specific)"
            echo "  scale [n]     Scale agents to n instances (default: 3)"
            echo "  update        Update to latest version"
            echo "  backup        Create backup of data and config"
            echo "  health        Run health check"
            echo "  help          Show this help"
            echo
            ;;
        *)
            # Default deployment
            check_requirements
            create_directories
            generate_secure_config
            
            # Try to pull images, build if fails
            if ! pull_images; then
                build_images
            fi
            
            start_services
            wait_for_services
            init_database
            show_deployment_info
            ;;
    esac
}

# Error handling
trap 'print_error "❌ Deployment failed! Check logs for details."; exit 1' ERR

# Run main function
main "$@"
