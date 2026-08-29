package update_history

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// UpdateHistoryService fetches and parses the Office for Mac update history table
// from the Microsoft Learn documentation site.
//
// Data is sourced from:
//
//	https://learn.microsoft.com/en-us/officeupdates/update-history-office-for-mac
//
// The page is scraped using golang.org/x/net/html. Each table row becomes an
// UpdateHistoryEntry with the release date, version, and download URLs for the
// suite installer and individual app updaters.
type UpdateHistoryService struct {
	client client.Client
}

// NewService creates a new UpdateHistoryService.
func NewService(c client.Client) *UpdateHistoryService {
	return &UpdateHistoryService{client: c}
}

// GetUpdateHistoryV1 fetches and parses the complete Office for Mac update history.
// It returns entries ordered as they appear on the Microsoft Learn page (newest first).
//
// GET https://learn.microsoft.com/en-us/officeupdates/update-history-office-for-mac
func (s *UpdateHistoryService) GetUpdateHistoryV1(ctx context.Context) (*UpdateHistoryResponse, error) {
	_, body, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextHTML).
		GetBytes(constants.UpdateHistoryURL)
	if err != nil {
		return nil, fmt.Errorf("fetch update history: %w", err)
	}

	entries, err := parseUpdateHistory(body)
	if err != nil {
		return nil, fmt.Errorf("parse update history: %w", err)
	}

	return &UpdateHistoryResponse{Entries: entries}, nil
}
