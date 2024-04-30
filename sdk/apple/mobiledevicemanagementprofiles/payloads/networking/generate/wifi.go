package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/networking/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// WiFiMobileDeviceManagementProfileConfigOptions holds customizable options for creating a WiFi profile.
type WiFiMobileDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName                 string
	SSID                               string
	AutoJoin                           bool
	CaptiveBypass                      bool
	DisableAssociationMACRandomization bool
	DisplayedOperatorName              string
	DomainName                         string
	EnableIPv6                         bool
	EncryptionType                     string
	HIDDEN_NETWORK                     bool
	IsHotspot                          bool
	MCCAndMNCs                         []string
	NAIRealmNames                      []string
	Password                           string
	PayloadCertificateUUID             string
	ProxyPACURL                        string
	ProxyServer                        string
	ProxyServerPort                    int
	ProxyType                          string
	RoamingConsortiumOIs               []string
	ServiceProviderRoamingEnabled      bool
	SetupModes                         []string
	TLSCertificateRequired             bool
	TLSMinimumVersion                  string
	TLSMaximumVersion                  string
	TLSTrustedCertificates             []string
	TLSTrustedServerNames              []string
	UserName                           string
	UserPassword                       string
	EAPClientConfig                    EAPClientConfigOptions
	QoSMarking                         QoSMarkingConfigOptions
}

// EAPClientConfigOptions defines the configuration options for EAP Client.
type EAPClientConfigOptions struct {
	AcceptEAPTypes                 []int
	EAPFASTProvisionPAC            bool
	EAPFASTProvisionPACAnonymously bool
	EAPFASTUsePAC                  bool
	EAPSIMNumberOfRANDs            int
	OneTimeUserPassword            bool
	OuterIdentity                  string
	TTLSInnerAuthentication        string
}

// QoSMarkingConfigOptions defines the configuration options for QoS marking policy.
type QoSMarkingConfigOptions struct {
	QoSMarkingAppleAudioVideoCalls    bool
	QoSMarkingEnabled                 bool
	QoSMarkingAllowListAppIdentifiers []string
}

// CreateWiFiMobileDeviceManagementProfile creates a device management profile for WiFi settings using provided configuration options.
func CreateWiFiMobileDeviceManagementProfile(options WiFiMobileDeviceManagementProfileConfigOptions) models.WiFiMobileDeviceManagementProfile {
	// Mapping the WiFi-specific configurations
	return models.WiFiMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.WiFiMobileDeviceManagementProfileSubsetPayload{
			{
				SSID_STR:                           options.SSID,
				AutoJoin:                           &options.AutoJoin,
				CaptiveBypass:                      &options.CaptiveBypass,
				DisableAssociationMACRandomization: &options.DisableAssociationMACRandomization,
				DisplayedOperatorName:              options.DisplayedOperatorName,
				DomainName:                         options.DomainName,
				EnableIPv6:                         &options.EnableIPv6,
				EncryptionType:                     options.EncryptionType,
				HIDDEN_NETWORK:                     &options.HIDDEN_NETWORK,
				IsHotspot:                          &options.IsHotspot,
				MCCAndMNCs:                         options.MCCAndMNCs,
				NAIRealmNames:                      options.NAIRealmNames,
				Password:                           options.Password,
				PayloadCertificateUUID:             options.PayloadCertificateUUID,
				ProxyPACURL:                        options.ProxyPACURL,
				ProxyServer:                        options.ProxyServer,
				ProxyServerPort:                    options.ProxyServerPort,
				ProxyType:                          options.ProxyType,
				RoamingConsortiumOIs:               options.RoamingConsortiumOIs,
				ServiceProviderRoamingEnabled:      &options.ServiceProviderRoamingEnabled,
				SetupModes:                         options.SetupModes,
				TLSCertificateRequired:             &options.TLSCertificateRequired,
				TLSMinimumVersion:                  options.TLSMinimumVersion,
				TLSMaximumVersion:                  options.TLSMaximumVersion,
				TLSTrustedCertificates:             options.TLSTrustedCertificates,
				TLSTrustedServerNames:              options.TLSTrustedServerNames,
				UserName:                           options.UserName,
				UserPassword:                       options.UserPassword,
				EAPClientConfiguration:             mapEAPClientConfig(options.EAPClientConfig),
				QoSMarkingPolicy:                   mapQoSMarking(options.QoSMarking),
				PayloadIdentifier:                  "com.gomdmsdkapple.wifipayload",
				PayloadType:                        "com.apple.wifi.managed",
				PayloadUUID:                        helpers.GenerateUUID(),
				PayloadVersion:                     1,
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.wifipayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}

func mapEAPClientConfig(options EAPClientConfigOptions) models.EAPClientConfiguration {
	return models.EAPClientConfiguration{
		AcceptEAPTypes:                 options.AcceptEAPTypes,
		EAPFASTProvisionPAC:            &options.EAPFASTProvisionPAC,
		EAPFASTProvisionPACAnonymously: &options.EAPFASTProvisionPACAnonymously,
		EAPFASTUsePAC:                  &options.EAPFASTUsePAC,
		EAPSIMNumberOfRANDs:            options.EAPSIMNumberOfRANDs,
		OneTimeUserPassword:            &options.OneTimeUserPassword,
		OuterIdentity:                  options.OuterIdentity,
		TTLSInnerAuthentication:        options.TTLSInnerAuthentication,
	}
}

func mapQoSMarking(options QoSMarkingConfigOptions) models.QoSMarkingPolicy {
	return models.QoSMarkingPolicy{
		QoSMarkingAppleAudioVideoCalls:    &options.QoSMarkingAppleAudioVideoCalls,
		QoSMarkingEnabled:                 &options.QoSMarkingEnabled,
		QoSMarkingAllowListAppIdentifiers: options.QoSMarkingAllowListAppIdentifiers,
	}
}
