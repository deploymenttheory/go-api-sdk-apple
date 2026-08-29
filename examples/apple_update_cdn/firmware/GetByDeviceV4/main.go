package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
)

func main() {
	fmt.Println("=== Apple Update CDN - Get Firmware by Device V4 ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	// Examples: "Mac14,3" (Mac mini M2), "iPhone15,2" (iPhone 14 Pro), "iPad14,4" (iPad mini 6th gen)
	identifier := "Mac14,3"

	response, _, err := c.AppleUpdateCDNAPI.Firmware.GetByDeviceV4(ctx, identifier)
	if err != nil {
		log.Fatalf("Error getting firmware for device %s: %v", identifier, err)
	}

	fmt.Printf("Device: %s (%s)\n", response.Identifier, response.Name)
	fmt.Printf("  Board Config: %s\n", response.BoardConfig)
	fmt.Printf("  Platform: %s\n", response.Platform)
	fmt.Printf("  Firmware versions: %d\n", len(response.Firmwares))

	for i, fw := range response.Firmwares {
		fmt.Printf("\nFirmware %d:\n", i+1)
		fmt.Printf("  Version: %s\n", fw.Version)
		fmt.Printf("  Build ID: %s\n", fw.BuildID)
		fmt.Printf("  Signed: %v\n", fw.Signed)
		fmt.Printf("  File Size: %d bytes (%.2f GB)\n", fw.FileSize, float64(fw.FileSize)/1e9)
		fmt.Printf("  Released: %s\n", fw.ReleaseDate.Format("2006-01-02"))
		fmt.Printf("  Uploaded: %s\n", fw.UploadDate.Format("2006-01-02"))
		fmt.Printf("  SHA-1: %s\n", fw.SHA1Sum)
		fmt.Printf("  SHA-256: %s\n", fw.SHA256Sum)
		fmt.Printf("  MD5: %s\n", fw.MD5Sum)
		fmt.Printf("  URL: %s\n", fw.URL)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
