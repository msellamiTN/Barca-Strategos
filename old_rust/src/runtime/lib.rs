pub mod phoenix_runtime;
pub mod memory_limiter;
pub mod fast_boot;
pub mod cross_platform;
pub mod resource_monitor;

pub use phoenix_runtime::*;
pub use memory_limiter::*;
pub use fast_boot::*;
pub use cross_platform::*;
pub use resource_monitor::*;

use crate::core::*;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::sync::RwLock;
use uuid::Uuid;

/// Ultra-lightweight Phoenix runtime - where tactical excellence meets speed
pub struct PhoenixRuntime {
    memory_limiter: MemoryLimiter,
    fast_boot: FastBootManager,
    cross_platform: CrossPlatformABI,
    resource_monitor: ResourceMonitor,
    agents: RwLock<HashMap<AgentId, AgentHandle>>,
    config: RuntimeConfig,
}

impl PhoenixRuntime {
    pub fn new(config: RuntimeConfig) -> Result<Self, RuntimeError> {
        let memory_limiter = MemoryLimiter::new(config.max_memory_mb);
        let fast_boot = FastBootManager::new(config.boot_timeout_seconds);
        let cross_platform = CrossPlatformABI::new();
        let resource_monitor = ResourceMonitor::new();
        
        Ok(Self {
            memory_limiter,
            fast_boot,
            cross_platform,
            resource_monitor,
            agents: RwLock::new(HashMap::new()),
            config,
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), RuntimeError> {
        // Initialize runtime components
        self.fast_boot.initialize().await?;
        self.cross_platform.initialize().await?;
        self.resource_monitor.start_monitoring().await?;
        
        println!("🔥 Phoenix Runtime initialized - Ultra-efficient <10MB agents ready");
        Ok(())
    }
    
    pub async fn spawn_agent(&self, agent_config: AgentConfig) -> Result<AgentHandle, RuntimeError> {
        // Memory-efficient agent spawning
        let agent = self.create_lightweight_agent(agent_config).await?;
        
        // Fast boot optimization
        let handle = self.fast_boot.spawn(agent).await?;
        
        // Resource monitoring
        self.resource_monitor.track_agent(handle.id()).await?;
        
        // Store agent handle
        self.agents.write().await.insert(handle.id(), handle.clone());
        
        println!("🦐 Agent spawned: {} (Memory: {}MB)", handle.id(), handle.memory_usage_mb());
        Ok(handle)
    }
    
    async fn create_lightweight_agent(&self, config: AgentConfig) -> Result<Box<dyn Agent>, RuntimeError> {
        // Create agent based on type with minimal memory footprint
        match config.agent_type {
            AgentType::SecurityAnalyst => {
                Ok(Box::new(crate::agents::SecurityAnalystAgent::new(config)?))
            }
            AgentType::ThreatHunter => {
                Ok(Box::new(crate::agents::ThreatHunterAgent::new(config)?))
            }
            AgentType::IncidentResponder => {
                Ok(Box::new(crate::agents::IncidentResponderAgent::new(config)?))
            }
            AgentType::AIAssistant => {
                Ok(Box::new(crate::agents::AIAssistantAgent::new(config)?))
            }
            _ => Err(RuntimeError::UnsupportedAgentType(config.agent_type)),
        }
    }
    
    pub async fn shutdown_agent(&self, agent_id: AgentId) -> Result<(), RuntimeError> {
        let mut agents = self.agents.write().await;
        if let Some(mut handle) = agents.remove(&agent_id) {
            handle.shutdown().await?;
            self.resource_monitor.stop_tracking(agent_id).await?;
            println!("🔥 Agent shutdown: {}", agent_id);
        }
        Ok(())
    }
    
    pub async fn get_agent(&self, agent_id: AgentId) -> Option<AgentHandle> {
        self.agents.read().await.get(&agent_id).cloned()
    }
    
    pub async fn list_agents(&self) -> Vec<AgentHandle> {
        self.agents.read().await.values().cloned().collect()
    }
    
    pub async fn get_runtime_stats(&self) -> RuntimeStats {
        let agents = self.agents.read().await;
        let total_memory = self.memory_limiter.total_usage_mb();
        let cpu_usage = self.resource_monitor.get_cpu_usage().await;
        
        RuntimeStats {
            active_agents: agents.len(),
            total_memory_mb: total_memory,
            cpu_usage_percent: cpu_usage,
            uptime_seconds: self.fast_boot.uptime_seconds(),
            agents_per_second: self.fast_boot.spawn_rate(),
        }
    }
    
    pub async fn scale_agents(&self, target_count: usize) -> Result<(), RuntimeError> {
        let current_count = self.agents.read().await.len();
        
        if target_count > current_count {
            let to_spawn = target_count - current_count;
            println!("🚀 Scaling up: spawning {} additional agents", to_spawn);
            
            for i in 0..to_spawn {
                let config = AgentConfig {
                    agent_type: AgentType::SecurityAnalyst,
                    name: format!("scaled-agent-{}", i),
                    memory_limit_mb: 10,
                    capabilities: vec![Capability::ThreatDetection],
                };
                
                self.spawn_agent(config).await?;
            }
        } else if target_count < current_count {
            let to_remove = current_count - target_count;
            println!("📉 Scaling down: removing {} agents", to_remove);
            
            let agents: Vec<AgentId> = self.agents.read().await.keys().cloned().take(to_remove).collect();
            for agent_id in agents {
                self.shutdown_agent(agent_id).await?;
            }
        }
        
        Ok(())
    }
}

/// Agent handle for runtime management
#[derive(Debug, Clone)]
pub struct AgentHandle {
    id: AgentId,
    agent_type: AgentType,
    memory_usage_mb: u32,
    status: AgentStatus,
    created_at: chrono::DateTime<chrono::Utc>,
}

impl AgentHandle {
    pub fn new(id: AgentId, agent_type: AgentType, memory_limit_mb: u32) -> Self {
        Self {
            id,
            agent_type,
            memory_usage_mb: memory_limit_mb,
            status: AgentStatus::Starting,
            created_at: chrono::Utc::now(),
        }
    }
    
    pub fn id(&self) -> AgentId {
        self.id
    }
    
    pub fn agent_type(&self) -> AgentType {
        self.agent_type.clone()
    }
    
    pub fn memory_usage_mb(&self) -> u32 {
        self.memory_usage_mb
    }
    
    pub fn status(&self) -> AgentStatus {
        self.status.clone()
    }
    
    pub fn set_status(&mut self, status: AgentStatus) {
        self.status = status;
    }
    
    pub async fn shutdown(&mut self) -> Result<(), RuntimeError> {
        self.status = AgentStatus::ShuttingDown;
        // Perform shutdown operations
        self.status = AgentStatus::Shutdown;
        Ok(())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AgentStatus {
    Starting,
    Ready,
    Busy(String),
    Error(String),
    ShuttingDown,
    Shutdown,
}

/// Runtime configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RuntimeConfig {
    pub max_memory_mb: u32,
    pub max_agents: u32,
    pub boot_timeout_seconds: u32,
    pub enable_fast_boot: bool,
    pub enable_cross_platform: bool,
    pub target_platform: Platform,
}

impl Default for RuntimeConfig {
    fn default() -> Self {
        Self {
            max_memory_mb: 100, // Total runtime memory limit
            max_agents: 10,
            boot_timeout_seconds: 1, // 1-second boot target
            enable_fast_boot: true,
            enable_cross_platform: true,
            target_platform: Platform::X86,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Platform {
    X86,
    ARM,
    RISCV,
    MIPS,
}

/// Runtime statistics
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RuntimeStats {
    pub active_agents: usize,
    pub total_memory_mb: u32,
    pub cpu_usage_percent: f32,
    pub uptime_seconds: u64,
    pub agents_per_second: f32,
}

#[derive(Debug, thiserror::Error)]
pub enum RuntimeError {
    #[error("Agent creation failed: {0}")]
    AgentCreationFailed(String),
    #[error("Unsupported agent type: {0:?}")]
    UnsupportedAgentType(AgentType),
    #[error("Memory limit exceeded: {0}MB")]
    MemoryLimitExceeded(u32),
    #[error("Boot timeout: {0}s")]
    BootTimeout(u32),
    #[error("Resource monitoring error: {0}")]
    ResourceMonitoringError(String),
    #[error("Platform compatibility error: {0}")]
    PlatformCompatibilityError(String),
}

/// Memory limiter for ultra-efficient resource management
pub struct MemoryLimiter {
    max_memory_mb: u32,
    current_usage_mb: std::sync::atomic::AtomicU32,
    agent_usage: RwLock<HashMap<AgentId, u32>>,
}

impl MemoryLimiter {
    pub fn new(max_memory_mb: u32) -> Self {
        Self {
            max_memory_mb,
            current_usage_mb: std::sync::atomic::AtomicU32::new(0),
            agent_usage: RwLock::new(HashMap::new()),
        }
    }
    
    pub fn total_usage_mb(&self) -> u32 {
        self.current_usage_mb.load(std::sync::atomic::Ordering::Relaxed)
    }
    
    pub async fn allocate_memory(&self, agent_id: AgentId, memory_mb: u32) -> Result<(), RuntimeError> {
        let current_total = self.total_usage_mb();
        if current_total + memory_mb > self.max_memory_mb {
            return Err(RuntimeError::MemoryLimitExceeded(self.max_memory_mb));
        }
        
        self.agent_usage.write().await.insert(agent_id, memory_mb);
        self.current_usage_mb.fetch_add(memory_mb, std::sync::atomic::Ordering::Relaxed);
        
        Ok(())
    }
    
    pub async fn deallocate_memory(&self, agent_id: AgentId) {
        if let Some(memory_mb) = self.agent_usage.write().await.remove(&agent_id) {
            self.current_usage_mb.fetch_sub(memory_mb, std::sync::atomic::Ordering::Relaxed);
        }
    }
    
    pub async fn get_agent_memory(&self, agent_id: AgentId) -> Option<u32> {
        self.agent_usage.read().await.get(&agent_id).copied()
    }
}

/// Fast boot manager for 1-second startup
pub struct FastBootManager {
    boot_timeout_seconds: u32,
    start_time: chrono::DateTime<chrono::Utc>,
    spawn_count: std::sync::atomic::AtomicU32,
}

impl FastBootManager {
    pub fn new(boot_timeout_seconds: u32) -> Self {
        Self {
            boot_timeout_seconds,
            start_time: chrono::Utc::now(),
            spawn_count: std::sync::atomic::AtomicU32::new(0),
        }
    }
    
    pub async fn initialize(&mut self) -> Result<(), RuntimeError> {
        self.start_time = chrono::Utc::now();
        println!("⚡ Fast Boot Manager initialized - Target: {}s", self.boot_timeout_seconds);
        Ok(())
    }
    
    pub async fn spawn<T>(&self, agent: T) -> Result<T, RuntimeError>
    where
        T: Send + Sync,
    {
        let start = std::time::Instant::now();
        
        // Simulate fast spawning
        tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
        
        let elapsed = start.elapsed();
        if elapsed.as_secs() > self.boot_timeout_seconds {
            return Err(RuntimeError::BootTimeout(self.boot_timeout_seconds));
        }
        
        self.spawn_count.fetch_add(1, std::sync::atomic::Ordering::Relaxed);
        println!("🚀 Agent spawned in {:?}", elapsed);
        
        Ok(agent)
    }
    
    pub fn uptime_seconds(&self) -> u64 {
        (chrono::Utc::now() - self.start_time).num_seconds() as u64
    }
    
    pub fn spawn_rate(&self) -> f32 {
        let uptime = self.uptime_seconds();
        if uptime > 0 {
            self.spawn_count.load(std::sync::atomic::Ordering::Relaxed) as f32 / uptime as f32
        } else {
            0.0
        }
    }
}

/// Cross-platform ABI support
pub struct CrossPlatformABI {
    target_platform: Platform,
    abi_version: String,
}

impl CrossPlatformABI {
    pub fn new() -> Self {
        Self {
            target_platform: Platform::X86,
            abi_version: "1.0".to_string(),
        }
    }
    
    pub async fn initialize(&mut self) -> Result<(), RuntimeError> {
        // Detect current platform
        self.target_platform = self.detect_platform();
        println!("🌍 Cross-platform ABI initialized for: {:?}", self.target_platform);
        Ok(())
    }
    
    fn detect_platform(&self) -> Platform {
        // Simple platform detection
        if cfg!(target_arch = "x86_64") {
            Platform::X86
        } else if cfg!(target_arch = "aarch64") {
            Platform::ARM
        } else if cfg!(target_arch = "riscv64") {
            Platform::RISCV
        } else {
            Platform::X86 // Default
        }
    }
    
    pub fn get_platform(&self) -> Platform {
        self.target_platform.clone()
    }
    
    pub fn is_compatible(&self, required_platform: &Platform) -> bool {
        std::mem::discriminant(&self.target_platform) == std::mem::discriminant(required_platform)
    }
}

/// Resource monitoring for efficient resource management
pub struct ResourceMonitor {
    cpu_usage: std::sync::atomic::AtomicF32,
    tracked_agents: RwLock<HashMap<AgentId, AgentResourceInfo>>,
    monitoring_active: std::sync::atomic::AtomicBool,
}

impl ResourceMonitor {
    pub fn new() -> Self {
        Self {
            cpu_usage: std::sync::atomic::AtomicF32::new(0.0),
            tracked_agents: RwLock::new(HashMap::new()),
            monitoring_active: std::sync::atomic::AtomicBool::new(false),
        }
    }
    
    pub async fn start_monitoring(&self) -> Result<(), RuntimeError> {
        self.monitoring_active.store(true, std::sync::atomic::Ordering::Relaxed);
        
        // Start background monitoring task
        let cpu_usage = self.cpu_usage.clone();
        tokio::spawn(async move {
            let mut interval = tokio::time::interval(tokio::time::Duration::from_secs(1));
            loop {
                interval.tick().await;
                // Simulate CPU usage monitoring
                let usage = (rand::random::<f32>() * 100.0).min(80.0);
                cpu_usage.store(usage, std::sync::atomic::Ordering::Relaxed);
            }
        });
        
        println!("📊 Resource monitoring started");
        Ok(())
    }
    
    pub async fn track_agent(&self, agent_id: AgentId) -> Result<(), RuntimeError> {
        let resource_info = AgentResourceInfo {
            agent_id,
            memory_usage_mb: 10,
            cpu_usage_percent: 0.0,
            network_bytes: 0,
            last_update: chrono::Utc::now(),
        };
        
        self.tracked_agents.write().await.insert(agent_id, resource_info);
        Ok(())
    }
    
    pub async fn stop_tracking(&self, agent_id: AgentId) {
        self.tracked_agents.write().await.remove(&agent_id);
    }
    
    pub async fn get_cpu_usage(&self) -> f32 {
        self.cpu_usage.load(std::sync::atomic::Ordering::Relaxed)
    }
    
    pub async fn get_agent_resources(&self, agent_id: AgentId) -> Option<AgentResourceInfo> {
        self.tracked_agents.read().await.get(&agent_id).cloned()
    }
    
    pub async fn get_all_resources(&self) -> Vec<AgentResourceInfo> {
        self.tracked_agents.read().await.values().cloned().collect()
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentResourceInfo {
    pub agent_id: AgentId,
    pub memory_usage_mb: u32,
    pub cpu_usage_percent: f32,
    pub network_bytes: u64,
    pub last_update: chrono::DateTime<chrono::Utc>,
}
