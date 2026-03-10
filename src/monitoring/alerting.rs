// use crate::core::*;
// use crate::security::*;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::collections::{HashMap, VecDeque};
use std::sync::Arc;
use tokio::sync::{RwLock, mpsc};
use chrono::{DateTime, Utc, Duration};
use uuid::Uuid;

/// Real-time alerting system for Barca-Strategos Phoenix
/// Provides intelligent alerting with escalation and notification capabilities

pub struct AlertingSystem {
    alert_config: AlertConfig,
    alert_rules: Arc<RwLock<Vec<AlertRule>>>,
    alert_queue: Arc<RwLock<VecDeque<Alert>>>,
    escalation_engine: EscalationEngine,
    notification_service: NotificationService,
    alert_history: Arc<RwLock<Vec<Alert>>>,
    metrics_collector: Arc<RwLock<AlertMetrics>>,
}

impl AlertingSystem {
    pub fn new(config: AlertConfig) -> Self {
        Self {
            alert_config: config.clone(),
            alert_rules: Arc::new(RwLock::new(Vec::new())),
            alert_queue: Arc::new(RwLock::new(VecDeque::with_capacity(1000))),
            escalation_engine: EscalationEngine::new(&config.escalation_config),
            notification_service: NotificationService::new(&config.notification_config),
            alert_history: Arc::new(RwLock::new(Vec::with_capacity(10000))),
            metrics_collector: Arc::new(RwLock::new(AlertMetrics::default())),
        }
    }
    
    /// Initialize alerting system
    pub async fn initialize(&mut self) -> Result<(), AlertingError> {
        // Load default alert rules
        self.load_default_rules().await?;
        
        // Initialize escalation engine
        self.escalation_engine.initialize().await?;
        
        // Initialize notification service
        self.notification_service.initialize().await?;
        
        // Start background processing
        self.start_background_processing().await?;
        
        Ok(())
    }
    
    /// Process security event and generate alerts
    pub async fn process_security_event(&self, event: &SecurityEvent) -> Result<Vec<Alert>, AlertingError> {
        let mut alerts = Vec::new();
        let rules = self.alert_rules.read().await;
        
        // Evaluate against all alert rules
        for rule in rules.iter() {
            if let Some(alert) = rule.evaluate_security_event(event).await? {
                alerts.push(alert);
            }
        }
        
        // Process generated alerts
        for alert in alerts {
            self.process_alert(alert).await?;
        }
        
        Ok(alerts)
    }
    
    /// Process network event and generate alerts
    pub async fn process_network_event(&self, event: &NetworkSecurityEvent) -> Result<Vec<Alert>, AlertingError> {
        let mut alerts = Vec::new();
        let rules = self.alert_rules.read().await;
        
        // Evaluate against all alert rules
        for rule in rules.iter() {
            if let Some(alert) = rule.evaluate_network_event(event).await? {
                alerts.push(alert);
            }
        }
        
        // Process generated alerts
        for alert in alerts {
            self.process_alert(alert).await?;
        }
        
        Ok(alerts)
    }
    
    /// Process system event and generate alerts
    pub async fn process_system_event(&self, event: &SystemEvent) -> Result<Vec<Alert>, AlertingError> {
        let mut alerts = Vec::new();
        let rules = self.alert_rules.read().await;
        
        // Evaluate against all alert rules
        for rule in rules.iter() {
            if let Some(alert) = rule.evaluate_system_event(event).await? {
                alerts.push(alert);
            }
        }
        
        // Process generated alerts
        for alert in alerts {
            self.process_alert(alert).await?;
        }
        
        Ok(alerts)
    }
    
    /// Create custom alert rule
    pub async fn create_alert_rule(&self, rule: AlertRule) -> Result<(), AlertingError> {
        let mut rules = self.alert_rules.write().await;
        rules.push(rule);
        Ok(())
    }
    
