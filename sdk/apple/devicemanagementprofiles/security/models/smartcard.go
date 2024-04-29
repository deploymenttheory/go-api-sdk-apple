/*
Device Management Profile: Smart Card
Version: 	macOS 10.12.4+

Ref: 			https://developer.apple.com/documentation/devicemanagement/smartcard

Profile Example:
<?xml version=”1.0” encoding=”UTF-8”?>
<!DOCTYPE plist PUBLIC “-//Apple//DTD PLIST 1.0//EN” “http://www.apple.com/DTDs/PropertyList-1.0.dtd”>
<plist version=”1.0”>
<dict>

	<key>PayloadContent</key>
	<array>
	    <dict>
	        <key>UserPairing</key>
	        <false/>
	        <key>allowSmartCard</key>
	        <false/>
	        <key>checkCertificateTrust</key>
	        <false/>
	        <key>oneCardPerUser</key>
	        <false/>
	        <key>tokenRemovalAction</key>
	        <integer>1</integer>
	        <key>enforceSmartCard</key>
	        <true/>
	        <key>PayloadIdentifier</key>
	        <string>com.example.mysmartcardpayload</string>
	        <key>PayloadType</key>
	        <string>com.apple.security.smartcard</string>
	        <key>PayloadUUID</key>
	        <string>88f7336c-d9f6-44d1-b486-11e4080e2223</string>
	        <key>PayloadVersion</key>
	        <integer>1</integer>
	    </dict>
	</array>
	<key>PayloadDisplayName</key>
	<string>SmartCard</string>
	<key>PayloadIdentifier</key>
	<string>com.example.myprofile</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>85091214-a32f-4131-8b03-0045e5d81c42</string>
	<key>PayloadVersion</key>
	<integer>1</integer>

</dict>
</plist>
*/
package models

// ResourceSmartCardDeviceManagementProfile defines the structure for managing smart card configurations on devices.
type ResourceSmartCardDeviceManagementProfile struct {
	Version                  string                                          `plist:"version,attr"`
	PayloadContent           []SmartCardDeviceManagementProfileSubsetPayload `plist:"PayloadContent"`
	PayloadDescription       string                                          `plist:"PayloadDescription,omitempty"`
	PayloadDisplayName       string                                          `plist:"PayloadDisplayName,omitempty"`
	PayloadEnabled           string                                          `plist:"PayloadEnabled,omitempty"`
	PayloadIdentifier        string                                          `plist:"PayloadIdentifier,omitempty"`
	PayloadOrganization      string                                          `plist:"PayloadOrganization,omitempty"`
	PayloadRemovalDisallowed string                                          `plist:"PayloadRemovalDisallowed,omitempty"`
	PayloadScope             string                                          `plist:"PayloadScope,omitempty"`
	PayloadType              string                                          `plist:"PayloadType,omitempty"`
	PayloadUUID              string                                          `plist:"PayloadUUID,omitempty"`
	PayloadVersion           int                                             `plist:"PayloadVersion,omitempty"`
}

// SmartCardDeviceManagementProfileSubsetPayload represents the smart card specific settings in the profile.
type SmartCardDeviceManagementProfileSubsetPayload struct {
	AllowSmartCard        *bool  `plist:"allowSmartCard,omitempty"`
	CheckCertificateTrust int    `plist:"checkCertificateTrust"`
	EnforceSmartCard      *bool  `plist:"enforceSmartCard,omitempty"`
	OneCardPerUser        *bool  `plist:"oneCardPerUser,omitempty"`
	TokenRemovalAction    int    `plist:"tokenRemovalAction,omitempty"`
	UserPairing           *bool  `plist:"UserPairing,omitempty"`
	PayloadIdentifier     string `plist:"PayloadIdentifier"`
	PayloadType           string `plist:"PayloadType"`
	PayloadUUID           string `plist:"PayloadUUID"`
	PayloadVersion        int    `plist:"PayloadVersion"`
}
