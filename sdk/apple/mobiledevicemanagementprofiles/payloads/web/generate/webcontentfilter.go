package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/web/common"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/web/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// WebContentFilterMobileDeviceManagementProfileConfigOptions holds the configuration options for creating a Web Content Filter profile.
type WebContentFilterMobileDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName                        string
	AutoFilterEnabled                         bool
	FilterType                                string
	PermittedURLs                             []string
	DenylistURLs                              []string
	FilterBrowsers                            bool
	FilterSockets                             bool
	FilterPackets                             bool
	ContentFilterUUID                         string
	FilterDataProviderBundleIdentifier        string
	FilterDataProviderDesignatedRequirement   string
	FilterPacketProviderBundleIdentifier      string
	FilterPacketProviderDesignatedRequirement string
	Organization                              string
	ServerAddress                             string
	UserDefinedName                           string
	UserName                                  string
	Password                                  string
	VendorConfig                              common.VendorConfig
}

// CreateWebContentFilterMobileDeviceManagementProfile creates a device management profile for web content filtering using provided configuration options.
func CreateWebContentFilterMobileDeviceManagementProfile(options WebContentFilterMobileDeviceManagementProfileConfigOptions) models.WebContentFilterMobileDeviceManagementProfile {
	return models.WebContentFilterMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.WebContentFilterMobileDeviceManagementProfileSubsetPayload{
			{
				AutoFilterEnabled:                         &options.AutoFilterEnabled,
				FilterType:                                options.FilterType,
				PermittedURLs:                             options.PermittedURLs,
				DenylistURLs:                              options.DenylistURLs,
				FilterBrowsers:                            &options.FilterBrowsers,
				FilterSockets:                             &options.FilterSockets,
				FilterPackets:                             &options.FilterPackets,
				ContentFilterUUID:                         options.ContentFilterUUID,
				FilterDataProviderBundleIdentifier:        options.FilterDataProviderBundleIdentifier,
				FilterDataProviderDesignatedRequirement:   options.FilterDataProviderDesignatedRequirement,
				FilterPacketProviderBundleIdentifier:      options.FilterPacketProviderBundleIdentifier,
				FilterPacketProviderDesignatedRequirement: options.FilterPacketProviderDesignatedRequirement,
				Organization:                              options.Organization,
				ServerAddress:                             options.ServerAddress,
				UserDefinedName:                           options.UserDefinedName,
				UserName:                                  options.UserName,
				Password:                                  options.Password,
				VendorConfig:                              options.VendorConfig,
				PayloadIdentifier:                         "com.gomdmsdkapple.webcontentfilterpayload",
				PayloadType:                               "com.apple.webcontent-filter",
				PayloadUUID:                               helpers.GenerateUUID(),
				PayloadVersion:                            1,
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.webcontentfilterpayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}
