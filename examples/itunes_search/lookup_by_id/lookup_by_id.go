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

	fmt.Println("=== Lookup by iTunes ID Example ===")
	fmt.Println("Looking up Jack Johnson by iTunes artist ID (909253):")

	result, _, err := c.ItunesAPI.Search.LookupV1(ctx, &search.LookupOptions{
		ID: 909253,
	})
	if err != nil {
		log.Fatalf("Error looking up ID: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
