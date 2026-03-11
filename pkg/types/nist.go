package types

import "time"

// NIST CSF represents the NIST Cybersecurity Framework implementation
type NISTCSF struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Status       NISTStatus  `json:"status"`
	Functions    []Function  `json:"functions"`
	Profiles     []Profile   `json:"profiles"`
	Assessments  []Assessment `json:"assessments"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// NISTStatus represents implementation status
type NISTStatus string

const (
	NISTStatusDraft        NISTStatus = "draft"
	NISTStatusActive       NISTStatus = "active"
	NISTStatusUnderReview  NISTStatus = "under_review"
	NISTStatusApproved     NISTStatus = "approved"
	NISTStatusArchived     NISTStatus = "archived"
)

// Function represents a NIST CSF function (Identify, Protect, Detect, Respond, Recover)
type Function struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Category    FunctionType  `json:"category"`
	Categories  []SubCategory `json:"categories"`
	Status      FunctionStatus `json:"status"`
	Progress    float64       `json:"progress"`
	Owner       string        `json:"owner"`
	LastReview  *time.Time    `json:"last_review,omitempty"`
	NextReview  time.Time     `json:"next_review"`
}

// FunctionType represents the five core functions
type FunctionType string

const (
	FunctionTypeIdentify FunctionType = "identify"
	FunctionTypeProtect  FunctionType = "protect"
	FunctionTypeDetect   FunctionType = "detect"
	FunctionTypeRespond  FunctionType = "respond"
	FunctionTypeRecover  FunctionType = "recover"
)

// FunctionStatus represents function implementation status
type FunctionStatus string

const (
	FunctionStatusNotStarted     FunctionStatus = "not_started"
	FunctionStatusInProgress     FunctionStatus = "in_progress"
	FunctionStatusPartiallyComplete FunctionStatus = "partially_complete"
	FunctionStatusComplete        FunctionStatus = "complete"
	FunctionStatusOptimized       FunctionStatus = "optimized"
)

// SubCategory represents NIST CSF subcategories
type SubCategory struct {
	ID          string                 `json:"id"`
	Number      string                 `json:"number"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	InformativeReferences []string    `json:"informative_references"`
	ImplementationStatus ImplementationStatus `json:"implementation_status"`
	Measurements []Measurement         `json:"measurements"`
	Resources   []Resource             `json:"resources"`
	Priority    Priority               `json:"priority"`
	Owner       string                 `json:"owner"`
	DueDate     time.Time              `json:"due_date"`
}

// ImplementationStatus represents subcategory implementation
type ImplementationStatus string

const (
	ImplementationStatusNotImplemented ImplementationStatus = "not_implemented"
	ImplementationStatusPartiallyImplemented ImplementationStatus = "partially_implemented"
	ImplementationStatusImplemented    ImplementationStatus = "implemented"
	ImplementationStatusTested         ImplementationStatus = "tested"
	ImplementationStatusValidated      ImplementationStatus = "validated"
)

// Measurement represents metrics for subcategories
type Measurement struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Method      MethodType   `json:"method"`
	Frequency   Frequency    `json:"frequency"`
	Target      float64      `json:"target"`
	Current     float64      `json:"current"`
	Unit        string       `json:"unit"`
	Status      MetricStatus `json:"status"`
	LastUpdated time.Time    `json:"last_updated"`
}

// MethodType represents measurement methods
type MethodType string

const (
	MethodTypeAutomated   MethodType = "automated"
	MethodTypeManual       MethodType = "manual"
	MethodTypeHybrid       MethodType = "hybrid"
	MethodTypeQualitative  MethodType = "qualitative"
)

// Frequency represents measurement frequency
type Frequency string

const (
	FrequencyRealtime  Frequency = "realtime"
	FrequencyHourly    Frequency = "hourly"
	FrequencyDaily     Frequency = "daily"
	FrequencyWeekly    Frequency = "weekly"
	FrequencyMonthly   Frequency = "monthly"
	FrequencyQuarterly Frequency = "quarterly"
	FrequencyAnnually  Frequency = "annually"
)

// MetricStatus represents measurement status
type MetricStatus string

const (
	MetricStatusNotMeasured MetricStatus = "not_measured"
	MetricStatusBelowTarget MetricStatus = "below_target"
	MetricStatusAtTarget    MetricStatus = "at_target"
	MetricStatusAboveTarget MetricStatus = "above_target"
)

