package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Organization Devices ===")

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

	opts := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldProductFamily,
			devices.FieldStatus,
			devices.FieldAddedToOrgDateTime,
			devices.FieldUpdatedDateTime,
		},
		Limit: 100,
	}

	response, _, err := c.AXMAPI.Devices.GetV1(ctx, opts)
	if err != nil {
		log.Fatalf("Error getting organization devices: %v", err)
	}

	fmt.Printf("Found %d devices\n", len(response.Data))

	for i, device := range response.Data {
		fmt.Printf("\nDevice %d:\n", i+1)
		fmt.Printf("  ID: %s\n", device.ID)
		fmt.Printf("  Type: %s\n", device.Type)
		if device.Attributes != nil {
			fmt.Printf("  Serial: %s\n", device.Attributes.SerialNumber)
			fmt.Printf("  Model: %s\n", device.Attributes.DeviceModel)
			fmt.Printf("  Family: %s\n", device.Attributes.ProductFamily)
			fmt.Printf("  Status: %s\n", device.Attributes.Status)
			if device.Attributes.AddedToOrgDateTime != nil {
				fmt.Printf("  Added: %s\n", device.Attributes.AddedToOrgDateTime.Format(time.RFC3339))
			}
			if device.Attributes.UpdatedDateTime != nil {
				fmt.Printf("  Updated: %s\n", device.Attributes.UpdatedDateTime.Format(time.RFC3339))
			}
		}
	}

	if response.Links != nil && response.Links.Next != "" {
		fmt.Printf("\nNext page: %s\n", response.Links.Next)
	}

	if response.Meta != nil && response.Meta.Paging != nil {
		fmt.Printf("Pagination - Limit: %d\n", response.Meta.Paging.Limit)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
