package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/apple_update_cdn_api/gdmf"
)

func main() {
	fmt.Println("=== Apple GDMF - Get Public Firmware Versions V2 ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	response, _, err := c.AppleUpdateCDNAPI.GDMF.GetPublicVersionsV2(ctx)
	if err != nil {
		log.Fatalf("Error getting public firmware versions: %v", err)
	}

	if response.PublicAssetSets != nil {
		fmt.Println("\n--- PublicAssetSets ---")
		printPlatformAssets("macOS", response.PublicAssetSets.MacOS)
		printPlatformAssets("iOS", response.PublicAssetSets.IOS)
		printPlatformAssets("visionOS", response.PublicAssetSets.VisionOS)
	}

	if response.AssetSets != nil {
		fmt.Println("\n--- AssetSets (seed/additional) ---")
		printPlatformAssets("macOS", response.AssetSets.MacOS)
		printPlatformAssets("iOS", response.AssetSets.IOS)
		printPlatformAssets("visionOS", response.AssetSets.VisionOS)
	}

	if response.PublicBackgroundSecurityImprovements != nil {
		fmt.Println("\n--- PublicBackgroundSecurityImprovements (Rapid Security Responses) ---")
		printPlatformAssets("macOS", response.PublicBackgroundSecurityImprovements.MacOS)
		printPlatformAssets("iOS", response.PublicBackgroundSecurityImprovements.IOS)
		printPlatformAssets("visionOS", response.PublicBackgroundSecurityImprovements.VisionOS)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}

func printPlatformAssets(platform string, entries []*gdmf.AssetEntry) {
	if len(entries) == 0 {
		return
	}
	fmt.Printf("\n  %s (%d versions):\n", platform, len(entries))
	for _, entry := range entries {
		fmt.Printf("    Version: %s  Build: %s  Posted: %s  Expires: %s  Devices: %d\n",
			entry.ProductVersion,
			entry.Build,
			entry.PostingDate,
			entry.ExpirationDate,
			len(entry.SupportedDevices),
		)
	}
}
