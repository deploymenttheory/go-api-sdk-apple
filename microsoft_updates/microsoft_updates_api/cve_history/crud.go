package cve_history

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// CVEHistoryService fetches and parses Office for Mac security release notes
// from the Microsoft Learn documentation site.
//
// Data is sourced from:
//
//	https://learn.microsoft.com/en-us/officeupdates/release-notes-office-for-mac
//
// The page is scraped using golang.org/x/net/html. Each release heading becomes
// a CVEEntry with the release date, version, and the list of CVE identifiers.
type CVEHistoryService struct {
	client client.Client
}

// NewService creates a new CVEHistoryService.
func NewService(c client.Client) *CVEHistoryService {
	return &CVEHistoryService{client: c}
}

// GetCVEHistoryV1 fetches and parses the complete Office for Mac CVE/security
// release history. Entries are ordered as they appear on the page (newest first).
//
// GET https://learn.microsoft.com/en-us/officeupdates/release-notes-office-for-mac
func (s *CVEHistoryService) GetCVEHistoryV1(ctx context.Context) (*CVEHistoryResponse, error) {
	_, body, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextHTML).
		GetBytes(constants.CVEHistoryURL)
	if err != nil {
		return nil, fmt.Errorf("fetch CVE history: %w", err)
	}

	entries, err := parseCVEHistory(body)
	if err != nil {
		return nil, fmt.Errorf("parse CVE history: %w", err)
	}

	return &CVEHistoryResponse{Entries: entries}, nil
}