    /// Update alert rule
    pub async fn update_alert_rule(&self, rule_id: &str, updates: AlertRule) -> Result<(), AlertingError> {
        let mut rules = self.alert_rules.write().await;
        if let Some(rule) = rules.iter_mut().find(|r| r.id == rule_id) {
            *rule = updates;
        } else {
            return Err(AlertingError::RuleNotFound(rule_id.to_string()));
        }
        Ok(())
    }
    
    /// Delete alert rule
    pub async fn delete_alert_rule(&self, rule_id: &str) -> Result<(), AlertingError> {
        let mut rules = self.alert_rules.write().await;
        rules.retain(|r| r.id != rule_id);
        Ok(())
    }
    
    /// Get active alerts
    pub async fn get_active_alerts(&self) -> Vec<Alert> {
        let queue = self.alert_queue.read().await;
        queue.iter().cloned().collect()
    }
    
    /// Get alert history
    pub async fn get_alert_history(&self, filter: &AlertFilter) -> Result<Vec<Alert>, AlertingError> {
        let history = self.alert_history.read().await;
        
        let filtered_alerts: Vec<Alert> = history.iter()
            .filter(|alert| self.matches_filter(alert, filter))
            .cloned()
            .collect();
        
        Ok(filtered_alerts)
    }
    
    /// Acknowledge alert
    pub async fn acknowledge_alert(&self, alert_id: &str, user: &str) -> Result<(), AlertingError> {
        let mut queue = self.alert_queue.write().await;
        if let Some(alert) = queue.iter_mut().find(|a| a.id == alert_id) {
            alert.status = AlertStatus::Acknowledged;
            alert.acknowledged_by = Some(user.to_string());
            alert.acknowledged_at = Some(Utc::now());
            
            // Update escalation engine
            self.escalation_engine.handle_acknowledgment(alert).await?;
        } else {
            return Err(AlertingError::AlertNotFound(alert_id.to_string()));
        }
        
        Ok(())
    }
    
    /// Resolve alert
    pub async fn resolve_alert(&self, alert_id: &str, resolution: &AlertResolution) -> Result<(), AlertingError> {
        let mut queue = self.alert_queue.write().await;
        if let Some(alert) = queue.iter_mut().find(|a| a.id == alert_id) {
            alert.status = AlertStatus::Resolved;
            alert.resolution = Some(resolution.clone());
            alert.resolved_at = Some(Utc::now());
            
            // Update escalation engine
            self.escalation_engine.handle_resolution(alert, resolution).await?;
        } else {
            return Err(AlertingError::AlertNotFound(alert_id.to_string()));
        }
        
        Ok(())
    }
    
    /// Get alerting statistics
    pub async fn get_alerting_stats(&self) -> AlertingStats {
        let metrics = self.metrics_collector.read().await;
        let queue = self.alert_queue.read().await;
        let history = self.alert_history.read().await;
        
        AlertingStats {
            total_alerts: metrics.total_alerts,
            active_alerts: queue.len(),
            resolved_alerts: history.iter().filter(|a| a.status == AlertStatus::Resolved).count(),
            acknowledged_alerts: history.iter().filter(|a| a.status == AlertStatus::Acknowledged).count(),
            average_resolution_time_minutes: metrics.average_resolution_time_minutes,
            escalation_rate: metrics.escalation_rate,
            false_positive_rate: metrics.false_positive_rate,
        }
    }
    
    // Private methods
    
