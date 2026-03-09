pub mod quest_system;
pub mod achievements;
pub mod leaderboard;
pub mod rewards;
pub mod challenges;

pub use quest_system::*;
pub use achievements::*;
pub use leaderboard::*;
pub use rewards::*;
pub use challenges::*;

use crate::core::*;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::sync::RwLock;
use uuid::Uuid;

/// Gamified security operations system
pub struct SecurityQuestEngine {
    quest_generator: AIQuestGenerator,
    narrative_engine: NarrativeEngine,
    reward_system: TokenizedRewardSystem,
    achievement_tracker: AchievementTracker,
    leaderboard: LeaderboardManager,
    player_progress: RwLock<HashMap<UserId, PlayerProgress>>,
    config: GamificationConfig,
}

impl SecurityQuestEngine {
    pub fn new(config: GamificationConfig) -> Result<Self, GamificationError> {
        let quest_generator = AIQuestGenerator::new(&config)?;
        let narrative_engine = NarrativeEngine::new(&config)?;
        let reward_system = TokenizedRewardSystem::new(&config)?;
        let achievement_tracker = AchievementTracker::new(&config)?;
        let leaderboard = LeaderboardManager::new(&config)?;
        
        Ok(Self {
            quest_generator,
            narrative_engine,
            reward_system,
            achievement_tracker,
            leaderboard,
            player_progress: RwLock::new(HashMap::new()),
            config,
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        self.quest_generator.initialize().await?;
        self.narrative_engine.initialize().await?;
        self.reward_system.initialize().await?;
        self.achievement_tracker.initialize().await?;
        self.leaderboard.initialize().await?;
        
        println!("🎮 Security Quest Engine initialized - Gamification ready!");
        Ok(())
    }
    
    pub async fn generate_personalized_quest(&self, analyst: SecurityAnalyst) -> Result<SecurityQuest, GamificationError> {
        // Analyze analyst skills and current threats
        let skill_profile = self.analyze_analyst_skills(&analyst).await?;
        let threat_landscape = self.get_current_threat_landscape().await?;
        
        // Generate personalized quest with AI narration
        let quest = self.quest_generator.create_quest(skill_profile, threat_landscape).await?;
        let narrative = self.narrative_engine.create_quest_narrative(&quest, &analyst).await?;
        
        let personalized_quest = SecurityQuest {
            quest,
            narrative,
            personal_difficulty: self.calculate_personal_difficulty(&skill_profile, &quest),
            adaptive_hints: self.generate_adaptive_hints(&skill_profile, &quest),
        };
        
        println!("🎯 Generated quest for {}: {}", analyst.name, personalized_quest.quest.title);
        Ok(personalized_quest)
    }
    
    pub async fn create_daily_quests(&self) -> Result<Vec<SecurityQuest>, GamificationError> {
        let mut quests = Vec::new();
        
        // Threat hunting quest
        quests.push(SecurityQuest {
            quest: Quest {
                id: QuestId::new(),
                title: "Threat Hunter: Find the Hidden IOC".to_string(),
                description: "Analyze today's logs and identify at least 3 new indicators of compromise".to_string(),
                quest_type: QuestType::ThreatHunting,
                difficulty: QuestDifficulty::Medium,
                requirements: vec![
                    QuestRequirement::AnalyzeLogs,
                    QuestRequirement::IdentifyIOC,
                    QuestRequirement::CreateAlert,
                ],
                rewards: vec![
                    Reward::ExperiencePoints(100),
                    Reward::Badge("Threat Hunter".to_string()),
                    Reward::Tokens(50),
                ],
                time_limit: Some(chrono::Duration::hours(8)),
                created_at: chrono::Utc::now(),
            },
            narrative: QuestNarrative {
                title: "The Digital Shadow".to_string(),
                story: "A mysterious threat lurks in the network shadows. As a elite threat hunter, you must uncover the hidden indicators of compromise before they strike again. The clock is ticking...".to_string(),
                objectives: vec![
                    "Scan network logs for suspicious patterns".to_string(),
                    "Identify at least 3 unique IOCs".to_string(),
                    "Create detailed alerts for the security team".to_string(),
                ],
            },
            personal_difficulty: PersonalDifficulty::Adaptive,
            adaptive_hints: vec![
                "Check for unusual outbound connections".to_string(),
                "Look for fileless malware indicators".to_string(),
                "Analyze timing patterns in log data".to_string(),
            ],
        });
        
        // Incident response quest
        quests.push(SecurityQuest {
            quest: Quest {
                id: QuestId::new(),
                title: "Incident Responder: Critical Containment".to_string(),
                description: "Respond to and contain a simulated security incident within 30 minutes".to_string(),
                quest_type: QuestType::IncidentResponse,
                difficulty: QuestDifficulty::Hard,
                requirements: vec![
                    QuestRequirement::AnalyzeIncident,
                    QuestRequirement::ContainThreat,
                    QuestRequirement::DocumentResponse,
                ],
                rewards: vec![
                    Reward::ExperiencePoints(200),
                    Reward::Badge("Incident Master".to_string()),
                    Reward::Tokens(100),
                ],
                time_limit: Some(chrono::Duration::minutes(30)),
                created_at: chrono::Utc::now(),
            },
            narrative: QuestNarrative {
                title: "Race Against Time".to_string(),
                story: "ALERT: Critical security incident detected! The network is under attack and every second counts. Your mission: contain the threat before it spreads. The fate of the digital kingdom rests in your hands...".to_string(),
                objectives: vec![
                    "Triage the incident within 5 minutes".to_string(),
                    "Isolate affected systems".to_string(),
                    "Document your response strategy".to_string(),
                ],
            },
            personal_difficulty: PersonalDifficulty::Adaptive,
            adaptive_hints: vec![
                "Prioritize critical assets first".to_string(),
                    "Use the incident response playbook".to_string(),
                    "Document every action for post-mortem".to_string(),
            ],
        });
        
        // Risk assessment quest
        quests.push(SecurityQuest {
            quest: Quest {
                id: QuestId::new(),
                title: "Risk Analyst: Vulnerability Discovery".to_string(),
                description: "Identify and assess at least 5 new vulnerabilities in the infrastructure".to_string(),
                quest_type: QuestType::RiskAssessment,
                difficulty: QuestDifficulty::Medium,
                requirements: vec![
                    QuestRequirement::ScanVulnerabilities,
                    QuestRequirement::AssessRisk,
                    QuestRequirement::RecommendMitigation,
                ],
                rewards: vec![
                    Reward::ExperiencePoints(150),
                    Reward::Badge("Risk Expert".to_string()),
                    Reward::Tokens(75),
                ],
                time_limit: Some(chrono::Duration::hours(6)),
                created_at: chrono::Utc::now(),
            },
            narrative: QuestNarrative {
                title: "The Hidden Weaknesses".to_string(),
                description: "Like a master locksmith, you must find the hidden weaknesses in our digital fortress. Each vulnerability discovered makes us stronger. Your analytical skills will be our shield...".to_string(),
                objectives: vec![
                    "Perform comprehensive vulnerability scan".to_string(),
                    "Assess risk levels for each finding".to_string(),
                    "Prioritize remediation efforts".to_string(),
                ],
            },
            personal_difficulty: PersonalDifficulty::Adaptive,
            adaptive_hints: vec![
                "Focus on critical infrastructure first".to_string(),
                "Consider both technical and business impact".to_string(),
                    "Use CVSS scoring for consistency".to_string(),
            ],
        });
        
        println!("🎯 Generated {} daily quests", quests.len());
        Ok(quests)
    }
    
    async fn analyze_analyst_skills(&self, analyst: &SecurityAnalyst) -> Result<SkillProfile, GamificationError> {
        let mut progress = self.player_progress.write().await;
        let player_progress = progress.entry(analyst.id.clone()).or_insert_with(PlayerProgress::new);
        
        // Calculate skill levels based on completed quests and achievements
        let threat_hunting_skill = player_progress.calculate_skill_level(SkillType::ThreatHunting);
        let incident_response_skill = player_progress.calculate_skill_level(SkillType::IncidentResponse);
        let risk_assessment_skill = player_progress.calculate_skill_level(SkillType::RiskAssessment);
        let communication_skill = player_progress.calculate_skill_level(SkillType::Communication);
        
        Ok(SkillProfile {
            analyst_id: analyst.id.clone(),
            threat_hunting: threat_hunting_skill,
            incident_response: incident_response_skill,
            risk_assessment: risk_assessment_skill,
            communication: communication_skill,
            overall_level: player_progress.overall_level,
            preferred_difficulty: player_progress.preferred_difficulty,
        })
    }
    
    async fn get_current_threat_landscape(&self) -> Result<ThreatLandscape, GamificationError> {
        // Simulate getting current threat intelligence
        Ok(ThreatLandscape {
            active_threats: vec![
                Threat {
                    id: "APT-001".to_string(),
                    name: "Advanced Persistent Threat Group".to_string(),
                    severity: ThreatSeverity::High,
                    category: ThreatCategory::Espionage,
                    tactics: vec!["Initial Access".to_string(), "Lateral Movement".to_string()],
                },
                Threat {
                    id: "MALWARE-042".to_string(),
                    name: "Ransomware Variant".to_string(),
                    severity: ThreatSeverity::Critical,
                    category: ThreatCategory::Ransomware,
                    tactics: vec!["Encryption".to_string(), "Exfiltration".to_string()],
                },
            ],
            vulnerability_trends: vec![
                "Web application vulnerabilities".to_string(),
                "Misconfigured cloud services".to_string(),
            ],
            industry_specific_threats: vec![
                "Financial sector targeting".to_string(),
            ],
        })
    }
    
    fn calculate_personal_difficulty(&self, skill_profile: &SkillProfile, quest: &Quest) -> PersonalDifficulty {
        let base_difficulty = quest.difficulty.clone();
        
        // Adjust based on analyst skill level
        match skill_profile.overall_level {
            Level::Beginner => PersonalDifficulty::Easier,
            Level::Intermediate => PersonalDifficulty::Balanced,
            Level::Advanced => PersonalDifficulty::Harder,
            Level::Expert => PersonalDifficulty::Expert,
        }
    }
    
    fn generate_adaptive_hints(&self, skill_profile: &SkillProfile, quest: &Quest) -> Vec<String> {
        let mut hints = Vec::new();
        
        match quest.quest_type {
            QuestType::ThreatHunting => {
                if skill_profile.threat_hunting.level < 3 {
                    hints.push("Start with known IOC patterns".to_string());
                    hints.push("Use timeline analysis techniques".to_string());
                }
            }
            QuestType::IncidentResponse => {
                if skill_profile.incident_response.level < 3 {
                    hints.push("Follow the incident response playbook".to_string());
                    hints.push("Prioritize containment over eradication".to_string());
                }
            }
            QuestType::RiskAssessment => {
                if skill_profile.risk_assessment.level < 3 {
                    hints.push("Use CVSS scoring framework".to_string());
                    hints.push("Consider business impact in assessment".to_string());
                }
            }
        }
        
        hints
    }
    
    pub async fn complete_quest(&self, analyst_id: &UserId, quest_id: &QuestId, completion_data: QuestCompletion) -> Result<QuestResult, GamificationError> {
        let mut progress = self.player_progress.write().await;
        let player_progress = progress.entry(analyst_id.clone()).or_insert_with(PlayerProgress::new);
        
        // Update quest progress
        let quest_result = player_progress.complete_quest(quest_id, completion_data).await?;
        
        // Check for achievements
        let new_achievements = self.achievement_tracker.check_achievements(&player_progress).await?;
        
        // Update leaderboard
        self.leaderboard.update_score(analyst_id, quest_result.experience_gained).await?;
        
        // Process rewards
        for reward in &quest_result.rewards {
            self.reward_system.grant_reward(analyst_id, reward).await?;
        }
        
        println!("🎉 Quest completed! {} gained {} XP", analyst_id.0, quest_result.experience_gained);
        
        Ok(quest_result)
    }
    
    pub async fn get_leaderboard(&self, time_period: LeaderboardPeriod) -> Result<Vec<LeaderboardEntry>, GamificationError> {
        self.leaderboard.get_leaderboard(time_period).await
    }
    
    pub async fn get_player_progress(&self, analyst_id: &UserId) -> Option<PlayerProgress> {
        self.player_progress.read().await.get(analyst_id).cloned()
    }
}

/// AI-powered quest generation
pub struct AIQuestGenerator {
    config: GamificationConfig,
    quest_templates: Vec<QuestTemplate>,
}

impl AIQuestGenerator {
    pub fn new(config: &GamificationConfig) -> Result<Self, GamificationError> {
        let quest_templates = vec![
            QuestTemplate {
                name: "Threat Investigation".to_string(),
                base_difficulty: QuestDifficulty::Medium,
                quest_type: QuestType::ThreatHunting,
                requirements: vec![QuestRequirement::AnalyzeLogs, QuestRequirement::IdentifyIOC],
                base_rewards: vec![Reward::ExperiencePoints(100), Reward::Tokens(50)],
            },
            QuestTemplate {
                name: "Incident Response".to_string(),
                base_difficulty: QuestDifficulty::Hard,
                quest_type: QuestType::IncidentResponse,
                requirements: vec![QuestRequirement::AnalyzeIncident, QuestRequirement::ContainThreat],
                base_rewards: vec![Reward::ExperiencePoints(200), Reward::Tokens(100)],
            },
        ];
        
        Ok(Self {
            config: config.clone(),
            quest_templates,
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        println!("🎯 AI Quest Generator initialized");
        Ok(())
    }
    
    pub async fn create_quest(&self, skill_profile: SkillProfile, threat_landscape: ThreatLandscape) -> Result<Quest, GamificationError> {
        // Select appropriate template based on skills and threats
        let template = self.select_quest_template(&skill_profile, &threat_landscape)?;
        
        // Generate quest with AI
        let quest = Quest {
            id: QuestId::new(),
            title: self.generate_quest_title(&template, &threat_landscape).await?,
            description: self.generate_quest_description(&template, &skill_profile, &threat_landscape).await?,
            quest_type: template.quest_type,
            difficulty: self.adjust_difficulty(&template.base_difficulty, &skill_profile),
            requirements: template.requirements,
            rewards: self.adjust_rewards(&template.base_rewards, &skill_profile),
            time_limit: Some(self.calculate_time_limit(&template.quest_type, &skill_profile)),
            created_at: chrono::Utc::now(),
        };
        
        Ok(quest)
    }
    
    fn select_quest_template(&self, skill_profile: &SkillProfile, threat_landscape: &ThreatLandscape) -> Result<&QuestTemplate, GamificationError> {
        // Simple selection logic - in production, use more sophisticated AI
        if skill_profile.threat_hunting.level > 3 && !threat_landscape.active_threats.is_empty() {
            Ok(&self.quest_templates[0]) // Threat hunting
        } else {
            Ok(&self.quest_templates[1]) // Incident response
        }
    }
    
    async fn generate_quest_title(&self, template: &QuestTemplate, threat_landscape: &ThreatLandscape) -> Result<String, GamificationError> {
        match template.quest_type {
            QuestType::ThreatHunting => {
                if let Some(threat) = threat_landscape.active_threats.first() {
                    Ok(format!("Hunt the {}: {}", threat.name, template.name))
                } else {
                    Ok(format!("Digital Detective: {}", template.name))
                }
            }
            QuestType::IncidentResponse => {
                Ok(format!("Emergency Response: {}", template.name))
            }
            _ => Ok(template.name.clone()),
        }
    }
    
    async fn generate_quest_description(&self, template: &QuestTemplate, skill_profile: &SkillProfile, threat_landscape: &ThreatLandscape) -> Result<String, GamificationError> {
        match template.quest_type {
            QuestType::ThreatHunting => {
                Ok(format!(
                    "Based on your expertise level ({:?}) and current threat landscape, investigate potential security breaches and identify indicators of compromise. Focus on: {}",
                    skill_profile.overall_level,
                    threat_landscape.active_threats.first().map_or("general threats", |t| &t.name)
                ))
            }
            QuestType::IncidentResponse => {
                Ok(format!(
                    "A security incident requires your immediate attention. Use your {} skills to analyze, contain, and resolve the situation efficiently.",
                    skill_profile.get_primary_skill_name()
                ))
            }
            _ => Ok(template.name.clone()),
        }
    }
    
    fn adjust_difficulty(&self, base_difficulty: &QuestDifficulty, skill_profile: &SkillProfile) -> QuestDifficulty {
        match skill_profile.overall_level {
            Level::Beginner => QuestDifficulty::Easy,
            Level::Intermediate => base_difficulty.clone(),
            Level::Advanced => QuestDifficulty::Hard,
            Level::Expert => QuestDifficulty::Expert,
        }
    }
    
    fn adjust_rewards(&self, base_rewards: &[Reward], skill_profile: &SkillProfile) -> Vec<Reward> {
        let mut adjusted_rewards = base_rewards.to_vec();
        
        // Increase rewards for higher skill levels
        let multiplier = match skill_profile.overall_level {
            Level::Beginner => 1.0,
            Level::Intermediate => 1.2,
            Level::Advanced => 1.5,
            Level::Expert => 2.0,
        };
        
        for reward in &mut adjusted_rewards {
            if let Reward::ExperiencePoints(ref mut xp) = reward {
                *xp = (*xp as f64 * multiplier) as u32;
            }
            if let Reward::Tokens(ref mut tokens) = reward {
                *tokens = (*tokens as f64 * multiplier) as u32;
            }
        }
        
        adjusted_rewards
    }
    
    fn calculate_time_limit(&self, quest_type: &QuestType, skill_profile: &SkillProfile) -> chrono::Duration {
        let base_time = match quest_type {
            QuestType::ThreatHunting => chrono::Duration::hours(8),
            QuestType::IncidentResponse => chrono::Duration::hours(2),
            QuestType::RiskAssessment => chrono::Duration::hours(6),
            _ => chrono::Duration::hours(4),
        };
        
        // Adjust based on skill level
        let multiplier = match skill_profile.overall_level {
            Level::Beginner => 1.5,
            Level::Intermediate => 1.2,
            Level::Advanced => 1.0,
            Level::Expert => 0.8,
        };
        
        chrono::Duration::milliseconds((base_time.num_milliseconds() as f64 * multiplier) as i64)
    }
}

/// Narrative engine for immersive quest storytelling
pub struct NarrativeEngine {
    story_templates: HashMap<QuestType, Vec<StoryTemplate>>,
    config: GamificationConfig,
}

impl NarrativeEngine {
    pub fn new(config: &GamificationConfig) -> Result<Self, GamificationError> {
        let mut story_templates = HashMap::new();
        
        // Threat hunting stories
        story_templates.insert(QuestType::ThreatHunting, vec![
            StoryTemplate {
                title: "The Digital Shadow".to_string(),
                theme: "mystery".to_string(),
                tone: "suspenseful".to_string(),
                elements: vec!["clues".to_string(), "investigation".to_string(), "discovery".to_string()],
            },
            StoryTemplate {
                title: "Code Breaker".to_string(),
                theme: "puzzle".to_string(),
                tone: "intellectual".to_string(),
                elements: vec!["patterns".to_string(), "logic".to_string(), "breakthrough".to_string()],
            },
        ]);
        
        // Incident response stories
        story_templates.insert(QuestType::IncidentResponse, vec![
            StoryTemplate {
                title: "Race Against Time".to_string(),
                theme: "emergency".to_string(),
                tone: "urgent".to_string(),
                elements: vec!["pressure".to_string(), "quick thinking".to_string(), "heroism".to_string()],
            },
            StoryTemplate {
                title: "Crisis Commander".to_string(),
                theme: "leadership".to_string(),
                tone: "authoritative".to_string(),
                elements: vec!["strategy".to_string(), "coordination".to_string(), "resolution".to_string()],
            },
        ]);
        
        Ok(Self {
            story_templates,
            config: config.clone(),
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        println!("📚 Narrative Engine initialized");
        Ok(())
    }
    
    pub async fn create_quest_narrative(&self, quest: &Quest, analyst: &SecurityAnalyst) -> Result<QuestNarrative, GamificationError> {
        let templates = self.story_templates.get(&quest.quest_type)
            .ok_or(GamificationError::NarrativeError("No templates for quest type".to_string()))?;
        
        let template = templates.first().unwrap(); // Simple selection - could be more sophisticated
        
        let narrative = QuestNarrative {
            title: template.title.clone(),
            story: self.generate_story(&template, quest, analyst).await?,
            objectives: self.generate_objectives(&quest, &template),
        };
        
        Ok(narrative)
    }
    
    async fn generate_story(&self, template: &StoryTemplate, quest: &Quest, analyst: &SecurityAnalyst) -> Result<String, GamificationError> {
        let story = format!(
            "🎭 {}\n\n{}, your expertise is needed! {}\n\nYour mission: {}\n\nGood luck, security champion! 🦐",
            template.title,
            analyst.name,
            self.get_story_intro(&template.theme),
            quest.description
        );
        
        Ok(story)
    }
    
    fn get_story_intro(&self, theme: &str) -> String {
        match theme {
            "mystery" => "A digital mystery unfolds in the network shadows. Clues are scattered, patterns are hidden, and only your keen analytical eye can piece together the truth.".to_string(),
            "emergency" => "ALERT! The digital fortress is under attack. Every second counts as you race against time to defend our systems and protect our data.".to_string(),
            "puzzle" => "A complex digital puzzle awaits your solution. Cryptic patterns and hidden meanings guard the secrets you seek. Your logical mind is the key.".to_string(),
            "leadership" => "The team looks to you for guidance. In this moment of crisis, your leadership and strategic thinking will make the difference between success and failure.".to_string(),
            _ => "A new security challenge emerges, calling upon your skills and expertise to protect our digital assets.".to_string(),
        }
    }
    
    fn generate_objectives(&self, quest: &Quest, template: &StoryTemplate) -> Vec<String> {
        quest.requirements.iter().map(|req| {
            match req {
                QuestRequirement::AnalyzeLogs => format!("🔍 {} - Analyze system logs for suspicious activity", template.elements.get(0).unwrap_or(&"investigate".to_string())),
                QuestRequirement::IdentifyIOC => format!("🎯 {} - Identify indicators of compromise", template.elements.get(2).unwrap_or(&"discover".to_string())),
                QuestRequirement::ContainThreat => format!("🛡️ {} - Contain the identified threat", template.elements.get(1).unwrap_or(&"protect".to_string())),
                QuestRequirement::DocumentResponse => format!("📝 {} - Document your response strategy", template.elements.get(1).unwrap_or(&"coordinate".to_string())),
                _ => format!("🔧 Complete security task"),
            }
        }).collect()
    }
}

/// Tokenized reward system
pub struct TokenizedRewardSystem {
    token_economy: TokenEconomy,
    reward_history: RwLock<HashMap<UserId, Vec<RewardRecord>>>,
    config: GamificationConfig,
}

impl TokenizedRewardSystem {
    pub fn new(config: &GamificationConfig) -> Result<Self, GamificationError> {
        let token_economy = TokenEconomy::new(&config)?;
        
        Ok(Self {
            token_economy,
            reward_history: RwLock::new(HashMap::new()),
            config: config.clone(),
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        self.token_economy.initialize().await?;
        println!("💰 Tokenized Reward System initialized");
        Ok(())
    }
    
    pub async fn grant_reward(&self, user_id: &UserId, reward: &Reward) -> Result<(), GamificationError> {
        let record = RewardRecord {
            id: RewardId::new(),
            user_id: user_id.clone(),
            reward: reward.clone(),
            granted_at: chrono::Utc::now(),
            status: RewardStatus::Granted,
        };
        
        // Update reward history
        let mut history = self.reward_history.write().await;
        let user_rewards = history.entry(user_id.clone()).or_insert_with(Vec::new);
        user_rewards.push(record);
        
        // Process token rewards
        if let Reward::Tokens(amount) = reward {
            self.token_economy.mint_tokens(user_id, *amount).await?;
        }
        
        println!("🎁 Reward granted to {}: {:?}", user_id.0, reward);
        Ok(())
    }
    
    pub async fn get_user_balance(&self, user_id: &UserId) -> Result<TokenBalance, GamificationError> {
        self.token_economy.get_balance(user_id).await
    }
}

/// Achievement tracking system
pub struct AchievementTracker {
    achievements: HashMap<AchievementId, Achievement>,
    user_achievements: RwLock<HashMap<UserId, HashMap<AchievementId, AchievementProgress>>>,
    config: GamificationConfig,
}

impl AchievementTracker {
    pub fn new(config: &GamificationConfig) -> Result<Self, GamificationError> {
        let mut achievements = HashMap::new();
        
        // Define achievements
        achievements.insert(AchievementId("first_quest".to_string()), Achievement {
            id: AchievementId("first_quest".to_string()),
            name: "First Steps".to_string(),
            description: "Complete your first security quest".to_string(),
            category: AchievementCategory::Milestone,
            rarity: AchievementRarity::Common,
            points: 50,
            icon: "🎯".to_string(),
        });
        
        achievements.insert(AchievementId("threat_hunter_10".to_string()), Achievement {
            id: AchievementId("threat_hunter_10".to_string()),
            name: "Novice Hunter".to_string(),
            description: "Complete 10 threat hunting quests".to_string(),
            category: AchievementCategory::Skill,
            rarity: AchievementRarity::Uncommon,
            points: 200,
            icon: "🕵️".to_string(),
        });
        
        Ok(Self {
            achievements,
            user_achievements: RwLock::new(HashMap::new()),
            config: config.clone(),
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        println!("🏆 Achievement Tracker initialized");
        Ok(())
    }
    
    pub async fn check_achievements(&self, player_progress: &PlayerProgress) -> Result<Vec<Achievement>, GamificationError> {
        let mut new_achievements = Vec::new();
        
        // Check various achievement conditions
        if player_progress.total_quests_completed == 1 {
            if let Some(achievement) = self.achievements.get(&AchievementId("first_quest".to_string())) {
                new_achievements.push(achievement.clone());
            }
        }
        
        if player_progress.quests_by_type.get(&QuestType::ThreatHunting).unwrap_or(&0) >= &10 {
            if let Some(achievement) = self.achievements.get(&AchievementId("threat_hunter_10".to_string())) {
                new_achievements.push(achievement.clone());
            }
        }
        
        Ok(new_achievements)
    }
}

/// Supporting types and structures

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityQuest {
    pub quest: Quest,
    pub narrative: QuestNarrative,
    pub personal_difficulty: PersonalDifficulty,
    pub adaptive_hints: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Quest {
    pub id: QuestId,
    pub title: String,
    pub description: String,
    pub quest_type: QuestType,
    pub difficulty: QuestDifficulty,
    pub requirements: Vec<QuestRequirement>,
    pub rewards: Vec<Reward>,
    pub time_limit: Option<chrono::Duration>,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QuestId(pub String);

impl QuestId {
    pub fn new() -> Self {
        Self(uuid::Uuid::new_v4().to_string())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum QuestType {
    ThreatHunting,
    IncidentResponse,
    RiskAssessment,
    VulnerabilityScanning,
    Forensics,
    Compliance,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum QuestDifficulty {
    Easy,
    Medium,
    Hard,
    Expert,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum QuestRequirement {
    AnalyzeLogs,
    IdentifyIOC,
    CreateAlert,
    AnalyzeIncident,
    ContainThreat,
    DocumentResponse,
    ScanVulnerabilities,
    AssessRisk,
    RecommendMitigation,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Reward {
    ExperiencePoints(u32),
    Tokens(u32),
    Badge(String),
    Title(String),
    Custom(String, serde_json::Value),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QuestNarrative {
    pub title: String,
    pub story: String,
    pub objectives: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PersonalDifficulty {
    Easier,
    Balanced,
    Harder,
    Expert,
    Adaptive,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityAnalyst {
    pub id: UserId,
    pub name: String,
    pub role: String,
    pub experience_level: u32,
    pub skills: Vec<Skill>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserId(pub String);

impl UserId {
    pub fn new() -> Self {
        Self(uuid::Uuid::new_v4().to_string())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Skill {
    pub name: String,
    pub level: u32,
    pub experience: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SkillProfile {
    pub analyst_id: UserId,
    pub threat_hunting: SkillLevel,
    pub incident_response: SkillLevel,
    pub risk_assessment: SkillLevel,
    pub communication: SkillLevel,
    pub overall_level: Level,
    pub preferred_difficulty: QuestDifficulty,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SkillLevel {
    pub level: u32,
    pub experience: u32,
    pub rank: String,
}

impl SkillLevel {
    pub fn new(level: u32) -> Self {
        let rank = match level {
            1..=3 => "Novice".to_string(),
            4..=6 => "Apprentice".to_string(),
            7..=9 => "Journeyman".to_string(),
            10..=12 => "Expert".to_string(),
            _ => "Master".to_string(),
        };
        
        Self { level, experience: 0, rank }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Level {
    Beginner,
    Intermediate,
    Advanced,
    Expert,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum SkillType {
    ThreatHunting,
    IncidentResponse,
    RiskAssessment,
    Communication,
}

impl SkillProfile {
    pub fn get_primary_skill_name(&self) -> String {
        let skills = vec![
            (&self.threat_hunting.level, "Threat Hunting"),
            (&self.incident_response.level, "Incident Response"),
            (&self.risk_assessment.level, "Risk Assessment"),
            (&self.communication.level, "Communication"),
        ];
        
        skills.into_iter()
            .max_by_key(|(level, _)| *level)
            .map(|(_, name)| *name)
            .unwrap_or("General Security")
            .to_string()
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ThreatLandscape {
    pub active_threats: Vec<Threat>,
    pub vulnerability_trends: Vec<String>,
    pub industry_specific_threats: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Threat {
    pub id: String,
    pub name: String,
    pub severity: ThreatSeverity,
    pub category: ThreatCategory,
    pub tactics: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ThreatSeverity {
    Low,
    Medium,
    High,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ThreatCategory {
    Malware,
    Espionage,
    Ransomware,
    Phishing,
    DDoS,
    InsiderThreat,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlayerProgress {
    pub user_id: UserId,
    pub total_quests_completed: u32,
    pub quests_by_type: HashMap<QuestType, u32>,
    pub overall_level: Level,
    pub preferred_difficulty: QuestDifficulty,
    pub achievements: Vec<AchievementId>,
    pub current_streak: u32,
    pub best_streak: u32,
}

impl PlayerProgress {
    pub fn new() -> Self {
        Self {
            user_id: UserId::new(),
            total_quests_completed: 0,
            quests_by_type: HashMap::new(),
            overall_level: Level::Beginner,
            preferred_difficulty: QuestDifficulty::Medium,
            achievements: Vec::new(),
            current_streak: 0,
            best_streak: 0,
        }
    }
    
    pub fn calculate_skill_level(&self, skill_type: SkillType) -> SkillLevel {
        let completed = match skill_type {
            SkillType::ThreatHunting => self.quests_by_type.get(&QuestType::ThreatHunting).unwrap_or(&0),
            SkillType::IncidentResponse => self.quests_by_type.get(&QuestType::IncidentResponse).unwrap_or(&0),
            SkillType::RiskAssessment => self.quests_by_type.get(&QuestType::RiskAssessment).unwrap_or(&0),
            SkillType::Communication => &(self.total_quests_completed / 2),
        };
        
        SkillLevel::new(*completed)
    }
    
    pub async fn complete_quest(&mut self, quest_id: &QuestId, completion_data: QuestCompletion) -> Result<QuestResult, GamificationError> {
        self.total_quests_completed += 1;
        
        // Update quest type counter
        *self.quests_by_type.entry(completion_data.quest_type.clone()).or_insert(0) += 1;
        
        // Update streak
        if completion_data.success {
            self.current_streak += 1;
            if self.current_streak > self.best_streak {
                self.best_streak = self.current_streak;
            }
        } else {
            self.current_streak = 0;
        }
        
        // Calculate experience gained
        let experience_gained = self.calculate_experience_gained(&completion_data);
        
        Ok(QuestResult {
            quest_id: quest_id.clone(),
            success: completion_data.success,
            experience_gained,
            level_up: self.check_level_up(),
            rewards: completion_data.rewards.clone(),
            completion_time: completion_data.completion_time,
        })
    }
    
    fn calculate_experience_gained(&self, completion: &QuestCompletion) -> u32 {
        let base_xp = match completion.difficulty {
            QuestDifficulty::Easy => 50,
            QuestDifficulty::Medium => 100,
            QuestDifficulty::Hard => 200,
            QuestDifficulty::Expert => 400,
        };
        
        let time_bonus = if completion.completion_time < chrono::Duration::minutes(30) {
            1.5
        } else if completion.completion_time < chrono::Duration::hours(1) {
            1.2
        } else {
            1.0
        };
        
        let streak_bonus = 1.0 + (self.current_streak as f64 * 0.1);
        
        (base_xp as f64 * time_bonus * streak_bonus) as u32
    }
    
    fn check_level_up(&self) -> bool {
        // Simple level calculation
        let new_level = match self.total_quests_completed {
            0..=4 => Level::Beginner,
            5..=14 => Level::Intermediate,
            15..=29 => Level::Advanced,
            _ => Level::Expert,
        };
        
        !std::mem::discriminant(&self.overall_level) == std::mem::discriminant(&new_level)
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QuestCompletion {
    pub quest_type: QuestType,
    pub difficulty: QuestDifficulty,
    pub success: bool,
    pub completion_time: chrono::Duration,
    pub rewards: Vec<Reward>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QuestResult {
    pub quest_id: QuestId,
    pub success: bool,
    pub experience_gained: u32,
    pub level_up: bool,
    pub rewards: Vec<Reward>,
    pub completion_time: chrono::Duration,
}

// Additional supporting types

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StoryTemplate {
    pub title: String,
    pub theme: String,
    pub tone: String,
    pub elements: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QuestTemplate {
    pub name: String,
    pub base_difficulty: QuestDifficulty,
    pub quest_type: QuestType,
    pub requirements: Vec<QuestRequirement>,
    pub base_rewards: Vec<Reward>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Achievement {
    pub id: AchievementId,
    pub name: String,
    pub description: String,
    pub category: AchievementCategory,
    pub rarity: AchievementRarity,
    pub points: u32,
    pub icon: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AchievementId(pub String);

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AchievementCategory {
    Milestone,
    Skill,
    Special,
    Seasonal,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AchievementRarity {
    Common,
    Uncommon,
    Rare,
    Epic,
    Legendary,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AchievementProgress {
    pub achievement_id: AchievementId,
    pub progress: f32,
    pub completed: bool,
    pub completed_at: Option<chrono::DateTime<chrono::Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RewardRecord {
    pub id: RewardId,
    pub user_id: UserId,
    pub reward: Reward,
    pub granted_at: chrono::DateTime<chrono::Utc>,
    pub status: RewardStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RewardId(pub String);

impl RewardId {
    pub fn new() -> Self {
        Self(uuid::Uuid::new_v4().to_string())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum RewardStatus {
    Granted,
    Claimed,
    Expired,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenEconomy {
    balances: RwLock<HashMap<UserId, TokenBalance>>,
}

impl TokenEconomy {
    pub fn new(_config: &GamificationConfig) -> Result<Self, GamificationError> {
        Ok(Self {
            balances: RwLock::new(HashMap::new()),
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        println!("💰 Token Economy initialized");
        Ok(())
    }
    
    pub async fn mint_tokens(&self, user_id: &UserId, amount: u32) -> Result<(), GamificationError> {
        let mut balances = self.balances.write().await;
        let balance = balances.entry(user_id.clone()).or_insert(TokenBalance::new());
        balance.add_tokens(amount);
        Ok(())
    }
    
    pub async fn get_balance(&self, user_id: &UserId) -> Result<TokenBalance, GamificationError> {
        Ok(self.balances.read().await.get(user_id).cloned().unwrap_or(TokenBalance::new()))
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenBalance {
    pub total_earned: u32,
    pub total_spent: u32,
    pub current_balance: u32,
}

impl TokenBalance {
    pub fn new() -> Self {
        Self {
            total_earned: 0,
            total_spent: 0,
            current_balance: 0,
        }
    }
    
    pub fn add_tokens(&mut self, amount: u32) {
        self.total_earned += amount;
        self.current_balance += amount;
    }
    
    pub fn spend_tokens(&mut self, amount: u32) -> bool {
        if self.current_balance >= amount {
            self.total_spent += amount;
            self.current_balance -= amount;
            true
        } else {
            false
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LeaderboardManager {
    entries: RwLock<HashMap<LeaderboardPeriod, Vec<LeaderboardEntry>>>,
}

impl LeaderboardManager {
    pub fn new(_config: &GamificationConfig) -> Result<Self, GamificationError> {
        Ok(Self {
            entries: RwLock::new(HashMap::new()),
        })
    }
    
    pub async fn initialize(&mut self) -> Result<(), GamificationError> {
        println!("🏆 Leaderboard Manager initialized");
        Ok(())
    }
    
    pub async fn update_score(&self, user_id: &UserId, score: u32) -> Result<(), GamificationError> {
        let mut entries = self.entries.write().await;
        let daily_entries = entries.entry(LeaderboardPeriod::Daily).or_insert_with(Vec::new);
        
        // Update or add entry
        if let Some(entry) = daily_entries.iter_mut().find(|e| e.user_id == *user_id) {
            entry.score += score;
            entry.last_updated = chrono::Utc::now();
        } else {
            daily_entries.push(LeaderboardEntry {
                user_id: user_id.clone(),
                score,
                rank: 0, // Will be calculated when retrieving
                last_updated: chrono::Utc::now(),
            });
        }
        
        Ok(())
    }
    
    pub async fn get_leaderboard(&self, period: LeaderboardPeriod) -> Result<Vec<LeaderboardEntry>, GamificationError> {
        let entries = self.entries.read().await;
        let mut period_entries = entries.get(&period).cloned().unwrap_or_default();
        
        // Sort by score
        period_entries.sort_by(|a, b| b.score.cmp(&a.score));
        
        // Update ranks
        for (index, entry) in period_entries.iter_mut().enumerate() {
            entry.rank = (index + 1) as u32;
        }
        
        Ok(period_entries)
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum LeaderboardPeriod {
    Daily,
    Weekly,
    Monthly,
    AllTime,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LeaderboardEntry {
    pub user_id: UserId,
    pub score: u32,
    pub rank: u32,
    pub last_updated: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GamificationConfig {
    pub enable_quests: bool,
    pub enable_achievements: bool,
    pub enable_leaderboard: bool,
    pub enable_tokens: bool,
    pub quest_refresh_interval_hours: u32,
}

impl Default for GamificationConfig {
    fn default() -> Self {
        Self {
            enable_quests: true,
            enable_achievements: true,
            enable_leaderboard: true,
            enable_tokens: true,
            quest_refresh_interval_hours: 24,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum GamificationError {
    QuestGenerationError(String),
    NarrativeError(String),
    RewardError(String),
    AchievementError(String),
    ConfigurationError(String),
}
