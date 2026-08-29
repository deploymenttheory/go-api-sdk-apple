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
	fmt.Println("=== Apple Business Manager - Get Device Management Service by ID ===")

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

	serverID := "1F97349736CF4614A94F624E705841AD"

	opts := &devicemanagement.RequestQueryOptions{
		Fields: []string{
			devicemanagement.FieldServerName,
			devicemanagement.FieldServerType,
			devicemanagement.FieldEnableMdmDisownFlag,
			devicemanagement.FieldDefaultProductFamilies,
			devicemanagement.FieldStatus,
			devicemanagement.FieldDeviceCount,
			devicemanagement.FieldLastConnectedDateTime,
			devicemanagement.FieldLastConnectedIp,
			devicemanagement.FieldCreatedDateTime,
			devicemanagement.FieldUpdatedDateTime,
		},
	}

	response, _, err := c.AXMAPI.DeviceManagement.GetByMDMServerIDV1(ctx, serverID, opts)
	if err != nil {
		log.Fatalf("Error getting MDM server: %v", err)
	}

	server := response.Data
	fmt.Printf("MDM Server:\n")
	fmt.Printf("  ID: %s\n", server.ID)
	fmt.Printf("  Type: %s\n", server.Type)
	if server.Attributes != nil {
		fmt.Printf("  Name: %s\n", server.Attributes.ServerName)
		fmt.Printf("  Server Type: %s\n", server.Attributes.ServerType)
		fmt.Printf("  Status: %s\n", server.Attributes.Status)
		fmt.Printf("  Device Count: %d\n", server.Attributes.DeviceCount)
		fmt.Printf("  Enable MDM Disown Flag: %v\n", server.Attributes.EnableMdmDisownFlag)
		fmt.Printf("  Default Product Families: %v\n", server.Attributes.DefaultProductFamilies)
		if server.Attributes.LastConnectedDateTime != nil {
			fmt.Printf("  Last Connected: %s\n", server.Attributes.LastConnectedDateTime.Format(time.RFC3339))
		}
		fmt.Printf("  Last Connected IP: %s\n", server.Attributes.LastConnectedIp)
		if server.Attributes.CreatedDateTime != nil {
			fmt.Printf("  Created: %s\n", server.Attributes.CreatedDateTime.Format(time.RFC3339))
		}
		if server.Attributes.UpdatedDateTime != nil {
			fmt.Printf("  Updated: %s\n", server.Attributes.UpdatedDateTime.Format(time.RFC3339))
		}
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
