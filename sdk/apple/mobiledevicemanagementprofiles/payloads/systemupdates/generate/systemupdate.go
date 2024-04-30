package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/systemupdates/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// SystemUpdateDeviceManagementProfileConfigOptions holds customizable options for the system update profile.
type SystemUpdateMobileDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName                          string
	AutomaticallyInstallAppUpdates              bool
	AutomaticallyInstallMacOSUpdates            bool
	AutomaticCheckEnabled                       bool
	AutomaticDownload                           bool
	CatalogURL                                  string
	ConfigDataInstall                           bool
	CriticalUpdateInstall                       bool
	RestrictSoftwareUpdateRequireAdminToInstall bool
	AllowPreReleaseInstallation                 bool
}

// CreateSystemUpdateMobileDeviceManagementProfile creates a device management profile for system updates using provided configuration options.
func CreateSystemUpdateMobileDeviceManagementProfile(options SystemUpdateMobileDeviceManagementProfileConfigOptions) models.SystemUpdateMobileDeviceManagementProfile {
	return models.SystemUpdateMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.SystemUpdateMobileDeviceManagementProfileSubsetPayload{
			{
				AutomaticallyInstallAppUpdates:              &options.AutomaticallyInstallAppUpdates,
				AutomaticallyInstallMacOSUpdates:            &options.AutomaticallyInstallMacOSUpdates,
				AutomaticCheckEnabled:                       &options.AutomaticCheckEnabled,
				AutomaticDownload:                           &options.AutomaticDownload,
				CatalogURL:                                  options.CatalogURL,
				ConfigDataInstall:                           &options.ConfigDataInstall,
				CriticalUpdateInstall:                       &options.CriticalUpdateInstall,
				RestrictSoftwareUpdateRequireAdminToInstall: &options.RestrictSoftwareUpdateRequireAdminToInstall,
				AllowPreReleaseInstallation:                 &options.AllowPreReleaseInstallation,
				PayloadIdentifier:                           "com.gomdmsdkapple.softwareupdatepayload",
				PayloadType:                                 "com.apple.SoftwareUpdate",
				PayloadUUID:                                 helpers.GenerateUUID(),
				PayloadVersion:                              1,
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.softwareupdatepayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}
