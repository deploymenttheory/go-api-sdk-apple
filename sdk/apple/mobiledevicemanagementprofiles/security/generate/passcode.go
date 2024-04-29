package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/security/models"

	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// PasscodeDeviceMobileManagementProfileConfigOptions holds customizable options for the passcode profile.
type PasscodeMobileDeviceManagementProfileConfigOptions struct {
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

// CreatePasscodeMobileDeviceManagementProfile creates a device management profile for passcodes using provided configuration options.
func CreatePasscodeMobileDeviceManagementProfile(options PasscodeMobileDeviceManagementProfileConfigOptions) models.ObjectPasscodeMobileDeviceManagementProfile {
	descriptions := []models.Description{
		{Locale: "en-US", Description: "Passcode must meet complexity requirements."},
	}

	return models.ObjectPasscodeMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.PasscodeMobileDeviceManagementProfileSubsetPayload{
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
