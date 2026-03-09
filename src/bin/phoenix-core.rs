use barca_strategos::core::*;
use barca_strategos::runtime::*;
use barca_strategos::collaboration::*;
use barca_strategos::ai::*;
use barca_strategos::gamification::*;
use tokio;
use std::env;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logging
    env_logger::init();
    
    println!("🔥 Starting Barca-Strategos Phoenix Core");
    println!("🦐 Ultra-Efficient AI Security Framework");
    
    // Load configuration
    let config = PhoenixConfig::from_env()?;
    
    // Initialize runtime
    let mut runtime = PhoenixRuntime::new(config.runtime)?;
    runtime.initialize().await?;
    
    // Initialize collaboration hub
    let mut collab_hub = UnifiedCollaborationHub::new(config.collaboration)?;
    collab_hub.initialize().await?;
    
    // Initialize AI assistant
    let mut ai_assistant = BarcaAIAssistant::new(&config.ai)?;
    ai_assistant.initialize().await?;
    
    // Initialize gamification
    let mut quest_engine = SecurityQuestEngine::new(config.gamification)?;
    quest_engine.initialize().await?;
    
    println!("✅ Phoenix Core initialized successfully");
    println!("🚀 All systems ready - Ultra-efficient agents deployed");
    
    // Start main application loop
    start_phoenix_server(runtime, collab_hub, ai_assistant, quest_engine).await?;
    
    Ok(())
}

async fn start_phoenix_server(
    runtime: PhoenixRuntime,
    collab_hub: UnifiedCollaborationHub,
    ai_assistant: BarcaAIAssistant,
    quest_engine: SecurityQuestEngine,
) -> Result<(), Box<dyn std::error::Error>> {
    println!("🌐 Starting Phoenix web server on http://0.0.0.0:8080");
    
    // Create sample agents for demonstration
    for i in 0..3 {
        let agent_config = AgentConfig {
            agent_type: AgentType::SecurityAnalyst,
            name: format!("phoenix-agent-{}", i),
            memory_limit_mb: 10,
            capabilities: vec![Capability::ThreatDetection],
        };
        
        match runtime.spawn_agent(agent_config).await {
            Ok(handle) => {
                println!("🦐 Spawned agent: {} (Memory: {}MB)", handle.id(), handle.memory_usage_mb());
            }
            Err(e) => {
                eprintln!("❌ Failed to spawn agent {}: {}", i, e);
            }
        }
    }
    
    // Show runtime stats
    let stats = runtime.get_runtime_stats().await;
    println!("📊 Runtime Stats:");
    println!("  • Active Agents: {}", stats.active_agents);
    println!("  • Total Memory: {}MB", stats.total_memory_mb);
    println!("  • CPU Usage: {:.1}%", stats.cpu_usage_percent);
    println!("  • Uptime: {}s", stats.uptime_seconds);
    
    // Keep the server running
    loop {
        tokio::time::sleep(tokio::time::Duration::from_secs(10)).await;
        
        // Update stats periodically
        let stats = runtime.get_runtime_stats().await;
        println!("🔄 Runtime Update: {} agents active, {}MB memory used", 
                stats.active_agents, stats.total_memory_mb);
    }
}
