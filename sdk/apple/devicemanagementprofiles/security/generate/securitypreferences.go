package generate

import (
	"github.com/deploymenttheory/go-mdm-sdk-apple/sdk/apple/devicemanagementprofiles/security/models"
	"github.com/deploymenttheory/go-mdm-sdk-apple/sdk/helpers"
)

// SecurityPreferencesDeviceManagementProfileConfigOptions holds customizable options for the security preferences profile.
type SecurityPreferencesDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName       string
	DontAllowFireWallUI      bool
	DontAllowLockMessageUI   bool
	DontAllowPasswordResetUI bool
}

// CreateSecurityPreferencesDeviceManagementProfile creates a device management profile for security preferences using provided configuration options.
func CreateSecurityPreferencesDeviceManagementProfile(options SecurityPreferencesDeviceManagementProfileConfigOptions) models.ResourceSecurityPreferencesDeviceManagementProfile {
	return models.ResourceSecurityPreferencesDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.SecurityPreferencesDeviceManagementProfileSubsetPayload{
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
