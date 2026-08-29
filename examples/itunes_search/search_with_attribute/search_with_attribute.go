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

	fmt.Println("=== Search with Attribute Example ===")
	fmt.Println("Searching for music tracks by artist 'Taylor Swift' using artistTerm attribute:")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:      "Taylor Swift",
		Media:     search.MediaMusic,
		Entity:    search.EntityMusicTrack,
		Attribute: search.AttributeArtistTerm,
		Country:   "US",
		Limit:     25,
		Lang:      "en_us",
		Explicit:  search.ExplicitNo,
	})
	if err != nil {
		log.Fatalf("Error searching with attribute: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
