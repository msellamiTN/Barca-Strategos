pub mod dashboard;
pub mod cognitive_collaboration;
pub mod compliance_center;
pub mod security_operations;
pub mod risk_management;
pub mod agent_interaction;
pub mod real_time_monitoring;

pub use dashboard::*;
pub use cognitive_collaboration::*;
pub use compliance_center::*;
pub use security_operations::*;
pub use risk_management::*;
pub use agent_interaction::*;
pub use real_time_monitoring::*;

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::sync::RwLock;
use uuid::Uuid;
use crate::common::{UserId, User, UserPreferences, Workspace, CollaborationContext, AnalyticsIntelligenceEngine, SettingsConfigurationManager, AnalyticsInsights, InterfaceAdaptation};

/// Main GUI System for Barca-Strategos Phoenix
/// Integrates all system features with optimized intuitive interfaces

pub struct PhoenixGUISystem {
    dashboard: MainDashboard,
    cognitive_workspace: CognitiveCollaborationWorkspace,
    compliance_center: ComplianceManagementCenter,
    security_ops: SecurityOperationsCenter,
    risk_workspace: RiskManagementWorkspace,
    agent_hub: AgentInteractionHub,
    monitoring_center: RealTimeMonitoringCenter,
    analytics_engine: AnalyticsIntelligenceEngine,
    config_manager: SettingsConfigurationManager,
    user_sessions: RwLock<HashMap<UserId, UserSession>>,
    system_state: RwLock<GUISystemState>,
}

impl PhoenixGUISystem {
    pub fn new(config: GUIConfig) -> Self {
        Self {
            dashboard: MainDashboard::new(&config.dashboard),
            cognitive_workspace: CognitiveCollaborationWorkspace::new(&config.cognitive),
            compliance_center: ComplianceManagementCenter::new(&config.compliance),
            security_ops: SecurityOperationsCenter::new(&config.security),
            risk_workspace: RiskManagementWorkspace::new(&config.risk),
            agent_hub: AgentInteractionHub::new(&config.agents),
            monitoring_center: RealTimeMonitoringCenter::new(&config.monitoring),
            analytics_engine: AnalyticsIntelligenceEngine::new(&config.analytics),
            config_manager: SettingsConfigurationManager::new(&config.settings),
            user_sessions: RwLock::new(HashMap::new()),
            system_state: RwLock::new(GUISystemState::new()),
        }
    }

    pub async fn initialize(&mut self) -> Result<(), GUIError> {
        // Initialize all GUI components
        self.dashboard.initialize().await?;
        self.cognitive_workspace.initialize().await?;
        self.compliance_center.initialize().await?;
        self.security_ops.initialize().await?;
        self.risk_workspace.initialize().await?;
        self.agent_hub.initialize().await?;
        self.monitoring_center.initialize().await?;
        self.analytics_engine.initialize().await?;
        self.config_manager.initialize().await?;

        // Start background services
        self.start_background_services().await?;

        println!("🎨 Phoenix GUI System initialized with all integrated features");
        Ok(())
    }

    /// Create user session with personalized interface
    pub async fn create_user_session(&self, user: User, preferences: UserPreferences) -> Result<UserSession, GUIError> {
        let session = UserSession {
            id: Uuid::new_v4().to_string(),
            user_id: user.id.clone(),
            user,
            preferences,
            active_workspaces: Vec::new(),
            cognitive_load: CognitiveLoad::Normal,
            collaboration_context: None,
            created_at: chrono::Utc::now(),
        };

        let mut sessions = self.user_sessions.write().await;
        sessions.insert(session.user_id.clone(), session.clone());

        Ok(session)
    }

