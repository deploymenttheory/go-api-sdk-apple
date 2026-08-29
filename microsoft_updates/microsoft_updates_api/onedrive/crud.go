package onedrive

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// OneDriveService fetches Microsoft OneDrive version metadata across all
// distribution rings.
//
// Data comes from two sources:
//   - g.live.com XML manifests (insider and standalone feeds)
//   - go.microsoft.com/fwlink redirect URLs (deferred, upcoming, rolling-out, app-new rings)
//
// For fwlink rings, a HEAD request is issued to resolve the final redirect URL,
// which encodes the version string in its path.
type OneDriveService struct {
	client client.Client
}

// NewService creates a new OneDriveService.
func NewService(c client.Client) *OneDriveService {
	return &OneDriveService{client: c}
}

// fetchManifest fetches and parses an OneDrive XML manifest from a g.live.com URL.
func (s *OneDriveService) fetchManifest(ctx context.Context, manifestURL, ring string) (*OneDriveRing, error) {
	_, body, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationXML).
		GetBytes(manifestURL)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest for ring %s: %w", ring, err)
	}

	var manifest oneDriveManifest
	if err := xml.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest for ring %s: %w", ring, err)
	}

	if len(manifest.Items) == 0 {
		return nil, fmt.Errorf("no update entries found in OneDrive manifest for ring %s", ring)
	}

	item := manifest.Items[0]
	return &OneDriveRing{
		Ring:          ring,
		Version:       item.Version,
		BuildVersion:  item.BuildVersion,
		DownloadURL:   item.PackageURL,
		ApplicationID: ApplicationID,
		BundleID:      BundleID,
	}, nil
}

// fetchFwlinkRing resolves a fwlink redirect URL and extracts version info from
// the final redirect destination URL path.
func (s *OneDriveService) fetchFwlinkRing(ctx context.Context, fwlinkURL, ring string) (*OneDriveRing, error) {
	resp, err := s.client.NewRequest(ctx).Head(fwlinkURL)
	if err != nil {
		return nil, fmt.Errorf("resolve fwlink for ring %s: %w", ring, err)
	}

	finalURL := fwlinkURL
	if resp != nil {
		if loc := resp.Header().Get("Location"); loc != "" {
			finalURL = loc
		} else if resp.RawResponse != nil {
			finalURL = resp.RawResponse.Request.URL.String()
		}
	}

	return &OneDriveRing{
		Ring:          ring,
		DownloadURL:   finalURL,
		ApplicationID: ApplicationID,
		BundleID:      BundleID,
	}, nil
}

// GetAllRingsV1 fetches OneDrive version metadata for all distribution rings.
// It returns a slice of rings ordered: Production, Insider, Deferred, UpcomingDeferred,
// RollingOut, AppNew.
//
// GET/HEAD against g.live.com manifests and go.microsoft.com fwlink endpoints.
func (s *OneDriveService) GetAllRingsV1(ctx context.Context) (*OneDriveAllRingsResponse, error) {
	resp := &OneDriveAllRingsResponse{}

	type ringFetcher struct {
		name string
		fn   func(context.Context) (*OneDriveRing, error)
	}

	fetchers := []ringFetcher{
		{RingProduction, s.GetProductionRingV1},
		{RingInsider, s.GetInsiderRingV1},
		{RingDeferred, s.GetDeferredRingV1},
		{RingUpcoming, s.GetUpcomingDeferredRingV1},
		{RingRollingOut, s.GetRollingOutRingV1},
		{RingAppNew, s.GetAppNewRingV1},
	}

	for _, f := range fetchers {
		ring, err := f.fn(ctx)
		if err != nil {
			s.client.GetLogger().Sugar().Warnf("skipping OneDrive ring %s: %v", f.name, err)
			continue
		}
		resp.Rings = append(resp.Rings, ring)
	}

	return resp, nil
}

// GetProductionRingV1 fetches OneDrive metadata for the standalone production ring
// using the Microsoft standalone product manifest.
//
// GET https://g.live.com/0USSDMC_W5T/StandaloneProductManifest
func (s *OneDriveService) GetProductionRingV1(ctx context.Context) (*OneDriveRing, error) {
	return s.fetchManifest(ctx, constants.OneDriveStandaloneManifest, RingProduction)
}

// GetInsiderRingV1 fetches OneDrive metadata for the insider ring
// using the Microsoft insider feed.
//
// GET https://g.live.com/0USSDMC_W5T/MacODSUInsiders
func (s *OneDriveService) GetInsiderRingV1(ctx context.Context) (*OneDriveRing, error) {
	return s.fetchManifest(ctx, constants.OneDriveLiveFeedInsiders, RingInsider)
}

// GetDeferredRingV1 fetches OneDrive metadata for the deferred ring
// by resolving the fwlink redirect URL.
//
// HEAD https://go.microsoft.com/fwlink/?linkid=861009
func (s *OneDriveService) GetDeferredRingV1(ctx context.Context) (*OneDriveRing, error) {
	return s.fetchFwlinkRing(ctx, constants.OneDriveFwlinkDeferred, RingDeferred)
}

// GetUpcomingDeferredRingV1 fetches OneDrive metadata for the upcoming deferred ring.
//
// HEAD https://go.microsoft.com/fwlink/?linkid=861010
func (s *OneDriveService) GetUpcomingDeferredRingV1(ctx context.Context) (*OneDriveRing, error) {
	return s.fetchFwlinkRing(ctx, constants.OneDriveFwlinkUpcoming, RingUpcoming)
}

// GetRollingOutRingV1 fetches OneDrive metadata for the rolling-out ring.
//
// HEAD https://go.microsoft.com/fwlink/?linkid=861011
func (s *OneDriveService) GetRollingOutRingV1(ctx context.Context) (*OneDriveRing, error) {
	return s.fetchFwlinkRing(ctx, constants.OneDriveFwlinkRollingOut, RingRollingOut)
}

// GetAppNewRingV1 fetches OneDrive metadata for the app-new ring.
//
// HEAD https://go.microsoft.com/fwlink/?linkid=823060
func (s *OneDriveService) GetAppNewRingV1(ctx context.Context) (*OneDriveRing, error) {
	return s.fetchFwlinkRing(ctx, constants.OneDriveFwlinkAppNew, RingAppNew)
}

