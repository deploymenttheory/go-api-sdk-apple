package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
)

func main() {
	fmt.Println("=== Apple Update CDN - List All Firmware V3 (All Platforms) ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	response, _, err := c.AppleUpdateCDNAPI.Firmware.ListAllFirmwareV3(ctx)
	if err != nil {
		log.Fatalf("Error listing all firmware: %v", err)
	}

	fmt.Printf("Found %d device models (all platforms)\n", len(response.Devices))

	// Tally by platform prefix for a quick overview.
	counts := map[string]int{
		"Mac":     0,
		"iPhone":  0,
		"iPad":    0,
		"Other":   0,
	}
	for id := range response.Devices {
		switch {
		case len(id) >= 3 && id[:3] == "Mac":
			counts["Mac"]++
		case len(id) >= 6 && id[:6] == "iPhone":
			counts["iPhone"]++
		case len(id) >= 4 && id[:4] == "iPad":
			counts["iPad"]++
		default:
			counts["Other"]++
		}
	}

	fmt.Printf("\nDevice breakdown:\n")
	fmt.Printf("  Mac:     %d models\n", counts["Mac"])
	fmt.Printf("  iPhone:  %d models\n", counts["iPhone"])
	fmt.Printf("  iPad:    %d models\n", counts["iPad"])
	fmt.Printf("  Other:   %d models\n", counts["Other"])

	// Print details for the first 5 devices as a sample.
	i := 0
	for id, device := range response.Devices {
		if i >= 5 {
			break
		}
		fmt.Printf("\nDevice %d:\n", i+1)
		fmt.Printf("  Identifier: %s\n", id)
		fmt.Printf("  Name: %s\n", device.Name)
		fmt.Printf("  Firmware versions: %d\n", len(device.Firmwares))
		if len(device.Firmwares) > 0 {
			latest := device.Firmwares[0]
			fmt.Printf("  Latest: %s (%s) signed=%v\n", latest.Version, latest.BuildID, latest.Signed)
		}
		i++
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
