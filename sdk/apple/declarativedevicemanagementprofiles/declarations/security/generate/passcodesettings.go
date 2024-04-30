package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/declarativedevicemanagementprofiles/declarations/security/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// PasscodeDeclarativeDeviceManagementProfileConfigOptions holds customizable options for the passcode profile.
type PasscodeDeclarativeDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName           string
	ChangeAtNextAuth             bool
	FailedAttemptsResetInMinutes int
	MaximumFailedAttempts        int
	MaximumGracePeriodInMinutes  int
	MaximumInactivityInMinutes   int
	MaximumPasscodeAgeInDays     int
	MinimumComplexCharacters     int
	MinimumLength                int
	PasscodeReuseLimit           int
	RequireAlphanumericPasscode  bool
	RequireComplexPasscode       bool
	RequirePasscode              bool
	CustomRegexPattern           string
	CustomRegexDescription       map[string]string
}

// CreatePasscodeSettingsDeclarativeDeviceManagementProfile creates a device management profile for passcodes using provided configuration options.
func CreatePasscodeSettingsDeclarativeDeviceManagementProfile(options PasscodeDeclarativeDeviceManagementProfileConfigOptions) models.ObjectPasscodeSettingsDeclarativeDeviceManagementProfile {
	return models.ObjectPasscodeSettingsDeclarativeDeviceManagementProfile{
		Type:        "com.apple.configuration.passcode.settings",
		Identifier:  helpers.GenerateUUID(),
		ServerToken: helpers.GenerateUUID(),
		Payload: models.PasscodeSettingsDeclarativeDeviceManagementProfileSubsetPayload{
			ChangeAtNextAuth:             &options.ChangeAtNextAuth,
			FailedAttemptsResetInMinutes: options.FailedAttemptsResetInMinutes,
			MaximumFailedAttempts:        options.MaximumFailedAttempts,
			MaximumGracePeriodInMinutes:  options.MaximumGracePeriodInMinutes,
			MaximumInactivityInMinutes:   options.MaximumInactivityInMinutes,
			MaximumPasscodeAgeInDays:     options.MaximumPasscodeAgeInDays,
			MinimumComplexCharacters:     options.MinimumComplexCharacters,
			MinimumLength:                options.MinimumLength,
			PasscodeReuseLimit:           options.PasscodeReuseLimit,
			RequireAlphanumericPasscode:  &options.RequireAlphanumericPasscode,
			RequireComplexPasscode:       &options.RequireComplexPasscode,
			RequirePasscode:              &options.RequirePasscode,
			CustomRegex: models.PasscodeSettingsCustomRegexObject{
				Description: models.PasscodeSettingsCustomRegexDescriptionObject(options.CustomRegexDescription),
				Regex:       options.CustomRegexPattern,
			},
		},
	}
}