// Resource represents implementation resources
type Resource struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Type        ResourceType `json:"type"`
	Description string       `json:"description"`
	URL         string       `json:"url"`
	Category    ResourceCategory `json:"category"`
	Priority    Priority     `json:"priority"`
	Status      ResourceStatus `json:"status"`
	Owner       string       `json:"owner"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ResourceType represents resource types
type ResourceType string

const (
	ResourceTypeDocument     ResourceType = "document"
	ResourceTypeTool         ResourceType = "tool"
	ResourceTypeService      ResourceType = "service"
	ResourceTypeTraining     ResourceType = "training"
	ResourceTypeProcedure    ResourceType = "procedure"
	ResourceTypePolicy       ResourceType = "policy"
	ResourceTypeGuidance     ResourceType = "guidance"
)

// ResourceCategory represents resource categories
type ResourceCategory string

const (
	ResourceCategoryPolicy     ResourceCategory = "policy"
	ResourceCategoryProcedure  ResourceCategory = "procedure"
	ResourceCategoryGuidance   ResourceCategory = "guidance"
	ResourceCategoryTool       ResourceCategory = "tool"
	ResourceCategoryTraining   ResourceCategory = "training"
	ResourceCategoryReference  ResourceCategory = "reference"
)

// ResourceStatus represents resource status
type ResourceStatus string

const (
	ResourceStatusDraft        ResourceStatus = "draft"
	ResourceStatusUnderReview  ResourceStatus = "under_review"
	ResourceStatusApproved     ResourceStatus = "approved"
	ResourceStatusPublished    ResourceStatus = "published"
	ResourceStatusArchived     ResourceStatus = "archived"
)

// Priority represents priority levels
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// Profile represents NIST CSF implementation profiles
type Profile struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        ProfileType   `json:"type"`
	Functions   []Function    `json:"functions"`
	Scopes      []Scope       `json:"scopes"`
	Targets     []Target      `json:"targets"`
	Status      ProfileStatus `json:"status"`
	Version     string        `json:"version"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ProfileType represents profile types
type ProfileType string

const (
	ProfileTypeGeneric   ProfileType = "generic"
	ProfileTypeSector    ProfileType = "sector"
	ProfileTypeOrganization ProfileType = "organization"
	ProfileTypeCustom    ProfileType = "custom"
)

// Scope represents implementation scope
type Scope struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        ScopeType    `json:"type"`
	Entities    []Entity     `json:"entities"`
	Exclusions  []string     `json:"exclusions"`
}

// ScopeType represents scope types
type ScopeType string

const (
	ScopeTypeOrganization ScopeType = "organization"
	ScopeTypeBusinessUnit ScopeType = "business_unit"
	ScopeTypeSystem       ScopeType = "system"
	ScopeTypeProcess      ScopeType = "process"
	ScopeTypeLocation     ScopeType = "location"
)

// Entity represents scoped entities
type Entity struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        EntityType `json:"type"`
	Description string     `json:"description"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// EntityType represents entity types
type EntityType string

const (
	EntityTypeSystem       EntityType = "system"
	EntityTypeApplication EntityType = "application"
	EntityTypeNetwork     EntityType = "network"
	EntityTypeProcess      EntityType = "process"
	EntityTypeFacility    EntityType = "facility"
	EntityTypeData        EntityType = "data"
)

// Target represents implementation targets
type Target struct {
	ID          string       `json:"id"`
	SubCategoryID string    `json:"sub_category_id"`
	Description string       `json:"description"`
	Target      TargetLevel  `json:"target"`
	Current     TargetLevel  `json:"current"`
	Metric      string       `json:"metric"`
	DueDate     time.Time    `json:"due_date"`
	Status       TargetStatus `json:"status"`
}

// TargetLevel represents target maturity levels
type TargetLevel string

const (
	TargetLevelPartial TargetLevel = "partial"
	TargetLevelRiskInformed TargetLevel = "risk_informed"
	TargetLevelRepeatable TargetLevel = "repeatable"
	TargetLevelAdaptive   TargetLevel = "adaptive"
)

// TargetStatus represents target achievement status
type TargetStatus string

const (
	TargetStatusNotStarted TargetStatus = "not_started"
	TargetStatusInProgress TargetStatus = "in_progress"
	TargetStatusAchieved   TargetStatus = "achieved"
	TargetStatusExceeded   TargetStatus = "exceeded"
	TargetStatusMissed     TargetStatus = "missed"
)

// ProfileStatus represents profile status
type ProfileStatus string

const (
	ProfileStatusDraft        ProfileStatus = "draft"
	ProfileStatusUnderReview  ProfileStatus = "under_review"
	ProfileStatusApproved     ProfileStatus = "approved"
	ProfileStatusActive       ProfileStatus = "active"
	ProfileStatusDeprecated   ProfileStatus = "deprecated"
)

// Assessment represents NIST CSF assessments
type Assessment struct {
	ID          string         `json:"id"`
	ProfileID   string         `json:"profile_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        AssessmentType `json:"type"`
	Method      AssessmentMethod `json:"method"`
	Scope       string         `json:"scope"`
	Assessor    string         `json:"assessor"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Status      AssessmentStatus `json:"status"`
	Results     []AssessmentResult `json:"results"`
	Score       float64        `json:"score"`
	Findings    []Finding      `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// AssessmentType represents assessment types
