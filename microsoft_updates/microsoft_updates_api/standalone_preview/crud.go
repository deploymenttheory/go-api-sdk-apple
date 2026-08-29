package standalone_preview

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// StandalonePreviewService fetches macOS standalone application metadata from the
// Microsoft Office CDN preview (Insider Slow) channel.
type StandalonePreviewService struct {
	client  client.Client
	baseURL string
}

// NewService creates a new StandalonePreviewService targeting the preview CDN channel.
func NewService(c client.Client) *StandalonePreviewService {
	return &StandalonePreviewService{
		client:  c,
		baseURL: constants.StandalonePreviewCDNBaseURL,
	}
}

func (s *StandalonePreviewService) fetchPackage(ctx context.Context, appID string) (*Package, error) {
	url := s.baseURL + appID + ".xml"

	_, body, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationXML).
		GetBytes(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", appID, err)
	}

	var plist plistArray
	if err := xml.Unmarshal(body, &plist); err != nil {
		return nil, fmt.Errorf("parse plist for %s: %w", appID, err)
	}

	if len(plist.Items) == 0 {
		return nil, fmt.Errorf("no update entries found for %s", appID)
	}

	return plist.Items[0].toPackage(appID), nil
}

// GetLatestV1 fetches the latest preview metadata for all known standalone applications.
//
// GET https://officecdnmac.microsoft.com/pr/{previewChannelUUID}/MacAutoupdate/{AppID}.xml
func (s *StandalonePreviewService) GetLatestV1(ctx context.Context) (*StandalonePreviewResponse, error) {
	resp := &StandalonePreviewResponse{}
	for _, appID := range AllAppIDs {
		pkg, err := s.fetchPackage(ctx, appID)
		if err != nil {
			s.client.GetLogger().Sugar().Warnf("skipping preview %s: %v", appID, err)
			continue
		}
		resp.Packages = append(resp.Packages, pkg)
	}
	return resp, nil
}

// GetPackageByApplicationIDV1 fetches preview metadata for a single application by ID.
//
// GET https://officecdnmac.microsoft.com/pr/{previewChannelUUID}/MacAutoupdate/{AppID}.xml
func (s *StandalonePreviewService) GetPackageByApplicationIDV1(ctx context.Context, appID string) (*Package, error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	return s.fetchPackage(ctx, appID)
}

// GetPackageByNameV1 fetches preview metadata for the application with the given display name.
func (s *StandalonePreviewService) GetPackageByNameV1(ctx context.Context, name string) (*Package, error) {
	if name == "" {
		return nil, fmt.Errorf("application name is required")
	}

	for appID, appName := range AppNames {
		if appName == name {
			return s.fetchPackage(ctx, appID)
		}
	}

	return nil, fmt.Errorf("application %q not found in preview app list", name)
}
