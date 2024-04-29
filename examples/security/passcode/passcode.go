package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/go-mdm-sdk-apple/examples/helpers"
	"github.com/deploymenttheory/go-mdm-sdk-apple/sdk/apple/devicemanagementprofiles/security/generate"

	"github.com/deploymenttheory/go-mdm-sdk-apple/sdk/apple/serialize"
)

// main generates a passcode profile with the provided options and saves it to a .mobileconfig file.
func main() {
	options := GeneratePasscodeProfile()

	passcodeProfile := generate.CreatePasscodeDeviceManagementProfile(options)

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

// GeneratePasscodeProfile generates a passcode profile with the provided options.
func GeneratePasscodeProfile() generate.PasscodeDeviceManagementProfileConfigOptions {
	return generate.PasscodeDeviceManagementProfileConfigOptions{
		PayloadDisplayName:  "passcode",
		AllowSimple:         true,
		ForcePIN:            true,
		MaxFailedAttempts:   5,
		MaxGracePeriod:      1,
		MaxInactivity:       2,
		MaxPINAgeInDays:     30,
		MinLength:           8,
		PinHistory:          2,
		RequireAlphanumeric: false,
	}
}
