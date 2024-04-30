/*
Device Management Profile: WiFi
Version:
- iOS 4.0+
- iPadOS 4.0+
- macOS 10.7+
- tvOS 9.0+
- watchOS 3.2+
- visionOS 1.0+

Ref: 			https://developer.apple.com/documentation/devicemanagement/wifi

Profile Example:
<?xml version=”1.0” encoding=”UTF-8”?>
<!DOCTYPE plist PUBLIC “-//Apple//DTD PLIST 1.0//EN” “http://www.apple.com/DTDs/PropertyList-1.0.dtd”>
<plist version=”1.0”>
<dict>

	<key>New item</key>
	<dict>
	    <key>PayloadContent</key>
	    <array>
	        <dict>
	            <key>AutoJoin</key>
	            <true/>
	            <key>CaptiveBypass</key>
	            <false/>
	            <key>DisableAssociationMACRandomization</key>
	            <false/>
	            <key>EncryptionType</key>
	            <string>WPA</string>
	            <key>HIDDEN_NETWORK</key>
	            <true/>
	            <key>IsHotspot</key>
	            <false/>
	            <key>Password</key>
	            <string>Password123</string>
	            <key>PayloadDisplayName</key>
	            <string>Wi-Fi</string>
	            <key>ProxyType</key>
	            <string>None</string>
	            <key>SSID_STR</key>
	            <string>Example</string>
	            <key>PayloadIdentifier</key>
	            <string>com.example.mywifipayload</string>
	            <key>PayloadType</key>
	            <string>com.apple.wifi.managed</string>
	            <key>PayloadUUID</key>
	            <string>94c487e0-d6f8-41e3-b66d-a89994e6919b</string>
	            <key>PayloadVersion</key>
	            <integer>1</integer>
	        </dict>
	    </array>
	    <key>PayloadDisplayName</key>
	    <string>Wi-Fi</string>
	    <key>PayloadIdentifier</key>
	    <string>com.example.myprofile</string>
	    <key>PayloadType</key>
	    <string>Configuration</string>
	    <key>PayloadUUID</key>
	    <string>71e9b0f7-02f8-4aea-b365-b381d872909a</string>
	    <key>PayloadVersion</key>
	    <integer>1</integer>
	</dict>

</dict>
</plist>
*/
package models

// WiFiMobileDeviceManagementProfile defines the structure for managing WiFi configurations on devices.
type WiFiMobileDeviceManagementProfile struct {
	Version                  string                                           `plist:"version,attr"`
	PayloadContent           []WiFiMobileDeviceManagementProfileSubsetPayload `plist:"PayloadContent"`
	PayloadDescription       string                                           `plist:"PayloadDescription,omitempty"`
	PayloadDisplayName       string                                           `plist:"PayloadDisplayName,omitempty"`
	PayloadEnabled           string                                           `plist:"PayloadEnabled,omitempty"`
	PayloadIdentifier        string                                           `plist:"PayloadIdentifier,omitempty"`
	PayloadOrganization      string                                           `plist:"PayloadOrganization,omitempty"`
	PayloadRemovalDisallowed string                                           `plist:"PayloadRemovalDisallowed,omitempty"`
	PayloadScope             string                                           `plist:"PayloadScope,omitempty"`
	PayloadType              string                                           `plist:"PayloadType,omitempty"`
	PayloadUUID              string                                           `plist:"PayloadUUID,omitempty"`
	PayloadVersion           int                                              `plist:"PayloadVersion,omitempty"`
}

// WiFiMobileDeviceManagementProfileSubsetPayload represents the WiFi configuration settings in the profile.
type WiFiMobileDeviceManagementProfileSubsetPayload struct {
	SSID_STR                           string                 `plist:"SSID_STR"`
	AutoJoin                           *bool                  `plist:"AutoJoin"`
	CaptiveBypass                      *bool                  `plist:"CaptiveBypass"`
	DisableAssociationMACRandomization *bool                  `plist:"DisableAssociationMACRandomization"`
	DisplayedOperatorName              string                 `plist:"DisplayedOperatorName"`
	DomainName                         string                 `plist:"DomainName"`
	EAPClientConfiguration             EAPClientConfiguration `plist:"EAPClientConfiguration"`
	EnableIPv6                         *bool                  `plist:"EnableIPv6"`
	EncryptionType                     string                 `plist:"EncryptionType"`
	HESSID                             string                 `plist:"HESSID"`
	HIDDEN_NETWORK                     *bool                  `plist:"HIDDEN_NETWORK"`
	IsHotspot                          *bool                  `plist:"IsHotspot"`
	MCCAndMNCs                         []string               `plist:"MCCAndMNCs"`
	NAIRealmNames                      []string               `plist:"NAIRealmNames"`
	Password                           string                 `plist:"Password"`
	PayloadCertificateUUID             string                 `plist:"PayloadCertificateUUID"`
	ProxyPACURL                        string                 `plist:"ProxyPACURL"`
	ProxyServer                        string                 `plist:"ProxyServer"`
	ProxyServerPort                    int                    `plist:"ProxyServerPort"`
	ProxyType                          string                 `plist:"ProxyType"`
	QoSMarkingPolicy                   QoSMarkingPolicy       `plist:"QoSMarkingPolicy"`
	RoamingConsortiumOIs               []string               `plist:"RoamingConsortiumOIs"`
	ServiceProviderRoamingEnabled      *bool                  `plist:"ServiceProviderRoamingEnabled"`
	SetupModes                         []string               `plist:"SetupModes"`
	TLSCertificateRequired             *bool                  `plist:"TLSCertificateRequired"`
	TLSMinimumVersion                  string                 `plist:"TLSMinimumVersion"`
	TLSMaximumVersion                  string                 `plist:"TLSMaximumVersion"`
	TLSTrustedCertificates             []string               `plist:"TLSTrustedCertificates"`
	TLSTrustedServerNames              []string               `plist:"TLSTrustedServerNames"`
	UserName                           string                 `plist:"UserName"`
	UserPassword                       string                 `plist:"UserPassword"`
	PayloadIdentifier                  string                 `plist:"PayloadIdentifier"`
	PayloadType                        string                 `plist:"PayloadType"`
	PayloadUUID                        string                 `plist:"PayloadUUID"`
	PayloadVersion                     int                    `plist:"PayloadVersion"`
}

// EAPClientConfiguration structures the necessary elements for enterprise network configurations.
type EAPClientConfiguration struct {
	AcceptEAPTypes                 []int  `plist:"AcceptEAPTypes"`
	EAPFASTProvisionPAC            *bool  `plist:"EAPFASTProvisionPAC"`
	EAPFASTProvisionPACAnonymously *bool  `plist:"EAPFASTProvisionPACAnonymously"`
	EAPFASTUsePAC                  *bool  `plist:"EAPFASTUsePAC"`
	EAPSIMNumberOfRANDs            int    `plist:"EAPSIMNumberOfRANDs"`
	OneTimeUserPassword            *bool  `plist:"OneTimeUserPassword"`
	OuterIdentity                  string `plist:"OuterIdentity"`
	TTLSInnerAuthentication        string `plist:"TTLSInnerAuthentication"`
}

// QoSMarkingPolicy defines the settings for Quality of Service marking in a WiFi network.
type QoSMarkingPolicy struct {
	QoSMarkingAppleAudioVideoCalls    *bool    `plist:"QoSMarkingAppleAudioVideoCalls"`
	QoSMarkingEnabled                 *bool    `plist:"QoSMarkingEnabled"`
	QoSMarkingAllowListAppIdentifiers []string `plist:"QoSMarkingAllowListAppIdentifiers"`
}
