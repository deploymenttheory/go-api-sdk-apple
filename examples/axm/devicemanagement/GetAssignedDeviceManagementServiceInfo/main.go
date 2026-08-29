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
	fmt.Println("=== Apple Business Manager - Get Assigned Device Management Service Information ===")

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

	deviceID := "XABC123X0ABC123X0"

	opts := &devicemanagement.RequestQueryOptions{
		Fields: []string{
			devicemanagement.FieldServerName,
			devicemanagement.FieldServerType,
			devicemanagement.FieldCreatedDateTime,
			devicemanagement.FieldUpdatedDateTime,
		},
	}

	response, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerInfoByDeviceIDV1(ctx, deviceID, opts)
	if err != nil {
		log.Fatalf("Error getting assigned server information for device %s: %v", deviceID, err)
	}

	fmt.Printf("Device ID: %s\n", deviceID)
	fmt.Printf("Assigned Server Information:\n")
	fmt.Printf("  Server ID: %s\n", response.Data.ID)
	fmt.Printf("  Type: %s\n", response.Data.Type)

	if response.Data.Attributes != nil {
		fmt.Printf("  Server Name: %s\n", response.Data.Attributes.ServerName)
		fmt.Printf("  Server Type: %s\n", response.Data.Attributes.ServerType)
		if response.Data.Attributes.CreatedDateTime != nil {
			fmt.Printf("  Created: %s\n", response.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
		}
		if response.Data.Attributes.UpdatedDateTime != nil {
			fmt.Printf("  Updated: %s\n", response.Data.Attributes.UpdatedDateTime.Format(time.RFC3339))
		}
	}

	if response.Data.Relationships != nil && response.Data.Relationships.Devices != nil {
		fmt.Printf("  Devices Relationship: %+v\n", response.Data.Relationships.Devices)
	}

	if response.Links != nil {
		fmt.Printf("  Self Link: %s\n", response.Links.Self)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
