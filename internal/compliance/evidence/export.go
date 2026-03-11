package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Service generates compliance evidence.
type Service struct{}

// New creates an Evidence service.
func New() *Service {
	return &Service{}
}

// ExportJSON exports compliance evidence as JSON.
func (s *Service) ExportJSON(ctx context.Context) ([]byte, error) {
	evidence := map[string]interface{}{
		"framework": "SOC2",
		"timestamp": time.Now().UTC(),
		"findings": []types.FindingSeverity{
			types.FindingSeverityLow,
			types.FindingSeverityMedium,
			types.FindingSeverityHigh,
			types.FindingSeverityCritical,
		},
		"recommendations": []types.RecommendationPriority{
			types.RecommendationPriorityLow,
			types.RecommendationPriorityMedium,
			types.RecommendationPriorityHigh,
			types.RecommendationPriorityCritical,
		},
	}
	return json.MarshalIndent(evidence, "", "  ")
}

// ExportPDF exports compliance evidence as PDF.
func (s *Service) ExportPDF(ctx context.Context) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "SOC2 Evidence Report")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Generated: %s", time.Now().UTC().Format(time.RFC3339)))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Findings:")
	pdf.Ln(8)
	findings := []string{"Low", "Medium", "High", "Critical"}
	for _, f := range findings {
		pdf.Cell(40, 8, fmt.Sprintf("- %s", f))
		pdf.Ln(6)
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	log.Println("evidence: PDF exported")
	return buf.Bytes(), nil
}
