# 🐳 Docker Deployment Guide

## Overview

This guide covers deploying the Barca-Strategos Phoenix GUI with Docker and Docker Compose. The application provides a comprehensive web-based cognitive collaboration platform with real-time monitoring, compliance management, and AI agent interaction.

## 🚀 Quick Start

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- 2GB+ RAM
- 2+ CPU cores

### One-Command Deployment

```bash
# Clone and start the application
git clone <repository-url>
cd Barca-Strategos
make quick-start
```

This will:

1. Build the Docker image
2. Start all services with Docker Compose  
3. Launch the web interface at `http://localhost:8080`

## 📋 Services Overview

### Core Services

| Service      | Port  | Description |
|--------------|-------|-------------|
| `phoenix-gui`| 8080  | Main web application with GUI |
| `postgres`   | 5432  | PostgreSQL database |
| `redis`      | 6379  | Redis cache and session storage |

### Optional Services

| Service     | Port    | Profile     | Description |
|-------------|---------|-------------|-------------|
| `nginx`     | 80, 443 | `production`| Reverse proxy and SSL termination |
| `prometheus`| 9090    | `monitoring` | Metrics collection |
| `grafana`   | 3000    | `monitoring` | Visualization dashboard |
| `jaeger`    | 16686   | `monitoring` | Distributed tracing |

## 🔧 Configuration

### Environment Variables

Create a `.env` file:

```bash
# Database Configuration
POSTGRES_DB=phoenix
POSTGRES_USER=phoenix
POSTGRES_PASSWORD=your_secure_password

# Redis Configuration
REDIS_PASSWORD=your_redis_password

# Application Configuration
PHOENIX_HOST=0.0.0.0
PHOENIX_PORT=8080
PHOENIX_MAX_CONNECTIONS=1000
PHOENIX_CORS_ENABLED=true

# Security
JWT_SECRET=your_jwt_secret_here
RUST_LOG=info

# Optional: AI Integration
API_KEY=your_ai_api_key
MODEL_PROVIDER=openai
```

### Production Configuration

For production, create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  phoenix-gui:
    environment:
      - RUST_LOG=warn
      - PHOENIX_MAX_CONNECTIONS=5000
    deploy:
      replicas: 3
      resources:
        limits:
          memory: 512M
          cpus: '2.0'

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - phoenix-gui
```

## 🏗️ Build and Deployment

### Development Build

```bash
# Build for development
make docker-build

# Run development environment
make docker-run
```

### Production Build

```bash
# Build optimized production image
docker build -t phoenix-gui:prod --target production .

# Deploy with production configuration
make prod
```

### Custom Build

```bash
# Build with custom tag
docker build -t phoenix-gui:custom .

# Build with build arguments
docker build \
  --build-arg RUST_VERSION=1.75 \
  --build-arg APP_VERSION=1.0.0 \
  -t phoenix-gui:custom .
```

## 📊 Monitoring and Logging

### Viewing Logs

```bash
# View all logs
make docker-logs

# View specific service logs
docker-compose logs -f phoenix-gui
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Health Checks

```bash
# Check system status
make status

# Manual health check
curl http://localhost:8080/api/system/health
```

### Monitoring Stack

Enable monitoring with:

```bash
# Start with monitoring services
docker-compose --profile monitoring up -d

# Access Grafana (admin/admin)
open http://localhost:3000

# Access Prometheus
open http://localhost:9090

# Access Jaeger
open http://localhost:16686
```

## 🗄️ Database Management

### Database Migrations

```bash
# Run migrations
make db-migrate

# Reset database
make db-reset
```

### Backup and Restore

1. Create backup

```bash
make backup
```

2. Restore from backup

```bash
make restore BACKUP_FILE=backup_20240309_120000.sql
```

### Manual Database Access

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U phoenix -d phoenix

# Connect to Redis
docker-compose exec redis redis-cli
```

## 🔒 Security Configuration

### SSL/TLS Setup

1. Create SSL certificates:

```bash
mkdir -p nginx/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/key.pem \
  -out nginx/ssl/cert.pem
