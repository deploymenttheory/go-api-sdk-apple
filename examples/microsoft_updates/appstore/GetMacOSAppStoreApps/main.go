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

	resp, err := c.MicrosoftUpdatesAPI.AppStoreMacOS.GetAllAppsV1(context.Background())
	if err != nil {
		log.Fatalf("AppStoreMacOS.GetAllAppsV1: %v", err)
	}

	fmt.Printf("macOS App Store — Microsoft apps (%d found)\n", resp.ResultCount)
	fmt.Println("===========================================")
	for _, app := range resp.Results {
		fmt.Printf("%-45s version=%-15s bundle=%s\n",
			app.TrackName,
			app.Version,
			app.BundleID,
		)
	}
}
