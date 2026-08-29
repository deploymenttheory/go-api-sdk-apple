package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
)

func main() {
	fmt.Println("=== Apple Update CDN - List All iOS (iPhone) Firmware V3 ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	response, _, err := c.AppleUpdateCDNAPI.Firmware.ListAllIOSFirmwareV3(ctx)
	if err != nil {
		log.Fatalf("Error listing iOS firmware: %v", err)
	}

	fmt.Printf("Found %d iPhone device models\n", len(response.Devices))

	for id, device := range response.Devices {
		fmt.Printf("\nDevice: %s\n", id)
		fmt.Printf("  Name: %s\n", device.Name)
		fmt.Printf("  Board Config: %s\n", device.BoardConfig)
		fmt.Printf("  Platform: %s\n", device.Platform)
		fmt.Printf("  Firmware versions: %d\n", len(device.Firmwares))

		for _, fw := range device.Firmwares {
			fmt.Printf("    - %s (%s) signed=%v size=%d bytes\n",
				fw.Version, fw.BuildID, fw.Signed, fw.Size)
		}
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
