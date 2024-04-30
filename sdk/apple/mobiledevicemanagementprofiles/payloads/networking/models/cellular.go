/*
Device Management Profile: Cellular
Version:
- iOS 7.0+
- iPadOS 7.0+
- watchOS 3.2+

Ref: 			https://developer.apple.com/documentation/devicemanagement/cellular

Profile Example:
<?xml version=”1.0” encoding=”UTF-8”?>
<!DOCTYPE plist PUBLIC “-//Apple//DTD PLIST 1.0//EN” “http://www.apple.com/DTDs/PropertyList-1.0.dtd”>
<plist version=”1.0”>
<dict>

	<key>PayloadContent</key>
	<array>
	    <dict>
	        <key>AttachAPN</key>
	        <dict>
	            <key>Name</key>
	            <string>example.com</string>
	        </dict>
	        <key>PayloadIdentifier</key>
	        <string>com.example.mycellularnetworkpayload</string>
	        <key>PayloadType</key>
	        <string>com.apple.cellular</string>
	        <key>PayloadUUID</key>
	        <string>5a024a67-119f-4b38-8648-4c28a054ec5f</string>
	        <key>PayloadVersion</key>
	        <real>1</real>
	    </dict>
	</array>
	<key>PayloadDisplayName</key>
	<string>Cellular</string>
	<key>PayloadIdentifier</key>
	<string>com.example.myprofile</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>07eeff13-902a-408b-9bec-2228b86f944f</string>
	<key>PayloadVersion</key>
	<integer>1</integer>

</dict>
</plist>
*/
package models

// CellularMobileDeviceManagementProfile defines the structure for managing cellular configurations on devices.
type CellularMobileDeviceManagementProfile struct {
	Version                  string                                               `plist:"version,attr"`
	PayloadContent           []CellularMobileDeviceManagementProfileSubsetPayload `plist:"PayloadContent"`
	PayloadDescription       string                                               `plist:"PayloadDescription,omitempty"`
	PayloadDisplayName       string                                               `plist:"PayloadDisplayName,omitempty"`
	PayloadEnabled           string                                               `plist:"PayloadEnabled,omitempty"`
	PayloadIdentifier        string                                               `plist:"PayloadIdentifier,omitempty"`
	PayloadOrganization      string                                               `plist:"PayloadOrganization,omitempty"`
	PayloadRemovalDisallowed string                                               `plist:"PayloadRemovalDisallowed,omitempty"`
	PayloadScope             string                                               `plist:"PayloadScope,omitempty"`
	PayloadType              string                                               `plist:"PayloadType,omitempty"`
	PayloadUUID              string                                               `plist:"PayloadUUID,omitempty"`
	PayloadVersion           int                                                  `plist:"PayloadVersion,omitempty"`
}

// CellularMobileDeviceManagementProfilePayload represents the cellular configuration settings in the profile.
type CellularMobileDeviceManagementProfileSubsetPayload struct {
	APNs              []CellularMobileDeviceManagementProfileSubsetCellularAPNsItem `plist:"APNs"`
	AttachAPN         CellularMobileDeviceManagementProfileSubsetCellularAttachAPN  `plist:"AttachAPN"`
	PayloadIdentifier string                                                        `plist:"PayloadIdentifier"`
	PayloadType       string                                                        `plist:"PayloadType"`
	PayloadUUID       string                                                        `plist:"PayloadUUID"`
	PayloadVersion    int                                                           `plist:"PayloadVersion"`
}

// CellularMobileDeviceManagementProfileSubsetCellularAPNsItem defines individual APN settings.
type CellularMobileDeviceManagementProfileSubsetCellularAPNsItem struct {
	AllowedProtocolMask                  int    `plist:"AllowedProtocolMask"`
	AllowedProtocolMaskInDomesticRoaming int    `plist:"AllowedProtocolMaskInDomesticRoaming"`
	AllowedProtocolMaskInRoaming         int    `plist:"AllowedProtocolMaskInRoaming"`
	AuthenticationType                   string `plist:"AuthenticationType"`
	EnableXLAT464                        *bool  `plist:"EnableXLAT464"`
	Name                                 string `plist:"Name"`
	Password                             string `plist:"Password,omitempty"`
	ProxyPort                            int    `plist:"ProxyPort,omitempty"`
	ProxyServer                          string `plist:"ProxyServer,omitempty"`
	Username                             string `plist:"Username,omitempty"`
}

// CellularMobileDeviceManagementProfileSubsetCellularAttachAPN defines the attach APN configuration.
type CellularMobileDeviceManagementProfileSubsetCellularAttachAPN struct {
	AllowedProtocolMask int    `plist:"AllowedProtocolMask"`
	AuthenticationType  string `plist:"AuthenticationType"`
	Name                string `plist:"Name"`
	Password            string `plist:"Password,omitempty"`
	Username            string `plist:"Username,omitempty"`
}
