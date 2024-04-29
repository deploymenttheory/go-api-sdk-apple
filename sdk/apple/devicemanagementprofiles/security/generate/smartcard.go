package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/devicemanagementprofiles/security/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// SmartCardDeviceManagementProfileConfigOptions holds customizable options for the smart card profile.
type SmartCardDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName    string
	UserPairing           bool
	AllowSmartCard        bool
	CheckCertificateTrust int
	OneCardPerUser        bool
	TokenRemovalAction    int
	EnforceSmartCard      bool
}

// CreateSmartCardDeviceManagementProfile creates a device management profile for smart card using provided configuration options.
func CreateSmartCardDeviceManagementProfile(options SmartCardDeviceManagementProfileConfigOptions) models.ResourceSmartCardDeviceManagementProfile {
	return models.ResourceSmartCardDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.SmartCardDeviceManagementProfileSubsetPayload{
			{
				UserPairing:           &options.UserPairing,
				AllowSmartCard:        &options.AllowSmartCard,
				CheckCertificateTrust: options.CheckCertificateTrust,
				OneCardPerUser:        &options.OneCardPerUser,
				TokenRemovalAction:    options.TokenRemovalAction,
				EnforceSmartCard:      &options.EnforceSmartCard,
				PayloadIdentifier:     "com.gomdmsdkapple.smartcardpayload",
				PayloadType:           "com.apple.security.smartcard",
				PayloadUUID:           helpers.GenerateUUID(),
				PayloadVersion:        1,
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.smartcardpayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}
