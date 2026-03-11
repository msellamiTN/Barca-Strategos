package v1

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"barca-strategos/internal/compliance/policy"
	"barca-strategos/pkg/types"
)

func RegisterPolicyRoutes(app *fiber.App, policySvc *policy.Service, broadcast func([]byte)) {
	policySvc.SetBroadcast(broadcast)

	// Policy Management Routes
	app.Post("/api/v1/compliance/policy", func(c *fiber.Ctx) error {
		var request types.PolicyRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		policy, err := policySvc.CreatePolicy(c.Context(), &request)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(policy)
	})

	app.Get("/api/v1/compliance/policy", func(c *fiber.Ctx) error {
		policies, err := policySvc.ListPolicies(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(policies)
	})

	app.Get("/api/v1/compliance/policy/:id", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		policy, err := policySvc.GetPolicy(c.Context(), policyID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(policy)
	})

	app.Put("/api/v1/compliance/policy/:id", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		
		var request types.PolicyUpdateRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		err := policySvc.UpdatePolicy(c.Context(), policyID, &request)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Policy updated successfully",
		})
	})

	app.Post("/api/v1/compliance/policy/:id/approve", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		
		var request struct {
			ApproverID string  `json:"approver_id"`
			Comments   *string `json:"comments,omitempty"`
		}
		
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		err := policySvc.ApprovePolicy(c.Context(), policyID, request.ApproverID, request.Comments)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Policy approval processed successfully",
		})
	})

	app.Post("/api/v1/compliance/policy/:id/submit", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		
		var request struct {
			Approvers []string `json:"approvers"`
		}
		
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		err := policySvc.SubmitForApproval(c.Context(), policyID, request.Approvers)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Policy submitted for approval successfully",
		})
	})

	app.Post("/api/v1/compliance/policy/:id/publish", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		
		err := policySvc.PublishPolicy(c.Context(), policyID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Policy published successfully",
		})
	})

	app.Post("/api/v1/compliance/policy/:id/archive", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		
		err := policySvc.ArchivePolicy(c.Context(), policyID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Policy archived successfully",
		})
	})

	// Compliance Management Routes
	app.Get("/api/v1/compliance/policy/:id/compliance", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		compliance, err := policySvc.GetPolicyCompliance(c.Context(), policyID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(compliance)
	})

	app.Post("/api/v1/compliance/policy/report", func(c *fiber.Ctx) error {
		var scope types.PolicyScope
		if err := json.Unmarshal(c.Body(), &scope); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		report, err := policySvc.GenerateComplianceReport(c.Context(), &scope)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(report)
	})

	// Statistics Routes
	app.Get("/api/v1/compliance/policy/stats", func(c *fiber.Ctx) error {
		stats, err := policySvc.GetPolicyStats(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(stats)
	})

	// Category Management Routes
	app.Get("/api/v1/compliance/policy/categories", func(c *fiber.Ctx) error {
		// Return predefined categories
		categories := []types.PolicyCategory{
			{
				ID:               "SECURITY",
				Name:             "Security Policies",
				Description:      "Policies related to information security",
				Subcategories:    []string{"Access Control", "Data Protection", "Network Security", "Incident Response"},
				ApprovalRequired: true,
				ReviewPeriodDays: 365,
			},
			{
				ID:               "COMPLIANCE",
				Name:             "Compliance Policies",
				Description:      "Policies related to regulatory compliance",
				Subcategories:    []string{"Regulatory Compliance", "Audit Requirements", "Reporting Standards", "Documentation"},
				ApprovalRequired: true,
				ReviewPeriodDays: 180,
			},
			{
				ID:               "OPERATIONAL",
				Name:             "Operational Policies",
				Description:      "Policies related to business operations",
				Subcategories:    []string{"Business Continuity", "Disaster Recovery", "Change Management", "Service Management"},
				ApprovalRequired: false,
				ReviewPeriodDays: 730,
			},
			{
				ID:               "HR",
				Name:             "HR Policies",
				Description:      "Policies related to human resources",
				Subcategories:    []string{"Employee Conduct", "Training and Awareness", "Privacy and Data Handling", "Remote Work"},
				ApprovalRequired: true,
				ReviewPeriodDays: 365,
			},
		}

		return c.JSON(categories)
	})

	// Workflow Management Routes
	app.Get("/api/v1/compliance/policy/:id/workflows", func(c *fiber.Ctx) error {
		policyID := c.Params("id")
		
		// Mock workflow data - in real implementation, this would query the service
		workflows := []types.ApprovalWorkflow{
			{
				ID:        fmt.Sprintf("wf-%s", policyID),
				PolicyID:  policyID,
				Approvers: []types.Approver{
					{
						ID:    "approver1",
						Name:  "John Doe",
						Email: "john.doe@company.com",
						Role:  "Policy Committee",
					},
					{
						ID:    "approver2",
						Name:  "Jane Smith",
						Email: "jane.smith@company.com",
						Role:  "Compliance Officer",
					},
				},
				CreatedAt: c.Context().Time(),
				ExpiresAt: c.Context().Time().Add(7 * 24 * 3600 * 1000000000), // 7 days in nanoseconds
				Status:    "pending",
			},
		}

		return c.JSON(workflows)
	})

	// Search and Filter Routes
	app.Post("/api/v1/compliance/policy/search", func(c *fiber.Ctx) error {
		var request struct {
			Query      string   `json:"query"`
			Categories []string `json:"categories"`
			Status     []string `json:"status"`
			Tags       []string `json:"tags"`
		}
		
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		policies, err := policySvc.ListPolicies(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Filter policies based on search criteria
		var filteredPolicies []*types.Policy
		for _, policy := range policies {
			match := true
			
			// Filter by query (title and description)
			if request.Query != "" {
				titleMatch := contains(policy.Title, request.Query)
				descMatch := contains(policy.Description, request.Query)
				if !titleMatch && !descMatch {
					match = false
				}
			}
			
			// Filter by categories
			if len(request.Categories) > 0 {
				categoryMatch := false
				for _, category := range request.Categories {
					if policy.Category.ID == category {
						categoryMatch = true
						break
					}
				}
				if !categoryMatch {
					match = false
				}
			}
			
			// Filter by status
			if len(request.Status) > 0 {
				statusMatch := false
				for _, status := range request.Status {
					if string(policy.Status) == status {
						statusMatch = true
						break
					}
				}
				if !statusMatch {
					match = false
				}
			}
			
			// Filter by tags
			if len(request.Tags) > 0 {
				tagMatch := false
				for _, tag := range request.Tags {
					for _, policyTag := range policy.Tags {
						if policyTag == tag {
							tagMatch = true
							break
						}
					}
					if tagMatch {
						break
					}
				}
				if !tagMatch {
					match = false
				}
			}
			
			if match {
				filteredPolicies = append(filteredPolicies, policy)
			}
		}

		return c.JSON(filteredPolicies)
	})

	// Bulk Operations Routes
	app.Post("/api/v1/compliance/policy/bulk", func(c *fiber.Ctx) error {
		var request struct {
			Action     string   `json:"action"`
			PolicyIDs  []string `json:"policy_ids"`
			Parameters map[string]interface{} `json:"parameters"`
		}
		
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		var results []map[string]interface{}
		
		for _, policyID := range request.PolicyIDs {
			result := map[string]interface{}{
				"policy_id": policyID,
				"success":   false,
				"message":   "",
			}
			
			switch request.Action {
			case "archive":
				err := policySvc.ArchivePolicy(c.Context(), policyID)
				if err != nil {
					result["message"] = err.Error()
				} else {
					result["success"] = true
					result["message"] = "Policy archived successfully"
				}
			case "publish":
				err := policySvc.PublishPolicy(c.Context(), policyID)
				if err != nil {
					result["message"] = err.Error()
				} else {
					result["success"] = true
					result["message"] = "Policy published successfully"
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
