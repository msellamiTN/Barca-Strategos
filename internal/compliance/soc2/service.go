package soc2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"barca-strategos/pkg/types"
)

type Service struct {
	controls             map[string]*types.SOC2Control
	assessments          []types.SOC2Assessment
	reports              map[string]types.SOC2Report
	stats                *types.SOC2Stats
	broadcast            func([]byte)
}

func New() *Service {
	return &Service{
		controls: make(map[string]*types.SOC2Control),
		assessments: make([]types.SOC2Assessment, 0),
		reports:   make(map[string]types.SOC2Report),
		stats:     &types.SOC2Stats{},
	}
}

func (s *Service) SetBroadcast(broadcast func([]byte)) {
	s.broadcast = broadcast
}

func (s *Service) Initialize(ctx context.Context) error {
	// Initialize SOC 2 controls
	s.loadSOC2Controls()
	
	// Start background monitoring
	go s.startBackgroundMonitoring()
	
	return nil
}

func (s *Service) AssessSOC2Compliance(ctx context.Context, scope *types.SOC2Scope) (*types.SOC2Assessment, error) {
	// Get all relevant SOC 2 controls
	controls := s.getControlsByScope(scope)
	
	// Assess each control
	var controlAssessments []types.SOC2ControlAssessment
	for _, control := range controls {
		assessment := s.assessSOC2Control(control)
		controlAssessments = append(controlAssessments, *assessment)
	}
	
	// Calculate overall SOC 2 score
	overallScore := s.calculateSOC2Score(controlAssessments)
	
	// Generate findings
	findings := s.generateSOC2Findings(controlAssessments)
	
	assessment := &types.SOC2Assessment{
		AssessmentID:        uuid.New().String(),
		Timestamp:           time.Now().UTC(),
		Framework:           "SOC 2",
		Version:             "2017",
		Scope:               *scope,
		OverallScore:        overallScore,
		ControlAssessments: controlAssessments,
		Findings:            findings,
		Recommendations:     s.generateSOC2Recommendations(findings),
		LastAssessed:        time.Now().UTC(),
	}
	
	// Store assessment
	s.assessments = append(s.assessments, *assessment)
	s.updateStats()
	s.broadcastSOC2Update()
	
	return assessment, nil
}

func (s *Service) GetControlStatus(ctx context.Context, controlID string) (*types.SOC2ControlStatus, error) {
	control, exists := s.controls[controlID]
	if !exists {
		return nil, types.SOC2Error{
			Code:    "CONTROL_NOT_FOUND",
			Message: fmt.Sprintf("SOC 2 control %s not found", controlID),
		}
	}
	return &control.Status, nil
}

func (s *Service) UpdateControl(ctx context.Context, controlID string, update *types.SOC2ControlUpdate) error {
	control, exists := s.controls[controlID]
	if !exists {
		return types.SOC2Error{
			Code:    "CONTROL_NOT_FOUND",
			Message: fmt.Sprintf("SOC 2 control %s not found", controlID),
		}
	}
	
	// Update control based on update type
	switch update.UpdateType {
	case types.UpdateTypeAssessment:
		// Update assessment information
		control.LastReviewDate = &update.Timestamp
	case types.UpdateTypeImplementation:
		// Update implementation details
		if update.Timestamp.After(control.ImplementationDate.UnixNano()) {
			control.ImplementationDate = &update.Timestamp
		}
		control.Evidence = append(control.Evidence, update.Evidence...)
	case types.UpdateTypeStatus:
		// Update control status
		// Parse status from notes if provided
		if update.Notes != "" {
			switch update.Notes {
			case "not_implemented":
				control.Status = types.SOC2ControlStatusNotImplemented
			case "partially_implemented":
				control.Status = types.SOC2ControlStatusPartiallyImplemented
			case "implemented":
				control.Status = types.SOC2ControlStatusImplemented
			case "compliant":
				control.Status = types.SOC2ControlStatusCompliant
			}
		}
	}
	
	s.controls[controlID] = control
	s.broadcastSOC2Update()
	
	return nil
}

