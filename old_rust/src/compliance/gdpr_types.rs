/// GDPR-specific types
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GDPRMetrics {
    pub total_requests: u64,
    pub processed_requests: u64,
    pub average_processing_time_ms: u64,
    pub compliance_score: f64,
    pub data_breaches: u64,
    pub last_updated: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GDPRControlUpdate {
    pub control_id: String,
    pub status: GDPRControlStatus,
    pub evidence: Vec<String>,
    pub updated_by: String,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum GDPRControlStatus {
    NotImplemented,
    PartiallyImplemented,
    Implemented,
    Compliant,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BreachNotificationStats {
    pub total_breaches: u64,
    pub notifications_sent: u64,
    pub average_notification_time_hours: f64,
    pub last_breach: Option<DateTime<Utc>>,
}
