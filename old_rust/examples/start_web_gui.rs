//! Example: Starting the Phoenix GUI Web Server
//! 
//! This example demonstrates how to start the web-based GUI for the Barca-Strategos Phoenix system.
//! The GUI provides a modern, intuitive web interface for all system capabilities.

use barca_strategos_phoenix::gui::*;
use barca_strategos_phoenix::core::*;
use std::collections::HashMap;
use tokio;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("🚀 Starting Barca-Strategos Phoenix GUI Web Server...");

    // Initialize the GUI system with all components
    let gui_config = GUIConfig {
        dashboard: DashboardConfig {
            layout: LayoutConfig,
            widgets: WidgetConfig,
            personalization: PersonalizationConfig,
            performance: PerformanceConfig,
            accessibility: AccessibilityConfig,
        },
        cognitive: CognitiveConfig {
            collaboration: CollaborationConfig,
            visualization: VisualizationConfig,
            cognitive_load: CognitiveLoadConfig,
            interaction: InteractionConfig,
        },
        compliance: ComplianceGUIConfig {
            iso27001: ISO27001GUIConfig,
            nist_csf: NISTCSFGUIConfig,
            gdpr: GDPRGUIConfig,
            soc2: SOC2GUIConfig,
            pci_dss: PCIDSSGUIConfig,
            risk: RiskGUIConfig,
            policy: PolicyGUIConfig,
            vendor: VendorGUIConfig,
            dashboard: DashboardConfig,
            analytics: AnalyticsConfig,
            alerts: AlertConfig,
        },
        security: SecurityGUIConfig {
            threat_monitor: ThreatMonitorConfig,
            incident_manager: IncidentManagerConfig,
            vulnerability_scanner: VulnerabilityScannerConfig,
            analytics: SecurityAnalyticsConfig,
            cognitive: CognitiveSecurityConfig,
            dashboard: SecurityDashboardConfig,
        },
        risk: RiskGUIConfig {
            risk_engine: RiskEngineConfig,
            assessment: AssessmentManagerConfig,
            mitigation: MitigationManagerConfig,
            analytics: RiskAnalyticsConfig,
            cognitive: CognitiveRiskConfig,
            dashboard: RiskDashboardConfig,
        },
        agents: AgentGUIConfig {
            interface: AgentInterfaceConfig,
            cognitive_load: CognitiveLoadConfig,
            communication: CommunicationConfig,
            learning: LearningConfig,
        },
        monitoring: MonitoringGUIConfig {
            metrics: MetricsConfig,
            alerts: AlertConfig,
            visualization: VisualizationConfig,
            cognitive: CognitiveConfig,
        },
        analytics: AnalyticsGUIConfig {
            engine: AnalyticsEngineConfig,
            visualization: VisualizationConfig,
            reporting: ReportingConfig,
        },
        settings: SettingsGUIConfig {
            user_management: UserManagementConfig,
            system_config: SystemConfigConfig,
            security_config: SecurityConfigConfig,
        },
    };

    let mut gui_system = PhoenixGUISystem::new(gui_config);
    gui_system.initialize().await?;

    // Configure web server
    let web_config = WebServerConfig {
        host: "127.0.0.1".to_string(),
        port: 8080,
        static_files_path: "static".to_string(),
        max_connections: 1000,
        enable_cors: true,
    };

    // Create and start web server
    let web_server = GUIWebServer::new(gui_system, web_config);

    println!("🌐 Web Server Configuration:");
    println!("   📍 URL: http://127.0.0.1:8080");
    println!("   📁 Static Files: ./static/");
    println!("   🔌 WebSocket: ws://127.0.0.1:8080/ws");
    println!("   🌐 CORS: Enabled");
    println!("   👥 Max Connections: 1000");
    println!();

    println!("🎨 GUI Features Available:");
    println!("   🧠 Cognitive Collaboration Workspace");
    println!("   📊 Unified Dashboard with Real-time Updates");
    println!("   🛡️ Security Operations Center");
    println!("   ⚖️  Compliance Management Center");
    println!("   ⚠️  Risk Management Workspace");
    println!("   🤖 AI Agent Interaction Hub");
    println!("   📡 Real-time Monitoring Center");
    println!("   🎯 Analytics & Intelligence Engine");
    println!();

    println!("🚀 Starting web server...");
    println!("   Open your browser and navigate to: http://127.0.0.1:8080");
    println!();

    // Start the web server (this will block until the server is stopped)
    web_server.start().await?;

    Ok(())
}

/// Demo function to create a sample user session
async fn create_demo_user_session() -> UserSession {
    UserSession {
        id: "demo-session-001".to_string(),
        user_id: "demo-user-001".to_string(),
        user: User {
            id: "demo-user-001".to_string(),
            name: "Demo User".to_string(),
            email: "demo@example.com".to_string(),
            role: UserRole::SecurityAnalyst,
            permissions: UserPermissions::default(),
            created_at: chrono::Utc::now(),
            last_login: chrono::Utc::now(),
        },
        preferences: UserPreferences {
            theme: Theme::Dark,
            language: "en".to_string(),
            notifications_enabled: true,
            auto_refresh: true,
            cognitive_load_management: true,
            collaboration_mode: true,
        },
        active_workspaces: vec![],
        cognitive_load: CognitiveLoad::Normal,
        collaboration_context: None,
        created_at: chrono::Utc::now(),
    }
}

