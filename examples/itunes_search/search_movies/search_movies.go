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

	fmt.Println("=== Movie Search Example ===")
	fmt.Println("Searching for movies with term 'star wars':")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:     "star wars",
		Media:    search.MediaMovie,
		Entity:   search.EntityMovie,
		Country:  "US",
		Limit:    10,
		Lang:     "en_us",
		Explicit: search.ExplicitYes,
	})
	if err != nil {
		log.Fatalf("Error searching for movies: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
