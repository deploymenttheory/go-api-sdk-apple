package apple_update_cdn

import (
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/apple_update_cdn_api/cdn"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/apple_update_cdn_api/firmware"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/apple_update_cdn_api/gdmf"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/client"
)

// Client is the main entry point for the Apple Update CDN SDK.
// It provides access to three services:
//   - Firmware: discovers macOS IPSW restore firmware via the ipsw.me API
//   - GDMF: queries Apple's official signed-version feed (gdmf.apple.com)
//   - CDN: parses Apple CDN URLs and resolves file metadata via HEAD requests
type Client struct {
	transport         *client.Transport
	AppleUpdateCDNAPI *AppleUpdateCDNAPIClient
}

// AppleUpdateCDNAPIClient groups all Apple Update CDN services.
type AppleUpdateCDNAPIClient struct {
	Firmware *firmware.FirmwareService
	GDMF     *gdmf.GDMFService
	CDN      *cdn.CDNService
}

// NewClient creates a new Apple Update CDN client with optional configuration.
//
// Example:
//
//	c, err := apple_update_cdn.NewClient(
//	    apple_update_cdn.WithTimeout(15 * time.Second),
//	    apple_update_cdn.WithLogger(logger),
//	)
func NewClient(options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransport(options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		AppleUpdateCDNAPI: &AppleUpdateCDNAPIClient{
			Firmware: firmware.NewService(transport),
			GDMF:     gdmf.NewService(transport),
			CDN:      cdn.NewService(transport),
		},
	}, nil
}

// NewDefaultClient creates a new Apple Update CDN client with default settings
// (30s timeout, 3 retries, no logging).
func NewDefaultClient() (*Client, error) {
	return NewClient()
}

// Close releases resources held by the client.
func (c *Client) Close() error {
	return c.transport.Close()
}