    async fn load_default_rules(&self) -> Result<(), AlertingError> {
        let mut rules = self.alert_rules.write().await;
        
        // Default security alert rules
        rules.push(AlertRule {
            id: "critical_security_event".to_string(),
            name: "Critical Security Event".to_string(),
            description: "Alert on critical security events".to_string(),
            enabled: true,
            severity: AlertSeverity::Critical,
            conditions: vec![
                AlertCondition::EventType("SECURITY_EVENT".to_string()),
                AlertCondition::SeverityThreshold(80.0),
            ],
            actions: vec![
                AlertAction::ImmediateNotification,
                AlertAction::EscalateToLevel(2),
                AlertAction::CreateIncident,
            ],
            cooldown_minutes: 5,
            rate_limit_per_hour: 10,
        });
        
        // Default network alert rules
        rules.push(AlertRule {
            id: "ddos_attack_detected".to_string(),
            name: "DDoS Attack Detected".to_string(),
            description: "Alert on DDoS attack patterns".to_string(),
            enabled: true,
            severity: AlertSeverity::High,
            conditions: vec![
                AlertCondition::EventType("NETWORK_EVENT".to_string()),
                AlertCondition::ThresholdExceeded("requests_per_second".to_string(), 1000.0),
            ],
            actions: vec![
                AlertAction::ImmediateNotification,
                AlertAction::EscalateToLevel(1),
                AlertAction::BlockSource,
            ],
            cooldown_minutes: 2,
            rate_limit_per_hour: 5,
        });
        
        // Default system alert rules
        rules.push(AlertRule {
            id: "system_resource_exhaustion".to_string(),
            name: "System Resource Exhaustion".to_string(),
            description: "Alert on system resource exhaustion".to_string(),
            enabled: true,
            severity: AlertSeverity::Medium,
            conditions: vec![
                AlertCondition::EventType("SYSTEM_EVENT".to_string()),
                AlertCondition::ThresholdExceeded("cpu_usage".to_string(), 90.0),
            ],
            actions: vec![
                AlertAction::Notification,
                AlertAction::LogOnly,
            ],
            cooldown_minutes: 10,
            rate_limit_per_hour: 3,
        });
        
        Ok(())
    }
    
    async fn start_background_processing(&self) -> Result<(), AlertingError> {
        // Start background tasks
        tokio::spawn(self.background_alert_processor());
        tokio::spawn(self.background_escalation_checker());
        tokio::spawn(self.background_metrics_collector());
        Ok(())
    }
    
