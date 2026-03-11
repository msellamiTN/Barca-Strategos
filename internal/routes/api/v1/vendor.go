package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"barca-strategos/internal/compliance/vendor"
	"barca-strategos/pkg/types"
)

func RegisterVendorRoutes(app *fiber.App, vendorSvc *vendor.Service, broadcast func([]byte)) {
	vendorSvc.SetBroadcast(broadcast)

	// Vendor Management Routes
	app.Post("/api/v1/compliance/vendor", func(c *fiber.Ctx) error {
		var request types.VendorRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		vendor, err := vendorSvc.AddVendor(c.Context(), &request)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(vendor)
	})

	app.Get("/api/v1/compliance/vendor", func(c *fiber.Ctx) error {
		vendors, err := vendorSvc.ListVendors(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(vendors)
	})

	app.Get("/api/v1/compliance/vendor/:id", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		vendor, err := vendorSvc.GetVendor(c.Context(), vendorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(vendor)
	})

	app.Put("/api/v1/compliance/vendor/:id", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		var request types.VendorUpdateRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		err := vendorSvc.UpdateVendor(c.Context(), vendorID, &request)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Vendor updated successfully",
		})
	})

	// Assessment Management Routes
	app.Post("/api/v1/compliance/vendor/:id/assess", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		assessment, err := vendorSvc.ConductVendorAssessment(c.Context(), vendorID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(assessment)
	})

	app.Get("/api/v1/compliance/vendor/:id/assessments", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		// Get vendor to access assessments
		vendor, err := vendorSvc.GetVendor(c.Context(), vendorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Mock assessments - in real implementation, this would query the service
		assessments := []types.VendorAssessment{
			{
				AssessmentID: fmt.Sprintf("ass-%s-1", vendorID),
				VendorID:     vendorID,
				VendorName:   vendor.Name,
				Timestamp:    time.Now().UTC(),
				RiskAssessment: types.VendorRiskAssessment{
					RiskScore:             0.6,
					RiskFactors:           []types.RiskFactor{
						{
							ID:                "rf1",
							Name:              "Data Access Controls",
							Description:       "Vendor data access control assessment",
							Severity:          types.RiskFactorSeverityMedium,
							Score:             0.6,
							MitigationRequired: true,
						},
					},
					FinancialStability:    0.8,
					OperationalCapability: 0.85,
					SecurityPosture:       0.9,
					ReputationScore:       0.8,
				},
				ComplianceAssessment: types.VendorComplianceAssessment{
					ComplianceScore:     0.85,
					ComplianceIssues:    []types.ComplianceIssue{},
					Certifications:      []types.Certification{},
					RegulatoryAdherence: 0.9,
					PolicyCompliance:    0.8,
				},
				OverallRiskScore:   0.68,
				RiskLevel:          types.VendorRiskLevelMedium,
				Recommendations:    []types.VendorRecommendation{},
				NextAssessmentDate: time.Now().UTC().AddDate(0, 0, 90),
			},
		}

		return c.JSON(assessments)
	})

	// Risk Status Routes
	app.Get("/api/v1/compliance/vendor/:id/risk-status", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		riskStatus, err := vendorSvc.GetVendorRiskStatus(c.Context(), vendorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(riskStatus)
	})

	// Compliance Monitoring Routes
	app.Get("/api/v1/compliance/vendor/:id/compliance", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		complianceStatus, err := vendorSvc.MonitorVendorCompliance(c.Context(), vendorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(complianceStatus)
	})

	// Reporting Routes
	app.Post("/api/v1/compliance/vendor/report", func(c *fiber.Ctx) error {
		var scope types.VendorScope
		if err := json.Unmarshal(c.Body(), &scope); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		report, err := vendorSvc.GenerateVendorReport(c.Context(), &scope)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(report)
	})

	// Statistics Routes
	app.Get("/api/v1/compliance/vendor/stats", func(c *fiber.Ctx) error {
		stats, err := vendorSvc.GetVendorStats(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(stats)
	})

	// Category Management Routes
	app.Get("/api/v1/compliance/vendor/categories", func(c *fiber.Ctx) error {
		// Return predefined categories
		categories := []types.VendorCategory{
			{
				ID:                     "CLOUD",
				Name:                   "Cloud Service Providers",
				Description:            "Vendors providing cloud infrastructure and services",
				RiskFactors:            []string{"Data sovereignty", "Service availability", "Data access controls", "Encryption standards"},
				BaseRiskScore:          0.6,
				AssessmentFrequencyDays: 180,
			},
			{
				ID:                     "SOFTWARE",
				Name:                   "Software Vendors",
				Description:            "Vendors providing software products and solutions",
				RiskFactors:            []string{"Software security", "Update management", "License compliance", "Support availability"},
				BaseRiskScore:          0.4,
				AssessmentFrequencyDays: 365,
			},
			{
				ID:                     "CONSULTING",
				Name:                   "Consulting Services",
				Description:            "Vendors providing consulting and professional services",
				RiskFactors:            []string{"Access to sensitive data", "Consultant credentials", "Confidentiality agreements", "Quality of deliverables"},
				BaseRiskScore:          0.5,
				AssessmentFrequencyDays: 365,
			},
			{
				ID:                     "INFRASTRUCTURE",
				Name:                   "Infrastructure Providers",
				Description:            "Vendors providing physical infrastructure and facilities",
				RiskFactors:            []string{"Physical security", "Disaster recovery", "Environmental controls", "Access management"},
				BaseRiskScore:          0.7,
				AssessmentFrequencyDays: 180,
			},
			{
				ID:                     "FINANCIAL",
				Name:                   "Financial Services",
				Description:            "Vendors providing financial and payment services",
				RiskFactors:            []string{"PCI DSS compliance", "Financial stability", "Regulatory compliance", "Transaction security"},
				BaseRiskScore:          0.8,
				AssessmentFrequencyDays: 90,
			},
		}

		return c.JSON(categories)
	})

	// Search and Filter Routes
	app.Post("/api/v1/compliance/vendor/search", func(c *fiber.Ctx) error {
		var request struct {
			Query      string   `json:"query"`
			Categories []string `json:"categories"`
			RiskLevels []string `json:"risk_levels"`
			Status     []string `json:"status"`
		}
		
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		vendors, err := vendorSvc.ListVendors(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Filter vendors based on search criteria
		var filteredVendors []*types.Vendor
		for _, vendor := range vendors {
			match := true
			
			// Filter by query (name and description)
			if request.Query != "" {
				nameMatch := contains(vendor.Name, request.Query)
				descMatch := contains(vendor.Description, request.Query)
				if !nameMatch && !descMatch {
					match = false
				}
			}
			
			// Filter by categories
			if len(request.Categories) > 0 {
				categoryMatch := false
				for _, category := range request.Categories {
					if vendor.Category.ID == category {
						categoryMatch = true
						break
					}
				}
				if !categoryMatch {
					match = false
				}
			}
			
			// Filter by risk levels
			if len(request.RiskLevels) > 0 && vendor.RiskLevel != nil {
				riskMatch := false
				for _, riskLevel := range request.RiskLevels {
					if string(*vendor.RiskLevel) == riskLevel {
						riskMatch = true
						break
					}
				}
				if !riskMatch {
					match = false
				}
			}
			
			// Filter by status
			if len(request.Status) > 0 {
				statusMatch := false
				for _, status := range request.Status {
					if string(vendor.Status) == status {
						statusMatch = true
						break
					}
				}
				if !statusMatch {
					match = false
				}
			}
			
			if match {
				filteredVendors = append(filteredVendors, vendor)
			}
		}

		return c.JSON(filteredVendors)
	})

	// Bulk Operations Routes
	app.Post("/api/v1/compliance/vendor/bulk", func(c *fiber.Ctx) error {
		var request struct {
			Action     string   `json:"action"`
			VendorIDs  []string `json:"vendor_ids"`
			Parameters map[string]interface{} `json:"parameters"`
		}
		
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		var results []map[string]interface{}
		
		for _, vendorID := range request.VendorIDs {
			result := map[string]interface{}{
				"vendor_id": vendorID,
				"success":   false,
				"message":   "",
			}
			
			switch request.Action {
			case "assess":
				_, err := vendorSvc.ConductVendorAssessment(c.Context(), vendorID)
				if err != nil {
					result["message"] = err.Error()
				} else {
					result["success"] = true
					result["message"] = "Vendor assessment completed successfully"
				}
			case "update_status":
				status, ok := request.Parameters["status"].(string)
				if !ok {
					result["message"] = "Status parameter required"
				} else {
					// Mock status update - in real implementation, this would call the service
					result["success"] = true
					result["message"] = fmt.Sprintf("Vendor status updated to %s", status)
				}
			default:
				result["message"] = "Unknown action"
			}
			
			results = append(results, result)
		}

		return c.JSON(fiber.Map{
			"results": results,
		})
	})

	// Risk Assessment Templates Routes
	app.Get("/api/v1/compliance/vendor/assessment-templates", func(c *fiber.Ctx) error {
		templates := []map[string]interface{}{
			{
				"id":          "basic_assessment",
				"name":        "Basic Vendor Assessment",
				"description": "Standard vendor risk assessment template",
				"risk_factors": []string{
					"Financial stability",
					"Operational capability",
					"Security posture",
					"Compliance adherence",
				},
				"frequency_days": 365,
			},
			{
				"id":          "high_risk_assessment",
				"name":        "High-Risk Vendor Assessment",
				"description": "Comprehensive assessment for high-risk vendors",
				"risk_factors": []string{
					"Data sovereignty",
					"Service availability",
					"Data access controls",
					"Encryption standards",
					"Physical security",
					"Disaster recovery",
					"Environmental controls",
					"Access management",
				},
				"frequency_days": 90,
			},
		}

		return c.JSON(templates)
	})

	// Vendor Contract Management Routes
	app.Get("/api/v1/compliance/vendor/:id/contracts", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		vendor, err := vendorSvc.GetVendor(c.Context(), vendorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Mock contract data
		contracts := []map[string]interface{}{
			{
				"id":           fmt.Sprintf("contract-%s-1", vendorID),
				"vendor_id":    vendorID,
				"vendor_name":  vendor.Name,
				"start_date":   vendor.ContractStartDate,
				"end_date":     vendor.ContractEndDate,
				"services":     vendor.ServicesOffered,
				"value":        "100000",
				"currency":     "USD",
				"status":       "active",
				"renewal_date": vendor.ContractEndDate,
			},
		}

		return c.JSON(contracts)
	})

	// Vendor Document Management Routes
	app.Get("/api/v1/compliance/vendor/:id/documents", func(c *fiber.Ctx) error {
		vendorID := c.Params("id")
		
		_, err := vendorSvc.GetVendor(c.Context(), vendorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Mock document data
		documents := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("doc-%s-1", vendorID),
				"vendor_id":   vendorID,
				"name":        "Security Assessment Report",
				"type":        "assessment",
				"upload_date": time.Now().UTC().AddDate(0, -1, 0),
				"status":      "current",
				"size":        "2.5MB",
			},
			{
				"id":          fmt.Sprintf("doc-%s-2", vendorID),
				"vendor_id":   vendorID,
				"name":        "Compliance Certificate",
				"type":        "certificate",
				"upload_date": time.Now().UTC().AddDate(0, -2, 0),
				"status":      "current",
				"size":        "1.2MB",
			},
		}

		return c.JSON(documents)
	})
}

// Helper function for string contains
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    s[:len(substr)] == substr || 
		    s[len(s)-len(substr):] == substr ||
		    findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
