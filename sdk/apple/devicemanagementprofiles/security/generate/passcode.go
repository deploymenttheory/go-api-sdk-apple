package generate

import (
	"github.com/deploymenttheory/go-mdm-sdk-apple/sdk/apple/devicemanagementprofiles/security/models"

	"github.com/deploymenttheory/go-mdm-sdk-apple/sdk/helpers"
)

// DeviceManagementProfileConfigOptions holds customizable options for the passcode profile.
type DeviceManagementProfileConfigOptions struct {
	PayloadDisplayName  string
	AllowSimple         bool
	ForcePIN            bool
	MaxFailedAttempts   int
	MaxGracePeriod      int
	MaxInactivity       int
	MaxPINAgeInDays     int
	MinLength           int
	PinHistory          int
	RequireAlphanumeric bool
	CustomRegexPattern  string
}

// CreatePasscodeDeviceManagementProfile creates a device management profile for passcodes using provided configuration options.
func CreatePasscodeDeviceManagementProfile(options DeviceManagementProfileConfigOptions) models.ResourcePasscodeDeviceManagementProfile {
	descriptions := []models.Description{
		{Locale: "en-US", Description: "Passcode must meet complexity requirements."},
	}

	return models.ResourcePasscodeDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.PasscodeDeviceManagementProfileSubsetPayload{
			{
				AllowSimple:         &options.AllowSimple,
				ForcePIN:            &options.ForcePIN,
				MaxFailedAttempts:   options.MaxFailedAttempts,
				MaxGracePeriod:      options.MaxGracePeriod,
				MaxInactivity:       options.MaxInactivity,
				MaxPINAgeInDays:     options.MaxPINAgeInDays,
				MinLength:           options.MinLength,
				PinHistory:          options.PinHistory,
				RequireAlphanumeric: &options.RequireAlphanumeric,
				PayloadIdentifier:   "com.gomdmsdkapple.passcodepayload",
				PayloadType:         "com.apple.mobiledevice.passwordpolicy",
				PayloadUUID:         helpers.GenerateUUID(),
				PayloadVersion:      1,
				CustomRegex: &models.CustomRegex{
					PasswordContentDescriptions: descriptions,
					PasswordContentRegex:        options.CustomRegexPattern,
				},
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.passcodepayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}
