package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Device Serial Numbers for Device Management Service ===")

	keyID := "44f6a58a-xxxx-4cab-xxxx-d071a3c36a42"
	issuerID := "BUSINESSAPI.3bb3a62b-xxxx-4802-xxxx-a69b86201c5a"
	privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
your-abm-api-key
-----END EC PRIVATE KEY-----`

	privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	c, err := axm.NewClient(keyID, issuerID, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	mdmServerID := "1F97349736CF4614A94F624E705841AD"

	opts := &devicemanagement.RequestQueryOptions{
		Limit: 100,
	}

	response, _, err := c.AXMAPI.DeviceManagement.GetDeviceSerialNumbersByServerIDV1(ctx, mdmServerID, opts)
	if err != nil {
		log.Fatalf("Error getting device serial numbers: %v", err)
	}

	fmt.Printf("Found %d devices assigned to MDM server %s\n", len(response.Data), mdmServerID)

	for i, linkage := range response.Data {
		fmt.Printf("  %d. Device ID: %s (Type: %s)\n", i+1, linkage.ID, linkage.Type)
	}

	if response.Links != nil && response.Links.Next != "" {
		fmt.Printf("\nNext page: %s\n", response.Links.Next)
	}

	if response.Meta != nil && response.Meta.Paging != nil {
		fmt.Printf("\nPagination - Limit: %d", response.Meta.Paging.Limit)
		if response.Meta.Paging.NextCursor != "" {
			fmt.Printf(", Next Cursor: %s", response.Meta.Paging.NextCursor)
		}
		fmt.Println()
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
