package serialize

import (
	"fmt"

	"howett.net/plist"
)

// DeviceManagementProfile represents the interface for all device management profiles.
type DeviceManagementProfile interface{}

// SerializeDeviceManagementProfileToPLIST converts a DeviceManagementProfile into an XML formatted PLIST string.
func SerializeDeviceManagementProfileToPLIST(profile DeviceManagementProfile) (string, error) {
	// Serialize the configuration to plist in XML format
	output, err := plist.MarshalIndent(profile, plist.XMLFormat, "    ")
	if err != nil {
		return "", fmt.Errorf("error marshalling to PLIST XML: %v", err)
	}
	return string(output), nil
}
