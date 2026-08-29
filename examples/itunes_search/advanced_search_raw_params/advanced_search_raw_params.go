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

	fmt.Println("=== Advanced Search Example ===")
	fmt.Println("Searching for Jack Johnson music tracks with artist attribute, country, and language:")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:      "Jack Johnson",
		Media:     search.MediaMusic,
		Entity:    search.EntityMusicTrack,
		Attribute: search.AttributeArtistTerm,
		Limit:     25,
		Country:   "US",
		Lang:      "en_us",
		Explicit:  search.ExplicitNo,
		Version:   2,
	})
	if err != nil {
		log.Fatalf("Error performing advanced search: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
