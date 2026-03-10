/// Common types shared across the Barca-Strategos Phoenix GUI system
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use chrono::{DateTime, Utc};
use uuid::Uuid;
use base64;

// User and session management types
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct UserId(String);

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

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserPreferences {
    pub theme: String,
    pub language: String,
    pub notifications_enabled: bool,
    pub dashboard_layout: DashboardLayout,
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
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum AdaptationType {
    ColorBlindMode,
    HighContrast,
    LargeText,
    SimplifiedInterface,
    CustomLayout,
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
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsIntelligenceEngine {
    pub engine_id: String,
    pub capabilities: Vec<String>,
    pub status: EngineStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum EngineStatus {
    Active,
    Inactive,
    Maintenance,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SettingsConfigurationManager {
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
    pub alert_thresholds: HashMap<String, f64>,
    pub retention_days: u32,
    pub enabled_monitors: Vec<String>,
}

// Compliance types
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComplianceMonitor {
    pub monitor_id: String,
    pub framework: String,
    pub status: MonitorStatus,
    pub last_check: DateTime<Utc>,
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

// Cryptography constants
pub const BASE64: base64::engine::GeneralPurpose = base64::engine::general_purpose::STANDARD;
