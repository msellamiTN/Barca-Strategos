pub mod agent;
pub mod security;
pub mod communication;
pub mod storage;
pub mod compliance;
pub mod collaboration;
pub mod ai;
pub mod gui;
pub mod monitoring;
pub mod runtime;

pub use agent::*;
pub use security::*;
pub use communication::*;
pub use storage::*;
pub use analytics::*;
pub use orchestration::*;
pub use config::*;
pub use monitoring::*;
pub use compliance::*;
pub use collaboration::*;
pub use ai::*;
pub use gui::*;

use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

/// Core Phoenix Framework Types
pub type AgentId = Uuid;
pub type TaskId = Uuid;
pub type MessageId = Uuid;
pub type SessionId = Uuid;
pub type SecurityEventId = Uuid;

/// Core trait for all Phoenix components
pub trait PhoenixComponent: Send + Sync {
    fn component_id(&self) -> ComponentId;
    fn component_type(&self) -> ComponentType;
    async fn initialize(&mut self) -> Result<(), ComponentError>;
    async fn shutdown(&mut self) -> Result<(), ComponentError>;
    async fn health_check(&self) -> ComponentHealth;
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComponentId(pub String);

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ComponentType {
    Agent,
    Runtime,
    Security,
    Collaboration,
    AI,
    Storage,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ComponentHealth {
    Healthy,
    Degraded(String),
    Unhealthy(String),
}

#[derive(Debug, thiserror::Error)]
pub enum ComponentError {
    #[error("Initialization failed: {0}")]
    InitializationFailed(String),
    #[error("Configuration error: {0}")]
    ConfigurationError(String),
    #[error("Security violation: {0}")]
    SecurityViolation(String),
    #[error("Resource exhausted: {0}")]
    ResourceExhausted(String),
    #[error("Communication error: {0}")]
    CommunicationError(String),
    #[error("Database error: {0}")]
    DatabaseError(String),
    #[error("AI model error: {0}")]
    AIModelError(String),
}

/// Core configuration structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PhoenixConfig {
    pub framework: FrameworkConfig,
    pub security: SecurityConfig,
    pub agents: AgentConfig,
    pub collaboration: CollaborationConfig,
    pub ai: AIConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FrameworkConfig {
    pub max_memory_mb: u32,
    pub max_agents: u32,
    pub boot_timeout_seconds: u32,
    pub log_level: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityConfig {
    pub enable_sandbox: bool,
    pub enable_encryption: bool,
    pub jwt_secret: String,
    pub certificate_path: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentConfig {
    pub default_memory_limit_mb: u32,
    pub max_execution_time_seconds: u32,
    pub enable_prompt_injection_protection: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CollaborationConfig {
    pub platforms: Vec<PlatformConfig>,
    pub enable_cross_platform_sync: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlatformConfig {
    pub platform: String,
    pub enabled: bool,
    pub config: serde_json::Value,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AIConfig {
    pub model_provider: String,
    pub api_key: String,
    pub max_tokens: u32,
    pub temperature: f32,
}

impl Default for PhoenixConfig {
    fn default() -> Self {
        Self {
            framework: FrameworkConfig {
                max_memory_mb: 10, // Ultra-lightweight by design
                max_agents: 10,
                boot_timeout_seconds: 1, // 1-second boot target
                log_level: "info".to_string(),
            },
            security: SecurityConfig {
                enable_sandbox: true,
                enable_encryption: true,
                jwt_secret: std::env::var("PHOENIX_JWT_SECRET")
                    .expect("PHOENIX_JWT_SECRET must be set in environment variables"),
                certificate_path: std::env::var("PHOENIX_CERT_PATH").ok(),
            },
            agents: AgentConfig {
                default_memory_limit_mb: 10,
                max_execution_time_seconds: 300,
                enable_prompt_injection_protection: true,
            },
            collaboration: CollaborationConfig {
                platforms: vec![],
                enable_cross_platform_sync: true,
            },
            ai: AIConfig {
                model_provider: std::env::var("PHOENIX_MODEL_PROVIDER")
                    .unwrap_or_else(|_| "openai".to_string()),
                api_key: std::env::var("PHOENIX_API_KEY")
                    .expect("PHOENIX_API_KEY must be set in environment variables"),
                max_tokens: std::env::var("PHOENIX_MAX_TOKENS")
                    .ok()
                    .and_then(|s| s.parse().ok())
                    .unwrap_or(1000),
                temperature: std::env::var("PHOENIX_TEMPERATURE")
                    .ok()
                    .and_then(|s| s.parse().ok())
                    .unwrap_or(0.7),
            },
        }
    }
}
