package itunes

import (
	"github.com/deploymenttheory/go-sdk-appleservices/itunes/client"
	"github.com/deploymenttheory/go-sdk-appleservices/itunes/itunes_api/search"
)

// Client is the main entry point for the iTunes Search API SDK.
type Client struct {
	transport *client.Transport
	ItunesAPI *ItunesAPIClient
}

// ItunesAPIClient groups all iTunes API services.
type ItunesAPIClient struct {
	Search *search.SearchService
}

// NewClient creates a new iTunes Search API client with optional configuration.
//
// Example:
//
//	c, err := itunes.NewClient(
//	    itunes.WithTimeout(15 * time.Second),
//	    itunes.WithLogger(logger),
//	)
func NewClient(options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransport(options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		ItunesAPI: &ItunesAPIClient{
			Search: search.NewService(transport),
		},
	}, nil
}

// NewDefaultClient creates a new iTunes Search API client with default settings
// (30s timeout, 3 retries, no logging).
func NewDefaultClient() (*Client, error) {
	return NewClient()
}

// Close releases resources held by the client.
func (c *Client) Close() error {
	return c.transport.Close()
}
