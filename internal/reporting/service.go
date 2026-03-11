package reporting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/barca-strategos/phoenix/internal/case"
	"github.com/barca-strategos/phoenix/internal/risk"
	"github.com/barca-strategos/phoenix/internal/asset"
)

// Service generates scheduled and on-demand reports.
type Service struct {
	caseSvc  *case.Service
	riskSvc  *risk.Service
	assetSvc *asset.Service
}

// New creates a Reporting service.
func New(caseSvc *case.Service, riskSvc *risk.Service, assetSvc *asset.Service) *Service {
	return &Service{
		caseSvc:  caseSvc,
		riskSvc:  riskSvc,
		assetSvc: assetSvc,
	}
}

// GenerateExecutiveSummary creates a high-level report.
func (s *Service) GenerateExecutiveSummary(ctx context.Context) ([]byte, error) {
	cases, _ := s.caseSvc.ListCases(ctx)
	risks, _ := s.riskSvc.ListRisks(ctx)
	assets, _ := s.assetSvc.ListAssets(ctx)

	summary := map[string]interface{}{
		"generated_at": time.Now().UTC(),
		"cases": map[string]int{
			"total":    len(cases),
			"open":     countCasesByStatus(cases, "new") + countCasesByStatus(cases, "in_progress"),
			"resolved": countCasesByStatus(cases, "resolved") + countCasesByStatus(cases, "closed"),
		},
		"risks": map[string]interface{}{
			"total":      len(risks),
			"avg_score":  avgRiskScore(risks),
			"high_risk":  countRisksAbove(risks, 15),
		},
		"assets": map[string]interface{}{
			"total":        len(assets),
			"criticality": countAssetsByCriticality(assets),
		},
	}
	return json.MarshalIndent(summary, "", "  ")
}

// GeneratePDFReport creates a PDF report.
func (s *Service) GeneratePDFReport(ctx context.Context, reportType string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("%s Report", reportType))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Generated: %s", time.Now().UTC().Format(time.RFC3339)))
	pdf.Ln(10)

	switch reportType {
	case "cases":
		s.addCasesSection(pdf)
	case "risks":
		s.addRisksSection(pdf)
	case "assets":
		s.addAssetsSection(pdf)
	default:
		pdf.Cell(40, 10, "Unknown report type")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *Service) addCasesSection(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Cases")
	pdf.Ln(8)
	cases, _ := s.caseSvc.ListCases(context.Background())
	for _, c := range cases {
		pdf.Cell(40, 8, fmt.Sprintf("- %s (%s)", c.Title, c.Status))
		pdf.Ln(6)
	}
}

func (s *Service) addRisksSection(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Risks")
	pdf.Ln(8)
	risks, _ := s.riskSvc.ListRisks(context.Background())
	for _, r := range risks {
		pdf.Cell(40, 8, fmt.Sprintf("- %s (Score: %.1f)", r.Title, r.Score))
		pdf.Ln(6)
	}
}

func (s *Service) addAssetsSection(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Assets")
	pdf.Ln(8)
	assets, _ := s.assetSvc.ListAssets(context.Background())
	for _, a := range assets {
		pdf.Cell(40, 8, fmt.Sprintf("- %s (%s)", a.Name, a.Criticality))
		pdf.Ln(6)
	}
}

// Helper functions
func countCasesByStatus(cases []*case.Case, status string) int {
	count := 0
	for _, c := range cases {
		if string(c.Status) == status {
			count++
		}
	}
	return count
}

func avgRiskScore(risks []*risk.Risk) float64 {
	if len(risks) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range risks {
		sum += r.Score
	}
	return sum / float64(len(risks))
}

func countRisksAbove(risks []*risk.Risk, threshold float64) int {
	count := 0
	for _, r := range risks {
		if r.Score > threshold {
			count++
		}
	}
	return count
}

func countAssetsByCriticality(assets []*asset.Asset) map[string]int {
	m := make(map[string]int)
	for _, a := range assets {
		m[string(a.Criticality)]++
	}
	return m
}
