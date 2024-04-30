/*
Mobile Device Management Profile: Web Content Filter
Version: 	-iOS 7.0+
iPadOS 7.0+
macOS 10.15+
visionOS 1.1+

Ref: 			https://developer.apple.com/documentation/devicemanagement/webcontentfilter

Profile Example:
<?xml version=”1.0” encoding=”UTF-8”?>
<!DOCTYPE plist PUBLIC “-//Apple//DTD PLIST 1.0//EN” “http://www.apple.com/DTDs/PropertyList-1.0.dtd”>
<plist version=”1.0”>
<dict>

	<key>PayloadContent</key>
	<array>
	    <dict>
	        <key>AutoFilterEnabled</key>
	        <true/>
	        <key>DenylistURLs</key>
	        <array>
	            <string>https://notallowedname.company.com</string>
	        </array>
	        <key>FilterBrowsers</key>
	        <true/>
	        <key>FilterSockets</key>
	        <true/>
	        <key>FilterType</key>
	        <string>BuiltIn</string>
	        <key>PermittedURLs</key>
	        <array>
	            <string>https://example.company.com</string>
	        </array>
	        <key>PayloadIdentifier</key>
	        <string>com.example.mywebcontentfilterpayload</string>
	        <key>PayloadType</key>
	        <string>com.apple.webcontent-filter</string>
	        <key>PayloadUUID</key>
	        <string>fb5d598f-0a96-4b77-9702-9edfc3417601</string>
	        <key>PayloadVersion</key>
	        <integer>1</integer>
	    </dict>
	</array>
	<key>PayloadDisplayName</key>
	<string>Web Content Filter</string>
	<key>PayloadIdentifier</key>
	<string>com.example.myprofile</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>b510e0c6-dc81-4b62-88d0-6a3ef82925e7</string>
	<key>PayloadVersion</key>
	<integer>1</integer>

</dict>
</plist>
*/
package models

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/web/common"
)

// WebContentFilterMobileDeviceManagementProfile defines the structure for managing system update policies on devices.
type WebContentFilterMobileDeviceManagementProfile struct {
	Version                  string                                                       `plist:"version,attr"`
	PayloadContent           []WebContentFilterMobileDeviceManagementProfileSubsetPayload `plist:"PayloadContent"`
	PayloadDescription       string                                                       `plist:"PayloadDescription,omitempty"`
	PayloadDisplayName       string                                                       `plist:"PayloadDisplayName,omitempty"`
	PayloadEnabled           string                                                       `plist:"PayloadEnabled,omitempty"`
	PayloadIdentifier        string                                                       `plist:"PayloadIdentifier,omitempty"`
	PayloadOrganization      string                                                       `plist:"PayloadOrganization,omitempty"`
	PayloadRemovalDisallowed string                                                       `plist:"PayloadRemovalDisallowed,omitempty"`
	PayloadScope             string                                                       `plist:"PayloadScope,omitempty"`
	PayloadType              string                                                       `plist:"PayloadType,omitempty"`
	PayloadUUID              string                                                       `plist:"PayloadUUID,omitempty"`
	PayloadVersion           int                                                          `plist:"PayloadVersion,omitempty"`
}

// WebContentFilterMobileDeviceManagementProfileSubsetPayload represents the web content filter settings configured in the policy.
type WebContentFilterMobileDeviceManagementProfileSubsetPayload struct {
	AutoFilterEnabled                         *bool               `plist:"AutoFilterEnabled,omitempty"`                         // Enables automatic filtering if true.
	FilterType                                string              `plist:"FilterType"`                                          // Specifies the type of filter, BuiltIn or Plugin.
	PermittedURLs                             []string            `plist:"PermittedURLs,omitempty"`                             // URLs that are always accessible.
	DenylistURLs                              []string            `plist:"DenylistURLs,omitempty"`                              // URLs that are blocked.
	FilterBrowsers                            *bool               `plist:"FilterBrowsers,omitempty"`                            // Enables filtering of WebKit traffic if true.
	FilterSockets                             *bool               `plist:"FilterSockets,omitempty"`                             // Enables filtering of socket traffic if true.
	FilterPackets                             *bool               `plist:"FilterPackets,omitempty"`                             // Enables filtering of network packets if true.
	ContentFilterUUID                         string              `plist:"ContentFilterUUID,omitempty"`                         // UUID for the content filter configuration.
	FilterDataProviderBundleIdentifier        string              `plist:"FilterDataProviderBundleIdentifier,omitempty"`        // Bundle identifier for the data provider.
	FilterDataProviderDesignatedRequirement   string              `plist:"FilterDataProviderDesignatedRequirement,omitempty"`   // Designated requirement for the data provider.
	FilterPacketProviderBundleIdentifier      string              `plist:"FilterPacketProviderBundleIdentifier,omitempty"`      // Bundle identifier for the packet provider.
	FilterPacketProviderDesignatedRequirement string              `plist:"FilterPacketProviderDesignatedRequirement,omitempty"` // Designated requirement for the packet provider.
	Organization                              string              `plist:"Organization,omitempty"`                              // Organization name for Plugin type.
	ServerAddress                             string              `plist:"ServerAddress,omitempty"`                             // Server address for the service.
	UserDefinedName                           string              `plist:"UserDefinedName,omitempty"`                           // Display name for the filter configuration.
	UserName                                  string              `plist:"UserName,omitempty"`                                  // User name for the service.
	Password                                  string              `plist:"Password,omitempty"`                                  // Password for the service.
	VendorConfig                              common.VendorConfig `plist:"VendorConfig,omitempty"`                              // Custom configuration for the vendor's plugin.
	PayloadIdentifier                         string              `plist:"PayloadIdentifier"`                                   // Unique identifier for this payload.
	PayloadType                               string              `plist:"PayloadType"`                                         // Type identifier for this payload, should be com.apple.webcontent-filter.
	PayloadUUID                               string              `plist:"PayloadUUID"`                                         // Unique identifier for this payload instance.
	PayloadVersion                            int                 `plist:"PayloadVersion"`                                      // Version of the payload format.
}
