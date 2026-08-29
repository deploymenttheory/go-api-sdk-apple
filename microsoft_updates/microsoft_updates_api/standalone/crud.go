package standalone

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// StandaloneService fetches macOS standalone application metadata from the
// Microsoft Office CDN production channel.
//
// The CDN hosts per-application Apple plist XML files at:
//
//	https://officecdnmac.microsoft.com/pr/{channelUUID}/MacAutoupdate/{AppID}.xml
//
// Each plist contains version, download URL, minimum OS, and checksum data.
type StandaloneService struct {
	client  client.Client
	baseURL string
}

// NewService creates a new StandaloneService targeting the production CDN channel.
func NewService(c client.Client) *StandaloneService {
	return &StandaloneService{
		client:  c,
		baseURL: constants.StandaloneCDNBaseURL,
	}
}

// newServiceWithBaseURL creates a StandaloneService with a custom CDN base URL.
// Used internally by the beta and preview constructors.
func newServiceWithBaseURL(c client.Client, baseURL string) *StandaloneService {
	return &StandaloneService{client: c, baseURL: baseURL}
}

// fetchPackage fetches and parses the plist XML for a single application ID.
func (s *StandaloneService) fetchPackage(ctx context.Context, appID string) (*Package, error) {
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

// GetLatestV1 fetches the latest metadata for all known standalone applications
// from the Microsoft Office CDN production channel. It returns one Package per
// application ID, collected from individual per-app plist XML endpoints.
//
// GET https://officecdnmac.microsoft.com/pr/{channelUUID}/MacAutoupdate/{AppID}.xml
func (s *StandaloneService) GetLatestV1(ctx context.Context) (*StandaloneResponse, error) {
	resp := &StandaloneResponse{}
	for _, appID := range AllAppIDs {
		pkg, err := s.fetchPackage(ctx, appID)
		if err != nil {
			s.client.GetLogger().Sugar().Warnf("skipping %s: %v", appID, err)
			continue
		}
		resp.Packages = append(resp.Packages, pkg)
	}
	return resp, nil
}

// GetPackageByApplicationIDV1 fetches the latest metadata for a single application
// identified by its Microsoft CDN application ID (e.g. "MSWD2019").
//
// GET https://officecdnmac.microsoft.com/pr/{channelUUID}/MacAutoupdate/{AppID}.xml
func (s *StandaloneService) GetPackageByApplicationIDV1(ctx context.Context, appID string) (*Package, error) {
	if appID == "" {
		return nil, fmt.Errorf("application ID is required")
	}
	return s.fetchPackage(ctx, appID)
}

// GetPackageByNameV1 fetches the latest metadata for the application with the
// given human-readable display name (e.g. "Microsoft Word"). The name lookup
// is performed against the AppNames map in constants.go.
//
// Returns an error if no application with that name is found.
func (s *StandaloneService) GetPackageByNameV1(ctx context.Context, name string) (*Package, error) {
	if name == "" {
		return nil, fmt.Errorf("application name is required")
	}

	for appID, appName := range AppNames {
		if appName == name {
			return s.fetchPackage(ctx, appID)
		}
	}

	return nil, fmt.Errorf("application %q not found in standalone app list", name)
}
