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

	fmt.Println("=== Japanese Language Search Example ===")
	fmt.Println("Searching for music with Japanese language setting (term: '宇多田ヒカル'):")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:    "宇多田ヒカル",
		Media:   search.MediaMusic,
		Entity:  search.EntityMusicTrack,
		Country: "JP",
		Limit:   10,
		Lang:    "ja_jp",
	})
	if err != nil {
		log.Fatalf("Error searching with Japanese language: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
