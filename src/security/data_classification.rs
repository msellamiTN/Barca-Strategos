use serde::{Deserialize, Serialize};
use regex::Regex;
use std::collections::HashMap;

/// Data classification levels for Barca-Strategos Phoenix
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub enum DataClassification {
    /// Public data - can be freely shared
    Public,
    /// Internal data - for internal use only
    Internal,
    /// Confidential data - restricted access required
    Confidential,
    /// Restricted data - highest level of protection
    Restricted,
}

impl DataClassification {
    /// Get the security requirements for this classification level
    pub fn security_requirements(&self) -> SecurityRequirements {
        match self {
            DataClassification::Public => SecurityRequirements {
                encryption_required: false,
                access_logging: false,
                retention_days: 365,
                audit_required: false,
            },
            DataClassification::Internal => SecurityRequirements {
                encryption_required: true,
                access_logging: true,
                retention_days: 730,
                audit_required: true,
            },
            DataClassification::Confidential => SecurityRequirements {
                encryption_required: true,
                access_logging: true,
                retention_days: 2555, // 7 years
                audit_required: true,
            },
            DataClassification::Restricted => SecurityRequirements {
                encryption_required: true,
                access_logging: true,
                retention_days: 3650, // 10 years
                audit_required: true,
            },
        }
    }
}

/// Security requirements for data classification levels
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityRequirements {
    pub encryption_required: bool,
    pub access_logging: bool,
    pub retention_days: u32,
    pub audit_required: bool,
}

/// Data Loss Prevention (DLP) engine
pub struct DLPEngine {
    classification_patterns: Vec<ClassificationPattern>,
    custom_patterns: HashMap<String, Regex>,
}

