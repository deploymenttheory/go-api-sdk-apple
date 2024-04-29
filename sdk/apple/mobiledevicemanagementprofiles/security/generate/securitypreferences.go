package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/security/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// SecurityPreferencesMobileDeviceManagementProfileConfigOptions holds customizable options for the security preferences profile.
type SecurityPreferencesMobileDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName       string
	DontAllowFireWallUI      bool
	DontAllowLockMessageUI   bool
	DontAllowPasswordResetUI bool
}

// CreateSecurityPreferencesMobileDeviceManagementProfile creates a device management profile for security preferences using provided configuration options.
func CreateSecurityPreferencesMobileDeviceManagementProfile(options SecurityPreferencesMobileDeviceManagementProfileConfigOptions) models.ObjectSecurityPreferencesMobileDeviceManagementProfile {
	return models.ObjectSecurityPreferencesMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.SecurityPreferencesMobileDeviceManagementProfileSubsetPayload{
			{
				DontAllowFireWallUI:      &options.DontAllowFireWallUI,
				DontAllowLockMessageUI:   &options.DontAllowLockMessageUI,
				DontAllowPasswordResetUI: &options.DontAllowPasswordResetUI,
				PayloadIdentifier:        "com.gomdmsdkapple.securitypreferencespayload",
				PayloadType:              "com.apple.preference.security",
				PayloadUUID:              helpers.GenerateUUID(),
				PayloadVersion:           1,
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.securitypreferencespayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}
