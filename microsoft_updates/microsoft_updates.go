package microsoft_updates

import (
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/client"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/appstore_ios"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/appstore_macos"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/cve_history"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/edge"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/onedrive"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/standalone"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/standalone_beta"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/standalone_preview"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/update_history"
)

// Client is the main entry point for the Microsoft Updates SDK.
// It provides access to nine services spanning official Microsoft endpoints
// and the Apple App Store:
//
//   - Standalone: macOS standalone apps from the production Office CDN channel
//   - StandaloneBeta: beta (Insider Fast) channel builds
//   - StandalonePreview: preview (Insider Slow) channel builds
//   - Edge: Microsoft Edge across all four channels (stable/beta/dev/canary)
//   - OneDrive: OneDrive across all distribution rings
//   - AppStoreMacOS: Microsoft apps in the macOS App Store via iTunes Search API
//   - AppStoreIOS: Microsoft apps in the iOS App Store via iTunes Search API
//   - UpdateHistory: Office for Mac update history (HTML scrape)
//   - CVEHistory: Office for Mac CVE/security release notes (HTML scrape)
type Client struct {
	transport           *client.Transport
	MicrosoftUpdatesAPI *MicrosoftUpdatesAPIClient
}

// MicrosoftUpdatesAPIClient groups all Microsoft Updates sub-services.
type MicrosoftUpdatesAPIClient struct {
	Standalone        *standalone.StandaloneService
	StandaloneBeta    *standalone_beta.StandaloneBetaService
	StandalonePreview *standalone_preview.StandalonePreviewService
	Edge              *edge.EdgeService
	OneDrive          *onedrive.OneDriveService
	AppStoreMacOS     *appstore_macos.AppStoreMacOSService
	AppStoreIOS       *appstore_ios.AppStoreIOSService
	UpdateHistory     *update_history.UpdateHistoryService
	CVEHistory        *cve_history.CVEHistoryService
}

// NewClient creates a new Microsoft Updates client with optional configuration.
//
// Example:
//
//	c, err := microsoft_updates.NewClient(
//	    microsoft_updates.WithTimeout(15 * time.Second),
//	    microsoft_updates.WithLogger(logger),
//	)
func NewClient(options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransport(options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		MicrosoftUpdatesAPI: &MicrosoftUpdatesAPIClient{
			Standalone:        standalone.NewService(transport),
			StandaloneBeta:    standalone_beta.NewService(transport),
			StandalonePreview: standalone_preview.NewService(transport),
			Edge:              edge.NewService(transport),
			OneDrive:          onedrive.NewService(transport),
			AppStoreMacOS:     appstore_macos.NewService(transport),
			AppStoreIOS:       appstore_ios.NewService(transport),
			UpdateHistory:     update_history.NewService(transport),
			CVEHistory:        cve_history.NewService(transport),
		},
	}, nil
}

// NewDefaultClient creates a new Microsoft Updates client with default settings
// (30s timeout, 3 retries, no logging).
func NewDefaultClient() (*Client, error) {
	return NewClient()
}

// Close releases resources held by the client.
func (c *Client) Close() error {
	return c.transport.Close()
}
