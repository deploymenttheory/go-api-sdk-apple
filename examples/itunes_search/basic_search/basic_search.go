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

	fmt.Println("=== Basic Search Example ===")
	fmt.Println("Searching for 'Jack Johnson' music with limit 5:")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:  "Jack Johnson",
		Media: search.MediaMusic,
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Error searching for music: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
