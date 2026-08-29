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

	resp, err := c.MicrosoftUpdatesAPI.Standalone.GetLatestV1(context.Background())
	if err != nil {
		log.Fatalf("GetLatestV1: %v", err)
	}

	fmt.Printf("Found %d standalone packages\n\n", len(resp.Packages))
	for _, pkg := range resp.Packages {
		fmt.Printf("%-40s version=%-25s min_os=%s\n",
			pkg.Title,
			pkg.FullVersion,
			pkg.MinimumOS,
		)
		if pkg.Location != "" {
			fmt.Printf("  full installer:   %s\n", pkg.Location)
		}
		if pkg.AppOnlyLocation != "" {
			fmt.Printf("  app-only update:  %s\n", pkg.AppOnlyLocation)
		}
	}
}
