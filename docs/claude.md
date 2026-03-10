# Professional Project Analysis & Enhancement Prompt
## Multi-Expert Framework: Architecture | Security | GRC | Product Ownership

---

## EXECUTIVE SUMMARY TEMPLATE
Generate a comprehensive analysis including:
- **Project Health Status**: Current state assessment (Excellent/Good/Fair/Critical)
- **Risk Profile**: High-level risk exposure
- **Opportunity Score**: Refactoring & optimization potential (1-10)
- **Strategic Alignment**: Business objective fulfillment

---

## PHASE 1: DEEP PROJECT ANALYSIS

### 1.1 Project Overview Assessment
**Analyze and document:**
- Project name, scope, and business objectives
- Current technology stack and dependencies
- Team composition and skill levels
- Project timeline and delivery status
- Key stakeholders and business drivers
- Current system state (monolith/microservices/hybrid)

### 1.2 Business Context Analysis (Product Owner Perspective)
**Evaluate:**
- Market positioning and competitive advantage
- Revenue impact and business metrics (KPIs/OKRs)
- Customer satisfaction and NPS scores
- Feature adoption rates and usage analytics
- Product roadmap alignment
- Technical debt vs. business capability trade-offs
- User persona requirements and pain points
- Scalability demands and growth projections

---

## PHASE 2: DETAILED ARCHITECTURE ANALYSIS

### 2.1 Current Architecture Assessment
**Examine and document:**

#### A. System Design
- System architecture pattern (Layered, Microservices, Serverless, CQRS, Event-Driven, etc.)
- Component decomposition and boundaries
- Domain-driven design implementation
- Service mesh topology (if applicable)
- Message queue architecture (if applicable)
- Database architecture (polyglot persistence analysis)
- Cache layer strategy (Redis, Memcached, application-level)
- API design patterns and versioning strategy
- Frontend architecture (SPA, SSR, PWA considerations)

#### B. Data Flow & Integration
- Request-response flows (critical paths)
- Data synchronization mechanisms
- Third-party integrations and their reliability
- Event streaming architecture (Kafka, RabbitMQ, etc.)
- Data transformation pipelines
- Batch processing workflows
- Real-time vs. eventual consistency patterns

#### C. Infrastructure & Deployment
- Cloud provider analysis (AWS, Azure, GCP, on-premise)
- Containerization strategy (Docker, Kubernetes orchestration)
- Infrastructure-as-Code maturity (Terraform, CloudFormation, Ansible)
- CI/CD pipeline architecture
- Environment management (dev, staging, prod)
- Configuration management approach
- Secrets management implementation
- Disaster recovery and backup strategy

#### D. Performance & Scalability
- Current performance metrics (latency, throughput, response times)
- Load testing results and capacity planning
- Horizontal vs. vertical scaling strategy
- Database query optimization status
- Caching effectiveness
- CDN implementation (static assets, API edge caching)
- Connection pooling and resource optimization
- Bottleneck identification
- Traffic patterns and seasonal scaling needs

#### E. Reliability & Resilience
- Uptime/SLA metrics
- Single points of failure analysis
- Circuit breaker patterns implementation
- Retry logic and backoff strategies
- Health check mechanisms
- Graceful degradation capabilities
- Bulkhead isolation patterns
- Recovery Time Objective (RTO) and Recovery Point Objective (RPO)

---

### 2.2 Architecture Quality Evaluation (Expert Architect Assessment)

#### A. SOLID Principles Compliance
- **Single Responsibility**: Module/class cohesion analysis
- **Open/Closed**: Extension points and change impact zones
- **Liskov Substitution**: Interface design and inheritance hierarchy
- **Interface Segregation**: Dependency coupling analysis
- **Dependency Inversion**: Abstraction layer adequacy

#### B. Design Patterns Analysis
- **Creational Patterns**: Object creation strategy efficiency
- **Structural Patterns**: Component composition elegance
- **Behavioral Patterns**: Communication and state management
- **Anti-patterns Detected**: God objects, tight coupling, code smell analysis

#### C. Architectural Characteristics Scoring
```
Rate each 1-10 (10 = Excellent):
- Modularity: ___
- Testability: ___
- Deployability: ___
- Scalability: ___
- Maintainability: ___
- Interoperability: ___
- Fault Tolerance: ___
- Performance: ___
- Security: ___
- Usability: ___
```

