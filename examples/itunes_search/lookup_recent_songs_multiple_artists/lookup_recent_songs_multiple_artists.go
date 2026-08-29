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

	fmt.Println("=== Lookup Recent Songs for Multiple Artists Example ===")
	fmt.Println("Looking up 5 most recent songs for multiple artists by AMG artist IDs (468749, 5723):")

	result, _, err := c.ItunesAPI.Search.LookupV1(ctx, &search.LookupOptions{
		AMGArtistIDs: []string{"468749", "5723"},
		Entity:       search.EntitySong,
		Limit:        5,
		Sort:         search.SortRecent,
	})
	if err != nil {
		log.Fatalf("Error looking up recent songs for multiple AMG artist IDs: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
