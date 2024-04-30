package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/examples/helpers"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/networking/generate"

	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/serialize"
)

// main generates a passcode profile with the provided options and saves it to a .mobileconfig file.
func main() {
	options := GenerateCellularProfile()

	passcodeProfile := generate.CreateCellularMobileDeviceManagementProfile(options)

	// Serialize the profile to plist in XML format
	profileXML, err := serialize.SerializeDeviceManagementProfileToPLIST(passcodeProfile)
	if err != nil {
		fmt.Println("Error serializing profile to PLIST XML:", err)
		return
	}

	// Infer the filename from PayloadDisplayName and add the ".mobileconfig" suffix
	fileName := options.PayloadDisplayName + ".mobileconfig"

	// Rest of your file handling code...
	fmt.Print("Enter the directory to save the .mobileconfig file (e.g., /path/to/folder): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dirPath := scanner.Text()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Join directory path and filename
	filePath := filepath.Join(dirPath, fileName)

	// Save the profile to the specified file path
	if err := helpers.SaveToFile(filePath, profileXML); err != nil {
		fmt.Println("Error saving profile to file:", err)
		return
	}
}

// GenerateCellularProfile generates a cellular profile with the provided options.
func GenerateCellularProfile() generate.CellularMobileDeviceManagementProfileConfigOptions {
	return generate.CellularMobileDeviceManagementProfileConfigOptions{
		PayloadDisplayName: "cellular",
		AttachAPN: generate.AttachAPNConfigOptions{
			Name:                "example.com",
			AuthenticationType:  "PAP",  // Assuming default if not specified
			AllowedProtocolMask: 3,      // Assuming IPv4 & IPv6 are allowed if not specified
			Username:            "user", // Example values
			Password:            "password",
		},
		APNs: []generate.APNConfigOptions{
			// Example APN configuration, you can add multiple APN configurations
			{
				Name:                                 "internet",
				AllowedProtocolMask:                  3, // Example: IPv4 & IPv6
				AllowedProtocolMaskInDomesticRoaming: 3,
				AllowedProtocolMaskInRoaming:         3,
				AuthenticationType:                   "PAP",
				EnableXLAT464:                        false,
				Password:                             "password",
				ProxyPort:                            8080,
				ProxyServer:                          "proxy.example.com",
				Username:                             "user",
			},
		},
	}
}