    async fn background_alert_processor(&self) {
        let mut interval = tokio::time::interval(Duration::seconds(1));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.process_alert_queue().await {
                eprintln!("Alerting: Error processing alert queue: {}", e);
            }
        }
    }
    
    async fn background_escalation_checker(&self) {
        let mut interval = tokio::time::interval(Duration::seconds(30));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.check_escalations().await {
                eprintln!("Alerting: Error checking escalations: {}", e);
            }
        }
    }
    
    async fn background_metrics_collector(&self) {
        let mut interval = tokio::time::interval(Duration::minutes(1));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.collect_metrics().await {
                eprintln!("Alerting: Error collecting metrics: {}", e);
            }
        }
    }
    
    async fn process_alert(&self, mut alert: Alert) -> Result<(), AlertingError> {
        // Add to queue
        let mut queue = self.alert_queue.write().await;
        queue.push_back(alert.clone());
        
        // Add to history
        let mut history = self.alert_history.write().await;
        history.push(alert.clone());
        
        // Keep history size limited
        if history.len() > 10000 {
            history.remove(0);
        }
        
        // Update metrics
        let mut metrics = self.metrics_collector.write().await;
        metrics.total_alerts += 1;
        
        // Process alert actions
        for action in &alert.actions {
            self.execute_alert_action(action, &alert).await?;
        }
        
        Ok(())
    }
    
    async fn process_alert_queue(&self) -> Result<(), AlertingError> {
        let mut queue = self.alert_queue.write().await;
        
        // Process alerts in order
        while let Some(alert) = queue.pop_front() {
            // Send notifications
            self.notification_service.send_alert(&alert).await?;
            
            // Check for escalation
            self.escalation_engine.check_escalation(&alert).await?;
            
            // Re-queue if still active
            if alert.status == AlertStatus::Active {
                queue.push_back(alert);
            }
        }
        
        Ok(())
    }
    
    async fn check_escalations(&self) -> Result<(), AlertingError> {
        let queue = self.alert_queue.read().await;
        
        for alert in queue.iter() {
            if alert.status == AlertStatus::Active {
                self.escalation_engine.check_escalation(alert).await?;
            }
        }
        
        Ok(())
    }
    
    async fn collect_metrics(&self) -> Result<(), AlertingError> {
        let history = self.alert_history.read().await;
        let mut metrics = self.metrics_collector.write().await;
        
        // Calculate average resolution time
        let resolved_alerts: Vec<&Alert> = history.iter()
            .filter(|a| a.status == AlertStatus::Resolved && a.resolved_at.is_some())
            .collect();
        
        if !resolved_alerts.is_empty() {
            let total_resolution_time: Duration = resolved_alerts.iter()
                .map(|a| a.resolved_at.unwrap() - a.created_at)
                .sum();
            
            metrics.average_resolution_time_minutes = total_resolution_time.num_minutes() as f64 / resolved_alerts.len() as f64;
        }
        
        // Calculate escalation rate
        let total_alerts = history.len() as f64;
        let escalated_alerts = history.iter().filter(|a| a.escalation_level > 0).count() as f64;
        metrics.escalation_rate = if total_alerts > 0.0 { escalated_alerts / total_alerts } else { 0.0 };
        
        // Calculate false positive rate
        let false_positive_alerts = history.iter().filter(|a| {
            a.resolution.as_ref().map_or(false, |r| r.resolution_type == ResolutionType::FalsePositive)
        }).count() as f64;
        metrics.false_positive_rate = if total_alerts > 0.0 { false_positive_alerts / total_alerts } else { 0.0 };
        
        Ok(())
    }
    
    async fn execute_alert_action(&self, action: &AlertAction, alert: &Alert) -> Result<(), AlertingError> {
        match action {
            AlertAction::ImmediateNotification => {
                self.notification_service.send_immediate_alert(alert).await?;
            },
            AlertAction::Notification => {
                self.notification_service.send_alert(alert).await?;
            },
            AlertAction::EscalateToLevel(level) => {
                alert.escalation_level = *level;
            },
            AlertAction::CreateIncident => {
                self.create_incident_from_alert(alert).await?;
            },
            AlertAction::BlockSource => {
                self.block_alert_source(alert).await?;
            },
            AlertAction::LogOnly => {
                // Already logged by default
            },
        }
        Ok(())
    }
    
    async fn create_incident_from_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Create incident from alert
        eprintln!("Alerting: Creating incident from alert {}", alert.id);
        Ok(())
    }
    
    async fn block_alert_source(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Block source of alert
        eprintln!("Alerting: Blocking source for alert {}", alert.id);
        Ok(())
    }
    
    fn matches_filter(&self, alert: &Alert, filter: &AlertFilter) -> bool {
        // Check if alert matches filter criteria
        if let Some(severity) = &filter.severity {
            if alert.severity != *severity {
                return false;
            }
        }
        
        if let Some(status) = &filter.status {
            if alert.status != *status {
                return false;
            }
        }
        
        if let Some(start_time) = filter.start_time {
            if alert.created_at < start_time {
                return false;
            }
        }
        
        if let Some(end_time) = filter.end_time {
            if alert.created_at > end_time {
                return false;
            }
        }
        
        true
    }
}

/// Alert structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Alert {
    pub id: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub severity: AlertSeverity,
    pub status: AlertStatus,
    pub title: String,
    pub description: String,
    pub source: AlertSource,
    pub event_data: serde_json::Value,
    pub actions: Vec<AlertAction>,
    pub escalation_level: u8,
    pub acknowledged_by: Option<String>,
    pub acknowledged_at: Option<DateTime<Utc>>,
    pub resolved_at: Option<DateTime<Utc>>,
    pub resolution: Option<AlertResolution>,
}

