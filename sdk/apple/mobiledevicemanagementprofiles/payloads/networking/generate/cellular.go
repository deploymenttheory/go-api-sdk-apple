package generate

import (
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/networking/models"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/helpers"
)

// CellularMobileDeviceManagementProfileConfigOptions holds customizable options for creating a cellular profile.
type CellularMobileDeviceManagementProfileConfigOptions struct {
	PayloadDisplayName string
	AttachAPN          AttachAPNConfigOptions
	APNs               []APNConfigOptions
}

// AttachAPNConfigOptions defines the configuration options for the Attach APN.
type AttachAPNConfigOptions struct {
	Name                string
	AuthenticationType  string
	AllowedProtocolMask int
	Username            string
	Password            string
}

// APNConfigOptions defines the configuration options for individual APNs.
type APNConfigOptions struct {
	Name                                 string
	AllowedProtocolMask                  int
	AllowedProtocolMaskInDomesticRoaming int
	AllowedProtocolMaskInRoaming         int
	AuthenticationType                   string
	EnableXLAT464                        bool
	Password                             string
	ProxyPort                            int
	ProxyServer                          string
	Username                             string
}

// CreateCellularMobileDeviceManagementProfile creates a device management profile for cellular settings using provided configuration options.
func CreateCellularMobileDeviceManagementProfile(options CellularMobileDeviceManagementProfileConfigOptions) models.CellularMobileDeviceManagementProfile {
	// Map the AttachAPNConfigOptions and APNConfigOptions to their respective models
	attachAPN := models.CellularMobileDeviceManagementProfileSubsetCellularAttachAPN{
		AllowedProtocolMask: options.AttachAPN.AllowedProtocolMask,
		AuthenticationType:  options.AttachAPN.AuthenticationType,
		Name:                options.AttachAPN.Name,
		Password:            options.AttachAPN.Password,
		Username:            options.AttachAPN.Username,
	}

	apns := make([]models.CellularMobileDeviceManagementProfileSubsetCellularAPNsItem, len(options.APNs))
	for i, apnOption := range options.APNs {
		apns[i] = models.CellularMobileDeviceManagementProfileSubsetCellularAPNsItem{
			Name:                                 apnOption.Name,
			AllowedProtocolMask:                  apnOption.AllowedProtocolMask,
			AllowedProtocolMaskInDomesticRoaming: apnOption.AllowedProtocolMaskInDomesticRoaming,
			AllowedProtocolMaskInRoaming:         apnOption.AllowedProtocolMaskInRoaming,
			AuthenticationType:                   apnOption.AuthenticationType,
			EnableXLAT464:                        apnOption.EnableXLAT464,
			Password:                             apnOption.Password,
			ProxyPort:                            apnOption.ProxyPort,
			ProxyServer:                          apnOption.ProxyServer,
			Username:                             apnOption.Username,
		}
	}

	return models.CellularMobileDeviceManagementProfile{
		Version: "1.0",
		PayloadContent: []models.CellularMobileDeviceManagementProfileSubsetPayload{
			{
				AttachAPN:         attachAPN,
				APNs:              apns,
				PayloadIdentifier: "com.gomdmsdkapple.cellularpayload",
				PayloadType:       "com.apple.cellular",
				PayloadUUID:       helpers.GenerateUUID(),
				PayloadVersion:    1,
			},
		},
		PayloadDisplayName: options.PayloadDisplayName,
		PayloadIdentifier:  "com.gomdmsdkapple.cellularpayload",
		PayloadType:        "Configuration",
		PayloadUUID:        helpers.GenerateUUID(),
		PayloadVersion:     1,
	}
}
