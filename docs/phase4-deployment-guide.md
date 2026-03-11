# Phase 4: Integration Testing & Deployment Guide

## 🎯 Overview

This guide covers the complete deployment process for the Phoenix API with SOC 2 Type II and PCI DSS v4.0 compliance modules. The deployment includes comprehensive testing, containerization, monitoring, and operational procedures.

## 📋 Prerequisites

### System Requirements
- **Docker**: 20.10+ with Docker Compose
- **Memory**: Minimum 4GB RAM (8GB recommended)
- **Storage**: Minimum 10GB free space
- **Network**: Port 8080, 5432, 6379, 9090, 3000 available
- **OS**: Linux, macOS, or Windows with WSL2

### Development Tools
- **Go**: 1.22+
- **Node.js**: 18+ (for frontend)
- **PostgreSQL**: 15+ (for local development)
- **Redis**: 7+ (for local development)

## 🧪 Testing Framework

### Unit Tests
Unit tests cover individual service components and business logic:

```bash
# Run all unit tests
go test ./internal/compliance/soc2/...
go test ./internal/compliance/pci_dss/...

# Run specific test
go test ./internal/compliance/soc2/ -run TestSOC2Service_AssessSOC2Compliance
```

**Coverage Areas:**
- SOC 2 service layer (40+ test cases)
- PCI DSS service layer (40+ test cases)
- Control/requirement lifecycle management
- Assessment scoring algorithms
- Evidence gap identification
- Recommendation generation
- Background monitoring tasks

### Integration Tests
Integration tests verify API endpoints and service interactions:

```bash
# Run API integration tests
go test ./internal/routes/api/v1/...

# Run specific API tests
go test ./internal/routes/api/v1/ -run TestSOC2APITestSuite
```

**Coverage Areas:**
- All SOC 2 API endpoints (15+ tests)
- All PCI DSS API endpoints (18+ tests)
- Request/response validation
- Error handling
- Authentication/authorization
- WebSocket integration

### End-to-End Tests
E2E tests verify complete compliance workflows:

```bash
# Run workflow tests
go test ./test/integration/...

# Run specific workflow
go test ./test/integration/ -run TestComplianceWorkflowTestSuite.TestSOC2CompleteWorkflow
```

**Coverage Areas:**
- Complete SOC 2 compliance workflow
- Complete PCI DSS compliance workflow
- Bulk operations
- Evidence management
- Reporting workflows
- Search and filtering
- Error handling scenarios

## 🐳 Docker Containerization

### Multi-Stage Dockerfile
The Dockerfile uses multi-stage builds for optimal image size and security:

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
# ... build process ...

# Runtime stage
FROM alpine:latest
# ... runtime configuration ...
```

**Security Features:**
- Non-root user execution
- Health checks
- Minimal attack surface
- Security labels
- Volume mounting for logs and reports

### Docker Compose Services
Complete infrastructure setup with monitoring:

```yaml
services:
  phoenix-api:      # Main application
  postgres:         # Database
  redis:            # Cache
  nginx:            # Reverse proxy
  prometheus:       # Metrics collection
  grafana:          # Monitoring dashboard
```

**Service Configuration:**
- **Phoenix API**: Port 8080, health checks, volume mounts
- **PostgreSQL**: Port 5432, persistent storage, initialization scripts
- **Redis**: Port 6379, persistent storage
- **Nginx**: Ports 80/443, SSL termination
- **Prometheus**: Port 9090, metrics collection
- **Grafana**: Port 3000, dashboards and alerts

## 🚀 Deployment Process

### Quick Start
```bash
# Deploy with default settings
./scripts/deploy.sh

# Deploy with specific environment
./scripts/deploy.sh --env production

# Deploy without backup
./scripts/deploy.sh --no-backup
```

### Step-by-Step Deployment

#### 1. Environment Setup
```bash
# Set environment variables
export ENVIRONMENT=production
export PHOENIX_PORT=8080
export PHOENIX_DATABASE_URL=postgres://phoenix:phoenix@postgres:5432/phoenix?sslmode=disable
export PHOENIX_REDIS_URL=redis://redis:6379/0
```

#### 2. Build and Deploy
```bash
# Build Docker image
./scripts/deploy.sh build

# Deploy services
./scripts/deploy.sh deploy
```

#### 3. Database Initialization
```bash
# Run migrations
./scripts/deploy.sh migrate

# Verify database
docker exec phoenix-api-postgres-1 psql -U phoenix -c "\dt"
```

#### 4. Health Verification
```bash
# Wait for services to be healthy
./scripts/deploy.sh health

