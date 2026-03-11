/// Common types shared across the Barca-Strategos Phoenix GUI system
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use chrono::{DateTime, Utc};
use uuid::Uuid;
use crate::gui::{InterfaceComplexity, LayoutType, ThemeType};

// User and session management types
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct UserId(pub String);

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: UserId,
    pub username: String,
    pub email: String,
    pub role: UserRole,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum UserRole {
    Admin,
    SecurityAnalyst,
    ComplianceOfficer,
    RiskManager,
    Viewer,
}

impl UserRole {
    pub fn has_system_access(&self) -> bool {
        matches!(self, UserRole::Admin | UserRole::SecurityAnalyst)
    }
    
    pub fn has_security_access(&self) -> bool {
        matches!(self, UserRole::Admin | UserRole::SecurityAnalyst)
    }
    
    pub fn has_compliance_access(&self) -> bool {
        matches!(self, UserRole::Admin | UserRole::ComplianceOfficer | UserRole::RiskManager)
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserPreferences {
    pub theme: String,
    pub language: String,
    pub notifications_enabled: bool,
    pub dashboard_layout: DashboardLayout,
    pub interface_complexity: InterfaceComplexity,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DashboardLayout {
    pub widgets: Vec<WidgetPosition>,
    pub grid_size: (u32, u32),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WidgetPosition {
    pub widget_id: String,
    pub x: u32,
    pub y: u32,
    pub width: u32,
    pub height: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserSession {
    pub user: User,
    pub preferences: UserPreferences,
    pub active_workspaces: Vec<Workspace>,
    pub collaboration_context: Option<CollaborationContext>,
    pub session_id: String,
    pub created_at: DateTime<Utc>,
    pub last_activity: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Workspace {
    pub id: String,
    pub name: String,
    pub workspace_type: WorkspaceType,
    pub members: Vec<UserId>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum WorkspaceType {
    SecurityOperations,
    ComplianceManagement,
    RiskAssessment,
    IncidentResponse,
    Collaboration,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CollaborationContext {
    pub workspace_id: String,
    pub active_channels: Vec<String>,
    pub shared_documents: Vec<String>,
    pub participant_permissions: HashMap<UserId, PermissionLevel>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum PermissionLevel {
    Read,
    Write,
    Admin,
}

// GUI adaptation types
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InterfaceAdaptation {
    pub user_id: UserId,
    pub adaptation_type: AdaptationType,
    pub parameters: HashMap<String, String>,
    pub created_at: DateTime<Utc>,
    pub expires_at: Option<DateTime<Utc>>,
    pub adaptive_features: Vec<AdaptiveFeature>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum AdaptationType {
    ColorBlindMode,
    HighContrast,
    LargeText,
    SimplifiedInterface,
    CustomLayout,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum InterfaceComplexity {
    Simplified,
    Standard,
    Advanced,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum LayoutType {
    Default,
    Compact,
    Detailed,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum ThemeType {
    Light,
    Dark,
    Auto,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum AdaptiveFeature {
    SimplifiedInterface,
    CriticalInformationHighlighting,
    ReducedCognitiveLoad,
    ProgressiveDisclosure,
    EnhancedCommunication,
    SharedWorkspace,
}

// Agent types (separate from AlertType)
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum AgentType {
    ThreatAnalyst,
    ComplianceChecker,
    RiskAssessor,
    IncidentResponder,
    SecurityAdvisor,
    AutomationBot,
}

// Shared enums across modules
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub enum RiskLevel {
    Low,
    Medium,
    High,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub enum FindingSeverity {
    Low,
    Medium,
    High,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub enum RecommendationPriority {
    Low,
    Medium,
    High,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum UpdateType {
    Implementation,
    Assessment,
    Review,
    Evidence,
    Policy,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub enum ThreatSeverity {
    Info,
    Low,
    Medium,
    High,
    Critical,
}

// GUI component types
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum EngineStatus {
    Active,
    Inactive,
    Maintenance,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SettingsConfigurationManager {
    pub manager_id: String,
    pub config_version: String,
    pub settings: HashMap<String, String>,
    pub last_updated: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsInsights {
    pub insights: Vec<Insight>,
    pub generated_at: DateTime<Utc>,
    pub confidence_score: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Insight {
    pub id: String,
    pub insight_type: InsightType,
    pub description: String,
    pub confidence: f64,
    pub recommendations: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum InsightType {
    SecurityThreat,
    ComplianceGap,
    RiskTrend,
    PerformanceIssue,
}

// Monitoring and configuration types
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MonitoringConfig {
    pub check_interval_seconds: u64,
    pub enable_real_time: bool,
    pub alert_threshold: f64,
    pub alert_thresholds: HashMap<String, f64>,
    pub retention_days: u32,
    pub enabled_monitors: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RiskConfig {
    pub assessment_interval_hours: u64,
    pub risk_threshold: f64,
    pub enable_automated_assessment: bool,
}

// Compliance types
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComplianceMonitor {
    pub monitor_id: String,
    pub framework: String,
    pub status: MonitorStatus,
    pub last_check: DateTime<Utc>,
}

impl ComplianceMonitor {
    pub fn new(_config: &MonitoringConfig) -> Self {
        Self {
            monitor_id: "default".to_string(),
            framework: "multi".to_string(),
            status: MonitorStatus::Active,
            last_check: Utc::now(),
        }
    }
    
    pub async fn initialize(&mut self) -> Result<(), String> {
        // Initialize the monitor
        self.last_check = Utc::now();
        Ok(())
    }
    
    pub async fn log_function_update(&mut self, _function_id: &str, _update: &str) -> Result<(), String> {
        // Log function update
        self.last_check = Utc::now();
        Ok(())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum MonitorStatus {
    Active,
    Inactive,
    Error,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RiskAssessmentEngine {
    pub engine_id: String,
    pub assessment_methods: Vec<String>,
    pub risk_model: String,
}

impl RiskAssessmentEngine {
    pub fn new(_config: &RiskConfig) -> Self {
        Self {
            engine_id: "default".to_string(),
            assessment_methods: vec!["quantitative".to_string(), "qualitative".to_string()],
            risk_model: "standard".to_string(),
        }
    }
}

// Cryptography constants - commented out until base64 dependency is resolved
// pub const BASE64: base64::engine::GeneralPurpose = base64::engine::general_purpose::STANDARD;

// Analytics and intelligence types
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsIntelligenceEngine {
    pub engine_id: String,
    pub model_version: String,
    pub capabilities: Vec<String>,
    pub status: EngineStatus,
}

impl AnalyticsIntelligenceEngine {
    pub fn new(_config: &AnalyticsConfig) -> Self {
        Self {
            engine_id: "analytics_default".to_string(),
            model_version: "1.0".to_string(),
            capabilities: vec!["threat_detection".to_string(), "anomaly_detection".to_string()],
            status: EngineStatus::Active,
        }
    }
    
    pub async fn handle_analytics_update(&mut self, _data: serde_json::Value) -> Result<(), String> {
        // Handle analytics update
        Ok(())
    }
    
    pub async fn process_all_data(&self) -> Result<AnalyticsData, String> {
        Ok(AnalyticsData {
            insights: Vec::new(),
            recommendations: Vec::new(),
            trends: Vec::new(),
        })
    }
    
    pub async fn get_interface_adaptation(&self, _user_id: &str) -> Result<InterfaceAdaptation, String> {
        Ok(InterfaceAdaptation {
            complexity: InterfaceComplexity::Standard,
            layout: LayoutType::Default,
            theme: ThemeType::Light,
        })
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsData {
    pub insights: Vec<Insight>,
    pub recommendations: Vec<String>,
    pub trends: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsConfig;
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SettingsConfig;

impl SettingsConfigurationManager {
    pub fn new(_config: &SettingsConfig) -> Self {
        Self {
            manager_id: "settings_default".to_string(),
            config_version: "1.0".to_string(),
            settings: HashMap::new(),
            last_updated: Utc::now(),
        }
    }
}
