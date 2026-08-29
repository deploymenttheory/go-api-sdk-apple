package appstore_ios

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// AppStoreIOSService fetches Microsoft application metadata from the iOS App Store
// via the Apple iTunes Search API.
type AppStoreIOSService struct {
	client client.Client
}

// NewService creates a new AppStoreIOSService.
func NewService(c client.Client) *AppStoreIOSService {
	return &AppStoreIOSService{client: c}
}

// searchApp queries the iTunes Search API for a single iOS application by name.
func (s *AppStoreIOSService) searchApp(ctx context.Context, term string) (*AppStoreIOSResponse, error) {
	params := s.client.QueryBuilder().
		AddString("term", term).
		AddString("country", iTunesCountry).
		AddString("entity", iTunesEntityIOS).
		Build()

	var result AppStoreIOSResponse

	_, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetQueryParams(params).
		SetResult(&result).
		Get(constants.AppStoreSearchEndpoint)
	if err != nil {
		return nil, fmt.Errorf("ios app store search for %q: %w", term, err)
	}

	return &result, nil
}

// GetAllAppsV1 queries the iTunes Search API for all known Microsoft iOS applications
// and returns the combined deduplicated results.
//
// GET https://itunes.apple.com/search?term={app}&country=us&entity=software
func (s *AppStoreIOSService) GetAllAppsV1(ctx context.Context) (*AppStoreIOSResponse, error) {
	combined := &AppStoreIOSResponse{}

	for _, term := range AppSearchTerms {
		result, err := s.searchApp(ctx, term)
		if err != nil {
			s.client.GetLogger().Sugar().Warnf("skipping ios app store search for %q: %v", term, err)
			continue
		}
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

// GetAppByNameV1 queries the iOS App Store for a single Microsoft app by display name.
//
// GET https://itunes.apple.com/search?term={name}&country=us&entity=software
func (s *AppStoreIOSService) GetAppByNameV1(ctx context.Context, name string) (*AppEntry, error) {
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

	return nil, fmt.Errorf("app %q not found in iOS App Store results", name)
}

// GetAppByBundleIDV1 fetches iOS App Store metadata for the application with the
// given bundle identifier.
//
// GET https://itunes.apple.com/search?term={bundleId}&country=us&entity=software
func (s *AppStoreIOSService) GetAppByBundleIDV1(ctx context.Context, bundleID string) (*AppEntry, error) {
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

	return nil, fmt.Errorf("app with bundle ID %q not found in iOS App Store results", bundleID)
}
