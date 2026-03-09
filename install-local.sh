#!/bin/bash

# ──────────────────────────────────────────────────────────────────────────────────────
# Barca-Strategos Phoenix - Local Ubuntu Installation Script
# Installs Phoenix GUI system natively without Docker
# ──────────────────────────────────────────────────────────────────────────────────────

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print banner
print_banner() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                    BARCA-STRATEGOS PHOENIX                   ║"
    echo "║              Local Ubuntu Installation (No Docker)               ║"
    echo "║                     🚀 Native Setup 🚀                    ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        echo -e "${YELLOW}⚠️  Running as root. Some services will be created with proper users.${NC}"
    else
        echo -e "${RED}❌ This script needs to be run with sudo for system-wide installation${NC}"
        exit 1
    fi
}

# Check system requirements
check_requirements() {
    echo -e "${BLUE}🔍 Checking system requirements...${NC}"
    
    # Check memory
    MEMORY=$(free -m | awk 'NR==2{printf "%.0f", $2}')
    if [[ $MEMORY -lt 2048 ]]; then
        echo -e "${RED}❌ System requires at least 2GB RAM (found: ${MEMORY}MB)${NC}"
        exit 1
    fi
    
    # Check disk space
    DISK=$(df / | awk 'NR==2{print $4}')
    if [[ $DISK -lt 10485760 ]]; then
        echo -e "${RED}❌ System requires at least 10GB disk space (found: $(($DISK/1024/1024))GB)${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ System requirements met (${MEMORY}MB RAM, $(($DISK/1024/1024))GB disk)${NC}"
}

# Install system dependencies
install_system_deps() {
    echo -e "${BLUE}📦 Installing system dependencies...${NC}"
    
    apt-get update
    apt-get install -y \
        build-essential \
        pkg-config \
        libssl-dev \
        libpq-dev \
        postgresql \
        postgresql-contrib \
        redis-server \
        nginx \
        curl \
        wget \
        git \
        htop \
        unzip \
        supervisor \
        ufw
    
    echo -e "${GREEN}✅ System dependencies installed${NC}"
}

# Install Rust
install_rust() {
    echo -e "${BLUE}🦀 Installing Rust...${NC}"
    
    if ! command -v rustc &> /dev/null; then
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
        source "$HOME/.cargo/env"
        echo -e "${GREEN}✅ Rust installed${NC}"
    else
        echo -e "${YELLOW}⚠️  Rust already installed, updating...${NC}"
        rustup update
    fi
}

# Create Phoenix user
create_phoenix_user() {
    echo -e "${BLUE}👤 Creating Phoenix user...${NC}"
    
    if ! id "phoenix" &>/dev/null; then
        useradd -r -s /bin/bash -m -d /opt/phoenix phoenix
        echo -e "${GREEN}✅ Phoenix user created${NC}"
    else
        echo -e "${YELLOW}⚠️  Phoenix user already exists${NC}"
    fi
}

# Create directory structure
create_directories() {
    echo -e "${BLUE}📁 Creating directory structure...${NC}"
    
    mkdir -p /opt/phoenix/{bin,config,data,logs,static,backups}
    chown -R phoenix:phoenix /opt/phoenix
    
    echo -e "${GREEN}✅ Directory structure created${NC}"
}

# Clone and build Phoenix
build_phoenix() {
    echo -e "${BLUE}🔨 Building Phoenix...${NC}"
    
    # Copy source code to Phoenix directory
    cp -r . /opt/phoenix/src/
    chown -R phoenix:phoenix /opt/phoenix
    
    # Set up Rust environment for phoenix user
    sudo -u phoenix bash -c "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y"
    
    # Build as Phoenix user with proper environment (online build)
    sudo -u phoenix bash -c "cd /opt/phoenix/src && source /opt/phoenix/.cargo/env && cargo build --release"
    
    # Copy binaries
    cp /opt/phoenix/src/target/release/phoenix-core /opt/phoenix/bin/
    chmod +x /opt/phoenix/bin/phoenix-core
    
    echo -e "${GREEN}✅ Phoenix built successfully${NC}"
}

# Configure PostgreSQL
setup_postgresql() {
    echo -e "${BLUE}🐘 Setting up PostgreSQL...${NC}"
    
    # Start PostgreSQL
    systemctl start postgresql
    systemctl enable postgresql
    
    # Create database and user
    sudo -u postgres psql -c "CREATE USER phoenix WITH PASSWORD 'phoenix123';"
    sudo -u postgres psql -c "CREATE DATABASE phoenix OWNER phoenix;"
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE phoenix TO phoenix;"
    
    echo -e "${GREEN}✅ PostgreSQL configured${NC}"
}

