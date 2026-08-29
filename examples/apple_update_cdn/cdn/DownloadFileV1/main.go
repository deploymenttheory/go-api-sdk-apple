package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/tools/download_progress"
)

func main() {
	fmt.Println("=== Apple Update CDN - Download IPSW File V1 ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	// The IPSW URL to download. Obtain this from ListUniqueMacFirmwareVersionsV3,
	// ListUniqueIOSFirmwareVersionsV3, or any firmware listing method.
	ipswURL := "https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw"

	// Destination path — adjust to a location with sufficient free space.
	// macOS IPSW files are typically 15–22 GB.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	destPath := filepath.Join(homeDir, "Downloads", "UniversalMac_26.4.1_25E253_Restore.ipsw")

	fmt.Printf("Downloading to: %s\n\n", destPath)

	// Set up a terminal progress bar that writes to stderr.
	bar := download_progress.New(os.Stderr)
	filename := filepath.Base(destPath)
	progressFn := func(written, total int64) {
		bar.Callback(filename)(written, total)
	}

	result, _, err := c.AppleUpdateCDNAPI.CDN.DownloadFileV1(ctx, ipswURL, destPath, progressFn)
	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	fmt.Printf("\nDownload complete:\n")
	fmt.Printf("  Destination:   %s\n", result.DestPath)
	fmt.Printf("  Bytes written: %d (%.2f GB)\n", result.BytesWritten, float64(result.BytesWritten)/1e9)
	fmt.Printf("  Duration:      %s\n", result.Duration.Round(1e9))
	fmt.Printf("  SHA-1:         %s\n", result.SHA1)
	fmt.Printf("  SHA-256:       %s\n", result.SHA256)
	fmt.Printf("  Verified:      %v\n", result.Verified)

	if !result.Verified {
		log.Fatalf("Checksum verification failed — the downloaded file may be corrupt.")
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling result: %v", err)
	}
	fmt.Println("\nFull result:")
	fmt.Println(string(jsonData))
}
