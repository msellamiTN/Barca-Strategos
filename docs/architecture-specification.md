# Barca-Strategos Phoenix: Architecture Specification

## Document Overview

**Version**: 1.0  
**Date**: March 2026  
**Status**: Post-Refactor Baseline  

This document provides a comprehensive architectural overview of the Barca-Strategos Phoenix platform using Mermaid diagrams for visual representation of system components, data flows, and deployment topology.

---

## 1. High-Level System Architecture

The Phoenix platform follows a **Microservices + Multi-Agent** architecture pattern, enabling independent scaling and clear domain separation.

```mermaid
graph TB
    subgraph "Presentation Layer"
        WEB[Web GUI - React SPA]
        TG[Telegram Bot]
        DC[Discord Bot]
        SL[Slack Bot]
        TM[Teams Bot]
    end

    subgraph "API Gateway Layer"
        GW[Axum Web Server]
        AUTH[Auth Middleware]
        CORS[CORS Handler]
        WS[WebSocket Handler]
    end

    subgraph "Core Services"
        CORE[Phoenix Core]
        AI[AI Assistant]
        COMP[Compliance Engine]
        SEC[Security Module]
        RISK[Risk Management]
        MON[Monitoring]
    end

    subgraph "Data Layer"
        PG[(PostgreSQL)]
        RD[(Redis Cache)]
        KF[Kafka/Queue]
    end

    WEB --> GW
    TG --> GW
    DC --> GW
    SL --> GW
    TM --> GW

    GW --> AUTH
    AUTH --> CORS
    CORS --> WS

    WS --> CORE
    WS --> AI
    WS --> COMP
    WS --> SEC
    WS --> RISK
    WS --> MON

    CORE --> PG
    CORE --> RD
    CORE --> KF
    AI --> PG
    COMP --> PG
    SEC --> PG
    RISK --> PG
    MON --> RD
```

---

## 2. Module Architecture

### 2.1 Core Module Dependencies

```mermaid
graph LR
    subgraph "src/lib.rs"
        LIB[Library Root]
    end

    subgraph "Primary Modules"
        GUI[gui]
        SEC[security]
        COMP[compliance]
        MON[monitoring]
    end

    LIB --> GUI
    LIB --> SEC
    LIB --> COMP
    LIB --> MON

    subgraph "GUI Submodules"
        DASH[dashboard]
        SECOPS[security_operations]
        COMPCTR[compliance_center]
        COGCOL[cognitive_collaboration]
        AGINT[agent_interaction]
        RISKMGT[risk_management]
        RTM[real_time_monitoring]
        WEBSRV[web_server]
    end

    GUI --> DASH
    GUI --> SECOPS
    GUI --> COMPCTR
    GUI --> COGCOL
    GUI --> AGINT
    GUI --> RISKMGT
    GUI --> RTM
    GUI --> WEBSRV

    subgraph "Compliance Frameworks"
        ISO[iso27001]
        NIST[nist_csf]
        GDPR[gdpr]
        SOC2[soc2]
        PCI[pci_dss]
    end

    COMP --> ISO
    COMP --> NIST
    COMP --> GDPR
    COMP --> SOC2
    COMP --> PCI
```

### 2.2 Compliance Module Detail

```mermaid
classDiagram
    class ComplianceEngine {
        +assess_compliance()
        +generate_report()
        +track_controls()
    }

    class ISO27001Compliance {
        +controls: Vec~Control~
        +assess_isms()
        +get_control_status()
    }

    class NISTCSFCompliance {
        +functions: Vec~CSFFunction~
        +assess_csf()
        +get_maturity_level()
    }

    class GDPRCompliance {
        +articles: Vec~Article~
        +assess_gdpr()
        +track_data_subjects()
        +manage_consent()
        +handle_breach()
    }

    class SOC2Compliance {
        +trust_criteria: Vec~SOC2Control~
        +assess_soc2()
        +is_control_in_scope()
        +trigger_soc2_alert()
    }

    class PCIDSSCompliance {
        +requirements: Vec~Requirement~
        +assess_pci()
        +scan_cardholder_env()
    }

    ComplianceEngine --> ISO27001Compliance
    ComplianceEngine --> NISTCSFCompliance
    ComplianceEngine --> GDPRCompliance
    ComplianceEngine --> SOC2Compliance
    ComplianceEngine --> PCIDSSCompliance
```

---

## 3. Data Flow Architecture

### 3.1 Request Processing Flow

```mermaid
sequenceDiagram
    participant U as User/Bot
    participant GW as API Gateway
    participant AUTH as Auth Middleware
    participant SVC as Service Layer
    participant DB as PostgreSQL
    participant CACHE as Redis
    participant WS as WebSocket

    U->>GW: HTTP/WS Request
    GW->>AUTH: Validate JWT
    AUTH-->>GW: Auth Result

    alt Authenticated
        GW->>SVC: Route to Service
        SVC->>CACHE: Check Cache
        
        alt Cache Hit
            CACHE-->>SVC: Cached Data
        else Cache Miss
            SVC->>DB: Query Database
            DB-->>SVC: Data
            SVC->>CACHE: Update Cache
        end

        SVC-->>GW: Response
        GW-->>U: HTTP Response
        
        opt Real-time Update
            SVC->>WS: Broadcast Update
            WS-->>U: WebSocket Message
        end
    else Not Authenticated
        GW-->>U: 401 Unauthorized
    end
```