/// Alert severity levels
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum AlertSeverity {
    Info,
    Low,
    Medium,
    High,
    Critical,
}

/// Alert status
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum AlertStatus {
    Active,
    Acknowledged,
    Resolved,
    Suppressed,
}

/// Alert source
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AlertSource {
    Security,
    Network,
    System,
    Application,
    External,
}

/// Alert action
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AlertAction {
    ImmediateNotification,
    Notification,
    EscalateToLevel(u8),
    CreateIncident,
    BlockSource,
    LogOnly,
}

/// Alert rule
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlertRule {
    pub id: String,
    pub name: String,
    pub description: String,
    pub enabled: bool,
    pub severity: AlertSeverity,
    pub conditions: Vec<AlertCondition>,
    pub actions: Vec<AlertAction>,
    pub cooldown_minutes: u64,
    pub rate_limit_per_hour: u32,
}

impl AlertRule {
    pub async fn evaluate_security_event(&self, event: &SecurityEvent) -> Result<Option<Alert>, AlertingError> {
        if !self.enabled {
            return Ok(None);
        }
        
        // Check conditions
        for condition in &self.conditions {
            if !condition.matches_security_event(event) {
                return Ok(None);
            }
        }
        
        // Create alert
        Ok(Some(Alert {
            id: Uuid::new_v4().to_string(),
            created_at: Utc::now(),
            updated_at: Utc::now(),
            severity: self.severity.clone(),
            status: AlertStatus::Active,
            title: self.name.clone(),
            description: format!("Security event: {:?}", event.action),
            source: AlertSource::Security,
            event_data: serde_json::to_value(event).unwrap_or_default(),
            actions: self.actions.clone(),
            escalation_level: 0,
            acknowledged_by: None,
            acknowledged_at: None,
            resolved_at: None,
            resolution: None,
        }))
    }
    
    pub async fn evaluate_network_event(&self, event: &NetworkSecurityEvent) -> Result<Option<Alert>, AlertingError> {
        if !self.enabled {
            return Ok(None);
        }
        
        // Check conditions
        for condition in &self.conditions {
            if !condition.matches_network_event(event) {
                return Ok(None);
            }
        }
        
        // Create alert
        Ok(Some(Alert {
            id: Uuid::new_v4().to_string(),
            created_at: Utc::now(),
            updated_at: Utc::now(),
            severity: self.severity.clone(),
            status: AlertStatus::Active,
            title: self.name.clone(),
            description: format!("Network event: {}", event.action),
            source: AlertSource::Network,
            event_data: serde_json::to_value(event).unwrap_or_default(),
            actions: self.actions.clone(),
            escalation_level: 0,
            acknowledged_by: None,
            acknowledged_at: None,
            resolved_at: None,
            resolution: None,
        }))
    }
    
    pub async fn evaluate_system_event(&self, event: &SystemEvent) -> Result<Option<Alert>, AlertingError> {
        if !self.enabled {
            return Ok(None);
        }
        
        // Check conditions
        for condition in &self.conditions {
            if !condition.matches_system_event(event) {
                return Ok(None);
            }
        }
        
        // Create alert
        Ok(Some(Alert {
            id: Uuid::new_v4().to_string(),
            created_at: Utc::now(),
            updated_at: Utc::now(),
            severity: self.severity.clone(),
            status: AlertStatus::Active,
            title: self.name.clone(),
            description: format!("System event: {}", event.message),
            source: AlertSource::System,
            event_data: serde_json::to_value(event).unwrap_or_default(),
            actions: self.actions.clone(),
            escalation_level: 0,
            acknowledged_by: None,
            acknowledged_at: None,
            resolved_at: None,
            resolution: None,
        }))
    }
}

/// Alert condition
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AlertCondition {
    EventType(String),
    SeverityThreshold(f64),
    ThresholdExceeded(String, f64),
    FieldContains(String, String),
    TimeWindow(u64),
}