```

2. Configure nginx:

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    location / {
        proxy_pass http://phoenix-gui:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Security Best Practices

1. **Use strong passwords** for database and Redis
2. **Enable HTTPS** in production
3. **Limit network exposure** with proper firewall rules
4. **Regularly update** base images
5. **Monitor logs** for suspicious activity

## 📈 Performance Optimization

### Resource Allocation

Adjust resource limits in `docker-compose.yml`:

```yaml
services:
  phoenix-gui:
    deploy:
      resources:
        limits:
          memory: 512M    # Increase for high load
          cpus: '2.0'     # Add more CPU for concurrency
        reservations:
          memory: 256M
          cpus: '1.0'
```

### Caching Configuration

Optimize Redis for caching:

```yaml
services:
  redis:
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
```

### Database Optimization

PostgreSQL tuning for production:

```yaml
services:
  postgres:
    environment:
      - POSTGRES_SHARED_PRELOAD_LIBRARIES=pg_stat_statements
    command: >
      postgres
      -c shared_preload_libraries=pg_stat_statements
      -c max_connections=200
      -c shared_buffers=128MB
      -c effective_cache_size=512MB
```

## 🔄 Scaling and High Availability

### Horizontal Scaling

```yaml
services:
  phoenix-gui:
    deploy:
      replicas: 3  # Run 3 instances
    environment:
      - PHOENIX_INSTANCE_ID={{.Task.Slot}}
```

### Load Balancing with Nginx

```nginx
upstream phoenix_backend {
    server phoenix-gui-1:8080;
    server phoenix-gui-2:8080;
    server phoenix-gui-3:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://phoenix_backend;
    }
}
```

## 🛠️ Troubleshooting

### Common Issues

1. **Port conflicts**:
   ```bash
   # Check what's using port 8080
   lsof -i :8080
   # Kill the process
   kill -9 <PID>
   ```

2. **Memory issues**:
   ```bash
   # Check memory usage
   docker stats
   # Increase memory limits
   ```

3. **Database connection errors**:
   ```bash
   # Check database status
   docker-compose exec postgres pg_isready
   # Reset database if needed
   make db-reset
   ```

### Debug Mode

Enable debug logging:

```bash
# Run with debug logs
RUST_LOG=debug docker-compose up phoenix-gui

# View detailed logs
docker-compose logs -f phoenix-gui | grep DEBUG
```

### Performance Issues

1. **Monitor resource usage**:
   ```bash
   docker stats
   htop
   ```

2. **Check database performance**:
   ```sql
   SELECT * FROM pg_stat_activity;
   ```

3. **Profile the application**:
   ```bash
   make profile
   ```

## 📋 Maintenance

### Regular Tasks

```bash
# Update images
docker-compose pull
docker-compose up -d

# Clean up unused resources
docker system prune -f

# Backup data
make backup

# Check logs for errors
docker-compose logs --tail=100 phoenix-gui | grep ERROR
```

### Health Monitoring

Set up automated health checks:

```bash
#!/bin/bash
# health-check.sh
if ! curl -f http://localhost:8080/api/system/health; then
    echo "Health check failed!" | mail -s "Phoenix Alert" admin@example.com
    docker-compose restart phoenix-gui
fi
```

### Log Rotation

Configure log rotation in `docker-compose.yml`:

```yaml
services:
  phoenix-gui:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 🚀 Production Deployment

### Pre-deployment Checklist

- [ ] Environment variables configured
- [ ] SSL certificates installed
- [ ] Database backup created
- [ ] Resource limits set appropriately
- [ ] Monitoring configured
- [ ] Log rotation configured
- [ ] Security scan completed

### Deployment Steps

```bash
# 1. Prepare environment
make setup-prod

# 2. Build and test
make ci-build
make ci-test

# 3. Deploy
make prod

# 4. Verify deployment
make status
curl http://localhost:8080/api/system/health
```

### Rollback Procedure

```bash
# Stop current deployment
docker-compose down

# Switch to previous version
docker tag phoenix-gui:previous phoenix-gui:current

# Restart with previous version
docker-compose up -d

# Verify rollback
make status
```

## 📞 Support

For issues with Docker deployment:

1. Check the troubleshooting section
2. Review logs with `make docker-logs`
3. Verify configuration in `.env`
4. Check resource availability
5. Consult the main documentation

Additional resources:
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Application Documentation](README.md)
