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

	fmt.Println("=== Advanced Lookup Example ===")
	fmt.Println("Looking up recent albums for multiple artists (Jack Johnson + Weezer) in the US:")

	result, _, err := c.ItunesAPI.Search.LookupV1(ctx, &search.LookupOptions{
		AMGArtistIDs: []string{"468749", "5723"},
		Entity:       search.EntityAlbum,
		Limit:        10,
		Sort:         search.SortRecent,
		Country:      "US",
	})
	if err != nil {
		log.Fatalf("Error performing advanced lookup: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