# Check status
./scripts/deploy.sh status
```

#### 5. Generate Compliance Reports
```bash
# Generate initial compliance reports
./scripts/deploy.sh reports
```

## 📊 Monitoring and Observability

### Prometheus Metrics
The application exposes comprehensive metrics for compliance monitoring:

**Key Metrics:**
- `compliance_soc2_controls_total`: Total SOC 2 controls
- `compliance_soc2_controls_compliant`: Compliant SOC 2 controls
- `compliance_pci_dss_requirements_total`: Total PCI DSS requirements
- `compliance_pci_dss_requirements_compliant`: Compliant PCI DSS requirements
- `compliance_assessments_total`: Total compliance assessments
- `compliance_findings_total`: Total compliance findings

**Access Metrics:**
- Prometheus UI: http://localhost:9090
- Metrics endpoint: http://localhost:8080/metrics

### Grafana Dashboards
Pre-configured dashboards for compliance visualization:

**Dashboard Categories:**
- SOC 2 Compliance Overview
- PCI DSS Compliance Overview
- Assessment Trends
- Findings and Recommendations
- Evidence Coverage
- System Health

**Access Dashboards:**
- Grafana UI: http://localhost:3000
- Default credentials: admin/admin

### Log Management
Structured logging with compliance-specific information:

**Log Categories:**
- API requests/responses
- Compliance assessments
- Control/requirement updates
- Background tasks
- Error conditions

**Log Locations:**
- Application logs: `./logs/app.log`
- Access logs: `./logs/access.log`
- Error logs: `./logs/error.log`

## 🔧 Configuration Management

### Environment Variables
Key configuration options for compliance modules:

```bash
# Compliance Module Configuration
PHOENIX_SOC2_MONITORING_INTERVAL=6h          # SOC 2 monitoring frequency
PHOENIX_PCI_DSS_MONITORING_INTERVAL=6h      # PCI DSS monitoring frequency
PHOENIX_COMPLIANCE_ASSESSMENT_INTERVAL=24h  # Assessment frequency
PHOENIX_ENABLE_BACKGROUND_TASKS=true         # Background task enablement

# Database Configuration
PHOENIX_DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
PHOENIX_REDIS_URL=redis://redis:6379/0

# Application Configuration
PHOENIX_ENV=development
PHOENIX_PORT=8080
```

### Configuration Files
- **Docker Compose**: `docker-compose.yml`
- **Dockerfile**: `Dockerfile`
- **Database Init**: `scripts/init-db.sql`
- **Nginx Config**: `nginx/nginx.conf`
- **Prometheus Config**: `monitoring/prometheus.yml`

## 🔄 Operational Procedures

### Regular Operations

#### Daily Health Checks
```bash
# Check service status
./scripts/deploy.sh status

# Review logs
./scripts/deploy.sh logs

# Check compliance scores
curl http://localhost:8080/api/v1/compliance/soc2/stats
curl http://localhost:8080/api/v1/compliance/pci_dss/stats
```

#### Weekly Compliance Reviews
```bash
# Generate compliance reports
./scripts/deploy.sh reports

# Review assessment results
curl http://localhost:8080/api/v1/compliance/soc2/assessments
curl http://localhost:8080/api/v1/compliance/pci_dss/assessments

# Check for new findings
curl http://localhost:8080/api/v1/compliance/soc2/search -d '{"query": "critical"}'
```

#### Monthly Maintenance
```bash
# Update dependencies
go mod tidy
docker-compose pull

# Restart services
./scripts/deploy.sh restart

# Backup data
./scripts/deploy.sh backup
```

### Incident Response

#### Service Unavailability
```bash
# Check logs for errors
./scripts/deploy.sh logs | grep ERROR

# Restart services
./scripts/deploy.sh restart

# Verify health
./scripts/deploy.sh health
```

#### Database Issues
```bash
# Check database connectivity
docker exec phoenix-api-postgres-1 pg_isready -U phoenix

# Review database logs
docker logs phoenix-api-postgres-1

# Restart database
docker-compose restart postgres
```

#### Performance Issues
```bash
# Check system resources
docker stats

# Review Prometheus metrics
curl http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])

# Scale services if needed
./scripts/deploy.sh scale phoenix-api 3
```

## 📈 Scaling and Performance

### Horizontal Scaling
```bash
# Scale API services
./scripts/deploy.sh scale phoenix-api 3

