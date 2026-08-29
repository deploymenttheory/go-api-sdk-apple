package main

import (
	"context"
	"fmt"
	"log"

	microsoft_updates "github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates"
	"github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates/microsoft_updates_api/edge"
)

func printRelease(channel string, r *edge.EdgeRelease) {
	if r == nil {
		fmt.Printf("  %-10s  (no data)\n", channel)
		return
	}
	fmt.Printf("  %-10s version=%-25s published=%s\n", channel, r.Version, r.PublishedTime)
	for _, a := range r.Artifacts {
		fmt.Printf("             artifact=%-8s url=%s\n", a.ArtifactName, a.Location)
	}
}

func main() {
	c, err := microsoft_updates.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.MicrosoftUpdatesAPI.Edge.GetAllChannelsV1(context.Background())
	if err != nil {
		log.Fatalf("Edge.GetAllChannelsV1: %v", err)
	}

	fmt.Println("Microsoft Edge — all channels (macOS)")
	fmt.Println("=====================================")
	printRelease("stable", resp.Stable)
	printRelease("beta", resp.Beta)
	printRelease("dev", resp.Dev)
	printRelease("canary", resp.Canary)
}
