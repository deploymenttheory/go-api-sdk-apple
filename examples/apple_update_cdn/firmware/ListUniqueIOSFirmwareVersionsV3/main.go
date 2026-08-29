package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
)

func main() {
	fmt.Println("=== Apple Update CDN - List Unique iOS Firmware Versions V3 ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	// Returns one entry per unique iOS build ID, sorted newest-first.
	// Each iPhone model has its own device-specific IPSW (unlike Universal Mac IPSWs).
	versions, _, err := c.AppleUpdateCDNAPI.Firmware.ListUniqueIOSFirmwareVersionsV3(ctx)
	if err != nil {
		log.Fatalf("Error listing unique iOS firmware versions: %v", err)
	}

	fmt.Printf("Found %d unique iOS firmware versions\n", len(versions))

	for i, fw := range versions {
		fmt.Printf("\nVersion %d:\n", i+1)
		fmt.Printf("  Version: %s\n", fw.Version)
		fmt.Printf("  Build ID: %s\n", fw.BuildID)
		fmt.Printf("  Signed: %v\n", fw.Signed)
		fmt.Printf("  Size: %d bytes (%.2f GB)\n", fw.Size, float64(fw.Size)/1e9)
		fmt.Printf("  Released: %s\n", fw.ReleaseDate.Format("2006-01-02"))
		fmt.Printf("  SHA-1: %s\n", fw.SHA1Sum)
		fmt.Printf("  URL: %s\n", fw.URL)
	}

	jsonData, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
