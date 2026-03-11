package pci

import (
	"context"
	"log"
	"time"
)

// Service runs PCI DSS compliance assessments.
type Service struct{}

// New creates a PCI DSS service.
func New() *Service {
	return &Service{}
}

// ScanCardholderEnvironment scans cardholder data environments.
func (s *Service) ScanCardholderEnvironment(ctx context.Context) error {
	log.Println("pci: scanning cardholder environment")
	// TODO: run vulnerability scans, assess controls
	return nil
}