func (s *Service) GenerateSOC2Report(ctx context.Context, assessment *types.SOC2Assessment) (*types.SOC2Report, error) {
	report := &types.SOC2Report{
		ReportID:      uuid.New().String(),
		GeneratedAt:   time.Now().UTC(),
		Assessment:    *assessment,
		ReportContent: s.generateReportContent(assessment),
		Format:        "json",
	}
	
	s.reports[report.ReportID] = *report
	return report, nil
}

func (s *Service) GetSOC2Stats(ctx context.Context) (*types.SOC2Stats, error) {
	return s.stats, nil
}

// Private methods

func (s *Service) loadSOC2Controls() {
	// Load SOC 2 controls (Trust Services Criteria)
	controls := []*types.SOC2Control{
		// Common Criteria 1: Governance
		{
			ID:          "CC1.1",
			Title:       "Governance",
			Description: "Establish and communicate governance framework",
			Category:    types.SOC2ControlCategoryGovernance,
			Subcategories: []string{
				"Governance framework",
				"Board oversight",
				"Management direction",
				"Legal and compliance",
				"Risk management",
				"Ethics and compliance",
			},
			Objective:      "Establish and communicate governance framework",
			ControlType:    types.SOC2ControlTypeOrganizational,
			Status:         types.SOC2ControlStatusImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -30)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -7)}[0],
			Evidence:       []string{"Governance policy document", "Board meeting minutes"},
			Owner:          "Board of Directors",
			RiskLevel:      types.RiskLevelLow,
		},
		
		// Common Criteria 2: Asset Management
		{
			ID:          "CC2.1",
			Title:       "Asset Inventory",
			Description: "Maintain complete and accurate inventory of all hardware, software, and data assets",
			Category:    types.SOC2ControlCategoryAssetManagement,
			Subcategories: []string{
				"Hardware inventory",
				"Software inventory",
				"Data inventory",
				"Cloud assets",
				"Mobile devices",
			},
			Objective:      "Maintain complete and accurate inventory of all assets",
			ControlType:    types.SOC2ControlTypeOrganizational,
			Status:         types.SOC2ControlStatusImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -20)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -5)}[0],
			Evidence:       []string{"Asset registry", "CMDB", "Asset management system"},
			Owner:          "IT Asset Manager",
			RiskLevel:      types.RiskLevelMedium,
		},
		
		// Common Criteria 3: Identity Management and Access Control
		{
			ID:          "CC3.2",
			Title:       "Identity Management and Access Control",
			Description: "Identify, authenticate, and authorize access to systems",
			Category:    types.SOC2ControlCategoryAccessControl,
			Subcategories: []string{
				"User access management",
				"Remote access",
				"Multi-factor authentication",
				"Privileged access management",
				"Account lifecycle management",
				"Access certification",
				"Identity proofing",
			},
			Objective:      "Identify, authenticate, and authorize access to systems",
			ControlType:    types.SOC2ControlTypeTechnical,
			Status:         types.SOC2ControlStatusImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -15)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -3)}[0],
			Evidence:       []string{"User access policies", "Authentication system", "MFA system"},
			Owner:          "Identity Management Team",
			RiskLevel:      types.RiskLevelHigh,
		},
		
		// Common Criteria 4: Security Awareness and Training
		{
			ID:          "CC4.1",
			Title:       "Security Awareness and Training",
			Description: "Provide security awareness training to all personnel",
			Category:    types.SOC2ControlCategoryOperational,
			Subcategories: []string{
				"Security training program",
				"Phishing awareness",
				"Social engineering awareness",
				"Security culture",
				"Threat intelligence sharing",
			},
			Objective:      "Ensure all personnel understand their security responsibilities",
			ControlType:    types.SOC2ControlTypeOperational,
			Status:         types.SOC2ControlStatusPartiallyImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -10)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -2)}[0],
			Evidence:       []string{"Security training materials", "Phishing simulations"},
			Owner:          "Security Team",
			RiskLevel:      types.RiskLevelMedium,
		},
		
		// Common Criteria 5: Vulnerability Management
		{
			ID:          "CC5.1",
			Title:       "Vulnerability Management",
			Description: "Identify, assess, and remediate vulnerabilities",
			Category:    types.SOC2ControlCategoryOperational,
			Subcategories: []string{
				"Vulnerability scanning",
				"Penetration testing",
				"Vulnerability assessment",
				"Patch management",
				"CVE monitoring",
			},
			Objective:      "Continuously identify and remediate vulnerabilities",
			ControlType:    types.SOC2ControlTypeTechnical,
			Status:         types.SOC2ControlStatusNotImplemented,
			ImplementationDate: nil,
			LastReviewDate: nil,
			Evidence:       []string{},
			Owner:          "Security Team",
			RiskLevel:      types.RiskLevelHigh,
		},
		
		// Common Criteria 6: Incident Response
		{
			ID:          "CC6.1",
			Title:       "Incident Response",
			Description: "Establish and implement incident response capabilities",
			Category:    types.SOC2ControlCategoryOperational,
			Subcategories: []string{
				"Incident response planning",
				"Incident response playbooks",
				"Incident notification procedures",
				"Forensic capabilities",
				"Tabletop exercises",
				"Threat hunting",
			},
			Objective:      "Ensure timely and effective incident response",
			ControlType:    types.SOC2ControlTypeOperational,
			Status:         types.SOC2ControlStatusImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -30)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -7)}[0],
			Evidence:       []string{"Incident response plan", "Playbooks"},
			Owner:          "SOC Team",
			RiskLevel:      types.RiskLevelMedium,
		},
		
		// Common Criteria 7: Disaster Recovery Planning
		{
			ID:          "CC7.1",
			Title:       "Disaster Recovery Planning",
			Description: "Establish and test disaster recovery plans",
			Category:    types.SOC2ControlCategoryOperational,
			Subcategories: []string{
				"Business continuity planning",
				"Disaster recovery testing",
				"Backup and recovery procedures",
				"Alternative processing sites",
				"Crisis communication",
				"Tabletop exercises",
			},
			Objective:      "Ensure business continuity during disruptions",
			ControlType:    types.SOC2ControlTypeOperational,
			Status:         types.SOC2ControlStatusNotImplemented,
			ImplementationDate: nil,
			LastReviewDate: nil,
			Evidence:       []string{},
			Owner:          "Business Continuity Team",
			RiskLevel:      types.RiskLevelHigh,
		},
		
		// Common Criteria 8: Test and Evaluation
		{
			ID:          "CC8.1",
			Title:       "Penetration Testing",
			Description: "Conduct regular penetration testing",
			Category:    types.SOC2ControlCategoryTestEvaluation,
			Subcategories: []string{
				"External penetration testing",
				"Internal penetration testing",
				"Social engineering testing",
				"Application security testing",
			},
			Objective:      "Test security controls through penetration testing",
			ControlType:    types.SOC2ControlTypeTechnical,
			Status:         types.SOC2ControlStatusImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -25)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -3)}[0],
			Evidence:       []string{"Penetration test reports", "Security assessment reports"},
			Owner:          "Security Team",
			RiskLevel:      types.RiskLevelHigh,
		},
		
		// Common Criteria 9: Communications Security
		{
			ID:          "CC9.1",
			Title:       "Network Security Monitoring",
			Description: "Monitor network traffic for security events",
			Category:    types.SOC2ControlCategoryCommunicationsSecurity,
			Subcategories: []string{
				"Network intrusion detection",
				"Malware analysis",
				"Log analysis",
				"Network traffic analysis",
				"IDS integration",
				"Threat hunting",
				"Network device monitoring",
			},
			Objective:      "Detect and respond to network security incidents",
			ControlType:    types.SOC2ControlTypeTechnical,
			Status:         types.SOC2ControlStatusImplemented,
			ImplementationDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -25)}[0],
			LastReviewDate: &[]time.Time{time.Now().UTC().AddDate(0, 0, -3)}[0],
			Evidence:       []string{"Network logs", "IDS alerts", "Firewall logs"},
			Owner:          "Network Security Team",
			RiskLevel:      types.RiskLevelMedium,
		},
	}
	
	// Load controls into framework
	for _, control := range controls {
		s.controls[control.ID] = control
	}
}