    /// Get unified dashboard integrating all system features
    pub async fn get_unified_dashboard(&self, user_id: &str) -> Result<UnifiedDashboard, GUIError> {
        let sessions = self.user_sessions.read().await;
        let session = sessions.get(user_id).ok_or(GUIError::SessionNotFound(user_id.to_string()))?;

        let dashboard = UnifiedDashboard {
            user_session: session.clone(),
            compliance_overview: self.compliance_center.get_overview(user_id).await?,
            security_status: self.security_ops.get_status_summary(user_id).await?,
            risk_dashboard: self.risk_workspace.get_dashboard(user_id).await?,
            agent_status: self.agent_hub.get_agent_status(user_id).await?,
            cognitive_workspace: self.cognitive_workspace.get_workspace_state(user_id).await?,
            monitoring_metrics: self.monitoring_center.get_metrics(user_id).await?,
            analytics_insights: self.analytics_engine.get_insights(user_id).await?,
            system_health: self.get_system_health().await?,
            vulnerabilities_patched_today: self.get_vulnerabilities_patched_today().await?,
        };

        Ok(dashboard)
    }

    /// Handle real-time GUI updates
    pub async fn handle_gui_update(&self, update: GUIUpdate) -> Result<(), GUIError> {
        match update.update_type {
            GUIUpdateType::ComplianceChange => {
                self.compliance_center.handle_update(update.data).await?;
            },
            GUIUpdateType::SecurityEvent => {
                self.security_ops.handle_security_update(update.data).await?;
            },
            GUIUpdateType::RiskAssessment => {
                self.risk_workspace.handle_risk_update(update.data).await?;
            },
            GUIUpdateType::AgentInteraction => {
                self.agent_hub.handle_agent_update(update.data).await?;
            },
            GUIUpdateType::CognitiveCollaboration => {
                self.cognitive_workspace.handle_collaboration_update(update.data).await?;
            },
            GUIUpdateType::MonitoringAlert => {
                self.monitoring_center.handle_monitoring_update(update.data).await?;
            },
            GUIUpdateType::AnalyticsUpdate => {
                self.analytics_engine.handle_analytics_update(update.data).await?;
            },
        }

        // Broadcast update to relevant user sessions
        self.broadcast_update_to_users(update).await?;

        Ok(())
    }

    /// Get optimized interface based on user role and context
    pub async fn get_optimized_interface(&self, user_id: &str, context: InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        let sessions = self.user_sessions.read().await;
        let session = sessions.get(user_id).ok_or(GUIError::SessionNotFound(user_id.to_string()))?;

        let interface = match context.interface_type {
            InterfaceType::SecurityAnalyst => {
                self.create_security_analyst_interface(&session, &context).await?
            },
            InterfaceType::ComplianceManager => {
                self.create_compliance_manager_interface(&session, &context).await?
            },
            InterfaceType::RiskManager => {
                self.create_risk_manager_interface(&session, &context).await?
            },
            InterfaceType::IncidentResponder => {
                self.create_incident_responder_interface(&session, &context).await?
            },
            InterfaceType::Executive => {
                self.create_executive_interface(&session, &context).await?
            },
            InterfaceType::Collaborative => {
                self.create_collaborative_interface(&session, &context).await?
            },
        };

        Ok(interface)
    }

    // Private helper methods

    async fn start_background_services(&self) -> Result<(), GUIError> {
        // Start real-time monitoring
        tokio::spawn(self.background_monitoring_loop());
        
        // Start cognitive load optimization
        tokio::spawn(self.cognitive_optimization_loop());
        
        // Start analytics processing
        tokio::spawn(self.analytics_processing_loop());
        
        // Start interface adaptation
        tokio::spawn(self.interface_adaptation_loop());

        Ok(())
    }