type AssessmentType string

const (
	AssessmentTypeInitial   AssessmentType = "initial"
	AssessmentTypePeriodic  AssessmentType = "periodic"
	AssessmentTypeTargeted  AssessmentType = "targeted"
	AssessmentTypeContinuous AssessmentType = "continuous"
)

// AssessmentMethod represents assessment methods
type AssessmentMethod string

const (
	AssessmentMethodInterview     AssessmentMethod = "interview"
	AssessmentMethodDocumentation  AssessmentMethod = "documentation"
	AssessmentMethodObservation    AssessmentMethod = "observation"
	AssessmentMethodTesting        AssessmentMethod = "testing"
	AssessmentMethodQuestionnaire   AssessmentMethod = "questionnaire"
	AssessmentMethodHybrid         AssessmentMethod = "hybrid"
)

// AssessmentStatus represents assessment status
type AssessmentStatus string

const (
	AssessmentStatusPlanned      AssessmentStatus = "planned"
	AssessmentStatusInProgress    AssessmentStatus = "in_progress"
	AssessmentStatusCompleted     AssessmentStatus = "completed"
	AssessmentStatusUnderReview   AssessmentStatus = "under_review"
	AssessmentStatusApproved      AssessmentStatus = "approved"
	AssessmentStatusRejected      AssessmentStatus = "rejected"
)

// AssessmentResult represents assessment results
type AssessmentResult struct {
	SubCategoryID string        `json:"sub_category_id"`
	Rating        RatingLevel   `json:"rating"`
	Score         float64       `json:"score"`
	Evidence      []string      `json:"evidence"`
	Notes         string        `json:"notes"`
	Assessor      string        `json:"assessor"`
	AssessedDate  time.Time     `json:"assessed_date"`
}

// RatingLevel represents maturity rating levels
type RatingLevel string

const (
	RatingLevelLevel0 RatingLevel = "level_0"
	RatingLevelLevel1 RatingLevel = "level_1"
	RatingLevelLevel2 RatingLevel = "level_2"
	RatingLevelLevel3 RatingLevel = "level_3"
	RatingLevelLevel4 RatingLevel = "level_4"
	RatingLevelLevel5 RatingLevel = "level_5"
)

// Recommendation represents assessment recommendations
type Recommendation struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Priority    Priority            `json:"priority"`
	Category    RecommendationCategory `json:"category"`
	SubCategoryID string            `json:"sub_category_id"`
	Actions     []Action            `json:"actions"`
	Owner       string              `json:"owner"`
	DueDate     time.Time           `json:"due_date"`
	Status      RecommendationStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// RecommendationCategory represents recommendation categories
type RecommendationCategory string

const (
	RecommendationCategoryPolicy        RecommendationCategory = "policy"
	RecommendationCategoryProcedure     RecommendationCategory = "procedure"
	RecommendationCategoryTechnical     RecommendationCategory = "technical"
	RecommendationCategoryAdministrative RecommendationCategory = "administrative"
	RecommendationCategoryTraining      RecommendationCategory = "training"
	RecommendationCategoryResource      RecommendationCategory = "resource"
)

// RecommendationStatus represents recommendation status
type RecommendationStatus string

const (
	RecommendationStatusOpen        RecommendationStatus = "open"
	RecommendationStatusInProgress  RecommendationStatus = "in_progress"
	RecommendationStatusCompleted    RecommendationStatus = "completed"
	RecommendationStatusAccepted    RecommendationStatus = "accepted"
	RecommendationStatusRejected    RecommendationStatus = "rejected"
	RecommendationStatusDeferred    RecommendationStatus = "deferred"
)

// Action represents recommendation actions
type Action struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        ActionType `json:"type"`
	Priority    Priority   `json:"priority"`
	Owner       string     `json:"owner"`
	DueDate     time.Time  `json:"due_date"`
	Status      ActionStatus `json:"status"`
	Progress    float64    `json:"progress"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ActionType represents action types
type ActionType string

const (
	ActionTypePolicy      ActionType = "policy"
	ActionTypeProcedure   ActionType = "procedure"
	ActionTypeTechnical   ActionType = "technical"
	ActionTypeTraining    ActionType = "training"
	ActionTypeProcurement ActionType = "procurement"
	ActionTypeConfiguration ActionType = "configuration"
)

// ActionStatus represents action status
type ActionStatus string

const (
	ActionStatusNotStarted    ActionStatus = "not_started"
	ActionStatusInProgress    ActionStatus = "in_progress"
	ActionStatusCompleted     ActionStatus = "completed"
	ActionStatusBlocked       ActionStatus = "blocked"
	ActionStatusCancelled     ActionStatus = "cancelled"
)
