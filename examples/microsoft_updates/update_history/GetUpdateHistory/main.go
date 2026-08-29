package main

import (
	"context"
	"fmt"
	"log"

	microsoft_updates "github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates"
)

func main() {
	c, err := microsoft_updates.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.MicrosoftUpdatesAPI.UpdateHistory.GetUpdateHistoryV1(context.Background())
	if err != nil {
		log.Fatalf("UpdateHistory.GetUpdateHistoryV1: %v", err)
	}

	fmt.Printf("Office for Mac update history (%d entries)\n", len(resp.Entries))
	fmt.Println("==========================================")
	for _, entry := range resp.Entries {
		archived := ""
		if entry.Archived {
			archived = " [ARCHIVED]"
		}
		fmt.Printf("%-25s version=%-10s%s\n", entry.ReleaseDate, entry.Version, archived)
		if entry.BusinessProSuiteDownload != "" {
			fmt.Printf("  suite (w/Teams): %s\n", entry.BusinessProSuiteDownload)
		}
		if entry.SuiteDownload != "" {
			fmt.Printf("  suite (no Teams): %s\n", entry.SuiteDownload)
		}
	}
}