#### D. Technical Debt Assessment
- Code quality metrics (cyclomatic complexity, code coverage)
- Dependency version management (outdated packages)
- Documentation gap analysis
- Knowledge silos and bus factor
- Legacy code burden
- Refactoring priority matrix
- Estimated remediation timeline and cost

---

## PHASE 3: CYBERSECURITY EXPERT ANALYSIS

### 3.1 Threat Landscape Assessment

#### A. Threat Modeling (STRIDE/PASTA)
```
STRIDE Analysis:
- Spoofing: Authentication vulnerabilities
- Tampering: Data integrity risks
- Repudiation: Audit trail gaps
- Information Disclosure: Data exposure risks
- Denial of Service: Availability threats
- Elevation of Privilege: Authorization weaknesses
```

#### B. Vulnerability Assessment
- OWASP Top 10 alignment analysis
  - A01:2021 – Broken Access Control
  - A02:2021 – Cryptographic Failures
  - A03:2021 – Injection
  - A04:2021 – Insecure Design
  - A05:2021 – Security Misconfiguration
  - A06:2021 – Vulnerable and Outdated Components
  - A07:2021 – Identification and Authentication Failures
  - A08:2021 – Software and Data Integrity Failures
  - A09:2021 – Logging and Monitoring Failures
  - A10:2021 – Server-Side Request Forgery (SSRF)

- CWE (Common Weakness Enumeration) gap analysis
- Static Application Security Testing (SAST) results
- Dynamic Application Security Testing (DAST) results
- Dependency vulnerability scanning (SCA - Software Composition Analysis)
- Container image scanning (Trivy, Grype)
- Infrastructure security scanning (CloudSploit, Prowler)

#### C. Authentication & Authorization
- Authentication mechanisms (OAuth 2.0, SAML, OIDC, mTLS)
- MFA implementation status
- Password policy enforcement
- Session management and timeout handling
- JWT/Token security (exp, nbf, iat, signing algorithms)
- Authorization framework maturity (RBAC, ABAC, PBAC)
- Privilege escalation prevention
- API key management and rotation

#### D. Encryption & Data Protection
- Data classification (PII, PHI, Financial, Proprietary)
- Encryption at rest (AES-256, HSM usage)
- Encryption in transit (TLS 1.2+, certificate pinning)
- Key management (KMIP, HSM, cloud KMS)
- End-to-end encryption capabilities
- Database encryption (transparent encryption, column-level)
- Backup encryption and secure deletion
- GDPR/CCPA data handling compliance

#### E. Network Security
- Firewall rules and network segmentation (micro-segmentation)
- VPC/VNet architecture
- DDoS protection mechanisms
- Web Application Firewall (WAF) implementation
- API rate limiting and throttling
- DNS security (DNSSEC, DNS filtering)
- VPN/Bastion host access controls
- Network intrusion detection/prevention (IDS/IPS)

#### F. Application Security
- Input validation and sanitization
- Output encoding implementation
- SQL injection prevention (parameterized queries)
- Cross-Site Scripting (XSS) mitigation
- Cross-Site Request Forgery (CSRF) protection
- Security headers (CSP, HSTS, X-Frame-Options, etc.)
- API security (authentication, authorization, rate limiting)
- Secure error handling (information leakage prevention)
- Dependency management and vulnerability patching

#### G. Infrastructure Security
- Container security (runtime protection, pod security policies)
- Kubernetes RBAC and network policies
- Secrets management (HashiCorp Vault, AWS Secrets Manager)
- IAM policies (principle of least privilege)
- Security groups and NACLs configuration
- Bastion/Jump host implementation
- VPN and secure tunneling
- Physical security (data center access)

#### H. Monitoring & Incident Response
- Security information and event management (SIEM)
- Log aggregation and analysis (ELK stack, Splunk)
- Real-time alerting and anomaly detection
- Security incident response plan (SIRP)
- Forensics and evidence preservation
- Penetration testing frequency and scope
- Bug bounty program status
- Incident post-mortem process

