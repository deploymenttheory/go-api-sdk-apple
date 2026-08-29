package axm

import (
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/apps"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/auditevents"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/blueprints"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/configurations"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/organizationalunits"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/packages"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/usergroups"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/users"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
)

// Client is the main entry point for the Apple Business Manager API SDK.
type Client struct {
	transport *client.Transport
	AXMAPI    *AXMAPIClient
}

// AXMAPIClient groups all Apple Business Manager API services.
type AXMAPIClient struct {
	Devices             *devices.Devices
	DeviceManagement    *devicemanagement.DeviceManagement
	AuditEvents         *auditevents.AuditEvents
	Users               *users.Users
	UserGroups          *usergroups.UserGroups
	OrganizationalUnits *organizationalunits.OrganizationalUnits
	Apps                *apps.Apps
	Packages            *packages.Packages
	Configurations      *configurations.Configurations
	Blueprints          *blueprints.Blueprints
}

// NewClient creates a new Apple Business Manager client.
// Parameters:
//   - keyID: Your Apple Developer Key ID
//   - issuerID: Your Apple Developer Issuer ID (Team ID)
//   - privateKey: Your Apple Developer private key (*rsa.PrivateKey or *ecdsa.PrivateKey)
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClient(keyID, issuerID string, privateKey any, options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransport(keyID, issuerID, privateKey, options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		AXMAPI: &AXMAPIClient{
			Devices:             devices.NewService(transport),
			DeviceManagement:    devicemanagement.NewService(transport),
			AuditEvents:         auditevents.NewService(transport),
			Users:               users.NewService(transport),
			UserGroups:          usergroups.NewService(transport),
			OrganizationalUnits: organizationalunits.NewService(transport),
			Apps:                apps.NewService(transport),
			Packages:            packages.NewService(transport),
			Configurations:      configurations.NewService(transport),
			Blueprints:          blueprints.NewService(transport),
		},
	}, nil
}

// NewClientFromFile creates a client using private key from file.
// Parameters:
//   - keyID: Your Apple Developer Key ID
//   - issuerID: Your Apple Developer Issuer ID (Team ID)
//   - privateKeyPath: Path to your Apple Developer private key file (.p8)
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClientFromFile(keyID, issuerID, privateKeyPath string, options ...client.ClientOption) (*Client, error) {
	privateKey, err := client.LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return NewClient(keyID, issuerID, privateKey, options...)
}

// NewClientFromEnv creates a client using environment variables.
// Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH
// Parameters:
//   - options: Optional configuration options (WithLogger, WithTimeout, etc.)
func NewClientFromEnv(options ...client.ClientOption) (*Client, error) {
	transport, err := client.NewTransportFromEnv(options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		transport: transport,
		AXMAPI: &AXMAPIClient{
			Devices:             devices.NewService(transport),
			DeviceManagement:    devicemanagement.NewService(transport),
			AuditEvents:         auditevents.NewService(transport),
			Users:               users.NewService(transport),
			UserGroups:          usergroups.NewService(transport),
			OrganizationalUnits: organizationalunits.NewService(transport),
			Apps:                apps.NewService(transport),
			Packages:            packages.NewService(transport),
			Configurations:      configurations.NewService(transport),
			Blueprints:          blueprints.NewService(transport),
		},
	}, nil
}
