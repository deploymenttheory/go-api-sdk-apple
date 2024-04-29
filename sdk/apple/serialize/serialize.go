package serialize

import (
	"fmt"

	"howett.net/plist"
)

// SerializeDeviceManagementProfileToPLIST converts a DeviceManagementProfile into an XML string.
func SerializeDeviceManagementProfileToPLIST(profile string) (string, error) {
	// Serialize the configuration to plist in XML format
	output, err := plist.MarshalIndent(profile, plist.XMLFormat, "    ")
	if err != nil {
		return "", fmt.Errorf("error marshalling to PLIST XML: %v", err)
	}
	return string(output), nil
}