impl DLPEngine {
    pub fn new() -> Self {
        let classification_patterns = vec![
            // Personal Identifiable Information (PII)
            ClassificationPattern {
                classification: DataClassification::Restricted,
                patterns: vec![
                    Regex::new(r"\b\d{3}-\d{2}-\d{4}\b").unwrap(), // SSN
                    Regex::new(r"\b[A-Z]{2}\d{2}\s?\d{4}\s?\d{4}\s?\d{4}\s?\d{2}\b").unwrap(), // IBAN
                    Regex::new(r"\b\d{3}[-.\s]?\d{2}[-.\s]?\d{4}\b").unwrap(), // Driver's license
                ],
                description: "Personal Identifiable Information".to_string(),
            },
            
            // Financial Information
            ClassificationPattern {
                classification: DataClassification::Restricted,
                patterns: vec![
                    Regex::new(r"\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b").unwrap(), // Credit card
                    Regex::new(r"\b\d{3}-\d{6}\b").unwrap(), // Bank account
                    Regex::new(r"\b\d{9,}\b").unwrap(), // Account numbers (9+ digits)
                ],
                description: "Financial Information".to_string(),
            },
            
            // Health Information
            ClassificationPattern {
                classification: DataClassification::Confidential,
                patterns: vec![
                    Regex::new(r"(?i)medical\s+record").unwrap(),
                    Regex::new(r"(?i)health\s+information").unwrap(),
                    Regex::new(r"(?i)patient\s+data").unwrap(),
                    Regex::new(r"(?i)diagnosis").unwrap(),
                    Regex::new(r"(?i)treatment").unwrap(),
                ],
                description: "Protected Health Information".to_string(),
            },
            
            // Authentication Credentials
            ClassificationPattern {
                classification: DataClassification::Restricted,
                patterns: vec![
                    Regex::new(r"(?i)password\s*[:=]\s*\S+").unwrap(),
                    Regex::new(r"(?i)api[_-]?key\s*[:=]\s*\S+").unwrap(),
                    Regex::new(r"(?i)secret\s*[:=]\s*\S+").unwrap(),
                    Regex::new(r"(?i)token\s*[:=]\s*\S+").unwrap(),
                    Regex::new(r"(?i)auth\s*[:=]\s*\S+").unwrap(),
                ],
                description: "Authentication Credentials".to_string(),
            },
            
            // Confidential Business Information
            ClassificationPattern {
                classification: DataClassification::Confidential,
                patterns: vec![
                    Regex::new(r"(?i)confidential").unwrap(),
                    Regex::new(r"(?i)proprietary").unwrap(),
                    Regex::new(r"(?i)trade\s+secret").unwrap(),
                    Regex::new(r"(?i)internal\s+only").unwrap(),
                    Regex::new(r"(?i)company\s+confidential").unwrap(),
                ],
                description: "Confidential Business Information".to_string(),
            },
            
            // Contact Information
            ClassificationPattern {
                classification: DataClassification::Internal,
                patterns: vec![
                    Regex::new(r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b").unwrap(), // Email
                    Regex::new(r"\b\d{3}[-.\s]?\d{3}[-.\s]?\d{4}\b").unwrap(), // Phone
                    Regex::new(r"\b\d{5}\s?\d{5}\b").unwrap(), // ZIP code
                ],
                description: "Contact Information".to_string(),
            },
        ];
        
        Self {
            classification_patterns,
            custom_patterns: HashMap::new(),
        }
    }
    
    /// Classify data content
    pub fn classify_data(&self, content: &str) -> DataClassification {
        let mut highest_classification = DataClassification::Public;
        
        // Check against classification patterns
        for pattern in &self.classification_patterns {
            for regex in &pattern.patterns {
                if regex.is_match(content) {
                    if pattern.classification > highest_classification {
                        highest_classification = pattern.classification.clone();
                    }
                    break; // Found match in this pattern group
                }
            }
        }
        
        // Check custom patterns
        for (name, regex) in &self.custom_patterns {
            if regex.is_match(content) {
                // Custom patterns default to Confidential
                if DataClassification::Confidential > highest_classification {
                    highest_classification = DataClassification::Confidential;
                }
            }
        }
        
        highest_classification
    }
    
    /// Check if data can be exported based on classification
    pub fn can_export(&self, content: &str, target_classification: DataClassification) -> bool {
        let data_classification = self.classify_data(content);
        data_classification <= target_classification
    }
    
    /// Sanitize data by removing or redacting sensitive information
    pub fn sanitize_data(&self, content: &str, allowed_classification: DataClassification) -> String {
        let mut sanitized = content.to_string();
        
        for pattern in &self.classification_patterns {
            if pattern.classification > allowed_classification {
                for regex in &pattern.patterns {
                    sanitized = regex.replace_all(&sanitized, "[REDACTED]").to_string();
                }
            }
        }
        
        sanitized
    }
    
    /// Add custom classification pattern
    pub fn add_custom_pattern(&mut self, name: String, pattern: Regex, classification: DataClassification) {
        self.custom_patterns.insert(name, pattern);
    }
    
    /// Get data classification statistics
    pub fn get_classification_stats(&self, content: &str) -> ClassificationStats {
        let mut stats = ClassificationStats::default();
        
        for pattern in &self.classification_patterns {
            for regex in &pattern.patterns {
                if regex.is_match(content) {
                    match pattern.classification {
                        DataClassification::Public => stats.public_count += 1,
                        DataClassification::Internal => stats.internal_count += 1,
                        DataClassification::Confidential => stats.confidential_count += 1,
                        DataClassification::Restricted => stats.restricted_count += 1,
                    }
                    break;
                }
            }
        }
        
        stats
    }
}

/// Pattern for data classification
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ClassificationPattern {
    pub classification: DataClassification,
    pub patterns: Vec<Regex>,
    pub description: String,
}

/// Statistics for data classification
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct ClassificationStats {
    pub public_count: u32,
    pub internal_count: u32,
    pub confidential_count: u32,
    pub restricted_count: u32,
}

impl ClassificationStats {
    pub fn total_matches(&self) -> u32 {
        self.public_count + self.internal_count + self.confidential_count + self.restricted_count
    }
    
    pub fn highest_classification(&self) -> DataClassification {
        if self.restricted_count > 0 {
            DataClassification::Restricted
        } else if self.confidential_count > 0 {
            DataClassification::Confidential
        } else if self.internal_count > 0 {
            DataClassification::Internal
        } else {
            DataClassification::Public
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_classify_pii_data() {
        let dlp = DLPEngine::new();
        let ssn_content = "User SSN: 123-45-6789";
        assert_eq!(dlp.classify_data(ssn_content), DataClassification::Restricted);
    }
    
    #[test]
    fn test_classify_public_data() {
        let dlp = DLPEngine::new();
        let public_content = "This is public information";
        assert_eq!(dlp.classify_data(public_content), DataClassification::Public);
    }
    
    #[test]
    fn test_sanitize_data() {
        let dlp = DLPEngine::new();
        let sensitive_content = "Password: secret123 and email: user@example.com";
        let sanitized = dlp.sanitize_data(sensitive_content, DataClassification::Public);
        assert!(sanitized.contains("[REDACTED]"));
        assert!(!sanitized.contains("secret123"));
    }
}
