package edge

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/constants"
)

// EdgeService fetches Microsoft Edge release metadata for all distribution channels.
//
// Data is sourced from:
//
//	https://edgeupdates.microsoft.com/api/products/{channel}
//
// Each endpoint returns a JSON array of product release objects. This service
// filters for macOS (Platform == "MacOS") and returns the latest release per channel.
type EdgeService struct {
	client client.Client
}

// NewService creates a new EdgeService.
func NewService(c client.Client) *EdgeService {
	return &EdgeService{client: c}
}

// fetchChannel fetches and parses the Edge release list for a given channel endpoint,
// then returns the latest macOS release.
func (s *EdgeService) fetchChannel(ctx context.Context, endpoint, channel string) (*EdgeRelease, error) {
	var products []edgeAPIProduct

	_, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&products).
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetch edge %s: %w", channel, err)
	}

	// The API returns an array; the first element matches the channel.
	for _, product := range products {
		for _, rel := range product.Releases {
			if rel.Platform != PlatformMacOS {
				continue
			}
			release := &EdgeRelease{
				Channel:       channel,
				Version:       rel.ReleaseVersion,
				PublishedTime: rel.PublishedTime,
			}
			for _, a := range rel.Artifacts {
				artifact := EdgeArtifact{
					ArtifactName: a.ArtifactName,
					Location:     a.Location,
					SizeInBytes:  a.SizeInBytes,
				}
				if a.HashAlgorithm == "SHA256" {
					artifact.HashSHA256 = a.Hash
				}
				release.Artifacts = append(release.Artifacts, artifact)
			}
			return release, nil
		}
	}

	return nil, fmt.Errorf("no macOS release found for edge %s channel", channel)
}

// GetAllChannelsV1 fetches the latest macOS Edge release for all four channels
// (stable, beta, dev, canary) and returns them as an aggregated response.
//
// GET https://edgeupdates.microsoft.com/api/products/{channel}
func (s *EdgeService) GetAllChannelsV1(ctx context.Context) (*EdgeAllChannelsResponse, error) {
	stable, err := s.GetStableV1(ctx)
	if err != nil {
		return nil, fmt.Errorf("stable: %w", err)
	}

	beta, err := s.GetBetaV1(ctx)
	if err != nil {
		return nil, fmt.Errorf("beta: %w", err)
	}

	dev, err := s.GetDevV1(ctx)
	if err != nil {
		return nil, fmt.Errorf("dev: %w", err)
	}

	canary, err := s.GetCanaryV1(ctx)
	if err != nil {
		return nil, fmt.Errorf("canary: %w", err)
	}

	return &EdgeAllChannelsResponse{
		Stable: stable.Release,
		Beta:   beta.Release,
		Dev:    dev.Release,
		Canary: canary.Release,
	}, nil
}

// GetStableV1 fetches the latest macOS Edge stable release.
//
// GET https://edgeupdates.microsoft.com/api/products/stable
func (s *EdgeService) GetStableV1(ctx context.Context) (*EdgeChannelResponse, error) {
	rel, err := s.fetchChannel(ctx, constants.EdgeStableEndpoint, ChannelStable)
	if err != nil {
		return nil, err
	}
	return &EdgeChannelResponse{Channel: ChannelStable, Release: rel}, nil
}

// GetBetaV1 fetches the latest macOS Edge beta release.
//
// GET https://edgeupdates.microsoft.com/api/products/beta
func (s *EdgeService) GetBetaV1(ctx context.Context) (*EdgeChannelResponse, error) {
	rel, err := s.fetchChannel(ctx, constants.EdgeBetaEndpoint, ChannelBeta)
	if err != nil {
		return nil, err
	}
	return &EdgeChannelResponse{Channel: ChannelBeta, Release: rel}, nil
}

// GetDevV1 fetches the latest macOS Edge dev release.
//
// GET https://edgeupdates.microsoft.com/api/products/dev
func (s *EdgeService) GetDevV1(ctx context.Context) (*EdgeChannelResponse, error) {
	rel, err := s.fetchChannel(ctx, constants.EdgeDevEndpoint, ChannelDev)
	if err != nil {
		return nil, err
	}
	return &EdgeChannelResponse{Channel: ChannelDev, Release: rel}, nil
}

// GetCanaryV1 fetches the latest macOS Edge canary release.
//
// GET https://edgeupdates.microsoft.com/api/products/canary
func (s *EdgeService) GetCanaryV1(ctx context.Context) (*EdgeChannelResponse, error) {
	rel, err := s.fetchChannel(ctx, constants.EdgeCanaryEndpoint, ChannelCanary)
	if err != nil {
		return nil, err
	}
	return &EdgeChannelResponse{Channel: ChannelCanary, Release: rel}, nil
}
