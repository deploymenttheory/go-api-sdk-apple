package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/itunes"
	"github.com/deploymenttheory/go-sdk-appleservices/itunes/itunes_api/search"
)

func main() {
	c, err := itunes.NewClient(itunes.WithDebug())
	if err != nil {
		log.Fatalf("Error creating iTunes client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	fmt.Println("=== Lookup UPC with Tracks Example ===")
	fmt.Println("Looking up album by UPC (720642462928) including tracks:")

	result, _, err := c.ItunesAPI.Search.LookupV1(ctx, &search.LookupOptions{
		UPC:    "720642462928",
		Entity: search.EntitySong,
	})
	if err != nil {
		log.Fatalf("Error looking up UPC with tracks: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
