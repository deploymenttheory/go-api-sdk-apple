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

	fmt.Println("=== Lookup by ISBN Example ===")
	fmt.Println("Looking up book by 13-digit ISBN (9780316069359):")

	result, _, err := c.ItunesAPI.Search.LookupV1(ctx, &search.LookupOptions{
		ISBN: "9780316069359",
	})
	if err != nil {
		log.Fatalf("Error looking up ISBN: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
