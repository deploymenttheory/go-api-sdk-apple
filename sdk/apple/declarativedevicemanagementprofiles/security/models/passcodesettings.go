/*
Declarative Device Management Profile: Passcode Settings
Version:
  - iOS 15.0+
	- iPadOS 15.0+
	- macOS 13.0+
	- watchOS 10.0+

Ref: 			https://developer.apple.com/documentation/devicemanagement/passcodesettings

DDM Profile Example:
*/

package models

// ObjectPasscodeSettingsDeclarativeDeviceManagementProfile defines the structure for managing passcode policies on devices.
type ObjectPasscodeSettingsDeclarativeDeviceManagementProfile struct {
	Type        string                                                          `json:"Type"`        // The declaration type, a dot-separated sequence of tokens.
	Identifier  string                                                          `json:"Identifier"`  // The unique identifier, typically a UUID, within the set of all declarations sent to a device.
	ServerToken string                                                          `json:"ServerToken"` // An identifier for a unique revision of the declaration.
	Payload     PasscodeSettingsDeclarativeDeviceManagementProfileSubsetPayload `json:"Payload"`     // The data-specific part of the declaration containing keys and values.
}

// PasscodeSettingsDeclarativeDeviceManagementProfileSubsetPayload represents the passcode requirements set in the policy.
type PasscodeSettingsDeclarativeDeviceManagementProfileSubsetPayload struct {
	ChangeAtNextAuth             *bool                             `json:"ChangeAtNextAuth"`
	FailedAttemptsResetInMinutes int                               `json:"FailedAttemptsResetInMinutes"`
	MaximumFailedAttempts        int                               `json:"MaximumFailedAttempts"`
	MaximumGracePeriodInMinutes  int                               `json:"MaximumGracePeriodInMinutes"`
	MaximumInactivityInMinutes   int                               `json:"MaximumInactivityInMinutes"`
	MaximumPasscodeAgeInDays     int                               `json:"MaximumPasscodeAgeInDays"`
	MinimumComplexCharacters     int                               `json:"MinimumComplexCharacters"`
	MinimumLength                int                               `json:"MinimumLength"`
	PasscodeReuseLimit           int                               `json:"PasscodeReuseLimit"`
	RequireAlphanumericPasscode  *bool                             `json:"RequireAlphanumericPasscode"`
	RequireComplexPasscode       *bool                             `json:"RequireComplexPasscode"`
	RequirePasscode              *bool                             `json:"RequirePasscode"`
	CustomRegex                  PasscodeSettingsCustomRegexObject `json:"CustomRegex"`
	PayloadIdentifier            string                            `json:"PayloadIdentifier"`
	PayloadType                  string                            `json:"PayloadType"`
	PayloadUUID                  string                            `json:"PayloadUUID"`
	PayloadVersion               int                               `json:"PayloadVersion"`
}

// PasscodeSettingsCustomRegexObject defines the regular expression and its descriptions to enforce password compliance.
type PasscodeSettingsCustomRegexObject struct {
	Description PasscodeSettingsCustomRegexDescriptionObject `json:"Description"`
	Regex       string                                       `json:"Regex"`
}

// PasscodeSettingsCustomRegexDescriptionObject contains descriptions keyed by language IDs.
type PasscodeSettingsCustomRegexDescriptionObject map[string]string // Map of locale to description