    async fn background_monitoring_loop(&self) {
        let mut interval = tokio::time::interval(std::time::Duration::from_secs(5));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.update_monitoring_data().await {
                eprintln!("GUI: Error updating monitoring data: {}", e);
            }
        }
    }

    async fn cognitive_optimization_loop(&self) {
        let mut interval = tokio::time::interval(std::time::Duration::from_secs(30));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.optimize_cognitive_load().await {
                eprintln!("GUI: Error optimizing cognitive load: {}", e);
            }
        }
    }

    async fn analytics_processing_loop(&self) {
        let mut interval = tokio::time::interval(std::time::Duration::from_secs(60));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.process_analytics().await {
                eprintln!("GUI: Error processing analytics: {}", e);
            }
        }
    }

    async fn interface_adaptation_loop(&self) {
        let mut interval = tokio::time::interval(std::time::Duration::from_secs(120));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.adapt_interfaces().await {
                eprintln!("GUI: Error adapting interfaces: {}", e);
            }
        }
    }

    async fn update_monitoring_data(&self) -> Result<(), GUIError> {
        let monitoring_data = self.monitoring_center.collect_all_metrics().await?;
        
        let mut state = self.system_state.write().await;
        state.update_monitoring_data(monitoring_data);

        Ok(())
    }

    async fn optimize_cognitive_load(&self) -> Result<(), GUIError> {
        let mut sessions = self.user_sessions.write().await;
        
        for session in sessions.values_mut() {
            if session.cognitive_load == CognitiveLoad::High {
                // Simplify interface for high cognitive load
                session.preferences.interface_complexity = InterfaceComplexity::Simplified;
            }
        }

        Ok(())
    }

    async fn process_analytics(&self) -> Result<(), GUIError> {
        let analytics_data = self.analytics_engine.process_all_data().await?;
        
        // Update relevant GUI components with analytics insights
        self.dashboard.update_with_analytics(analytics_data.clone()).await?;
        self.compliance_center.update_with_analytics(analytics_data.clone()).await?;
        self.security_ops.update_with_analytics(analytics_data.clone()).await?;

        Ok(())
    }

    async fn adapt_interfaces(&self) -> Result<(), GUIError> {
        let sessions = self.user_sessions.read().await;
        
        for session in sessions.values() {
            let adaptation = self.analytics_engine.get_interface_adaptation(&session.user_id).await?;
            
            // Apply adaptation to user's active interfaces
            self.apply_interface_adaptation(&session.user_id, adaptation).await?;
        }

        Ok(())
    }

    async fn apply_interface_adaptation(&self, user_id: &str, adaptation: InterfaceAdaptation) -> Result<(), GUIError> {
        // Apply adaptation to all relevant GUI components
        self.dashboard.apply_adaptation(user_id, &adaptation).await?;
        self.cognitive_workspace.apply_adaptation(user_id, &adaptation).await?;
        self.compliance_center.apply_adaptation(user_id, &adaptation).await?;

        Ok(())
    }

    async fn broadcast_update_to_users(&self, update: GUIUpdate) -> Result<(), GUIError> {
        let sessions = self.user_sessions.read().await;
        
        for session in sessions.values() {
            if self.should_receive_update(session, &update) {
                // Send update to user's WebSocket connection
                self.send_update_to_session(session, &update).await?;
            }
        }

        Ok(())
    }

    fn should_receive_update(&self, session: &UserSession, update: &GUIUpdate) -> bool {
        // Check if user's role and context require this update
        match update.update_type {
            GUIUpdateType::ComplianceChange => session.user.role.has_compliance_access(),
            GUIUpdateType::SecurityEvent => session.user.role.has_security_access(),
            GUIUpdateType::RiskAssessment => session.user.role.has_risk_access(),
            GUIUpdateType::AgentInteraction => session.user.role.has_agent_access(),
            _ => true,
        }
    }

    async fn send_update_to_session(&self, session: &UserSession, update: &GUIUpdate) -> Result<(), GUIError> {
        // Implementation would send WebSocket message to user's client
        println!("📡 Sending GUI update to user {}: {:?}", session.user_id, update.update_type);
        Ok(())
    }

    async fn get_system_health(&self) -> Result<SystemHealth, GUIError> {
        let state = self.system_state.read().await;
        Ok(state.get_system_health())
    }

    async fn create_security_analyst_interface(&self, session: &UserSession, context: &InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        Ok(OptimizedInterface {
            interface_id: Uuid::new_v4().to_string(),
            user_id: session.user_id.clone(),
            interface_type: InterfaceType::SecurityAnalyst,
            layout: InterfaceLayout::SecurityFocused,
            components: vec![
                InterfaceComponent::SecurityDashboard,
                InterfaceComponent::ThreatMonitoring,
                InterfaceComponent::IncidentResponse,
                InterfaceComponent::AgentCollaboration,
            ],
            cognitive_optimizations: self.get_cognitive_optimizations_for_role("security_analyst").await?,
            adaptive_features: self.get_adaptive_features_for_context(context).await?,
        })
    }

    async fn create_compliance_manager_interface(&self, session: &UserSession, context: &InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        Ok(OptimizedInterface {
            interface_id: Uuid::new_v4().to_string(),
            user_id: session.user_id.clone(),
            interface_type: InterfaceType::ComplianceManager,
            layout: InterfaceLayout::ComplianceFocused,
            components: vec![
                InterfaceComponent::ComplianceDashboard,
                InterfaceComponent::PolicyManagement,
                InterfaceComponent::AuditTracking,
                InterfaceComponent::ReportingTools,
            ],
            cognitive_optimizations: self.get_cognitive_optimizations_for_role("compliance_manager").await?,
            adaptive_features: self.get_adaptive_features_for_context(context).await?,
        })
    }

    async fn create_risk_manager_interface(&self, session: &UserSession, context: &InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        Ok(OptimizedInterface {
            interface_id: Uuid::new_v4().to_string(),
            user_id: session.user_id.clone(),
            interface_type: InterfaceType::RiskManager,
            layout: InterfaceLayout::RiskFocused,
            components: vec![
                InterfaceComponent::RiskDashboard,
                InterfaceComponent::AssessmentTools,
                InterfaceComponent::MitigationPlanning,
                InterfaceComponent::VendorManagement,
            ],
            cognitive_optimizations: self.get_cognitive_optimizations_for_role("risk_manager").await?,
            adaptive_features: self.get_adaptive_features_for_context(context).await?,
        })
    }

    async fn create_incident_responder_interface(&self, session: &UserSession, context: &InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        Ok(OptimizedInterface {
            interface_id: Uuid::new_v4().to_string(),
            user_id: session.user_id.clone(),
            interface_type: InterfaceType::IncidentResponder,
            layout: InterfaceLayout::IncidentFocused,
            components: vec![
                InterfaceComponent::IncidentDashboard,
                InterfaceComponent::ResponsePlaybooks,
                InterfaceComponent::CommunicationTools,
                InterfaceComponent::EvidenceCollection,
            ],
            cognitive_optimizations: self.get_cognitive_optimizations_for_role("incident_responder").await?,
            adaptive_features: self.get_adaptive_features_for_context(context).await?,
        })
    }

    async fn create_executive_interface(&self, session: &UserSession, context: &InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        Ok(OptimizedInterface {
            interface_id: Uuid::new_v4().to_string(),
            user_id: session.user_id.clone(),
            interface_type: InterfaceType::Executive,
            layout: InterfaceLayout::ExecutiveFocused,
            components: vec![
                InterfaceComponent::ExecutiveDashboard,
                InterfaceComponent::StrategicReports,
                InterfaceComponent::RiskOverview,
                InterfaceComponent::ComplianceStatus,
            ],
            cognitive_optimizations: self.get_cognitive_optimizations_for_role("executive").await?,
            adaptive_features: self.get_adaptive_features_for_context(context).await?,
        })
    }

    async fn create_collaborative_interface(&self, session: &UserSession, context: &InterfaceContext) -> Result<OptimizedInterface, GUIError> {
        Ok(OptimizedInterface {
            interface_id: Uuid::new_v4().to_string(),
            user_id: session.user_id.clone(),
            interface_type: InterfaceType::Collaborative,
            layout: InterfaceLayout::CollaborationFocused,
            components: vec![
                InterfaceComponent::CognitiveWorkspace,
                InterfaceComponent::AgentInteraction,
                InterfaceComponent::SharedDocuments,
                InterfaceComponent::CommunicationTools,
            ],
            cognitive_optimizations: self.get_cognitive_optimizations_for_role("collaborative").await?,
            adaptive_features: self.get_adaptive_features_for_context(context).await?,
        })
    }

    async fn get_cognitive_optimizations_for_role(&self, role: &str) -> Result<Vec<CognitiveOptimization>, GUIError> {
        match role {
            "security_analyst" => Ok(vec![
                CognitiveOptimization::ThreatPrioritization,
                CognitiveOptimization::PatternRecognition,
                CognitiveOptimization::QuickResponseActions,
            ]),
            "compliance_manager" => Ok(vec![
                CognitiveOptimization::DocumentOrganization,
                CognitiveOptimization::RequirementTracking,
                CognitiveOptimization::AuditPreparation,
            ]),
            "risk_manager" => Ok(vec![
                CognitiveOptimization::RiskVisualization,
                CognitiveOptimization::ImpactAssessment,
                CognitiveOptimization::MitigationPrioritization,
            ]),
            "incident_responder" => Ok(vec![
                CognitiveOptimization::IncidentPrioritization,
                CognitiveOptimization::StepByStepGuidance,
                CognitiveOptimization::CommunicationOptimization,
            ]),
            "executive" => Ok(vec![
                CognitiveOptimization::ExecutiveSummaries,
                CognitiveOptimization::KPIHighlighting,
                CognitiveOptimization::TrendVisualization,
            ]),
            "collaborative" => Ok(vec![
                CognitiveOptimization::SharedUnderstanding,
                CognitiveOptimization::CommunicationClarity,
                CognitiveOptimization::DecisionSupport,
            ]),
            _ => Ok(vec![]),
        }
    }

    async fn get_adaptive_features_for_context(&self, context: &InterfaceContext) -> Result<Vec<AdaptiveFeature>, GUIError> {
        let mut features = Vec::new();

        if context.urgency == UrgencyLevel::High {
            features.push(AdaptiveFeature::SimplifiedInterface);
            features.push(AdaptiveFeature::CriticalInformationHighlighting);
        }

        if context.cognitive_load == CognitiveLoad::High {
            features.push(AdaptiveFeature::ReducedCognitiveLoad);
            features.push(AdaptiveFeature::ProgressiveDisclosure);
        }

        if context.collaboration_mode {
            features.push(AdaptiveFeature::EnhancedCommunication);
            features.push(AdaptiveFeature::SharedWorkspace);
        }

        Ok(features)
    }
}