func (s *Service) assessSOC2Control(control *types.SOC2Control) *types.SOC2ControlAssessment {
	// Assess control implementation
	implementationScore := s.assessSOC2Implementation(control)
	effectivenessScore := s.assessSOC2Effectiveness(control)
	complianceScore := (implementationScore + effectivenessScore) / 2.0
	
	// Determine status
	var status types.SOC2ControlStatus
	if complianceScore >= 0.9 {
		status = types.SOC2ControlStatusCompliant
	} else if complianceScore >= 0.7 {
		status = types.SOC2ControlStatusImplemented
	} else if complianceScore >= 0.5 {
		status = types.SOC2ControlStatusPartiallyImplemented
	} else {
		status = types.SOC2ControlStatusNotImplemented
	}
	
	return &types.SOC2ControlAssessment{
		ControlID:       control.ID,
		ControlTitle:    control.Title,
		Category:        control.Category,
		ComplianceScore: complianceScore,
		Status:          status,
		Findings:        s.assessSOC2Findings(control, complianceScore),
		Recommendations: s.assessSOC2Recommendations(control, complianceScore),
		LastAssessed:    time.Now().UTC(),
	}
}

func (s *Service) assessSOC2Implementation(control *types.SOC2Control) float64 {
	// Assess implementation based on evidence and status
	var baseScore float64
	switch control.Status {
	case types.SOC2ControlStatusCompliant:
		baseScore = 1.0
	case types.SOC2ControlStatusImplemented:
		baseScore = 0.8
	case types.SOC2ControlStatusPartiallyImplemented:
		baseScore = 0.6
	case types.SOC2ControlStatusNotImplemented:
		baseScore = 0.0
	}
	
	// Adjust based on evidence quality
	var evidenceScore float64
	if len(control.Evidence) >= 3 {
		evidenceScore = 0.9
	} else if len(control.Evidence) >= 1 {
		evidenceScore = 0.7
	} else {
		evidenceScore = 0.5
	}
	
	return baseScore * evidenceScore
}