# Scale database (read replicas)
# Add to docker-compose.yml:
# postgres-replica:
#   image: postgres:15-alpine
#   ...
```

### Performance Optimization

#### Database Optimization
- Connection pooling
- Query optimization
- Index tuning
- Partitioning for large tables

#### Application Optimization
- Caching strategies
- Background task optimization
- Memory management
- Concurrent processing

#### Monitoring Optimization
- Metric retention policies
- Alert tuning
- Dashboard optimization
- Log rotation

## 🔒 Security Considerations

### Container Security
- Non-root user execution
- Minimal base images
- Security scanning
- Vulnerability management

### Network Security
- SSL/TLS encryption
- Network segmentation
- Firewall rules
- Access control

### Data Security
- Encryption at rest
- Encryption in transit
- Access logging
- Data retention policies

### Compliance Security
- SOC 2 Type II controls
- PCI DSS v4.0 requirements
- Audit logging
- Evidence collection

## 🚨 Troubleshooting

### Common Issues

#### Container Startup Issues
```bash
# Check container logs
docker logs phoenix-api

# Verify image build
docker images | grep phoenix-api

# Rebuild if needed
./scripts/deploy.sh build
```

#### Database Connection Issues
```bash
# Check database status
docker exec phoenix-api-postgres-1 pg_isready -U phoenix

# Test connection
docker exec phoenix-api-postgres-1 psql -U phoenix -c "SELECT 1;"

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

#### Performance Issues
```bash
# Check resource usage
docker stats

# Review application logs
./scripts/deploy.sh logs | grep "slow"

# Monitor metrics
curl http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])
```

### Debug Mode
```bash
# Enable debug logging
export PHOENIX_LOG_LEVEL=debug

# Run with verbose output
./scripts/deploy.sh --env development

# Access debug endpoints
curl http://localhost:8080/api/system/health?debug=true
```

## 📚 API Reference

### SOC 2 Compliance API
- `GET /api/v1/compliance/soc2/controls` - List SOC 2 controls
- `POST /api/v1/compliance/soc2/assess` - Conduct SOC 2 assessment
- `GET /api/v1/compliance/soc2/stats` - Get SOC 2 statistics
- `POST /api/v1/compliance/soc2/report` - Generate SOC 2 report

### PCI DSS Compliance API
- `GET /api/v1/compliance/pci_dss/requirements` - List PCI DSS requirements
- `POST /api/v1/compliance/pci_dss/assess` - Conduct PCI DSS assessment
- `GET /api/v1/compliance/pci_dss/stats` - Get PCI DSS statistics
- `POST /api/v1/compliance/pci_dss/report` - Generate PCI DSS report

### System API
- `GET /api/system/health` - Health check
- `GET /api/metrics` - Prometheus metrics
- `GET /api/version` - Application version

## 📋 Deployment Checklist

### Pre-Deployment
- [ ] Review environment variables
- [ ] Verify database configuration
- [ ] Check SSL certificates
- [ ] Validate monitoring setup
- [ ] Test backup procedures

### Deployment
- [ ] Backup existing data
- [ ] Build new Docker image
- [ ] Stop existing services
- [ ] Start new services
- [ ] Run database migrations
- [ ] Verify health checks
- [ ] Generate compliance reports

### Post-Deployment
- [ ] Verify service status
- [ ] Check monitoring dashboards
- [ ] Review application logs
- [ ] Test API endpoints
- [ ] Validate compliance functionality
- [ ] Document deployment

## 🎯 Success Metrics

### Deployment Success
- All services healthy
- Database migrations complete
- API endpoints responding
- Monitoring functional
- Compliance reports generated

### Performance Metrics
- API response time < 200ms
- Database query time < 100ms
- Memory usage < 512MB
- CPU usage < 50%
- Error rate < 1%

### Compliance Metrics
- SOC 2 compliance score > 80%
- PCI DSS compliance score > 80%
- Evidence coverage > 90%
- Findings addressed within SLA
- Reports generated on schedule

## 📞 Support and Escalation

### Documentation
- API documentation: `/docs/api/`
- Architecture guide: `/docs/architecture/`
- Troubleshooting guide: `/docs/troubleshooting/`

### Monitoring
- Grafana dashboards: http://localhost:3000
- Prometheus metrics: http://localhost:9090
- Application logs: `./logs/`

### Contact
- Development team: dev@barca-strategos.com
- Operations team: ops@barca-strategos.com
- Security team: security@barca-strategos.com

---

**Phase 4 Status**: ✅ Complete

The Phoenix API with SOC 2 Type II and PCI DSS v4.0 compliance modules is now fully deployed and operational. The system includes comprehensive testing, monitoring, and operational procedures for enterprise-grade compliance management.
