package appstore_macos

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// AppStoreMacOSService fetches Microsoft application metadata from the macOS App Store
// via the Apple iTunes Search API.
//
// Each application is queried individually using its search term. Results include
// version, release date, minimum OS, release notes, and App Store URL.
type AppStoreMacOSService struct {
	client client.Client
}

// NewService creates a new AppStoreMacOSService.
func NewService(c client.Client) *AppStoreMacOSService {
	return &AppStoreMacOSService{client: c}
}

// searchApp queries the iTunes Search API for a single application by name,
// filtering for macOS software.
func (s *AppStoreMacOSService) searchApp(ctx context.Context, term string) (*AppStoreResponse, error) {
	params := s.client.QueryBuilder().
		AddString("term", term).
		AddString("country", iTunesCountry).
		AddString("entity", iTunesEntityMacOS).
		Build()

	var result AppStoreResponse

	_, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetQueryParams(params).
		SetResult(&result).
		Get(constants.AppStoreSearchEndpoint)
	if err != nil {
		return nil, fmt.Errorf("app store search for %q: %w", term, err)
	}

	return &result, nil
}

// GetAllAppsV1 queries the iTunes Search API for all known Microsoft macOS App Store
// applications and returns the combined results.
//
// GET https://itunes.apple.com/search?term={app}&country=us&entity=macSoftware
func (s *AppStoreMacOSService) GetAllAppsV1(ctx context.Context) (*AppStoreResponse, error) {
	combined := &AppStoreResponse{}

	for _, term := range AppSearchTerms {
		result, err := s.searchApp(ctx, term)
		if err != nil {
			s.client.GetLogger().Sugar().Warnf("skipping app store search for %q: %v", term, err)
			continue
		}
		// Deduplicate by BundleID.
		seen := make(map[string]bool, len(combined.Results))
		for _, existing := range combined.Results {
			seen[existing.BundleID] = true
		}
		for _, entry := range result.Results {
			if !seen[entry.BundleID] {
				combined.Results = append(combined.Results, entry)
				seen[entry.BundleID] = true
			}
		}
	}

	combined.ResultCount = len(combined.Results)
	return combined, nil
}

// GetAppByNameV1 queries the iTunes Search API for a single Microsoft macOS app
// by its display name (e.g. "Microsoft Word").
//
// GET https://itunes.apple.com/search?term={name}&country=us&entity=macSoftware
func (s *AppStoreMacOSService) GetAppByNameV1(ctx context.Context, name string) (*AppEntry, error) {
	if name == "" {
		return nil, fmt.Errorf("app name is required")
	}

	result, err := s.searchApp(ctx, name)
	if err != nil {
		return nil, err
	}

	for i := range result.Results {
		if result.Results[i].TrackName == name {
			return &result.Results[i], nil
		}
	}

	return nil, fmt.Errorf("app %q not found in macOS App Store results", name)
}

// GetAppByBundleIDV1 fetches macOS App Store metadata for the application with the
// given bundle identifier (e.g. "com.microsoft.Word").
//
// GET https://itunes.apple.com/search?term={bundleId}&country=us&entity=macSoftware
func (s *AppStoreMacOSService) GetAppByBundleIDV1(ctx context.Context, bundleID string) (*AppEntry, error) {
	if bundleID == "" {
		return nil, fmt.Errorf("bundle ID is required")
	}

	result, err := s.searchApp(ctx, bundleID)
	if err != nil {
		return nil, err
	}

	for i := range result.Results {
		if result.Results[i].BundleID == bundleID {
			return &result.Results[i], nil
		}
	}

	return nil, fmt.Errorf("app with bundle ID %q not found in macOS App Store results", bundleID)
}