#### I. Compliance & Standards
- Regulatory requirements (PCI-DSS, HIPAA, SOX, GDPR, CCPA)
- Industry standards (ISO 27001, NIST, CIS Benchmarks)
- Certification status (SOC 2 Type II, ISO 27001)
- Audit trails and evidence collection
- Compliance automation and continuous monitoring
- Data residency requirements
- Third-party vendor security assessments

### 3.2 Security Scoring & Risk Rating
```
Security Maturity Model (0-5 scale):
0 = Non-existent
1 = Ad-hoc/Informal
2 = Repeatable/Documented
3 = Defined/Standardized
4 = Managed/Measured
5 = Optimized/Continuous Improvement

Areas to Rate:
- Secure Development Lifecycle: ___
- Vulnerability Management: ___
- Incident Response: ___
- Access Control: ___
- Data Protection: ___
- Threat Detection: ___
- Compliance Management: ___
- Security Awareness: ___

Overall Security Risk: [Critical|High|Medium|Low]
```

---

## PHASE 4: GRC (GOVERNANCE, RISK & COMPLIANCE) ANALYSIS

### 4.1 Governance Assessment

#### A. Decision-Making Framework
- Architecture review board (ARB) structure
- Change management process (CAB)
- Architectural governance policies
- Standards and guidelines enforcement
- Technology selection criteria
- Approval workflows and SLAs
- Escalation procedures

#### B. Documentation & Knowledge Management
- Architecture Decision Records (ADRs)
- System documentation completeness
- API documentation (OpenAPI/Swagger)
- Runbook and playbook availability
- Knowledge base maturity
- Training program effectiveness
- Team skill inventory

#### C. Portfolio Management
- Project roadmap alignment
- Resource allocation optimization
- Budget vs. actual spending
- ROI tracking
- Technical debt tracking
- Prioritization framework
- Strategic alignment scorecard

### 4.2 Risk Management

#### A. Risk Identification & Assessment
```
Risk Register Template:
| ID | Risk | Probability | Impact | Score | Owner | Mitigation |
|----|------|-------------|--------|-------|-------|-----------|
```

Key Risk Categories:
- **Technical Risks**: Scalability, performance, availability
- **Operational Risks**: Process maturity, staffing, training
- **Business Risks**: Market changes, competitive threats
- **Compliance Risks**: Regulatory violations, audit failures
- **Security Risks**: Data breach, unauthorized access
- **Financial Risks**: Cost overruns, resource constraints

#### B. Risk Metrics & KPIs
- Mean Time to Detect (MTTD)
- Mean Time to Respond (MTTR)
- Mean Time to Resolve (MTTR)
- Incident frequency and severity
- Vulnerability age and remediation time
- Security patch compliance rate
- Audit findings and closure rate
- Change failure rate (CFR)
- Deployment frequency
- Lead time for changes

#### C. Business Continuity & Disaster Recovery
- RTO (Recovery Time Objective) targets
- RPO (Recovery Point Objective) targets
- Disaster recovery plan (DRP) documentation
- Business continuity plan (BCP)
- Backup and restore procedures
- Failover testing frequency
- Alternate site readiness
- Communication plan during incidents

### 4.3 Compliance Assessment

#### A. Regulatory Requirements Mapping
- Identify applicable regulations and standards
- Map current controls to requirements
- Identify compliance gaps
- Remediation roadmap
- Evidence collection process
- Audit readiness assessment

#### B. Audit & Assurance
- Internal audit schedule
- External audit findings
- SOC 2/SOC 3 readiness
- ISO 27001 alignment
- NIST Cybersecurity Framework alignment
- CIS Controls implementation level
- Compliance dashboard and metrics

#### C. Data Governance
- Data classification scheme
- Data ownership assignment
- Data retention policies
- Data disposal procedures
- Data lineage tracking
- Data quality metrics
- Metadata management

#### D. Third-Party Risk Management
- Vendor assessment criteria
- SLA monitoring
- Security assessment requirements
- Contract compliance tracking
- Incident reporting obligations
- Offboarding procedures

---

## PHASE 5: PRODUCT OWNERSHIP PERSPECTIVE

