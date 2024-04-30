package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/security/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// SmartCardMobileDeviceManagementProfileConfigOptions holds customizable options for the smart card profile.
type SmartCardMobileDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName    string
	UserPairing           bool
	AllowSmartCard        bool
	CheckCertificateTrust int
	OneCardPerUser        bool
	TokenRemovalAction    int
	EnforceSmartCard      bool
}

// CreateSmartCardMobileDeviceManagementProfile creates a device management profile for smart card using provided configuration options.
func CreateSmartCardMobileDeviceManagementProfile(options SmartCardMobileDeviceManagementProfileConfigOptions) models.ObjectSmartCardMobileDeviceManagementProfile {
	return models.ObjectSmartCardMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.SmartCardMobileDeviceManagementProfileSubsetPayload{
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