func (s *Service) assessSOC2Effectiveness(control *types.SOC2Control) float64 {
	// Assess effectiveness based on control type and risk level
	var baseScore float64
	switch control.RiskLevel {
	case types.RiskLevelLow:
		baseScore = 0.9
	case types.RiskLevelMedium:
		baseScore = 0.8
	case types.RiskLevelHigh:
		baseScore = 0.7
	case types.RiskLevelCritical:
		baseScore = 0.6
	default:
		baseScore = 0.8
	}
	
	// Adjust based on control type
	var typeAdjustment float64
	switch control.ControlType {
	case types.SOC2ControlTypeOrganizational:
		typeAdjustment = 0.0
	case types.SOC2ControlTypeTechnical:
		typeAdjustment = 0.1
	case types.SOC2ControlTypeOperational:
		typeAdjustment = 0.0
	}
	
	return baseScore + typeAdjustment
}

func (s *Service) assessSOC2Findings(control *types.SOC2Control, complianceScore float64) []types.SOC2Finding {
	var findings []types.SOC2Finding
	
	if complianceScore < 0.7 {
		findings = append(findings, types.SOC2Finding{
			Severity:       types.FindingSeverityHigh,
			Description:    fmt.Sprintf("SOC 2 control %s is not adequately implemented", control.ID),
			Recommendation: fmt.Sprintf("Implement %s control according to SOC 2 requirements", control.Title),
			EvidenceGaps: s.identifySOC2EvidenceGaps(control),
		})
	}
	
	return findings
}

func (s *Service) assessSOC2Recommendations(control *types.SOC2Control, complianceScore float64) []string {
	var recommendations []string
	
	if complianceScore < 0.5 {
		recommendations = append(recommendations, fmt.Sprintf("Implement %s control completely", control.Title))
	} else if complianceScore < 0.8 {
		recommendations = append(recommendations, fmt.Sprintf("Enhance %s control implementation", control.Title))
	}
	
	return recommendations
}