impl AlertCondition {
    pub fn matches_security_event(&self, event: &SecurityEvent) -> bool {
        match self {
            AlertCondition::EventType(event_type) => event_type == "SECURITY_EVENT",
            AlertCondition::SeverityThreshold(threshold) => event.risk_score >= *threshold,
            AlertCondition::ThresholdExceeded(field, threshold) => {
                // Check if field threshold is exceeded
                false // Placeholder
            },
            AlertCondition::FieldContains(field, value) => {
                // Check if field contains value
                false // Placeholder
            },
            AlertCondition::TimeWindow(minutes) => {
                // Check if event is within time window
                true // Placeholder
            },
        }
    }
    
    pub fn matches_network_event(&self, event: &NetworkSecurityEvent) -> bool {
        match self {
            AlertCondition::EventType(event_type) => event_type == "NETWORK_EVENT",
            AlertCondition::SeverityThreshold(threshold) => event.risk_score >= *threshold,
            AlertCondition::ThresholdExceeded(field, threshold) => {
                // Check if field threshold is exceeded
                false // Placeholder
            },
            AlertCondition::FieldContains(field, value) => {
                // Check if field contains value
                false // Placeholder
            },
            AlertCondition::TimeWindow(minutes) => {
                // Check if event is within time window
                true // Placeholder
            },
        }
    }
    
    pub fn matches_system_event(&self, event: &SystemEvent) -> bool {
        match self {
            AlertCondition::EventType(event_type) => event_type == "SYSTEM_EVENT",
            AlertCondition::SeverityThreshold(threshold) => (event.severity as f64) >= *threshold,
            AlertCondition::ThresholdExceeded(field, threshold) => {
                // Check if field threshold is exceeded
                false // Placeholder
            },
            AlertCondition::FieldContains(field, value) => {
                // Check if field contains value
                false // Placeholder
            },
            AlertCondition::TimeWindow(minutes) => {
                // Check if event is within time window
                true // Placeholder
            },
        }
    }
}

