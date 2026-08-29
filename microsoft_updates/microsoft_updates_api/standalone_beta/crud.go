package standalone_beta

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// StandaloneBetaService fetches macOS standalone application metadata from the
// Microsoft Office CDN beta (Insider Fast) channel.
//
// Beta builds are pre-release and may change more frequently than production.
type StandaloneBetaService struct {
	client  client.Client
	baseURL string
}

// NewService creates a new StandaloneBetaService targeting the beta CDN channel.
func NewService(c client.Client) *StandaloneBetaService {
	return &StandaloneBetaService{
		client:  c,
		baseURL: constants.StandaloneBetaCDNBaseURL,
	}
}

func (s *StandaloneBetaService) fetchPackage(ctx context.Context, appID string) (*Package, error) {
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

// GetLatestV1 fetches the latest beta metadata for all known standalone applications.
//
// GET https://officecdnmac.microsoft.com/pr/{betaChannelUUID}/MacAutoupdate/{AppID}.xml
func (s *StandaloneBetaService) GetLatestV1(ctx context.Context) (*StandaloneBetaResponse, error) {
	resp := &StandaloneBetaResponse{}
	for _, appID := range AllAppIDs {
		pkg, err := s.fetchPackage(ctx, appID)
		if err != nil {
			s.client.GetLogger().Sugar().Warnf("skipping beta %s: %v", appID, err)
			continue
		}
		resp.Packages = append(resp.Packages, pkg)
	}
	return resp, nil
}

// GetPackageByApplicationIDV1 fetches beta metadata for a single application by ID.
//
// GET https://officecdnmac.microsoft.com/pr/{betaChannelUUID}/MacAutoupdate/{AppID}.xml
func (s *StandaloneBetaService) GetPackageByApplicationIDV1(ctx context.Context, appID string) (*Package, error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	return s.fetchPackage(ctx, appID)
}

// GetPackageByNameV1 fetches beta metadata for the application with the given display name.
func (s *StandaloneBetaService) GetPackageByNameV1(ctx context.Context, name string) (*Package, error) {
	if name == "" {
		return nil, fmt.Errorf("application name is required")
	}

	for appID, appName := range AppNames {
		if appName == name {
			return s.fetchPackage(ctx, appID)
		}
	}

	return nil, fmt.Errorf("application %q not found in beta app list", name)
}
