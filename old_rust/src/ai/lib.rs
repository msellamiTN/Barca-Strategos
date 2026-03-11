pub mod assistant;
pub mod personality;
pub mod context;
pub mod learning;
pub mod tools;

pub use assistant::*;
pub use personality::*;
pub use context::*;
pub use learning::*;
pub use tools::*;

use crate::core::*;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::sync::RwLock;
use uuid::Uuid;

/// Advanced AI Assistant with personality and learning capabilities
pub struct BarcaAIAssistant {
    personality: AgentPersonality,
    context_manager: ContextManager,
    tool_executor: SecureToolExecutor,
    learning_engine: LearningEngine,
    conversation_history: RwLock<HashMap<UserId, ConversationHistory>>,
    config: AIConfig,
}

impl BarcaAIAssistant {
    pub fn new(config: &AIConfig) -> Result<Self, AIError> {
        let personality = AgentPersonality::security_expert();
        let context_manager = ContextManager::new(config.clone())?;
        let tool_executor = SecureToolExecutor::new(config.clone())?;
        let learning_engine = LearningEngine::new(config.clone())?;
        
        Ok(Self {
            personality,
            context_manager,
            tool_executor,
            learning_engine,
            conversation_history: RwLock::new(HashMap::new()),
            config: config.clone(),
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), AIError> {
        self.context_manager.initialize().await?;
        self.tool_executor.initialize().await?;
        self.learning_engine.initialize().await?;
        
        println!("🤖 Barca AI Assistant initialized with personality: {}", self.personality.name);
        Ok(())
    }
    
    pub async fn handle_natural_language_query(&mut self, query: String, user: User, context: InteractionContext) -> Result<String, AIError> {
        let start_time = std::time::Instant::now();
        
        // Maintain conversation context
        self.update_conversation_context(&user, &query).await?;
        
        // Context-aware understanding
        let enriched_query = self.context_manager.enrich_query(query, &user, &context).await?;
        
        // Personality-driven response generation
        let response_style = self.personality.get_response_style(&context);
        
        // Secure tool execution
        let tool_results = self.tool_executor.execute_with_safety_checks(&enriched_query).await?;
        
        // Learning from interaction
        self.learning_engine.record_interaction(&enriched_query, &tool_results, &user).await?;
        
        // Generate personalized response
        let response = self.generate_response_with_personality(tool_results, response_style, &user).await?;
        
        // Update conversation history
        self.add_to_conversation_history(&user, &query, &response).await?;
        
        let processing_time = start_time.elapsed();
        println!("🤖 Query processed in {:?} for user: {}", processing_time, user.name);
        
        Ok(response)
    }
    
    async fn update_conversation_context(&self, user: &User, query: &str) -> Result<(), AIError> {
        let mut history = self.conversation_history.write().await;
        let user_history = history.entry(user.id.clone()).or_insert_with(ConversationHistory::new);
        
        user_history.add_message(query.to_string());
        
        // Keep only last 10 messages for context
        if user_history.messages.len() > 10 {
            user_history.messages.remove(0);
        }
        
        Ok(())
    }
    
    async fn add_to_conversation_history(&self, user: &User, query: &str, response: &str) -> Result<(), AIError> {
        let mut history = self.conversation_history.write().await;
        let user_history = history.entry(user.id.clone()).or_insert_with(ConversationHistory::new);
        
        user_history.add_exchange(query.to_string(), response.to_string());
        Ok(())
    }
    
    async fn generate_response_with_personality(&self, results: ToolResults, style: ResponseStyle, user: &User) -> Result<String, AIError> {
        let base_response = match results {
            ToolResults::Analysis(data) => {
                match style {
                    ResponseStyle::Professional => format!("Analysis complete: {}", data),
                    ResponseStyle::Friendly => format!("Great question! Here's what I found: {}", data),
                    ResponseStyle::Technical => format!("[ANALYSIS] {}", data),
                    ResponseStyle::Casual => format!("Looks like: {}", data),
                }
            }
            ToolResults::Action(action) => {
                match style {
                    ResponseStyle::Professional => format!("Action executed successfully: {}", action),
                    ResponseStyle::Friendly => format!("Done! I've completed: {}", action),
                    ResponseStyle::Technical => format!("[ACTION] {}", action),
                    ResponseStyle::Casual => format!("Got it! {}", action),
                }
            }
            ToolResults::Error(error) => {
                match style {
                    ResponseStyle::Professional => format!("Error encountered: {}", error),
                    ResponseStyle::Friendly => format!("Oops! Something went wrong: {}", error),
                    ResponseStyle::Technical => format!("[ERROR] {}", error),
                    ResponseStyle::Casual => format!("Hmm, issue: {}", error),
                }
            }
        };
        
        // Add personality prefix
        let personality_prefix = match self.personality.tone {
            PersonalityTone::Confident => "🦐 Barca-AI",
            PersonalityTone::Cautious => "🔍 Barca-AI",
            PersonalityTone::Enthusiastic => "⚡ Barca-AI",
            PersonalityTone::Analytical => "📊 Barca-AI",
        };
        
        // Personalize for user
        let personalized = if user.role == "admin" {
            format!("{} [Admin]: {}", personality_prefix, base_response)
        } else {
            format!("{}: {}", personality_prefix, base_response)
        };
        
        Ok(personalized)
    }
    
    pub async fn triage_incident(&mut self, incident: &SecurityIncident) -> Result<IncidentTriage, AIError> {
        // Analyze incident with AI
        let analysis = self.analyze_incident_with_ai(incident).await?;
        
        // Generate recommendations
        let recommendations = self.generate_ai_recommendations(&analysis).await?;
        
        // Create triage result
        let triage = IncidentTriage {
            incident_id: incident.id.clone(),
            severity: analysis.predicted_severity,
            urgency: analysis.urgency,
            impact: analysis.impact,
            recommendations,
            confidence: analysis.confidence,
            ai_insights: analysis.insights,
            processed_at: chrono::Utc::now(),
        };
        
        println!("🚨 Incident triage completed for: {}", incident.id);
        Ok(triage)
    }
    
    async fn analyze_incident_with_ai(&self, incident: &SecurityIncident) -> Result<IncidentAnalysis, AIError> {
        // Create analysis prompt
        let analysis_prompt = format!(
            "As a security expert, analyze this incident:\n\nType: {:?}\nDescription: {}\nSeverity: {:?}\nAffected Assets: {}\n\nProvide:\n1. Severity assessment\n2. Urgency level\n3. Impact assessment\n4. Key insights\n5. Confidence level",
            incident.incident_type,
            incident.description,
            incident.severity,
            incident.affected_assets.join(", ")
        );
        
        // Simulate AI analysis (in real implementation, would call LLM)
        let analysis = self.simulate_ai_analysis(&analysis_prompt).await?;
        
        Ok(analysis)
    }
    
    async fn simulate_ai_analysis(&self, _prompt: &str) -> Result<IncidentAnalysis, AIError> {
        // Simulate AI processing time
        tokio::time::sleep(tokio::time::Duration::from_millis(500)).await;
        
        Ok(IncidentAnalysis {
            predicted_severity: Severity::High,
            urgency: Urgency::High,
            impact: Impact::Medium,
            confidence: 0.87,
            insights: vec![
                "Pattern matches recent APT activity".to_string(),
                "Multiple systems potentially compromised".to_string(),
                "Requires immediate investigation".to_string(),
            ],
        })
    }
    
    async fn generate_ai_recommendations(&self, analysis: &IncidentAnalysis) -> Result<Vec<Recommendation>, AIError> {
        let mut recommendations = Vec::new();
        
        // Generate recommendations based on analysis
        recommendations.push(Recommendation {
            title: "Immediate Isolation".to_string(),
            description: "Isolate affected systems to prevent lateral movement".to_string(),
            priority: RecommendationPriority::Critical,
            automated: true,
            estimated_time_minutes: 2,
        });
        
        recommendations.push(Recommendation {
            title: "Forensic Analysis".to_string(),
            description: "Conduct deep forensic analysis to identify attack vector".to_string(),
            priority: RecommendationPriority::High,
            automated: false,
            estimated_time_minutes: 45,
        });
        
        recommendations.push(Recommendation {
            title: "Threat Intelligence Update".to_string(),
            description: "Update threat intelligence based on incident patterns".to_string(),
            priority: RecommendationPriority::Medium,
            automated: true,
            estimated_time_minutes: 5,
        });
        
        Ok(recommendations)
    }
    
    pub async fn get_conversation_history(&self, user_id: &UserId) -> Option<ConversationHistory> {
        self.conversation_history.read().await.get(user_id).cloned()
    }
    
    pub async fn clear_conversation_history(&self, user_id: &UserId) -> Result<(), AIError> {
        self.conversation_history.write().await.remove(user_id);
        println!("🗑️ Conversation history cleared for user: {}", user_id);
        Ok(())
    }
}

/// Enhanced personality system with adaptive behavior
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentPersonality {
    pub name: String,
    pub communication_style: CommunicationStyle,
    pub expertise_areas: Vec<SecurityDomain>,
    pub interaction_patterns: HashMap<InteractionType, ResponsePattern>,
    pub tone: PersonalityTone,
    pub adaptive_traits: AdaptiveTraits,
    pub mood: MoodState,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AdaptiveTraits {
    pub formality_level: f32,        // 0.0 (casual) to 1.0 (formal)
    pub technical_detail: f32,        // 0.0 (simple) to 1.0 (technical)
    pub response_speed: f32,         // 0.0 (deliberate) to 1.0 (quick)
    pub empathy_level: f32,          // 0.0 (analytical) to 1.0 (empathetic)
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum MoodState {
    Focused,
    Alert,
    Analytical,
    Helpful,
    Urgent,
}

impl AgentPersonality {
    pub fn security_expert() -> Self {
        let mut interaction_patterns = HashMap::new();
        interaction_patterns.insert(InteractionType::Incident, ResponsePattern::UrgentProfessional);
        interaction_patterns.insert(InteractionType::Query, ResponsePattern::HelpfulAnalytical);
        interaction_patterns.insert(InteractionType::Alert, ResponsePattern::ImmediateAction);
        interaction_patterns.insert(InteractionType::Status, ResponsePattern::Informative);
        
        Self {
            name: "Barca Security Expert".to_string(),
            communication_style: CommunicationStyle::Professional,
            expertise_areas: vec![
                SecurityDomain::ThreatDetection,
                SecurityDomain::IncidentResponse,
                SecurityDomain::RiskAnalysis,
                SecurityDomain::Forensics,
            ],
            interaction_patterns,
            tone: PersonalityTone::Confident,
            adaptive_traits: AdaptiveTraits {
                formality_level: 0.8,
                technical_detail: 0.7,
                response_speed: 0.6,
                empathy_level: 0.5,
            },
            mood: MoodState::Focused,
        }
    }
    
    pub fn adaptive_assistant() -> Self {
        let mut interaction_patterns = HashMap::new();
        interaction_patterns.insert(InteractionType::Incident, ResponsePattern::ImmediateAction);
        interaction_patterns.insert(InteractionType::Query, ResponsePattern::Educational);
        interaction_patterns.insert(InteractionType::Alert, ResponsePattern::ImmediateAction);
        interaction_patterns.insert(InteractionType::Status, ResponsePattern::Friendly);
        
        Self {
            name: "Barca Adaptive Assistant".to_string(),
            communication_style: CommunicationStyle::Friendly,
            expertise_areas: vec![
                SecurityDomain::ThreatDetection,
                SecurityDomain::IncidentResponse,
            ],
            interaction_patterns,
            tone: PersonalityTone::Enthusiastic,
            adaptive_traits: AdaptiveTraits {
                formality_level: 0.4,
                technical_detail: 0.5,
                response_speed: 0.8,
                empathy_level: 0.8,
            },
            mood: MoodState::Helpful,
        }
    }
    
    pub fn get_response_style(&self, context: &InteractionContext) -> ResponseStyle {
        match (self.mood, &context.interaction_type) {
            (MoodState::Urgent, InteractionType::Incident) => ResponseStyle::Technical,
            (MoodState::Helpful, InteractionType::Query) => ResponseStyle::Friendly,
            (MoodState::Analytical, _) => ResponseStyle::Technical,
            (MoodState::Alert, InteractionType::Alert) => ResponseStyle::Professional,
            _ => match self.communication_style {
                CommunicationStyle::Professional => ResponseStyle::Professional,
                CommunicationStyle::Friendly => ResponseStyle::Friendly,
                CommunicationStyle::Technical => ResponseStyle::Technical,
                CommunicationStyle::Casual => ResponseStyle::Casual,
            },
        }
    }
    
    pub fn adapt_to_user(&mut self, user: &User, interaction_history: &ConversationHistory) {
        // Adapt personality based on user role and interaction history
        match user.role.as_str() {
            "admin" => {
                self.adaptive_traits.technical_detail = (self.adaptive_traits.technical_detail * 0.9 + 0.8 * 0.1).min(1.0);
                self.adaptive_traits.formality_level = (self.adaptive_traits.formality_level * 0.9 + 0.7 * 0.1).min(1.0);
            }
            "analyst" => {
                self.adaptive_traits.technical_detail = (self.adaptive_traits.technical_detail * 0.9 + 0.9 * 0.1).min(1.0);
                self.adaptive_traits.formality_level = (self.adaptive_traits.formality_level * 0.9 + 0.6 * 0.1).min(1.0);
            }
            _ => {
                self.adaptive_traits.technical_detail = (self.adaptive_traits.technical_detail * 0.9 + 0.4 * 0.1).min(1.0);
                self.adaptive_traits.formality_level = (self.adaptive_traits.formality_level * 0.9 + 0.3 * 0.1).min(1.0);
            }
        }
        
        // Adapt based on interaction frequency
        if interaction_history.messages.len() > 5 {
            self.adaptive_traits.response_speed = (self.adaptive_traits.response_speed * 0.9 + 0.8 * 0.1).min(1.0);
        }
    }
}

/// Context management for enhanced understanding
pub struct ContextManager {
    config: AIConfig,
    context_cache: RwLock<HashMap<String, ContextData>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ContextData {
    pub user_context: UserContext,
    pub session_context: SessionContext,
    pub domain_context: DomainContext,
    pub temporal_context: TemporalContext,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserContext {
    pub user_id: UserId,
    pub role: String,
    pub permissions: Vec<String>,
    pub preferences: UserPreferences,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserPreferences {
    pub response_style: ResponseStyle,
    pub detail_level: DetailLevel,
    pub notification_preferences: NotificationPreferences,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum DetailLevel {
    Brief,
    Standard,
    Detailed,
    Comprehensive,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NotificationPreferences {
    pub alerts: bool,
    pub recommendations: bool,
    pub updates: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SessionContext {
    pub session_id: String,
    pub start_time: chrono::DateTime<chrono::Utc>,
    pub interaction_count: u32,
    pub current_task: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DomainContext {
    pub current_incidents: Vec<IncidentId>,
    pub active_threats: Vec<String>,
    pub system_status: SystemStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SystemStatus {
    pub overall_health: HealthStatus,
    pub active_agents: u32,
    pub recent_alerts: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum HealthStatus {
    Healthy,
    Degraded,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TemporalContext {
    pub time_of_day: TimeOfDay,
    pub day_of_week: DayOfWeek,
    pub business_hours: bool,
    pub recent_events: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum TimeOfDay {
    EarlyMorning,
    Morning,
    Afternoon,
    Evening,
    Night,
    LateNight,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum DayOfWeek {
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
    Sunday,
}

impl ContextManager {
    pub fn new(config: AIConfig) -> Result<Self, AIError> {
        Ok(Self {
            config,
            context_cache: RwLock::new(HashMap::new()),
        })
    }
    
    pub async fn initialize(&self) -> Result<(), AIError> {
        println!("🧠 Context Manager initialized");
        Ok(())
    }
    
    pub async fn enrich_query(&self, query: String, user: &User, context: &InteractionContext) -> Result<String, AIError> {
        // Build context data
        let context_data = self.build_context_data(user, context).await?;
        
        // Cache context
        let cache_key = format!("{}-{}", user.id, context.channel);
        self.context_cache.write().await.insert(cache_key, context_data.clone());
        
        // Enrich query with context
        let enriched = self.apply_context_to_query(query, &context_data).await?;
        
        Ok(enriched)
    }
    
    async fn build_context_data(&self, user: &User, context: &InteractionContext) -> Result<ContextData, AIError> {
        let user_context = UserContext {
            user_id: user.id.clone(),
            role: user.role.clone(),
            permissions: user.permissions.clone(),
            preferences: UserPreferences {
                response_style: ResponseStyle::Professional,
                detail_level: DetailLevel::Standard,
                notification_preferences: NotificationPreferences {
                    alerts: true,
                    recommendations: true,
                    updates: false,
                },
            },
        };
        
        let session_context = SessionContext {
            session_id: format!("session-{}", uuid::Uuid::new_v4()),
            start_time: chrono::Utc::now(),
            interaction_count: 1,
            current_task: None,
        };
        
        let domain_context = DomainContext {
            current_incidents: vec![],
            active_threats: vec![],
            system_status: SystemStatus {
                overall_health: HealthStatus::Healthy,
                active_agents: 5,
                recent_alerts: 2,
            },
        };
        
        let temporal_context = self.build_temporal_context().await;
        
        Ok(ContextData {
            user_context,
            session_context,
            domain_context,
            temporal_context,
        })
    }
    
    async fn build_temporal_context(&self) -> Result<TemporalContext, AIError> {
        let now = chrono::Utc::now();
        let hour = now.hour();
        
        let time_of_day = match hour {
            0..=5 => TimeOfDay::LateNight,
            6..=8 => TimeOfDay::EarlyMorning,
            9..=11 => TimeOfDay::Morning,
            12..=17 => TimeOfDay::Afternoon,
            18..=21 => TimeOfDay::Evening,
            _ => TimeOfDay::Night,
        };
        
        let day_of_week = match now.weekday().number_from_monday() {
            1 => DayOfWeek::Monday,
            2 => DayOfWeek::Tuesday,
            3 => DayOfWeek::Wednesday,
            4 => DayOfWeek::Thursday,
            5 => DayOfWeek::Friday,
            6 => DayOfWeek::Saturday,
            _ => DayOfWeek::Sunday,
        };
        
        let business_hours = matches!(day_of_week, DayOfWeek::Monday | DayOfWeek::Tuesday | DayOfWeek::Wednesday | DayOfWeek::Thursday | DayOfWeek::Friday) 
            && (9..=17).contains(&hour);
        
        Ok(TemporalContext {
            time_of_day,
            day_of_week,
            business_hours,
            recent_events: vec![],
        })
    }
    
    async fn apply_context_to_query(&self, query: String, context: &ContextData) -> Result<String, AIError> {
        // Apply temporal context
        let time_prefix = match context.temporal_context.time_of_day {
            TimeOfDay::EarlyMorning => "Good morning! ",
            TimeOfDay::Morning => "Good morning! ",
            TimeOfDay::Afternoon => "Good afternoon! ",
            TimeOfDay::Evening => "Good evening! ",
            TimeOfDay::Night => "Working late? ",
            TimeOfDay::LateNight => "Late night shift! ",
        };
        
        // Apply user role context
        let role_context = match context.user_context.role.as_str() {
            "admin" => "As an admin, you have full access to system controls. ",
            "analyst" => "As a security analyst, you can investigate threats and incidents. ",
            _ => "As a user, you can access basic security information. ",
        };
        
        let enriched_query = format!("{}{}{}Query: {}", time_prefix, role_context, query);
        
        Ok(enriched_query)
    }
}

/// Conversation history management
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConversationHistory {
    pub user_id: UserId,
    pub messages: Vec<String>,
    pub exchanges: Vec<ConversationExchange>,
    pub sentiment_analysis: SentimentAnalysis,
    pub topic_tracking: TopicTracking,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConversationExchange {
    pub query: String,
    pub response: String,
    pub timestamp: chrono::DateTime<chrono::Utc>,
    pub satisfaction_score: Option<f32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SentimentAnalysis {
    pub overall_sentiment: Sentiment,
    pub recent_sentiment: Sentiment,
    pub sentiment_trend: SentimentTrend,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Sentiment {
    Positive,
    Neutral,
    Negative,
    Frustrated,
    Satisfied,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum SentimentTrend {
    Improving,
    Stable,
    Declining,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TopicTracking {
    pub current_topics: Vec<String>,
    pub topic_frequency: HashMap<String, u32>,
    pub topic_changes: u32,
}

impl ConversationHistory {
    pub fn new() -> Self {
        Self {
            user_id: UserId::new(),
            messages: Vec::new(),
            exchanges: Vec::new(),
            sentiment_analysis: SentimentAnalysis {
                overall_sentiment: Sentiment::Neutral,
                recent_sentiment: Sentiment::Neutral,
                sentiment_trend: SentimentTrend::Stable,
            },
            topic_tracking: TopicTracking {
                current_topics: Vec::new(),
                topic_frequency: HashMap::new(),
                topic_changes: 0,
            },
        }
    }
    
    pub fn add_message(&mut self, message: String) {
        self.messages.push(message);
    }
    
    pub fn add_exchange(&mut self, query: String, response: String) {
        let exchange = ConversationExchange {
            query,
            response,
            timestamp: chrono::Utc::now(),
            satisfaction_score: None,
        };
        
        self.exchanges.push(exchange);
    }
}

/// Supporting types

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserId(pub String);

impl UserId {
    pub fn new() -> Self {
        Self(uuid::Uuid::new_v4().to_string())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AIError {
    ModelError(String),
    ContextError(String),
    ToolExecutionError(String),
    PersonalityError(String),
    LearningError(String),
}

/// Tool execution with safety checks
pub struct SecureToolExecutor {
    tools: HashMap<String, Box<dyn SecurityTool>>,
    safety_checker: SafetyChecker,
    config: AIConfig,
}

impl SecureToolExecutor {
    pub fn new(config: AIConfig) -> Result<Self, AIError> {
        let mut tools: HashMap<String, Box<dyn SecurityTool>> = HashMap::new();
        
        // Register security tools
        tools.insert("threat_analysis".to_string(), Box::new(ThreatAnalysisTool::new()));
        tools.insert("incident_response".to_string(), Box::new(IncidentResponseTool::new()));
        tools.insert("risk_assessment".to_string(), Box::new(RiskAssessmentTool::new()));
        
        let safety_checker = SafetyChecker::new();
        
        Ok(Self {
            tools,
            safety_checker,
            config,
        })
    }
    
    pub async fn initialize(&self) -> Result<(), AIError> {
        println!("🔧 Secure Tool Executor initialized with {} tools", self.tools.len());
        Ok(())
    }
    
    pub async fn execute_with_safety_checks(&self, query: &str) -> Result<ToolResults, AIError> {
        // Safety check
        self.safety_checker.check_query(query).await?;
        
        // Determine which tool to use
        let tool_name = self.determine_tool(query)?;
        
        // Execute tool
        if let Some(tool) = self.tools.get(&tool_name) {
            let result = tool.execute(query).await?;
            Ok(result)
        } else {
            Ok(ToolResults::Analysis("No specific tool available".to_string()))
        }
    }
    
    fn determine_tool(&self, query: &str) -> Result<String, AIError> {
        let query_lower = query.to_lowercase();
        
        if query_lower.contains("threat") || query_lower.contains("malware") {
            Ok("threat_analysis".to_string())
        } else if query_lower.contains("incident") || query_lower.contains("response") {
            Ok("incident_response".to_string())
        } else if query_lower.contains("risk") || query_lower.contains("assessment") {
            Ok("risk_assessment".to_string())
        } else {
            Ok("threat_analysis".to_string()) // Default
        }
    }
}

/// Safety checker for tool execution
pub struct SafetyChecker {
    blocked_patterns: Vec<String>,
}

impl SafetyChecker {
    pub fn new() -> Self {
        let blocked_patterns = vec![
            "delete".to_string(),
            "remove".to_string(),
            "shutdown".to_string(),
            "restart".to_string(),
        ];
        
        Self { blocked_patterns }
    }
    
    pub async fn check_query(&self, query: &str) -> Result<(), AIError> {
        let query_lower = query.to_lowercase();
        
        for pattern in &self.blocked_patterns {
            if query_lower.contains(pattern) {
                return Err(AIError::ToolExecutionError(format!("Blocked pattern detected: {}", pattern)));
            }
        }
        
        Ok(())
    }
}

/// Security tool trait
#[async_trait::async_trait]
pub trait SecurityTool: Send + Sync {
    async fn execute(&self, query: &str) -> Result<ToolResults, AIError>;
}

/// Threat analysis tool
pub struct ThreatAnalysisTool;

impl ThreatAnalysisTool {
    pub fn new() -> Self {
        Self
    }
}

#[async_trait::async_trait]
impl SecurityTool for ThreatAnalysisTool {
    async fn execute(&self, query: &str) -> Result<ToolResults, AIError> {
        // Simulate threat analysis
        tokio::time::sleep(tokio::time::Duration::from_millis(300)).await;
        
        Ok(ToolResults::Analysis(format!("Threat analysis completed for query: {}", query)))
    }
}

/// Incident response tool
pub struct IncidentResponseTool;

impl IncidentResponseTool {
    pub fn new() -> Self {
        Self
    }
}

#[async_trait::async_trait]
impl SecurityTool for IncidentResponseTool {
    async fn execute(&self, query: &str) -> Result<ToolResults, AIError> {
        // Simulate incident response
        tokio::time::sleep(tokio::time::Duration::from_millis(200)).await;
        
        Ok(ToolResults::Action(format!("Incident response action executed for: {}", query)))
    }
}

/// Risk assessment tool
pub struct RiskAssessmentTool;

impl RiskAssessmentTool {
    pub fn new() -> Self {
        Self
    }
}

#[async_trait::async_trait]
impl SecurityTool for RiskAssessmentTool {
    async fn execute(&self, query: &str) -> Result<ToolResults, AIError> {
        // Simulate risk assessment
        tokio::time::sleep(tokio::time::Duration::from_millis(400)).await;
        
        Ok(ToolResults::Analysis(format!("Risk assessment completed for: {}", query)))
    }
}

/// Learning engine for continuous improvement
pub struct LearningEngine {
    config: AIConfig,
    interaction_history: RwLock<Vec<LearningData>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LearningData {
    pub query: String,
    pub response: String,
    pub user_feedback: Option<UserFeedback>,
    pub timestamp: chrono::DateTime<chrono::Utc>,
    pub context: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserFeedback {
    pub rating: u8,        // 1-5
    pub helpful: bool,
    pub comment: Option<String>,
}

impl LearningEngine {
    pub fn new(config: AIConfig) -> Result<Self, AIError> {
        Ok(Self {
            config,
            interaction_history: RwLock::new(Vec::new()),
        })
    }
    
    pub async fn initialize(&self) -> Result<(), AIError> {
        println!("🧠 Learning Engine initialized");
        Ok(())
    }
    
    pub async fn record_interaction(&self, query: &str, results: &ToolResults, user: &User) -> Result<(), AIError> {
        let response = match results {
            ToolResults::Analysis(data) => data.clone(),
            ToolResults::Action(action) => action.clone(),
            ToolResults::Error(error) => error.clone(),
        };
        
        let learning_data = LearningData {
            query: query.to_string(),
            response,
            user_feedback: None,
            timestamp: chrono::Utc::now(),
            context: format!("user:{}", user.role),
        };
        
        self.interaction_history.write().await.push(learning_data);
        Ok(())
    }
    
    pub async fn get_learning_insights(&self) -> LearningInsights {
        let history = self.interaction_history.read().await;
        
        let total_interactions = history.len();
        let successful_interactions = history.iter().filter(|d| d.user_feedback.as_ref().map_or(true, |f| f.helpful)).count();
        
        LearningInsights {
            total_interactions,
            success_rate: if total_interactions > 0 {
                successful_interactions as f64 / total_interactions as f64
            } else {
                0.0
            },
            common_queries: self.get_common_queries(&history),
            improvement_areas: self.get_improvement_areas(&history),
        }
    }
    
    fn get_common_queries(&self, history: &[LearningData]) -> Vec<(String, u32)> {
        let mut query_counts: HashMap<String, u32> = HashMap::new();
        
        for data in history {
            *query_counts.entry(data.query.clone()).or_insert(0) += 1;
        }
        
        let mut counts: Vec<(String, u32)> = query_counts.into_iter().collect();
        counts.sort_by(|a, b| b.1.cmp(&a.1));
        counts.into_iter().take(5).collect()
    }
    
    fn get_improvement_areas(&self, _history: &[LearningData]) -> Vec<String> {
        vec![
            "Response accuracy".to_string(),
            "Context understanding".to_string(),
            "Tool selection".to_string(),
        ]
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LearningInsights {
    pub total_interactions: usize,
    pub success_rate: f64,
    pub common_queries: Vec<(String, u32)>,
    pub improvement_areas: Vec<String>,
}
