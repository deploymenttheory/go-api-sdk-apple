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

	fmt.Println("=== Lookup Albums by AMG Artist ID Example ===")
	fmt.Println("Looking up all albums for Jack Johnson (AMG ID 468749):")

	result, _, err := c.ItunesAPI.Search.LookupV1(ctx, &search.LookupOptions{
		AMGArtistID: "468749",
		Entity:      search.EntityAlbum,
	})
	if err != nil {
		log.Fatalf("Error looking up albums for AMG artist ID: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