func (s *Service) identifySOC2EvidenceGaps(control *types.SOC2Control) []string {
	var gaps []string
	
	// Check for common evidence gaps based on control category
	switch control.Category {
	case types.SOC2ControlCategoryGovernance:
		if !s.hasEvidenceContaining(control, "policy") {
			gaps = append(gaps, "Missing governance evidence")
		}
	case types.SOC2ControlCategoryAccessControl:
		if !s.hasEvidenceContaining(control, "authentication") {
			gaps = append(gaps, "Missing authentication evidence")
		}
	case types.SOC2ControlCategoryOperational:
		if !s.hasEvidenceContaining(control, "penetration test") && control.Title == "Penetration Testing" {
			gaps = append(gaps, "Missing penetration test evidence")
		}
		if !s.hasEvidenceContaining(control, "incident response") && control.Title == "Incident Response" {
			gaps = append(gaps, "Missing incident response evidence")
		}
	}
	
	return gaps
}

func (s *Service) hasEvidenceContaining(control *types.SOC2Control, searchTerm string) bool {
	for _, evidence := range control.Evidence {
		if len(evidence) >= len(searchTerm) {
			// Simple substring check
			for i := 0; i <= len(evidence)-len(searchTerm); i++ {
				if evidence[i:i+len(searchTerm)] == searchTerm {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) calculateSOC2Score(assessments []types.SOC2ControlAssessment) float64 {
	if len(assessments) == 0 {
		return 1.0
	}
	
	var totalScore float64
	for _, assessment := range assessments {
		totalScore += assessment.ComplianceScore
	}
	
	return totalScore / float64(len(assessments))
}

func (s *Service) generateSOC2Findings(assessments []types.SOC2ControlAssessment) []types.SOC2Finding {
	var findings []types.SOC2Finding
	
	for _, assessment := range assessments {
		findings = append(findings, assessment.Findings...)
	}
	
	return findings
}

func (s *Service) generateSOC2Recommendations(findings []types.SOC2Finding) []types.SOC2Recommendation {
	var recommendations []types.SOC2Recommendation
	
	// Group findings by severity
	var criticalFindings, highFindings, mediumFindings, lowFindings []types.SOC2Finding
	
	for _, finding := range findings {
		switch finding.Severity {
		case types.FindingSeverityCritical:
			criticalFindings = append(criticalFindings, finding)
		case types.FindingSeverityHigh:
			highFindings = append(highFindings, finding)
		case types.FindingSeverityMedium:
			mediumFindings = append(mediumFindings, finding)
		case types.FindingSeverityLow:
			lowFindings = append(lowFindings, finding)
		}
	}
	
	// Generate recommendations for each priority level
	if len(criticalFindings) > 0 {
		recommendations = append(recommendations, types.SOC2Recommendation{
			Priority:        types.RecommendationPriorityCritical,
			Title:           "Address Critical SOC 2 Issues",
			Description:     "Immediate action required for critical SOC 2 gaps",
			Findings:        criticalFindings,
			EstimatedEffort: "2-4 weeks",
			Owner:           "CISO",
		})
	}
	
	if len(highFindings) > 0 {
		recommendations = append(recommendations, types.SOC2Recommendation{
			Priority:        types.RecommendationPriorityHigh,
			Title:           "Address High Priority SOC 2 Issues",
			Description:     "Address high priority SOC 2 issues within 30 days",
			Findings:        highFindings,
			EstimatedEffort: "1-2 weeks",
			Owner:           "Security Team",
		})
	}
	
	if len(mediumFindings) > 0 {
		recommendations = append(recommendations, types.SOC2Recommendation{
			Priority:        types.RecommendationPriorityMedium,
			Title:           "Address Medium Priority SOC 2 Issues",
			Description:     "Address medium priority SOC 2 issues within 60 days",
			Findings:        mediumFindings,
			EstimatedEffort: "2-3 weeks",
			Owner:           "Department Heads",
		})
	}
	
	if len(lowFindings) > 0 {
		recommendations = append(recommendations, types.SOC2Recommendation{
			Priority:        types.RecommendationPriorityLow,
			Title:           "Address Low Priority SOC 2 Issues",
			Description:     "Address low priority SOC 2 issues within 90 days",
			Findings:        lowFindings,
			EstimatedEffort: "1-2 weeks",
			Owner:           "Security Team",
		})
	}
	
	return recommendations
}

func (s *Service) getControlsByScope(scope *types.SOC2Scope) []*types.SOC2Control {
	var controls []*types.SOC2Control
	
	for _, control := range s.controls {
		// Simple scope filtering - in real implementation, this would be more sophisticated
		controls = append(controls, control)
	}
	
	return controls
}

func (s *Service) generateReportContent(assessment *types.SOC2Assessment) string {
	// Generate report content - in real implementation, this would be more sophisticated
	return fmt.Sprintf("SOC 2 Compliance Report\n\nAssessment ID: %s\nOverall Score: %.2f\nControls Assessed: %d\nFindings: %d\nRecommendations: %d\nGenerated: %s",
		assessment.AssessmentID,
		assessment.OverallScore,
		len(assessment.ControlAssessments),
		len(assessment.Findings),
		len(assessment.Recommendations),
		assessment.Timestamp.Format("2006-01-02 15:04:05"),
	)
}

func (s *Service) updateStats() {
	// Update statistics based on current controls and assessments
	s.stats.TotalControls = len(s.controls)
	s.stats.CompliantControls = 0
	s.stats.ImplementedControls = 0
	s.stats.PartiallyImplementedControls = 0
	s.stats.NotImplementedControls = 0
	
	var totalScore float64
	for _, control := range s.controls {
		switch control.Status {
		case types.SOC2ControlStatusCompliant:
			s.stats.CompliantControls++
			totalScore += 1.0
		case types.SOC2ControlStatusImplemented:
			s.stats.ImplementedControls++
			totalScore += 0.8
		case types.SOC2ControlStatusPartiallyImplemented:
			s.stats.PartiallyImplementedControls++
			totalScore += 0.6
		case types.SOC2ControlStatusNotImplemented:
			s.stats.NotImplementedControls++
			totalScore += 0.0
		}
	}
	
	if s.stats.TotalControls > 0 {
		s.stats.AverageComplianceScore = totalScore / float64(s.stats.TotalControls)
	}
}

func (s *Service) startBackgroundMonitoring() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		// Monitor SOC 2 compliance status
		s.monitorSOC2Status()
		
		// Perform periodic assessments
		if time.Now().Day()%14 == 0 { // Every 14 days
			s.performSOC2Assessments()
		}
		
		// Collect metrics
		s.collectSOC2Metrics()
	}
}

func (s *Service) monitorSOC2Status() {
	// Monitor overall SOC 2 compliance status
	scope := &types.SOC2Scope{
		Departments: []string{"IT", "Security"},
		Systems:     []string{"Phoenix Core", "Database"},
		Processes:   []string{"Incident Response"},
	}
	
	assessment, err := s.AssessSOC2Compliance(context.Background(), scope)
	if err == nil {
		// Check for compliance issues
		if assessment.OverallScore < 0.8 {
			s.triggerSOC2Alert(assessment)
		}
	}
}

func (s *Service) performSOC2Assessments() {
	// Perform periodic SOC 2 control assessments
	for _, control := range s.controls {
		assessment := s.assessSOC2Control(control)
		s.UpdateControl(context.Background(), control.ID, &types.SOC2ControlUpdate{
			UpdateType: types.UpdateTypeAssessment,
			UpdatedBy:  "system",
			Timestamp:  time.Now().UTC(),
			Notes:      fmt.Sprintf("Periodic assessment: score %.2f", assessment.ComplianceScore),
			Evidence:   []string{},
		})
	}
}

func (s *Service) collectSOC2Metrics() {
	// Collect SOC 2 metrics
	// For now, just log the collection
	fmt.Printf("SOC 2: Metrics collected successfully\n")
}

func (s *Service) triggerSOC2Alert(assessment *types.SOC2Assessment) {
	// Trigger SOC 2 alert
	// In real implementation, this would integrate with the alert system
	fmt.Printf("SOC 2: Compliance alert triggered - Score: %.2f\n", assessment.OverallScore)
}

func (s *Service) broadcastSOC2Update() {
	if s.broadcast != nil {
		controls := make([]*types.SOC2Control, 0, len(s.controls))
		for _, control := range s.controls {
			controls = append(controls, control)
		}
		
		data := map[string]interface{}{
			"type":     "soc2_update",
			"controls": controls,
			"stats":    s.stats,
		}
		
		if jsonData, err := json.Marshal(data); err == nil {
			s.broadcast(jsonData)
		}
	}
}