### 5.1 Product Strategy Alignment
- **Vision & Mission**: Long-term product direction
- **Market Positioning**: Competitive analysis
- **Customer Segmentation**: Target personas and value propositions
- **Feature Prioritization**: Business value vs. technical effort
- **Revenue Model**: Monetization strategy alignment
- **Growth Metrics**: User acquisition, retention, expansion
- **Product Roadmap**: 6-month, 12-month, 24-month plans
- **Go-to-Market Strategy**: Launch plans and distribution

### 5.2 Technical Requirements vs. Business Needs
- Core feature functionality assessment
- Performance expectations vs. reality
- Scalability roadmap alignment
- User experience and accessibility standards
- API design from customer perspective
- Integration capabilities required
- Reporting and analytics needs
- Customization and extensibility requirements

### 5.3 Quality & User Satisfaction Metrics
- Customer satisfaction scores (CSAT, NPS)
- Feature usage analytics
- Bug severity and resolution time
- Performance SLAs
- Uptime/reliability targets
- Support ticket volume and resolution time
- User feedback incorporation process
- Beta testing program maturity

### 5.4 Product-Technology Trade-offs
- Technical debt impact on feature velocity
- Refactoring vs. new feature development
- Platform stability vs. innovation
- Cost of ownership vs. feature richness
- Multi-tenancy vs. single-tenant considerations
- Localization and internationalization requirements
- Mobile vs. web platform priorities
- Migration strategy for platform transitions

---

## PHASE 6: COMPREHENSIVE REFACTORING RECOMMENDATIONS

### 6.1 Refactoring Roadmap (Prioritized)

#### Priority Tier 1: Critical (Next 0-3 Months)
```
Issue: [Description]
Impact: [High|Medium|Low]
Effort: [High|Medium|Low]
Business Value: [High|Medium|Low]
Security Risk: [Critical|High|Medium|Low]
Refactoring Strategy:
  - Current State: ...
  - Desired State: ...
  - Approach: ...
  - Dependencies: ...
  - Testing Strategy: ...
  - Rollback Plan: ...
  - Success Metrics: ...
```

#### Priority Tier 2: High (3-6 Months)
- Strategic improvements
- Tech debt with moderate impact
- Performance optimizations

#### Priority Tier 3: Medium (6-12 Months)
- Nice-to-have improvements
- Minor optimization opportunities
- Code quality enhancements

#### Priority Tier 4: Low (12+ Months)
- Future-proofing initiatives
- Technology evaluations
- Long-term architectural evolution

### 6.2 Component-Level Refactoring

#### For Each Major Component:
```
Component Name: [Name]
Current Implementation: [Architecture pattern, tech stack]
Issues Identified:
  1. [Issue with impact assessment]
  2. [Issue with impact assessment]
  
Refactoring Recommendation:
  Pattern Change: [From X to Y pattern]
  Technology Upgrade: [Version changes]
  API Redesign: [Breaking changes]
  Data Model Changes: [Schema evolution]
  Migration Path: [Step-by-step approach]
  
Testing Requirements:
  - Unit tests required: [Coverage %]
  - Integration tests: [Scope]
  - Performance tests: [Benchmarks]
  - Security tests: [Scenarios]
  
Timeline: [Estimated duration]
Resource Requirements: [Team composition]
Risk Assessment: [Potential issues]
```

### 6.3 Architectural Improvements

#### A. Modularity Enhancement
- Current coupling metrics and dependencies
- Module boundary optimization
- Interface clarity improvement
- Circular dependency elimination
- Service extraction opportunities

#### B. Scalability Improvements
- Horizontal scaling enablement
- Database sharding strategy
- Caching layer enhancement
- Load balancing optimization
- Resource pooling implementation

#### C. Performance Optimization
- Query optimization and indexing
- Algorithm efficiency review
- Memory usage optimization
- Network call reduction
- Asynchronous processing implementation

#### D. Maintainability Improvements
- Code readability enhancement
- Documentation improvement
- Test coverage increase
- Logging and observability
- Debugging capability enhancement

#### E. Security Hardening
- Authentication/authorization strengthening
- Encryption implementation
- Input validation enhancement
- API security hardening
- Infrastructure security improvements

---

## PHASE 7: COMPONENT REVIEW & ASSESSMENT

