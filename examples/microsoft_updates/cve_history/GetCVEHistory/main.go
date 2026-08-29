package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	microsoft_updates "github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates"
)

func main() {
	c, err := microsoft_updates.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.MicrosoftUpdatesAPI.CVEHistory.GetCVEHistoryV1(context.Background())
	if err != nil {
		log.Fatalf("CVEHistory.GetCVEHistoryV1: %v", err)
	}

	fmt.Printf("Office for Mac CVE history (%d entries)\n", len(resp.Entries))
	fmt.Println("=======================================")
	for _, entry := range resp.Entries {
		if len(entry.CVEs) == 0 {
			continue
		}
		fmt.Printf("%s  %s\n", entry.ReleaseDate, entry.Version)
		fmt.Printf("  CVEs: %s\n", strings.Join(entry.CVEs, ", "))
	}
}