### 3.2 Compliance Assessment Flow

```mermaid
flowchart TD
    START([Start Assessment]) --> SELECT{Select Framework}
    
    SELECT -->|ISO 27001| ISO[Load ISO Controls]
    SELECT -->|NIST CSF| NIST[Load CSF Functions]
    SELECT -->|GDPR| GDPR[Load GDPR Articles]
    SELECT -->|SOC 2| SOC2[Load Trust Criteria]
    SELECT -->|PCI DSS| PCI[Load Requirements]

    ISO --> EVAL[Evaluate Controls]
    NIST --> EVAL
    GDPR --> EVAL
    SOC2 --> EVAL
    PCI --> EVAL

    EVAL --> SCORE[Calculate Compliance Score]
    SCORE --> GAPS[Identify Gaps]
    GAPS --> REPORT[Generate Report]
    REPORT --> STORE[(Store in DB)]
    STORE --> NOTIFY[Notify Stakeholders]
    NOTIFY --> END([End])
```

---

## 4. Security Architecture

### 4.1 Security Operations Center Flow

```mermaid
flowchart LR
    subgraph "Threat Sources"
        NET[Network Traffic]
        LOG[System Logs]
        VULN[Vulnerability Scans]
        INTEL[Threat Intelligence]
    end

    subgraph "Detection Layer"
        IDS[Intrusion Detection]
        SIEM[SIEM Integration]
        ML[ML Anomaly Detection]
    end

    subgraph "Response Layer"
        ALERT[Alert Generation]
        TRIAGE[Threat Triage]
        INC[Incident Response]
        FORENSIC[Forensic Analysis]
    end

    subgraph "Output"
        DASH[Security Dashboard]
        REPORT[Security Reports]
        METRICS[Security Metrics]
    end

    NET --> IDS
    LOG --> SIEM
    VULN --> SIEM
    INTEL --> ML

    IDS --> ALERT
    SIEM --> ALERT
    ML --> ALERT

    ALERT --> TRIAGE
    TRIAGE --> INC
    INC --> FORENSIC

    TRIAGE --> DASH
    INC --> REPORT
    FORENSIC --> METRICS
```

### 4.2 Authentication & Authorization

```mermaid
flowchart TD
    REQ[Incoming Request] --> JWT{Has JWT?}
    
    JWT -->|No| LOGIN[Login Required]
    LOGIN --> CRED[Validate Credentials]
    CRED -->|Valid| ISSUE[Issue JWT]
    CRED -->|Invalid| DENY1[Access Denied]
    ISSUE --> CACHE[Cache Session in Redis]

    JWT -->|Yes| VERIFY[Verify JWT Signature]
    VERIFY -->|Invalid| DENY2[Access Denied]
    VERIFY -->|Valid| EXPIRE{Token Expired?}
    
    EXPIRE -->|Yes| REFRESH{Refresh Token?}
    REFRESH -->|Yes| RENEW[Renew JWT]
    REFRESH -->|No| LOGIN
    RENEW --> CONTINUE

    EXPIRE -->|No| RBAC[Check RBAC Permissions]
    RBAC -->|Authorized| CONTINUE[Continue to Service]
    RBAC -->|Unauthorized| DENY3[403 Forbidden]
```

---

## 5. Deployment Architecture

### 5.1 Docker Compose Topology

```mermaid
graph TB
    subgraph "External"
        USER[Users/Clients]
        BOTS[Chat Platforms]
    end

    subgraph "Load Balancer"
        NGINX[Nginx Reverse Proxy]
    end

    subgraph "Application Tier"
        PHOENIX[phoenix-core:8080]
        TG_BOT[telegram-bot]
        DC_BOT[discord-bot]
        SL_BOT[slack-bot]
        TM_BOT[teams-bot]
    end

    subgraph "Data Tier"
        PG_M[(PostgreSQL Master)]
        PG_S[(PostgreSQL Slave)]
        REDIS[(Redis Cluster)]
    end

    subgraph "Message Tier"
        ZK[Zookeeper]
        KAFKA[Kafka Broker]
    end

    subgraph "Observability"
        PROM[Prometheus]
        GRAF[Grafana]
        JAEGER[Jaeger Tracing]
    end

    USER --> NGINX
    BOTS --> TG_BOT
    BOTS --> DC_BOT
    BOTS --> SL_BOT
    BOTS --> TM_BOT

    NGINX --> PHOENIX
    
    TG_BOT --> PHOENIX
    DC_BOT --> PHOENIX
    SL_BOT --> PHOENIX
    TM_BOT --> PHOENIX

    PHOENIX --> PG_M
    PG_M --> PG_S
    PHOENIX --> REDIS
    PHOENIX --> KAFKA
    KAFKA --> ZK

    PHOENIX --> PROM
    PROM --> GRAF
    PHOENIX --> JAEGER
```

