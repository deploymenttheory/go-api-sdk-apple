package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/apple_update_cdn_api/cdn"
)

func main() {
	fmt.Println("=== Apple Update CDN - Parse IPSW URL ===")

	// ParseURL is a pure parsing operation — no HTTP request is made and no
	// client is required.
	rawURL := "https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw"

	info, err := cdn.ParseURL(rawURL)
	if err != nil {
		log.Fatalf("Failed to parse CDN URL: %v", err)
	}

	fmt.Printf("Parsed URL components:\n")
	fmt.Printf("  Catalog Release: %s\n", info.CatalogRelease)
	fmt.Printf("  Asset Type:      %s\n", info.AssetType)
	fmt.Printf("  Asset ID:        %s\n", info.AssetID)
	fmt.Printf("  UUID:            %s\n", info.UUID)
	fmt.Printf("  Filename:        %s\n", info.Filename)
	fmt.Printf("  Platform:        %s\n", info.Platform)
	fmt.Printf("  Version:         %s\n", info.Version)
	fmt.Printf("  Build:           %s\n", info.Build)
	fmt.Printf("  Restore Type:    %s\n", info.RestoreType)

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