### 7.1 Code Quality Metrics
```
Metric                          Current    Target    Gap
-----------------------------------------------------------
Code Coverage                   ___%       85%       ___%
Cyclomatic Complexity (Avg)     ___        5-10      ___
Technical Debt Ratio            ___%       5%        ___%
Comment Ratio                   ___%       20%       ___%
Duplication Index               ___%       3%        ___%
MAINTAINABILITY INDEX           ___/100    75-100    ___
```

### 7.2 Peer Review Checklist
- [ ] Code follows style guidelines
- [ ] No security vulnerabilities introduced
- [ ] Performance impact assessed
- [ ] Backward compatibility maintained
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] Error handling implemented
- [ ] Logging sufficient for debugging
- [ ] No unnecessary dependencies added
- [ ] Database migrations considered
- [ ] Deployment procedure documented
- [ ] Rollback procedure defined

### 7.3 Component Evaluation Matrix
```
Component         | Security | Performance | Maintainability | Scalability | Score
------------------|----------|-------------|-----------------|-------------|-------
[Component 1]     |    7/10  |     6/10    |       8/10      |     5/10    |  6.5/10
[Component 2]     |    8/10  |     7/10    |       7/10      |     8/10    |  7.5/10
[Component 3]     |    5/10  |     4/10    |       5/10      |     3/10    |  4.25/10
```

---

## PHASE 8: DELIVERABLES & RECOMMENDATIONS

### 8.1 Executive Summary Report
- Current state assessment
- Risk profile and opportunities
- Top 5 critical recommendations
- Investment requirements (time, resources, budget)
- Expected outcomes and metrics
- Timeline for implementation

### 8.2 Detailed Architecture Blueprint
- Current architecture diagram (C4 Model recommended)
- Target architecture design
- Transition architecture (if major changes needed)
- Component interaction diagrams
- Data flow diagrams
- Deployment architecture

### 8.3 Risk & Mitigation Plan
- Identified risks with probability and impact
- Risk mitigation strategies
- Contingency plans
- Monitoring and control measures
- Regular review schedule

### 8.4 Refactoring Roadmap
- Phased implementation plan
- Resource allocation
- Dependency mapping
- Critical path analysis
- Milestone definitions
- Success criteria and KPIs

### 8.5 Implementation Guides
- Design patterns to adopt
- Coding standards and guidelines
- Testing strategies
- Deployment procedures
- Rollback procedures
- Knowledge transfer plan

### 8.6 Monitoring & Measurement Framework
- KPIs and metrics to track
- Dashboards and reporting
- Alerting thresholds
- Regular review cadence
- Feedback loops and improvement process

---

## PHASE 9: ONGOING OPTIMIZATION

### 9.1 Continuous Improvement Process
- Regular architecture reviews (quarterly)
- Technology radar/evaluation process
- Proof of concept (PoC) pipeline
- Innovation lab or skunk works initiatives
- Experimentation framework
- Learning and development program

### 9.2 Metrics & Monitoring Dashboard
```
Business Metrics:
- Revenue impact
- Customer satisfaction (NPS)
- Feature adoption rate
- Time-to-market
- Cost per transaction

Technical Metrics:
- System reliability (uptime %)
- Performance (latency, throughput)
- Security incidents (count, severity)
- Technical debt ratio
- Deployment frequency
- Lead time for changes
- Change failure rate
- Mean time to recovery

Operational Metrics:
- Team velocity
- Code review cycle time
- Time-to-production
- Incident response time
- Support ticket resolution time
```

### 9.3 Architecture Governance Review Board (ARB)
- Meeting frequency and structure
- Decision-making authority
- Escalation procedures
- Architecture principles enforcement
- Standard library management
- Technology evaluation framework
- Compliance and security gate

---

## EXPERT ASSESSMENT TEMPLATE

### Summary
**Architecture Quality**: ___/10
**Security Posture**: ___/10
**Compliance Readiness**: ___/10
**Operational Maturity**: ___/10
**Product Alignment**: ___/10
**Overall Score**: ___/10

### Top 3 Critical Issues
1. [Issue] - Impact: [High] - Timeline: [Immediate]
2. [Issue] - Impact: [High] - Timeline: [0-3 months]
3. [Issue] - Impact: [Medium] - Timeline: [3-6 months]

