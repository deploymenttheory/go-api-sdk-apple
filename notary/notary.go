package notary

import (
	"github.com/deploymenttheory/go-sdk-appleservices/notary/client"
	"github.com/deploymenttheory/go-sdk-appleservices/notary/notary_api/submissions"
)

// Client is the main entry point for the Apple Notary API SDK.
type Client struct {
	transport *client.Transport
	NotaryAPI *NotaryAPIClient
}

// NotaryAPIClient groups all Apple Notary API services.
type NotaryAPIClient struct {
	Submissions *submissions.Submissions
}

// NewClient creates a new Apple Notary API client.
// Parameters:
//   - keyID: Your App Store Connect API Key ID
//   - issuerID: Your App Store Connect Issuer ID (Team ID)
//   - privateKey: Your App Store Connect private key (*rsa.PrivateKey or *ecdsa.PrivateKey)
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClient(keyID, issuerID string, privateKey any, options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransport(keyID, issuerID, privateKey, options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		NotaryAPI: &NotaryAPIClient{
			Submissions: submissions.NewService(transport),
		},
	}, nil
}

// NewClientFromFile creates a client using a private key from file.
// Parameters:
//   - keyID: Your App Store Connect API Key ID
//   - issuerID: Your App Store Connect Issuer ID (Team ID)
//   - privateKeyPath: Path to your App Store Connect private key file (.p8)
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClientFromFile(keyID, issuerID, privateKeyPath string, options ...client.ClientOption) (*Client, error) {
	privateKey, err := client.LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return NewClient(keyID, issuerID, privateKey, options...)
}

// NewClientFromEnv creates a client using environment variables.
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, and one of APPLE_PRIVATE_KEY_PEM or APPLE_PRIVATE_KEY_PATH.
// Parameters:
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClientFromEnv(options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransportFromEnv(options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		NotaryAPI: &NotaryAPIClient{
			Submissions: submissions.NewService(transport),
		},
	}, nil
}
