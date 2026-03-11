-- Phoenix API Database Initialization Script
-- This script initializes the database with tables and data for compliance modules

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create compliance tables
CREATE TABLE IF NOT EXISTS compliance_assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    framework VARCHAR(50) NOT NULL,
    version VARCHAR(20) NOT NULL,
    scope JSONB NOT NULL,
    overall_score DECIMAL(5,4) NOT NULL,
    findings JSONB NOT NULL,
    recommendations JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_compliance_assessments_framework ON compliance_assessments(framework);
CREATE INDEX IF NOT EXISTS idx_compliance_assessments_created_at ON compliance_assessments(created_at);

-- Create SOC2 specific tables
CREATE TABLE IF NOT EXISTS soc2_controls (
    id VARCHAR(20) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    subcategories JSONB NOT NULL,
    objective TEXT NOT NULL,
    control_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    implementation_date TIMESTAMP WITH TIME ZONE,
    last_review_date TIMESTAMP WITH TIME ZONE,
    evidence JSONB NOT NULL,
    owner VARCHAR(255) NOT NULL,
    risk_level VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS soc2_control_updates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    control_id VARCHAR(20) NOT NULL REFERENCES soc2_controls(id),
    update_type VARCHAR(50) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    evidence JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create PCI DSS specific tables
CREATE TABLE IF NOT EXISTS pci_dss_requirements (
    id VARCHAR(20) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    subcategories JSONB NOT NULL,
    objective TEXT NOT NULL,
    control_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    implementation_date TIMESTAMP WITH TIME ZONE,
    last_review_date TIMESTAMP WITH TIME ZONE,
    evidence JSONB NOT NULL,
    owner VARCHAR(255) NOT NULL,
    risk_level VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pci_dss_requirement_updates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requirement_id VARCHAR(20) NOT NULL REFERENCES pci_dss_requirements(id),
    update_type VARCHAR(50) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    evidence JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create reports tables
CREATE TABLE IF NOT EXISTS compliance_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id VARCHAR(255) NOT NULL UNIQUE,
    framework VARCHAR(50) NOT NULL,
    version VARCHAR(20) NOT NULL,
    assessment_id UUID NOT NULL REFERENCES compliance_assessments(id),
    report_content TEXT NOT NULL,
    format VARCHAR(20) NOT NULL DEFAULT 'json',
    file_path VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create audit logs table
CREATE TABLE IF NOT EXISTS compliance_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for audit logs
CREATE INDEX IF NOT EXISTS idx_compliance_audit_logs_entity ON compliance_audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_compliance_audit_logs_created_at ON compliance_audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_compliance_audit_logs_user_id ON compliance_audit_logs(user_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_compliance_assessments_updated_at
    BEFORE UPDATE ON compliance_assessments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_soc2_controls_updated_at
    BEFORE UPDATE ON soc2_controls
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pci_dss_requirements_updated_at
    BEFORE UPDATE ON pci_dss_requirements
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create trigger for audit logging
CREATE OR REPLACE FUNCTION log_compliance_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO compliance_audit_logs (action, entity_type, entity_id, old_values, new_values)
        VALUES ('DELETE', TG_TABLE_NAME, OLD.id::text, row_to_json(OLD), NULL);
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO compliance_audit_logs (action, entity_type, entity_id, old_values, new_values)
        VALUES ('UPDATE', TG_TABLE_NAME, NEW.id::text, row_to_json(OLD), row_to_json(NEW));
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO compliance_audit_logs (action, entity_type, entity_id, old_values, new_values)
        VALUES ('INSERT', TG_TABLE_NAME, NEW.id::text, NULL, row_to_json(NEW));
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create audit triggers
CREATE TRIGGER soc2_controls_audit
    AFTER INSERT OR UPDATE OR DELETE ON soc2_controls
    FOR EACH ROW
    EXECUTE FUNCTION log_compliance_changes();

CREATE TRIGGER pci_dss_requirements_audit
    AFTER INSERT OR UPDATE OR DELETE ON pci_dss_requirements
    FOR EACH ROW
    EXECUTE FUNCTION log_compliance_changes();

CREATE TRIGGER compliance_assessments_audit
    AFTER INSERT OR UPDATE OR DELETE ON compliance_assessments
    FOR EACH ROW
    EXECUTE FUNCTION log_compliance_changes();

CREATE TRIGGER compliance_reports_audit
    AFTER INSERT OR UPDATE OR DELETE ON compliance_reports
    FOR EACH ROW
    EXECUTE FUNCTION log_compliance_changes();

-- Insert default SOC2 controls
INSERT INTO soc2_controls (id, title, description, category, subcategories, objective, control_type, status, implementation_date, last_review_date, evidence, owner, risk_level) VALUES
('CC1.1', 'Governance', 'Establish and communicate governance framework', 'governance', '["Governance framework", "Board oversight", "Management direction", "Legal and compliance", "Risk management", "Ethics and compliance"]', 'Establish and communicate governance framework', 'organizational', 'implemented', NOW() - INTERVAL '30 days', NOW() - INTERVAL '7 days', '["Governance policy document", "Board meeting minutes"]', 'Board of Directors', 'low'),
('CC2.1', 'Asset Inventory', 'Maintain complete and accurate inventory of all hardware, software, and data assets', 'asset_management', '["Hardware inventory", "Software inventory", "Data inventory", "Cloud assets", "Mobile devices"]', 'Maintain complete and accurate inventory of all assets', 'organizational', 'implemented', NOW() - INTERVAL '20 days', NOW() - INTERVAL '5 days', '["Asset registry", "CMDB", "Asset management system"]', 'IT Asset Manager', 'medium'),
('CC3.2', 'Identity Management and Access Control', 'Identify, authenticate, and authorize access to systems', 'access_control', '["User access management", "Remote access", "Multi-factor authentication", "Privileged access management", "Account lifecycle management", "Access certification", "Identity proofing"]', 'Identify, authenticate, and authorize access to systems', 'technical', 'implemented', NOW() - INTERVAL '15 days', NOW() - INTERVAL '3 days', '["User access policies", "Authentication system", "MFA system"]', 'Identity Management Team', 'high'),
('CC4.1', 'Security Awareness and Training', 'Provide security awareness training to all personnel', 'operational', '["Security training program", "Phishing awareness", "Social engineering awareness", "Security culture", "Threat intelligence sharing"]', 'Ensure all personnel understand their security responsibilities', 'operational', 'partially_implemented', NOW() - INTERVAL '10 days', NOW() - INTERVAL '2 days', '["Security training materials", "Phishing simulations"]', 'Security Team', 'medium'),
('CC5.1', 'Vulnerability Management', 'Identify, assess, and remediate vulnerabilities', 'operational', '["Vulnerability scanning", "Penetration testing", "Vulnerability assessment", "Patch management", "CVE monitoring"]', 'Continuously identify and remediate vulnerabilities', 'technical', 'not_implemented', NULL, NULL, '[]', 'Security Team', 'high'),
('CC6.1', 'Incident Response', 'Establish and implement incident response capabilities', 'operational', '["Incident response planning", "Incident response playbooks", "Incident notification procedures", "Forensic capabilities", "Tabletop exercises", "Threat hunting"]', 'Ensure timely and effective incident response', 'operational', 'implemented', NOW() - INTERVAL '30 days', NOW() - INTERVAL '7 days', '["Incident response plan", "Playbooks"]', 'SOC Team', 'medium'),
('CC7.1', 'Disaster Recovery Planning', 'Establish and test disaster recovery plans', 'operational', '["Business continuity planning", "Disaster recovery testing", "Backup and recovery procedures", "Alternative processing sites", "Crisis communication", "Tabletop exercises"]', 'Ensure business continuity during disruptions', 'operational', 'not_implemented', NULL, NULL, '[]', 'Business Continuity Team', 'high'),
('CC8.1', 'Penetration Testing', 'Conduct regular penetration testing', 'test_evaluation', '["External penetration testing", "Internal penetration testing", "Social engineering testing", "Application security testing"]', 'Test security controls through penetration testing', 'technical', 'implemented', NOW() - INTERVAL '25 days', NOW() - INTERVAL '3 days', '["Penetration test reports", "Security assessment reports"]', 'Security Team', 'high'),
('CC9.1', 'Network Security Monitoring', 'Monitor network traffic for security events', 'communications_security', '["Network intrusion detection", "Malware analysis", "Log analysis", "Network traffic analysis", "IDS integration", "Threat hunting", "Network device monitoring"]', 'Detect and respond to network security incidents', 'technical', 'implemented', NOW() - INTERVAL '25 days', NOW() - INTERVAL '3 days', '["Network logs", "IDS alerts", "Firewall logs"]', 'Network Security Team', 'medium');

-- Insert default PCI DSS requirements
INSERT INTO pci_dss_requirements (id, title, description, category, subcategories, objective, control_type, status, implementation_date, last_review_date, evidence, owner, risk_level) VALUES
('1.1', 'Network Security Controls', 'Install and maintain network security controls', 'network_security', '["Firewall configuration", "Network segmentation", "Secure network architecture", "Restrict traffic", "Document network topology"]', 'Protect cardholder data', 'technical', 'implemented', NOW() - INTERVAL '60 days', NOW() - INTERVAL '10 days', '["Firewall rules", "Network diagrams"]', 'Network Security Team', 'critical'),
('2.1', 'Secure Configurations', 'Apply secure configurations to all system components', 'system_configuration', '["Secure configuration standards", "System hardening", "Patch management", "Configuration management", "Vulnerability management"]', 'Maintain secure systems', 'technical', 'partially_implemented', NOW() - INTERVAL '30 days', NOW() - INTERVAL '5 days', '["Configuration baselines", "Patch reports"]', 'System Administration', 'high'),
('3.1', 'Protect Stored Account Data', 'Protect stored account data', 'data_protection', '["Data encryption", "Key management", "Data masking", "Secure storage", "Data retention policies"]', 'Protect cardholder data', 'technical', 'implemented', NOW() - INTERVAL '45 days', NOW() - INTERVAL '7 days', '["Encryption certificates", "Key management logs"]', 'Security Team', 'critical'),
('4.1', 'Protect Cardholder Data in Transit', 'Protect cardholder data in transit', 'data_protection', '["Strong cryptography", "Secure protocols", "SSL/TLS configuration", "Certificate management", "Network encryption"]', 'Protect data in transit', 'technical', 'implemented', NOW() - INTERVAL '40 days', NOW() - INTERVAL '6 days', '["TLS certificates", "Encryption logs"]', 'Security Team', 'critical'),
('5.1', 'Malware Protection', 'Protect all systems against malicious software', 'malware_protection', '["Antivirus software", "Malware detection", "Regular updates", "System monitoring", "Incident response"]', 'Prevent malware infections', 'technical', 'implemented', NOW() - INTERVAL '50 days', NOW() - INTERVAL '8 days', '["Antivirus reports", "Malware scan logs"]', 'Security Operations', 'high'),
('6.1', 'Secure Development', 'Develop and maintain secure systems and software', 'secure_development', '["Secure coding practices", "Code reviews", "Security testing", "Vulnerability scanning", "Change management"]', 'Secure development lifecycle', 'operational', 'partially_implemented', NOW() - INTERVAL '25 days', NOW() - INTERVAL '3 days', '["Code review reports", "Security test results"]', 'Development Team', 'high'),
('7.1', 'Access Control', 'Restrict access to cardholder data', 'access_control', '["Least privilege principle", "User authentication", "Access reviews", "Role-based access", "Physical access controls"]', 'Restrict data access', 'technical', 'implemented', NOW() - INTERVAL '55 days', NOW() - INTERVAL '9 days', '["Access control policies", "User access logs"]', 'Identity Management', 'critical'),
('8.1', 'Authentication', 'Identify and authenticate access to system components', 'access_control', '["Strong authentication", "Multi-factor authentication", "Password policies", "Session management", "Account management"]', 'Authenticate users', 'technical', 'implemented', NOW() - INTERVAL '35 days', NOW() - INTERVAL '4 days', '["MFA logs", "Authentication policies"]', 'Identity Management', 'high'),
('9.1', 'Physical Access Control', 'Restrict physical access to cardholder data', 'physical_security', '["Physical security controls", "Visitor management", "Surveillance systems", "Secure facilities", "Media destruction"]', 'Physical security', 'physical', 'implemented', NOW() - INTERVAL '70 days', NOW() - INTERVAL '12 days', '["Access logs", "Security camera footage"]', 'Physical Security', 'medium'),
('10.1', 'Logging and Monitoring', 'Track and monitor all access to network resources and cardholder data', 'monitoring', '["Audit logging", "Security monitoring", "Log analysis", "Incident detection", "Log retention"]', 'Monitor and track access', 'technical', 'implemented', NOW() - INTERVAL '48 days', NOW() - INTERVAL '8 days', '["System logs", "Monitoring dashboards"]', 'Security Operations', 'high'),
('11.1', 'Security Testing', 'Regularly test security systems and processes', 'testing', '["Penetration testing", "Vulnerability scanning", "Security assessments", "Incident response testing", "Wireless testing"]', 'Test security controls', 'technical', 'not_implemented', NULL, NULL, '[]', 'Security Team', 'critical'),
('12.1', 'Security Policies', 'Support information security with organizational policies and programs', 'policy_management', '["Information security policy", "Risk assessment", "Security awareness training", "Incident response plan", "Vendor management"]', 'Security governance', 'organizational', 'partially_implemented', NOW() - INTERVAL '20 days', NOW() - INTERVAL '2 days', '["Security policies", "Training records"]', 'Security Management', 'medium');

-- Create default admin user for compliance management
INSERT INTO compliance_audit_logs (action, entity_type, entity_id, user_id, old_values, new_values, ip_address, user_agent) VALUES
('SYSTEM_INIT', 'database', 'database', 'system', NULL, '{"status": "initialized"}', '127.0.0.1', 'Phoenix API System');

-- Create view for compliance statistics
CREATE OR REPLACE VIEW compliance_statistics AS
SELECT 
    'SOC2' as framework,
    COUNT(*) as total_controls,
    COUNT(CASE WHEN status = 'compliant' THEN 1 END) as compliant_controls,
    COUNT(CASE WHEN status = 'implemented' THEN 1 END) as implemented_controls,
    COUNT(CASE WHEN status = 'partially_implemented' THEN 1 END) as partially_implemented_controls,
    COUNT(CASE WHEN status = 'not_implemented' THEN 1 END) as not_implemented_controls,
    AVG(CASE 
        WHEN status = 'compliant' THEN 1.0
        WHEN status = 'implemented' THEN 0.8
        WHEN status = 'partially_implemented' THEN 0.6
        WHEN status = 'not_implemented' THEN 0.0
        ELSE 0.5
    END) as average_compliance_score,
    MAX(updated_at) as last_updated
FROM soc2_controls
GROUP BY 'SOC2'

UNION ALL

SELECT 
    'PCI DSS' as framework,
    COUNT(*) as total_requirements,
    COUNT(CASE WHEN status = 'compliant' THEN 1 END) as compliant_requirements,
    COUNT(CASE WHEN status = 'implemented' THEN 1 END) as implemented_requirements,
    COUNT(CASE WHEN status = 'partially_implemented' THEN 1 END) as partially_implemented_requirements,
    COUNT(CASE WHEN status = 'not_implemented' THEN 1 END) as not_implemented_requirements,
    AVG(CASE 
        WHEN status = 'compliant' THEN 1.0
        WHEN status = 'implemented' THEN 0.8
        WHEN status = 'partially_implemented' THEN 0.6
        WHEN status = 'not_implemented' THEN 0.0
        ELSE 0.5
    END) as average_compliance_score,
    MAX(updated_at) as last_updated
FROM pci_dss_requirements
GROUP BY 'PCI DSS';

-- Grant necessary permissions (adjust as needed for your setup)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO phoenix;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO phoenix;

-- Create indexes for better performance on large datasets
CREATE INDEX IF NOT EXISTS idx_soc2_controls_status ON soc2_controls(status);
CREATE INDEX IF NOT EXISTS idx_soc2_controls_category ON soc2_controls(category);
CREATE INDEX IF NOT EXISTS idx_soc2_controls_risk_level ON soc2_controls(risk_level);
CREATE INDEX IF NOT EXISTS idx_soc2_controls_owner ON soc2_controls(owner);

CREATE INDEX IF NOT EXISTS idx_pci_dss_requirements_status ON pci_dss_requirements(status);
CREATE INDEX IF NOT EXISTS idx_pci_dss_requirements_category ON pci_dss_requirements(category);
CREATE INDEX IF NOT EXISTS idx_pci_dss_requirements_risk_level ON pci_dss_requirements(risk_level);
CREATE INDEX IF NOT EXISTS idx_pci_dss_requirements_owner ON pci_dss_requirements(owner);

-- Create full-text search indexes for compliance content (requires pg_trgm extension)
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX IF NOT EXISTS idx_soc2_controls_search ON soc2_controls USING gin(to_tsvector('english', title || ' ' || description));
-- CREATE INDEX IF NOT EXISTS idx_pci_dss_requirements_search ON pci_dss_requirements USING gin(to_tsvector('english', title || ' ' || description));

COMMIT;
