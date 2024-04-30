/*
Device Management Profile: Software Update
Version: 	macOS 10.7+

Ref: 			https://developer.apple.com/documentation/devicemanagement/softwareupdate

Profile Example:
<?xml version=”1.0” encoding=”UTF-8”?>
<!DOCTYPE plist PUBLIC “-//Apple//DTD PLIST 1.0//EN” “http://www.apple.com/DTDs/PropertyList-1.0.dtd”>
<plist version=”1.0”>
<dict>

	<key>PayloadContent</key>
	<array>
	    <dict>
	        <key>AutomaticallyInstallAppUpdates</key>
	        <false/>
	        <key>PayloadIdentifier</key>
	        <string>com.example.mysoftwareupdatepayload</string>
	        <key>PayloadType</key>
	        <string>com.apple.SoftwareUpdate</string>
	        <key>PayloadUUID</key>
	        <string>af3c6efa-0dd3-4021-814b-6f2dba91428b</string>
	        <key>PayloadVersion</key>
	        <integer>1</integer>
	    </dict>
	</array>
	<key>PayloadDisplayName</key>
	<string>Software Update</string>
	<key>PayloadIdentifier</key>
	<string>com.example.myprofile</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>8b6061ab-31c6-4eee-ba5b-8a09ea8f5fa7</string>
	<key>PayloadVersion</key>
	<integer>1</integer>

</dict>
</plist>
*/
package models

// SystemUpdateMobileDeviceManagementProfile defines the structure for managing system update policies on devices.
type SystemUpdateMobileDeviceManagementProfile struct {
	Version                  string                                                   `plist:"version,attr"`
	PayloadContent           []SystemUpdateMobileDeviceManagementProfileSubsetPayload `plist:"PayloadContent"`
	PayloadDescription       string                                                   `plist:"PayloadDescription,omitempty"`
	PayloadDisplayName       string                                                   `plist:"PayloadDisplayName,omitempty"`
	PayloadEnabled           string                                                   `plist:"PayloadEnabled,omitempty"`
	PayloadIdentifier        string                                                   `plist:"PayloadIdentifier,omitempty"`
	PayloadOrganization      string                                                   `plist:"PayloadOrganization,omitempty"`
	PayloadRemovalDisallowed string                                                   `plist:"PayloadRemovalDisallowed,omitempty"`
	PayloadScope             string                                                   `plist:"PayloadScope,omitempty"`
	PayloadType              string                                                   `plist:"PayloadType,omitempty"`
	PayloadUUID              string                                                   `plist:"PayloadUUID,omitempty"`
	PayloadVersion           int                                                      `plist:"PayloadVersion,omitempty"`
}

// SystemUpdateMobileDeviceManagementProfileSubsetPayload represents the system update settings configured in the policy.
type SystemUpdateMobileDeviceManagementProfileSubsetPayload struct {
	AllowPreReleaseInstallation                 *bool  `plist:"AllowPreReleaseInstallation"`
	AutomaticallyInstallAppUpdates              *bool  `plist:"AutomaticallyInstallAppUpdates"`
	AutomaticallyInstallMacOSUpdates            *bool  `plist:"AutomaticallyInstallMacOSUpdates"`
	AutomaticCheckEnabled                       *bool  `plist:"AutomaticCheckEnabled"`
	AutomaticDownload                           *bool  `plist:"AutomaticDownload"`
	CatalogURL                                  string `plist:"CatalogURL"`
	ConfigDataInstall                           *bool  `plist:"ConfigDataInstall"`
	CriticalUpdateInstall                       *bool  `plist:"CriticalUpdateInstall"`
	RestrictSoftwareUpdateRequireAdminToInstall *bool  `plist:"restrict-software-update-require-admin-to-install"`
	PayloadIdentifier                           string `plist:"PayloadIdentifier"`
	PayloadType                                 string `plist:"PayloadType"`
	PayloadUUID                                 string `plist:"PayloadUUID"`
	PayloadVersion                              int    `plist:"PayloadVersion"`
}
