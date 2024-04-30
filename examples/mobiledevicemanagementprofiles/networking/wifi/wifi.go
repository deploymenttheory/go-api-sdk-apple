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
	options := GenerateWiFiProfile()

	passcodeProfile := generate.CreateWiFiMobileDeviceManagementProfile(options)

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

// GenerateWiFiProfile generates a cellular profile with the provided options.
func GenerateWiFiProfile() generate.WiFiMobileDeviceManagementProfileConfigOptions {
	return generate.WiFiMobileDeviceManagementProfileConfigOptions{
		PayloadDisplayName:                 "wifi",
		SSID:                               "Example",
		AutoJoin:                           true,
		CaptiveBypass:                      false,
		DisableAssociationMACRandomization: false,
		DisplayedOperatorName:              "Example Network",
		DomainName:                         "",
		EnableIPv6:                         true,
		EncryptionType:                     "WPA2",
		HIDDEN_NETWORK:                     true,
		IsHotspot:                          false,
		MCCAndMNCs:                         []string{},
		NAIRealmNames:                      []string{},
		Password:                           "Password123",
		PayloadCertificateUUID:             "",
		ProxyPACURL:                        "",
		ProxyServer:                        "",
		ProxyServerPort:                    0,
		ProxyType:                          "None",
		RoamingConsortiumOIs:               []string{},
		ServiceProviderRoamingEnabled:      false,
		SetupModes:                         []string{"System"},
		TLSCertificateRequired:             false,
		TLSMinimumVersion:                  "1.2",
		TLSMaximumVersion:                  "1.3",
		TLSTrustedCertificates:             []string{},
		TLSTrustedServerNames:              []string{},
		UserName:                           "",
		UserPassword:                       "",
		EAPClientConfig: generate.EAPClientConfigOptions{
			AcceptEAPTypes:                 []int{},
			EAPFASTProvisionPAC:            false,
			EAPFASTProvisionPACAnonymously: false,
			EAPFASTUsePAC:                  false,
			EAPSIMNumberOfRANDs:            0,
			OneTimeUserPassword:            false,
			OuterIdentity:                  "",
			TTLSInnerAuthentication:        "MSCHAPv2",
		},
		QoSMarking: generate.QoSMarkingConfigOptions{
			QoSMarkingAppleAudioVideoCalls:    false,
			QoSMarkingEnabled:                 false,
			QoSMarkingAllowListAppIdentifiers: []string{},
		},
	}
}
