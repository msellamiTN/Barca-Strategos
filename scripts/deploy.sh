#!/bin/bash

# Phoenix API Deployment Script
# This script deploys the Phoenix API with SOC2 and PCI DSS compliance modules

set -e

# Configuration
DOCKER_IMAGE="barca-strategos/phoenix-api:latest"
CONTAINER_NAME="phoenix-api"
ENVIRONMENT=${ENVIRONMENT:-development}
BACKUP_ENABLED=${BACKUP_ENABLED:-true}
HEALTH_CHECK_TIMEOUT=${HEALTH_CHECK_TIMEOUT:-60}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed or not in PATH"
        exit 1
    fi
}

# Function to check if Docker Compose is running
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed or not in PATH"
        exit 1
    fi
}

# Function to backup existing data
backup_data() {
    if [ "$BACKUP_ENABLED" = "true" ]; then
        log_info "Creating backup of existing data..."
        
        BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        
        # Backup PostgreSQL data
        if docker ps -q -f phoenix-api-postgres-1 &>/dev/null; then
            log_info "Backing up PostgreSQL data..."
            docker exec phoenix-api-postgres-1 pg_dump -U phoenix phoenix > "$BACKUP_DIR/postgres_backup.sql"
        fi
        
        # Backup Redis data
        if docker ps -q -f phoenix-api-redis-1 &>/dev/null; then
            log_info "Backing up Redis data..."
            docker exec phoenix-api-redis-1 redis-cli BGSAVE
            docker cp phoenix-api-redis-1:/data/dump.rdb "$BACKUP_DIR/redis_backup.rdb"
        fi
        
        # Backup reports and logs
        if [ -d "./reports" ]; then
            log_info "Backing up reports..."
            cp -r ./reports "$BACKUP_DIR/reports"
        fi
        
        if [ -d "./logs" ]; then
            log_info "Backing up logs..."
            cp -r ./logs "$BACKUP_DIR/logs"
        fi
        
        log_success "Backup completed: $BACKUP_DIR"
    fi
}

# Function to build Docker image
build_image() {
    log_info "Building Docker image..."
    # Use default Dockerfile (busybox-based for restricted networks)
    docker build -t "$DOCKER_IMAGE" .
    log_success "Docker image built successfully"
}

# Function to stop existing containers
stop_containers() {
    log_info "Stopping existing containers..."
    
    if docker ps -q -f "$CONTAINER_NAME" &>/dev/null; then
        docker stop "$CONTAINER_NAME"
        docker rm "$CONTAINER_NAME"
        log_success "Stopped and removed existing container"
    fi
    
    if docker-compose ps -q &>/dev/null; then
        docker-compose down
        log_success "Stopped Docker Compose services"
    fi
}

# Function to start services
start_services() {
    log_info "Starting Phoenix API services..."
    
    # Set environment variables
    export PHOENIX_ENV="$ENVIRONMENT"
    export PHOENIX_PORT="${PHOENIX_PORT:-8080}"
    export PHOENIX_DATABASE_URL="${PHOENIX_DATABASE_URL:-postgres://phoenix:phoenix@postgres:5432/phoenix?sslmode=disable}"
    export PHOENIX_REDIS_URL="${PHOENIX_REDIS_URL:-redis://redis:6379/0}"
    export PHOENIX_SOC2_MONITORING_INTERVAL="${PHOENIX_SOC2_MONITORING_INTERVAL:-6h}"
    export PHOENIX_PCI_DSS_MONITORING_INTERVAL="${PHENIX_PCI_DSS_MONITORING_INTERVAL:-6h}"
    export PHOENIX_COMPLIANCE_ASSESSMENT_INTERVAL="${PHOENIX_COMPLIANCE_ASSESSMENT_INTERVAL:-24h}"
    export PHOENIX_ENABLE_BACKGROUND_TASKS="${PHOENIX_ENABLE_BACKGROUND_TASKS:-true}"
    
    # Start with Docker Compose
    docker-compose up -d
    
    log_success "Services started successfully"
}

# Function to wait for health check
wait_for_health() {
    log_info "Waiting for services to be healthy..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:8080/api/system/health &>/dev/null; then
            log_success "All services are healthy"
            return 0
        fi
        
        log_warning "Health check attempt $attempt/$max_attempts failed. Retrying in 2 seconds..."
        sleep 2
        ((attempt++))
    done
    
    log_error "Health check failed after $max_attempts attempts"
    return 1
}