### 5.2 Local Installation Architecture

```mermaid
graph TB
    subgraph "Ubuntu Server"
        subgraph "Systemd Services"
            PHX[phoenix-core.service]
        end

        subgraph "System Services"
            PG[PostgreSQL 15]
            RD[Redis Server]
            NGX[Nginx]
        end

        subgraph "Directories"
            OPT[/opt/phoenix]
            BIN[/opt/phoenix/bin]
            SRC[/opt/phoenix/src]
            LOG[/var/log/phoenix]
        end

        subgraph "User"
            USR[phoenix user]
        end
    end

    subgraph "External Access"
        HTTP[HTTP :80]
        HTTPS[HTTPS :443]
        API[API :8080]
    end

    HTTP --> NGX
    HTTPS --> NGX
    NGX --> PHX
    PHX --> PG
    PHX --> RD
    USR --> PHX
    PHX --> BIN
    PHX --> LOG
```

---

## 6. Technology Stack Summary

```mermaid
mindmap
    root((Phoenix Stack))
        Backend
            Rust 1.94+
            Axum Web Framework
            Tokio Async Runtime
            SQLx PostgreSQL
            Redis Client
        Frontend
            React SPA
            WebSocket
            Tailwind CSS
        Security
            Ring Crypto
            JWT Auth
            Argon2 Hashing
            rustls TLS
        AI/ML
            OpenAI API
            tiktoken-rs
        Integrations
            Telegram Bot
            Discord Bot
            Slack Bot
            Teams Bot
        Infrastructure
            Docker
            Nginx
            PostgreSQL
            Redis
            Kafka
        Observability
            Prometheus
            Grafana
            Jaeger
            Tracing
```

---

## 7. API Endpoints Architecture

```mermaid
graph LR
    subgraph "REST API"
        subgraph "Dashboard"
            D1[GET /api/dashboard]
            D2[GET /api/dashboard/metrics]
        end

        subgraph "Security"
            S1[GET /api/security/threats]
            S2[POST /api/security/incidents]
            S3[GET /api/security/vulnerabilities]
        end

        subgraph "Compliance"
            C1[GET /api/compliance/status]
            C2[POST /api/compliance/assess]
            C3[GET /api/compliance/reports]
        end

        subgraph "Risk"
            R1[GET /api/risk/items]
            R2[POST /api/risk/assess]
            R3[PUT /api/risk/mitigate]
        end

        subgraph "Agents"
            A1[GET /api/agents/status]
            A2[POST /api/agents/scale]
            A3[GET /api/agents/tasks]
        end
    end

    subgraph "WebSocket API"
        WS1[/ws - Main Connection]
        WS2[/ws/collab - Collaboration]
        WS3[/ws/metrics - Live Metrics]
        WS4[/ws/alerts - Security Alerts]
    end

    subgraph "Health"
        H1[GET /api/system/health]
    end
```

---

## 8. State Management

```mermaid
stateDiagram-v2
    [*] --> Initializing: System Start

    Initializing --> Ready: All Services Up
    Initializing --> Degraded: Partial Services

    Ready --> Processing: Request Received
    Processing --> Ready: Request Complete
    Processing --> Error: Processing Failed

    Ready --> Scaling: Load Threshold
    Scaling --> Ready: Scale Complete

    Ready --> Maintenance: Admin Action
    Maintenance --> Ready: Maintenance Complete

    Degraded --> Ready: Services Recovered
    Degraded --> Critical: More Failures

    Error --> Ready: Error Handled
    Error --> Critical: Unrecoverable

    Critical --> [*]: System Shutdown
```

---

## Appendix: File Structure

```
barca-strategos/
├── src/
│   ├── main.rs              # Application entry point
│   ├── lib.rs               # Library root (gui, security, compliance, monitoring)
│   ├── gui/
│   │   ├── mod.rs
│   │   ├── dashboard.rs
│   │   ├── security_operations.rs
│   │   ├── compliance_center.rs
│   │   ├── cognitive_collaboration.rs
│   │   ├── agent_interaction.rs
│   │   ├── risk_management.rs
│   │   ├── real_time_monitoring.rs
│   │   └── web_server.rs
│   ├── compliance/
│   │   ├── mod.rs
│   │   ├── iso27001.rs
│   │   ├── nist_csf.rs
│   │   ├── gdpr.rs
│   │   ├── soc2.rs
│   │   └── pci_dss.rs
│   ├── security/
│   │   └── mod.rs
│   └── monitoring/
│       └── mod.rs
├── Cargo.toml
├── Dockerfile
├── docker-compose.yml
├── install-local.sh
└── docs/
    ├── report.md
    ├── claude.md
    └── architecture-specification.md
```

---

*Document generated from Barca-Strategos Phoenix Comprehensive Analysis Report*
