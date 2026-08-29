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

	fmt.Println("=== Debug JSON Parsing ===")

	result, _, err := c.ItunesAPI.Search.SearchV1(ctx, &search.SearchOptions{
		Term:  "Jack Johnson",
		Media: search.MediaMusic,
		Limit: 2,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Printf("Parsed Response: %s\n", string(jsonData))

	// Verify manual JSON round-trip using the SDK model type.
	rawJSON := `{"resultCount":1,"results":[{"wrapperType":"track","kind":"song","artistId":909253,"artistName":"Jack Johnson","trackName":"Test Song","primaryGenreName":"Rock"}]}`

	var parsed search.SearchResponse
	if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
		log.Fatalf("JSON parse error: %v", err)
	}

	testJSON, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling test parse: %v", err)
	}
	fmt.Printf("Test Parse Response: %s\n", string(testJSON))
}
