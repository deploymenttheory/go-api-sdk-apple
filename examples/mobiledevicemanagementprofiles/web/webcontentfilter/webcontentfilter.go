package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/examples/helpers"
	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/mobiledevicemanagementprofiles/payloads/web/generate"

	"github.com/deploymenttheory/go-apple-sdk-devicemanagement/sdk/apple/serialize"
)

func main() {
	options := CreateWebContentFilterProfile()

	passcodeProfile := generate.CreateWebContentFilterMobileDeviceManagementProfile(options)

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

// CreateWebContentFilterProfile generates a profile with the provided options.
func CreateWebContentFilterProfile() generate.WebContentFilterMobileDeviceManagementProfileConfigOptions {
	return generate.WebContentFilterMobileDeviceManagementProfileConfigOptions{
		PayloadDisplayName:                        "webcontentfilter",                                  // Display name for the mobileconfig file.
		AutoFilterEnabled:                         true,                                                // Automatically enables filtering based on rules.
		FilterType:                                "Plugin",                                            // Use built-in type; alternate is 'Plugin'. macOS does not support "BuiltIn"
		PermittedURLs:                             []string{"https://example.company.com"},             // URLs that are explicitly allowed.
		DenylistURLs:                              []string{"https://notallowedname.company.com"},      // URLs that are explicitly blocked.
		FilterBrowsers:                            true,                                                // Filter traffic from web browsers.
		FilterSockets:                             true,                                                // Filter network socket traffic.
		FilterPackets:                             true,                                                // Filter network packets.
		ContentFilterUUID:                         "fb5d598f-0a96-4b77-9702-9edfc3417601",              // Unique identifier for the filter configuration.
		FilterDataProviderBundleIdentifier:        "com.example.filterservice",                         // Identifier for data provider bundle.
		FilterDataProviderDesignatedRequirement:   "signer=\"Example Corp\"",                           // Designated requirement for the data provider.
		FilterPacketProviderBundleIdentifier:      "com.example.packetfilterservice",                   // Identifier for packet filter provider.
		FilterPacketProviderDesignatedRequirement: "signer=\"Example Corp\"",                           // Designated requirement for the packet filter provider.
		Organization:                              "Example Corporation",                               // The organization managing the filter.
		ServerAddress:                             "https://example.server.com",                        // Address for the filtering service server.
		UserDefinedName:                           "Example Web Filter",                                // User-defined name for the filter setup.
		UserName:                                  "admin",                                             // Username for authentication with the filtering service.
		Password:                                  "password",                                          // Password for authentication with the filtering service.
		VendorConfig:                              map[string]interface{}{"additionalConfig": "value"}, // Additional vendor-specific configuration options.
	}
}