# Function to run database migrations
run_migrations() {
    log_info "Running database migrations..."
    
    # Wait for database to be ready
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker exec phoenix-api-postgres-1 pg_isready -U phoenix &>/dev/null; then
            log_success "Database is ready"
            break
        fi
        
        log_warning "Database not ready, attempt $attempt/$max_attempts. Retrying in 2 seconds..."
        sleep 2
        ((attempt++))
    done
    
    if [ $attempt -gt $max_attempts ]; then
        log_error "Database failed to become ready"
        return 1
    fi
    
    # Run initialization script if needed
    if docker exec phoenix-api-postgres-1 psql -U phoenix -c "\dt compliance_statistics;" &>/dev/null; then
        log_info "Database already initialized"
    else
        log_info "Initializing database..."
        docker exec phoenix-api-postgres-1 psql -U phoenix -f /docker-entrypoint-initdb.d/init-db.sql
        log_success "Database initialized successfully"
    fi
}

# Function to show service status
show_status() {
    log_info "Service Status:"
    
    echo "=== Docker Containers ==="
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    echo ""
    echo "=== Service Endpoints ==="
    
    if curl -f http://localhost:8080/api/system/health &>/dev/null; then
        echo -e "${GREEN}✓${NC} API Health Check: http://localhost:8080/api/system/health"
    else
        echo -e "${RED}✗${NC} API Health Check: http://localhost:8080/api/system/health"
    fi
    
    echo ""
    echo "=== Database Connection ==="
    if docker exec phoenix-api-postgres-1 pg_isready -U phoenix &>/dev/null; then
        echo -e "${GREEN}✓${NC} PostgreSQL: Connected"
    else
        echo -e "${RED}✗${NC} PostgreSQL: Not Connected"
    fi
    
    if docker exec phoenix-api-redis-1 redis-cli ping &>/dev/null; then
        echo -e "${GREEN}✓${NC} Redis: Connected"
    else
        echo -e "${RED}✗${NC} Redis: Not Connected"
    fi
    
    echo ""
    echo "=== Monitoring Services ==="
    if curl -f http://localhost:9090/targets &>/dev/null; then
        echo -e "${GREEN}✓${NC} Prometheus: http://localhost:9090"
    else
        echo -e "${YELLOW}⚠${NC} Prometheus: Not Available"
    fi
    
    if curl -f http://localhost:3000/api/health &>/dev/null; then
        echo -e "${GREEN}✓${NC} Grafana: http://localhost:3000"
    else
        echo -e "${YELLOW}⚠${NC} Grafana: Not Available"
    fi
}

# Function to show logs
show_logs() {
    log_info "Showing logs (press Ctrl+C to exit)..."
    docker-compose logs -f
}

# Function to clean up
cleanup() {
    log_info "Cleaning up..."
    
    stop_containers
    
    # Remove unused Docker images
    log_info "Removing unused Docker images..."
    docker image prune -f
    
    # Remove unused Docker volumes
    log_info "Removing unused Docker volumes..."
    docker volume prune -f
    
    log_success "Cleanup completed"
}

# Function to run tests
run_tests() {
    log_info "Running tests..."
    
    # Run unit tests
    log_info "Running unit tests..."
    go test ./internal/compliance/soc2/...
    go test ./internal/compliance/pci_dss/...
    
    # Run integration tests
    log_info "Running integration tests..."
    go test ./internal/routes/api/v1/...
    
    # Run end-to-end tests
    log_info "Running end-to-end tests..."
    go test ./test/integration/...
    
    log_success "All tests completed"
}