// Supporting structures

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GUIConfig {
    pub dashboard: DashboardConfig,
    pub cognitive: CognitiveConfig,
    pub compliance: ComplianceGUIConfig,
    pub security: SecurityGUIConfig,
    pub risk: RiskGUIConfig,
    pub agents: AgentGUIConfig,
    pub monitoring: MonitoringGUIConfig,
    pub analytics: AnalyticsGUIConfig,
    pub settings: SettingsGUIConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserSession {
    pub id: String,
    pub user_id: String,
    pub user: User,
    pub preferences: UserPreferences,
    pub active_workspaces: Vec<Workspace>,
    pub cognitive_load: CognitiveLoad,
    pub collaboration_context: Option<CollaborationContext>,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UnifiedDashboard {
    pub user_session: UserSession,
    pub compliance_overview: ComplianceOverview,
    pub security_status: SecurityStatusSummary,
    pub risk_dashboard: RiskDashboard,
    pub agent_status: AgentStatusSummary,
    pub cognitive_workspace: CognitiveWorkspaceState,
    pub monitoring_metrics: MonitoringMetrics,
    pub analytics_insights: AnalyticsInsights,
    pub system_health: SystemHealth,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GUIUpdate {
    pub update_id: String,
    pub update_type: GUIUpdateType,
    pub data: serde_json::Value,
    pub timestamp: chrono::DateTime<chrono::Utc>,
    pub target_users: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum GUIUpdateType {
    ComplianceChange,
    SecurityEvent,
    RiskAssessment,
    AgentInteraction,
    CognitiveCollaboration,
    MonitoringAlert,
    AnalyticsUpdate,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OptimizedInterface {
    pub interface_id: String,
    pub user_id: String,
    pub interface_type: InterfaceType,
    pub layout: InterfaceLayout,
    pub components: Vec<InterfaceComponent>,
    pub cognitive_optimizations: Vec<CognitiveOptimization>,
    pub adaptive_features: Vec<AdaptiveFeature>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum InterfaceType {
    SecurityAnalyst,
    ComplianceManager,
    RiskManager,
    IncidentResponder,
    Executive,
    Collaborative,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum InterfaceLayout {
    SecurityFocused,
    ComplianceFocused,
    RiskFocused,
    IncidentFocused,
    ExecutiveFocused,
    CollaborationFocused,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum InterfaceComponent {
    SecurityDashboard,
    ThreatMonitoring,
    IncidentResponse,
    AgentCollaboration,
    ComplianceDashboard,
    PolicyManagement,
    AuditTracking,
    ReportingTools,
    RiskDashboard,
    AssessmentTools,
    MitigationPlanning,
    VendorManagement,
    ExecutiveDashboard,
    StrategicReports,
    CognitiveWorkspace,
    SharedDocuments,
    CommunicationTools,
    EvidenceCollection,
    ResponsePlaybooks,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum CognitiveOptimization {
    ThreatPrioritization,
    PatternRecognition,
    QuickResponseActions,
    DocumentOrganization,
    RequirementTracking,
    AuditPreparation,
    RiskVisualization,
    ImpactAssessment,
    MitigationPrioritization,
    IncidentPrioritization,
    StepByStepGuidance,
    CommunicationOptimization,
    ExecutiveSummaries,
    KPIHighlighting,
    TrendVisualization,
    SharedUnderstanding,
    CommunicationClarity,
    DecisionSupport,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AdaptiveFeature {
    SimplifiedInterface,
    CriticalInformationHighlighting,
    ReducedCognitiveLoad,
    ProgressiveDisclosure,
    EnhancedCommunication,
    SharedWorkspace,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InterfaceContext {
    pub interface_type: InterfaceType,
    pub urgency: UrgencyLevel,
    pub cognitive_load: CognitiveLoad,
    pub collaboration_mode: bool,
    pub specific_task: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum UrgencyLevel {
    Low,
    Medium,
    High,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum CognitiveLoad {
    Normal,
    Elevated,
    High,
    Overloaded,
}

#[derive(Debug, Clone, Serialize, Default)]
pub struct GUISystemState {
    monitoring_data: Option<MonitoringMetrics>,
    user_sessions_count: usize,
    active_interfaces: usize,
    system_performance: SystemPerformance,
}

impl GUISystemState {
    pub fn new() -> Self {
        Self {
            monitoring_data: None,
            user_sessions_count: 0,
            active_interfaces: 0,
            system_performance: SystemPerformance::default(),
        }
    }

    pub fn update_monitoring_data(&mut self, data: MonitoringMetrics) {
        self.monitoring_data = Some(data);
    }

    pub fn get_system_health(&self) -> SystemHealth {
        SystemHealth {
            status: if self.system_performance.cpu_usage < 80.0 {
                HealthStatus::Healthy
            } else {
                HealthStatus::Degraded
            },
            performance: self.system_performance.clone(),
            active_sessions: self.user_sessions_count,
            last_updated: chrono::Utc::now(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SystemHealth {
    pub status: HealthStatus,
    pub performance: SystemPerformance,
    pub active_sessions: usize,
    pub last_updated: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum HealthStatus {
    Healthy,
    Degraded,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct SystemPerformance {
    pub cpu_usage: f64,
    pub memory_usage: f64,
    pub response_time_ms: u64,
    pub error_rate: f64,
}

// Error types

#[derive(Debug, thiserror::Error)]
pub enum GUIError {
    #[error("Session not found: {0}")]
    SessionNotFound(String),
    #[error("Component initialization failed: {0}")]
    ComponentInitializationFailed(String),
    #[error("Update processing failed: {0}")]
    UpdateProcessingFailed(String),
    #[error("Interface creation failed: {0}")]
    InterfaceCreationFailed(String),
    #[error("Configuration error: {0}")]
    ConfigurationError(String),
    #[error("Rendering error: {0}")]
    RenderingError(String),
}

// Placeholder configuration structures (would be defined in respective modules)

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DashboardConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CognitiveConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComplianceGUIConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityGUIConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RiskGUIConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentGUIConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MonitoringGUIConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsGUIConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SettingsGUIConfig;
