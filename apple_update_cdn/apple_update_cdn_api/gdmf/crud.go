package gdmf

import (
	"context"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/client"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/constants"
	"resty.dev/v3"
)

// GDMFService handles communication with Apple's Global Device Management Feed API.
//
// The GDMF API is Apple's authoritative source for currently-signed firmware
// versions. It covers macOS, iOS/iPadOS/watchOS, and visionOS, and includes
// posting dates, expiration dates, and the list of supported hardware for each
// version. It does not provide IPSW download URLs.
//
// GDMF API: https://gdmf.apple.com/v2/pmv
type GDMFService struct {
	client client.Client
}

// NewService creates a new GDMF service.
func NewService(c client.Client) *GDMFService {
	return &GDMFService{client: c}
}

// GetPublicVersionsV2 fetches all currently-signed public firmware versions
// from Apple's GDMF API. The response includes three sets:
//
//   - PublicAssetSets: publicly released, fully signed versions
//   - AssetSets: additional versions (may include seed programme entries)
//   - PublicBackgroundSecurityImprovements: rapid security responses
//
// Each set is broken down by platform (macOS, iOS, visionOS).
//
// GET https://gdmf.apple.com/v2/pmv
func (s *GDMFService) GetPublicVersionsV2(ctx context.Context) (*GDMFResponse, *resty.Response, error) {
	var result GDMFResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointGDMFVersions)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