# Function to generate compliance reports
generate_reports() {
    log_info "Generating compliance reports..."
    
    # Generate SOC2 report
    log_info "Generating SOC 2 Type II report..."
    curl -X POST http://localhost:8080/api/v1/compliance/soc2/assess \
        -H "Content-Type: application/json" \
        -d '{"departments": ["IT", "Security"], "systems": ["Phoenix Core", "Database"], "processes": ["Incident Response"]}' \
        | jq -r '.assessment_id' > /tmp/soc2_assessment_id
    
    if [ -s /tmp/soc2_assessment_id ]; then
        curl -X POST http://localhost:8080/api/v1/compliance/soc2/report \
            -H "Content-Type: application/json" \
            -d "{\"assessment_id\": \"$(cat /tmp/soc2_assessment_id)\"}" \
            -o reports/soc2_report_$(date +%Y%m%d).json
        log_success "SOC 2 report generated: reports/soc2_report_$(date +%Y%m%d).json"
    fi
    
    # Generate PCI DSS report
    log_info "Generating PCI DSS v4.0 report..."
    curl -X POST http://localhost:8080/api/v1/compliance/pci_dss/assess \
        -H "Content-Type: application/json" \
        -d '{"departments": ["Payment Processing", "Security"], "systems": ["Payment Gateway", "Database"], "processes": ["Card Processing"]}' \
        | jq -r '.assessment_id' > /tmp/pci_assessment_id
    
    if [s /tmp/pci_assessment_id ]; then
        curl -X POST http://localhost:8080/api/v1/compliance/pci_dss/report \
            -H "Content-Type: application/json" \
            -d "{\"assessment_id\": \"$(cat /tmp/pci_assessment_id)\"}" \
            -o reports/pci_dss_report_$(date +%Y%m%d).json
        log_success "PCI DSS report generated: reports/pci_dss_report_$(date +%Y%m%d).json"
    fi
    
    rm -f /tmp/soc2_assessment_id /tmp/pci_assessment_id
    log_success "Compliance reports generated successfully"
}

# Function to scale services
scale_services() {
    local service=${1:-phoenix-api}
    local replicas=${2:-2}
    
    log_info "Scaling $service to $replicas replicas..."
    docker-compose up -d --scale "$service=$replicas"
    log_success "Service scaled successfully"
}

# Function to show help
show_help() {
    echo "Phoenix API Deployment Script"
    echo ""
    "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    "Commands:"
    echo "  build              Build the Docker image"
    echo "  deploy             Deploy the application (default)"
    echo "  stop               Stop all services"
    echo "  restart            Restart all services"
    echo "  status             Show service status"
    "  logs               Show logs"
    "  health             Wait for health check"
    echo "  migrate            Run database migrations"
    "  backup             Backup existing data"
    "  cleanup            Clean up containers and images"
    echo "  test               Run all tests"
    "  reports            Generate compliance reports"
    "  scale              Scale services (usage: scale <service> <replicas>)"
    echo "  help               Show this help message"
    echo ""
    "Options:"
    "  --env ENVIRONMENT  Set environment (default: development)"
    "  --no-backup        Skip backup before deployment"
    "  --timeout SECONDS  Health check timeout (default: 60)"
    echo ""
    "Examples:"
    echo "  $0 deploy --env production"
    "  $0 deploy --no-backup --timeout 120"
    "  $0 scale phoenix-api 3"
    "  $0 status"
    echo "  $0 logs"
}

# Main deployment function
deploy() {
    log_info "Starting Phoenix API deployment..."
    log_info "Environment: $ENVIRONMENT"
    
    check_docker
    check_docker_compose
    
    # Backup existing data
    backup_data
    
    # Build image
    build_image
    
    # Stop existing containers
    stop_containers
    
    # Start services
    start_services
    
    # Run migrations
    run_migrations
    
    # Wait for health check
    if wait_for_health; then
        log_success "Deployment completed successfully!"
        show_status
    else
        log_error "Deployment failed - services not healthy"
        exit 1
    fi
}

# Parse command line arguments
case "${1:-deploy}" in
    build)
        build_image
        ;;
    deploy)
        deploy
        ;;
    stop)
        stop_containers
        ;;
    restart)
        stop_containers
        start_services
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    health)
        wait_for_health
        ;;
    migrate)
        run_migrations
        ;;
    backup)
        backup_data
        ;;
    cleanup)
        cleanup
        ;;
    test)
        run_tests
        ;;
    reports)
        generate_reports
        ;;
    scale)
        scale_services "${2:-phoenix-api}" "${3:-2}"
        ;;
    help)
        show_help
        ;;
    *)
        echo "Unknown command: ${1}"
        show_help
        exit 1
        ;;
esac

# Parse additional options
while [[ "$#" -gt 0 ]]; do
    case "${1}" in
        --env)
            ENVIRONMENT="${2}"
            shift
            ;;
        --no-backup)
            BACKUP_ENABLED=false
            shift
            ;;
        --timeout)
            HEALTH_CHECK_TIMEOUT="${2}"
            shift
            ;;
        *)
            shift
            ;;
    esac
done
