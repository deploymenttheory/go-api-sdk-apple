package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Device Management Services ===")

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

	opts := &devicemanagement.RequestQueryOptions{
		Fields: []string{
			devicemanagement.FieldServerName,
			devicemanagement.FieldServerType,
			devicemanagement.FieldCreatedDateTime,
			devicemanagement.FieldUpdatedDateTime,
		},
		Limit: 100,
	}

	response, _, err := c.AXMAPI.DeviceManagement.GetV1(ctx, opts)
	if err != nil {
		log.Fatalf("Error getting MDM servers: %v", err)
	}

	fmt.Printf("Found %d MDM servers\n", len(response.Data))

	for i, server := range response.Data {
		fmt.Printf("\nMDM Server %d:\n", i+1)
		fmt.Printf("  ID: %s\n", server.ID)
		fmt.Printf("  Type: %s\n", server.Type)
		if server.Attributes != nil {
			fmt.Printf("  Name: %s\n", server.Attributes.ServerName)
			fmt.Printf("  Server Type: %s\n", server.Attributes.ServerType)
			if server.Attributes.CreatedDateTime != nil {
				fmt.Printf("  Created: %s\n", server.Attributes.CreatedDateTime.Format(time.RFC3339))
			}
			if server.Attributes.UpdatedDateTime != nil {
				fmt.Printf("  Updated: %s\n", server.Attributes.UpdatedDateTime.Format(time.RFC3339))
			}
		}
		if server.Relationships != nil && server.Relationships.Devices != nil {
			fmt.Printf("  Devices Link: %s\n", server.Relationships.Devices.Links.Self)
		}
	}

	if response.Links != nil && response.Links.Next != "" {
		fmt.Printf("\nNext page: %s\n", response.Links.Next)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