# Configure Redis
setup_redis() {
    echo -e "${BLUE}🔴 Setting up Redis...${NC}"
    
    # Start Redis
    systemctl start redis-server
    systemctl enable redis-server
    
    echo -e "${GREEN}✅ Redis configured${NC}"
}

# Configure Nginx
setup_nginx() {
    echo -e "${BLUE}🌐 Setting up Nginx...${NC}"
    
    # Create Nginx config
    cat > /etc/nginx/sites-available/phoenix << 'EOF'
server {
    listen 80;
    server_name _;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF
    
    # Enable site
    ln -sf /etc/nginx/sites-available/phoenix /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default
    
    # Test and restart Nginx
    nginx -t
    systemctl restart nginx
    systemctl enable nginx
    
    echo -e "${GREEN}✅ Nginx configured${NC}"
}

# Create environment file
create_env_file() {
    echo -e "${BLUE}🔐 Creating environment configuration...${NC}"
    
    cat > /opt/phoenix/.env << 'EOF'
# Phoenix Configuration
PHOENIX_HOST=127.0.0.1
PHOENIX_PORT=8080
PHOENIX_STATIC_PATH=/opt/phoenix/static
PHOENIX_MAX_CONNECTIONS=1000

# Database Configuration
DATABASE_URL=postgresql://phoenix:phoenix123@localhost/phoenix
DB_PASSWORD=phoenix123
POSTGRES_DB=phoenix
POSTGRES_USER=phoenix

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=

# Security
PHOENIX_JWT_SECRET=$(openssl rand -base64 32)
RUST_LOG=info

# Paths
PHOENIX_DATA_DIR=/opt/phoenix/data
PHOENIX_LOG_DIR=/opt/phoenix/logs
PHOENIX_CONFIG_DIR=/opt/phoenix/config
EOF
    
    chown phoenix:phoenix /opt/phoenix/.env
    chmod 600 /opt/phoenix/.env
    
    echo -e "${GREEN}✅ Environment file created${NC}"
}

# Create systemd service
create_systemd_service() {
    echo -e "${BLUE}⚙️  Creating systemd service...${NC}"
    
    cat > /etc/systemd/system/phoenix.service << 'EOF'
[Unit]
Description=Barca-Strategos Phoenix GUI
After=network.target postgresql.service redis-server.service

[Service]
Type=simple
User=phoenix
Group=phoenix
WorkingDirectory=/opt/phoenix
Environment=PATH=/opt/phoenix/bin:/usr/local/bin:/usr/bin:/bin
EnvironmentFile=/opt/phoenix/.env
ExecStart=/opt/phoenix/bin/phoenix-core
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
    
    # Reload systemd and enable service
    systemctl daemon-reload
    systemctl enable phoenix
    
    echo -e "${GREEN}✅ Systemd service created${NC}"
}

# Configure firewall
setup_firewall() {
    echo -e "${BLUE}🔥 Configuring firewall...${NC}"
    
    ufw --force enable
    ufw allow ssh
    ufw allow 80/tcp
    ufw allow 443/tcp
    ufw allow 8080/tcp
    
    echo -e "${GREEN}✅ Firewall configured${NC}"
}

# Start Phoenix service
start_phoenix() {
    echo -e "${BLUE}🚀 Starting Phoenix...${NC}"
    
    systemctl start phoenix
    sleep 3
    
    if systemctl is-active --quiet phoenix; then
        echo -e "${GREEN}✅ Phoenix started successfully${NC}"
    else
        echo -e "${RED}❌ Failed to start Phoenix${NC}"
        systemctl status phoenix
        exit 1
    fi
}

# Show installation info
show_installation_info() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                    🎉 INSTALLATION COMPLETE 🎉                    ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    echo -e "${GREEN}🌐 Phoenix GUI: http://$(hostname -I | awk '{print $1}'):80${NC}"
    echo -e "${GREEN}🔧 Direct Access: http://$(hostname -I | awk '{print $1}'):8080${NC}"
    echo -e "${GREEN}📊 Status: sudo systemctl status phoenix${NC}"
    echo -e "${GREEN}📝 Logs: sudo journalctl -u phoenix -f${NC}"
    echo -e "${GREEN}🔄 Restart: sudo systemctl restart phoenix${NC}"
    
    echo -e "${YELLOW}"
    echo "📋 Services running:"
    echo "  • Phoenix GUI (port 8080)"
    echo "  • PostgreSQL (port 5432)"
    echo "  • Redis (port 6379)"
    echo "  • Nginx (port 80)"
    echo -e "${NC}"
}

# Main installation function
main() {
    print_banner
    check_root
    check_requirements
    install_system_deps
    install_rust
    create_phoenix_user
    create_directories
    build_phoenix
    setup_postgresql
    setup_redis
    setup_nginx
    create_env_file
    create_systemd_service
    setup_firewall
    start_phoenix
    show_installation_info
}

# Run main function
main "$@"
