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

	fmt.Println("=== Podcast Search Example ===")
	fmt.Println("Searching for podcasts with term 'tech talk':")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:    "tech talk",
		Media:   search.MediaPodcast,
		Entity:  search.EntityPodcast,
		Country: "US",
		Limit:   10,
		Lang:    "en_us",
	})
	if err != nil {
		log.Fatalf("Error searching for podcasts: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
