package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
)

func main() {
	fmt.Println("=== Apple Update CDN - Get File Metadata V1 (HEAD request) ===")

	c, err := apple_update_cdn.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	// Issues a HEAD request — the full IPSW file is not downloaded.
	// Apple's CDN returns checksums and file size in the response headers.
	ipswURL := "https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw"

	meta, _, err := c.AppleUpdateCDNAPI.CDN.GetFileMetadataV1(ctx, ipswURL)
	if err != nil {
		log.Fatalf("Error getting file metadata: %v", err)
	}

	fmt.Printf("File metadata for:\n  %s\n\n", meta.URL)
	fmt.Printf("  Content-Type:   %s\n", meta.ContentType)
	fmt.Printf("  Content-Length: %d bytes (%.2f GB)\n", meta.ContentLength, float64(meta.ContentLength)/1e9)
	fmt.Printf("  SHA-1:          %s\n", meta.SHA1)
	fmt.Printf("  SHA-256:        %s\n", meta.SHA256)
	fmt.Printf("  ETag:           %s\n", meta.ETag)
	if !meta.LastModified.IsZero() {
		fmt.Printf("  Last-Modified:  %s\n", meta.LastModified.Format("2006-01-02 15:04:05 UTC"))
	}

	jsonData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