### Top 5 Recommendations
1. [Recommendation with business/technical justification]
2. [Recommendation with business/technical justification]
3. [Recommendation with business/technical justification]
4. [Recommendation with business/technical justification]
5. [Recommendation with business/technical justification]

### Investment Summary
- **Total Effort**: X person-months
- **Budget Estimate**: $X
- **Timeline**: X quarters
- **Expected ROI**: X%
- **Risk Level**: [Low|Medium|High]

### Sign-off
- **Technical Architect**: _________________ Date: _______
- **Security Officer**: _________________ Date: _______
- **GRC Manager**: _________________ Date: _______
- **Product Owner**: _________________ Date: _______

---

## HOW TO USE THIS PROMPT

1. **Provide Project Details**: Share codebase access, documentation, architecture diagrams, deployment info
2. **Define Scope**: Specify areas of focus (full stack, specific components, or particular concerns)
3. **Set Context**: Provide business objectives, constraints, timeline, and success criteria
4. **Request Analysis**: Ask for specific outputs from the sections that matter most
5. **Iterate**: Use findings to refine strategy and create detailed implementation plans

## ANALYSIS TEMPLATE INSTANTIATION

When applying this prompt to a specific project:

```
PROJECT ANALYSIS REQUEST

Project Name: [Name]
Organization: [Name]
Scope: [Description]
Timeline: [Duration]
Team Size: [Number]
Budget: [Amount]
Key Stakeholders: [List]

Analysis Focus Areas (Priority):
1. [Area]
2. [Area]
3. [Area]

Specific Questions to Answer:
1. [Question]
2. [Question]
3. [Question]

Current Pain Points:
1. [Pain point]
2. [Pain point]
3. [Pain point]

Desired Outcomes:
1. [Outcome]
2. [Outcome]
3. [Outcome]
```

---

**Document Version**: 2.0
**Last Updated**: 2024
**Classification**: Professional/Enterprise
**Audience**: CTOs, Architects, Security Leaders, Product Leaders, Compliance Officers
---

## SECTION 8: SYSTEM ARCHITECTURE & IMPLEMENTATION DETAILS (POST-REFACTOR)

### 8.1 Backend Architecture Components

#### 8.1.1 Core GUI Framework (`src/gui/`)
The web interface module has been refined to eliminate ghost modules and stabilize its core components:
- `dashboard.rs`: Central unified view for security metrics.
- `security_operations.rs`: SecOps interfaces with corrected metric properties (e.g., `vulnerabilities_patched_today`).
- `compliance_center.rs`: Centralized hub for tracking multi-framework compliance.
- `cognitive_collaboration.rs`: Human-AI interaction systems.
- `agent_interaction.rs`: Multi-agent management interface.
- `risk_management.rs` & `real_time_monitoring.rs`: Critical tracking and alerting systems.

#### 8.1.2 Unified Compliance Engine (`src/compliance/`)
A significant portion of the Phoenix platform's value proposition is its automated compliance monitoring across major frameworks:
- **NIST CSF (`nist_csf.rs`)**: Fully implemented across the 5 core functions (Identify, Protect, Detect, Respond, Recover) with proper string identifiers and assessment logic.
- **GDPR (`gdpr.rs`)**: Features a robust engine for tracking data subject requests (`DataSubjectRequestStats`), managing consent, and handling breach notifications with strict structural typing.
- **SOC 2 Type II (`soc2.rs`)**: Comprehensive implementation of Trust Services Criteria with detailed operational and technical subcategories spanning logical access to network security monitoring.
- **ISO 27001 & PCI DSS**: Modular implementations supporting automated evidence collection and policy tracking.

#### 8.1.3 Rust Toolchain & Deployment Strategy
The project relies exclusively on the standard Rust toolchain:
- **Minimum Rust Version**: 1.85.0+ (Currently deployed on 1.94.0).
- **Dependency Management**: Cargo configuration (`Cargo.toml`) is streamlined to use `rustls` (pure Rust TLS) across components like `sqlx` and `reqwest`, ensuring high portability and bypassing the need for native OpenSSL bindings.
- **Deployment Mode**: Capable of full local deployments via `install-local.sh` for maximum bare-metal performance, alongside traditional Docker-based orchestration.