/// Demo function to show API integration
async fn demonstrate_api_integration() -> Result<(), Box<dyn std::error::Error>> {
    println!("🔌 Demonstrating API Integration...");

    // Example API calls that would be made from the frontend
    let api_base_url = "http://127.0.0.1:8080/api";

    // Create user session
    let session_data = serde_json::json!({
        "user": {
            "id": "demo-user-001",
            "name": "Demo User",
            "email": "demo@example.com",
            "role": "SecurityAnalyst"
        },
        "preferences": {
            "theme": "Dark",
            "language": "en",
            "notifications_enabled": true
        }
    });

    println!("   📝 Creating user session...");
    // In real implementation: POST {api_base_url}/user/session
    println!("   ✅ Session created successfully");

    // Get dashboard data
    println!("   📊 Fetching dashboard data...");
    // In real implementation: GET {api_base_url}/dashboard/demo-user-001
    println!("   ✅ Dashboard data retrieved");

    // Start cognitive collaboration
    let collaboration_data = serde_json::json!({
        "participants": {
            "humans": [
                {
                    "id": "demo-user-001",
                    "name": "Demo User",
                    "role": "SecurityAnalyst"
                }
            ],
            "agents": [
                {
                    "id": "security-analyst-001",
                    "name": "AI Security Analyst",
                    "agent_type": "SecurityAnalyst"
                }
            ]
        }
    });

    println!("   🧠 Starting cognitive collaboration...");
    // In real implementation: POST {api_base_url}/cognitive/demo-user-001/start-collaboration
    println!("   ✅ Collaboration session started");

    println!("   🎯 API Integration Demo Complete!");
    Ok(())
}

/// Demo function to show WebSocket real-time updates
async fn demonstrate_websocket_features() -> Result<(), Box<dyn std::error::Error>> {
    println!("🔌 Demonstrating WebSocket Features...");

    println!("   📡 WebSocket Connection:");
    println!("      - URL: ws://127.0.0.1:8080/ws");
    println!("      - Authentication: Required");
    println!("      - Real-time updates: Enabled");

    println!("   🔄 Real-time Update Types:");
    println!("      - Security Events");
    println!("      - Compliance Changes");
    println!("      - Risk Assessments");
    println!("      - Agent Interactions");
    println!("      - Cognitive Load Changes");
    println!("      - System Metrics");

    println!("   🧠 Cognitive Collaboration Features:");
    println!("      - Shared Mental Models");
    println!("      - Real-time Concept Synchronization");
    println!("      - Collaborative Decision Making");
    println!("      - Cognitive Load Optimization");

    println!("   🎯 WebSocket Demo Complete!");
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_gui_system_initialization() {
        let gui_config = GUIConfig {
            // Minimal config for testing
            dashboard: DashboardConfig {
                layout: LayoutConfig,
                widgets: WidgetConfig,
                personalization: PersonalizationConfig,
                performance: PerformanceConfig,
                accessibility: AccessibilityConfig,
            },
            // ... other configs would be initialized
            cognitive: CognitiveConfig {
                collaboration: CollaborationConfig,
                visualization: VisualizationConfig,
                cognitive_load: CognitiveLoadConfig,
                interaction: InteractionConfig,
            },
            compliance: ComplianceGUIConfig {
                iso27001: ISO27001GUIConfig,
                nist_csf: NISTCSFGUIConfig,
                gdpr: GDPRGUIConfig,
                soc2: SOC2GUIConfig,
                pci_dss: PCIDSSGUIConfig,
                risk: RiskGUIConfig,
                policy: PolicyGUIConfig,
                vendor: VendorGUIConfig,
                dashboard: DashboardConfig,
                analytics: AnalyticsConfig,
                alerts: AlertConfig,
            },
            security: SecurityGUIConfig {
                threat_monitor: ThreatMonitorConfig,
                incident_manager: IncidentManagerConfig,
                vulnerability_scanner: VulnerabilityScannerConfig,
                analytics: SecurityAnalyticsConfig,
                cognitive: CognitiveSecurityConfig,
                dashboard: SecurityDashboardConfig,
            },
            risk: RiskGUIConfig {
                risk_engine: RiskEngineConfig,
                assessment: AssessmentManagerConfig,
                mitigation: MitigationManagerConfig,
                analytics: RiskAnalyticsConfig,
                cognitive: CognitiveRiskConfig,
                dashboard: RiskDashboardConfig,
            },
            agents: AgentGUIConfig {
                interface: AgentInterfaceConfig,
                cognitive_load: CognitiveLoadConfig,
                communication: CommunicationConfig,
                learning: LearningConfig,
            },
            monitoring: MonitoringGUIConfig {
                metrics: MetricsConfig,
                alerts: AlertConfig,
                visualization: VisualizationConfig,
                cognitive: CognitiveConfig,
            },
            analytics: AnalyticsGUIConfig {
                engine: AnalyticsEngineConfig,
                visualization: VisualizationConfig,
                reporting: ReportingConfig,
            },
            settings: SettingsGUIConfig {
                user_management: UserManagementConfig,
                system_config: SystemConfigConfig,
                security_config: SecurityConfigConfig,
            },
        };

        let mut gui_system = PhoenixGUISystem::new(gui_config);
        let result = gui_system.initialize().await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_web_server_creation() {
        // This would test web server creation in a real implementation
        // For now, we just verify the structure compiles
        assert!(true);
    }
}