/// Alert resolution
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlertResolution {
    pub resolution_type: ResolutionType,
    pub resolved_by: String,
    pub notes: String,
    pub resolved_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ResolutionType {
    TruePositive,
    FalsePositive,
    Benign,
    Mitigated,
    Contained,
}

/// Alert filter
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlertFilter {
    pub severity: Option<AlertSeverity>,
    pub status: Option<AlertStatus>,
    pub start_time: Option<DateTime<Utc>>,
    pub end_time: Option<DateTime<Utc>>,
    pub source: Option<AlertSource>,
    pub limit: Option<u32>,
}

/// Alerting statistics
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct AlertingStats {
    pub total_alerts: u64,
    pub active_alerts: usize,
    pub resolved_alerts: usize,
    pub acknowledged_alerts: usize,
    pub average_resolution_time_minutes: f64,
    pub escalation_rate: f64,
    pub false_positive_rate: f64,
}

/// Alert metrics
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
struct AlertMetrics {
    pub total_alerts: u64,
    pub average_resolution_time_minutes: f64,
    pub escalation_rate: f64,
    pub false_positive_rate: f64,
}

/// Escalation engine
pub struct EscalationEngine {
    escalation_config: EscalationConfig,
    escalation_rules: Vec<EscalationRule>,
}

impl EscalationEngine {
    pub fn new(config: &EscalationConfig) -> Self {
        Self {
            escalation_config: config.clone(),
            escalation_rules: Vec::new(),
        }
    }
    
    pub async fn initialize(&mut self) -> Result<(), AlertingError> {
        // Load escalation rules
        self.escalation_rules = vec![
            EscalationRule {
                id: "critical_escalation".to_string(),
                name: "Critical Alert Escalation".to_string(),
                conditions: vec![
                    EscalationCondition::Severity(AlertSeverity::Critical),
                    EscalationCondition::TimeWithoutAcknowledgment(5), // 5 minutes
                ],
                actions: vec![
                    EscalationAction::NotifyManagement,
                    EscalationAction::CreateIncident,
                ],
            },
        ];
        
        Ok(())
    }
    
    pub async fn check_escalation(&self, alert: &Alert) -> Result<(), AlertingError> {
        for rule in &self.escalation_rules {
            if rule.matches(alert) {
                self.execute_escalation_actions(&rule.actions, alert).await?;
            }
        }
        Ok(())
    }
    
    pub async fn handle_acknowledgment(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Handle alert acknowledgment
        Ok(())
    }
    
    pub async fn handle_resolution(&self, alert: &Alert, resolution: &AlertResolution) -> Result<(), AlertingError> {
        // Handle alert resolution
        Ok(())
    }
    
    async fn execute_escalation_actions(&self, actions: &[EscalationAction], alert: &Alert) -> Result<(), AlertingError> {
        for action in actions {
            match action {
                EscalationAction::NotifyManagement => {
                    eprintln!("Escalation: Notifying management for alert {}", alert.id);
                },
                EscalationAction::CreateIncident => {
                    eprintln!("Escalation: Creating incident for alert {}", alert.id);
                },
                EscalationAction::EscalateToTeam(team) => {
                    eprintln!("Escalation: Escalating to team {} for alert {}", team, alert.id);
                },
            }
        }
        Ok(())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EscalationRule {
    pub id: String,
    pub name: String,
    pub conditions: Vec<EscalationCondition>,
    pub actions: Vec<EscalationAction>,
}

impl EscalationRule {
    pub fn matches(&self, alert: &Alert) -> bool {
        for condition in &self.conditions {
            if !condition.matches(alert) {
                return false;
            }
        }
        true
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum EscalationCondition {
    Severity(AlertSeverity),
    TimeWithoutAcknowledgment(u64), // minutes
    EscalationLevel(u8),
}

impl EscalationCondition {
    pub fn matches(&self, alert: &Alert) -> bool {
        match self {
            EscalationCondition::Severity(severity) => &alert.severity == severity,
            EscalationCondition::TimeWithoutAcknowledgment(minutes) => {
                if let Some(acknowledged_at) = alert.acknowledged_at {
                    Utc::now().signed_duration_since(acknowledged_at).num_minutes() > *minutes as i64
                } else {
                    Utc::now().signed_duration_since(alert.created_at).num_minutes() > *minutes as i64
                }
            },
            EscalationCondition::EscalationLevel(level) => alert.escalation_level >= *level,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum EscalationAction {
    NotifyManagement,
    CreateIncident,
    EscalateToTeam(String),
}

/// Notification service
pub struct NotificationService {
    notification_config: NotificationConfig,
    channels: Vec<Box<dyn NotificationChannel>>,
}

impl NotificationService {
    pub fn new(config: &NotificationConfig) -> Self {
        Self {
            notification_config: config.clone(),
            channels: Vec::new(),
        }
    }
    
    pub async fn initialize(&mut self) -> Result<(), AlertingError> {
        // Initialize notification channels
        if self.notification_config.enable_email {
            self.channels.push(Box::new(EmailNotificationChannel::new(&self.notification_config.email_config)));
        }
        
        if self.notification_config.enable_slack {
            self.channels.push(Box::new(SlackNotificationChannel::new(&self.notification_config.slack_config)));
        }
        
        if self.notification_config.enable_webhook {
            self.channels.push(Box::new(WebhookNotificationChannel::new(&self.notification_config.webhook_config)));
        }
        
        Ok(())
    }
    
    pub async fn send_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        for channel in &self.channels {
            if let Err(e) = channel.send_alert(alert).await {
                eprintln!("Notification: Failed to send via channel: {}", e);
            }
        }
        Ok(())
    }
    
    pub async fn send_immediate_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        for channel in &self.channels {
            if let Err(e) = channel.send_immediate_alert(alert).await {
                eprintln!("Notification: Failed to send immediate alert via channel: {}", e);
            }
        }
        Ok(())
    }
}

#[async_trait]
pub trait NotificationChannel: Send + Sync {
    async fn send_alert(&self, alert: &Alert) -> Result<(), AlertingError>;
    async fn send_immediate_alert(&self, alert: &Alert) -> Result<(), AlertingError>;
}

/// Email notification channel
pub struct EmailNotificationChannel {
    email_config: EmailConfig,
}

impl EmailNotificationChannel {
    pub fn new(config: &EmailConfig) -> Self {
        Self {
            email_config: config.clone(),
        }
    }
}

#[async_trait]
impl NotificationChannel for EmailNotificationChannel {
    async fn send_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Send email notification
        eprintln!("Email: Sending alert {}", alert.title);
        Ok(())
    }
    
    async fn send_immediate_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Send immediate email notification
        eprintln!("Email: Sending immediate alert {}", alert.title);
        Ok(())
    }
}

/// Slack notification channel
pub struct SlackNotificationChannel {
    slack_config: SlackConfig,
}

impl SlackNotificationChannel {
    pub fn new(config: &SlackConfig) -> Self {
        Self {
            slack_config: config.clone(),
        }
    }
}

#[async_trait]
impl NotificationChannel for SlackNotificationChannel {
    async fn send_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Send Slack notification
        eprintln!("Slack: Sending alert {}", alert.title);
        Ok(())
    }
    
    async fn send_immediate_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Send immediate Slack notification
        eprintln!("Slack: Sending immediate alert {}", alert.title);
        Ok(())
    }
}

/// Webhook notification channel
pub struct WebhookNotificationChannel {
    webhook_config: WebhookConfig,
}

impl WebhookNotificationChannel {
    pub fn new(config: &WebhookConfig) -> Self {
        Self {
            webhook_config: config.clone(),
        }
    }
}

#[async_trait]
impl NotificationChannel for WebhookNotificationChannel {
    async fn send_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Send webhook notification
        eprintln!("Webhook: Sending alert {}", alert.title);
        Ok(())
    }
    
    async fn send_immediate_alert(&self, alert: &Alert) -> Result<(), AlertingError> {
        // Send immediate webhook notification
        eprintln!("Webhook: Sending immediate alert {}", alert.title);
        Ok(())
    }
}

/// Configuration structures

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlertConfig {
    pub escalation_config: EscalationConfig,
    pub notification_config: NotificationConfig,
    pub max_alerts_per_minute: u32,
    pub alert_retention_days: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EscalationConfig {
    pub enable_escalation: bool,
    pub max_escalation_level: u8,
    pub escalation_intervals_minutes: Vec<u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NotificationConfig {
    pub enable_email: bool,
    pub enable_slack: bool,
    pub enable_webhook: bool,
    pub email_config: EmailConfig,
    pub slack_config: SlackConfig,
    pub webhook_config: WebhookConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EmailConfig {
    pub smtp_server: String,
    pub smtp_port: u16,
    pub username: String,
    pub password: String,
    pub from_address: String,
    pub to_addresses: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SlackConfig {
    pub webhook_url: String,
    pub channel: String,
    pub username: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WebhookConfig {
    pub url: String,
    pub headers: HashMap<String, String>,
    pub timeout_seconds: u64,
}

/// Alerting errors
#[derive(Debug, thiserror::Error)]
pub enum AlertingError {
    #[error("Configuration error: {0}")]
    ConfigurationError(String),
    
    #[error("Rule not found: {0}")]
    RuleNotFound(String),
    
    #[error("Alert not found: {0}")]
    AlertNotFound(String),
    
    #[error("Notification error: {0}")]
    NotificationError(String),
    
    #[error("Escalation error: {0}")]
    EscalationError(String),
    
    #[error("Rate limit exceeded")]
    RateLimitExceeded,
}
