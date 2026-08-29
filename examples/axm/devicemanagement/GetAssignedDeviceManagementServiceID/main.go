package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Assigned Device Management Service ID ===")

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

	response, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerIDByDeviceIDV1(ctx, deviceID)
	if err != nil {
		log.Fatalf("Error getting assigned server ID for device %s: %v", deviceID, err)
	}

	fmt.Printf("Device ID: %s\n", deviceID)
	fmt.Printf("Assigned Server Linkage:\n")
	fmt.Printf("  Type: %s\n", response.Data.Type)
	fmt.Printf("  Server ID: %s\n", response.Data.ID)

	if response.Links != nil {
		fmt.Printf("  Self Link: %s\n", response.Links.Self)
		fmt.Printf("  Related Link: %s\n", response.Links.Related)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
