package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
)

func main() {
	fmt.Println("=== Apple Business Manager - Unassign Devices from Server ===")

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
	deviceIDs := []string{
		"XABC123X0ABC123X0",
		"YDEF456Y1DEF456Y1",
	}

	response, _, err := c.AXMAPI.DeviceManagement.UnassignDevicesV1(ctx, mdmServerID, deviceIDs)
	if err != nil {
		log.Fatalf("Error unassigning devices from server: %v", err)
	}

	fmt.Printf("Unassignment activity created:\n")
	fmt.Printf("  Activity ID: %s\n", response.Data.ID)
	fmt.Printf("  Type: %s\n", response.Data.Type)

	if response.Data.Attributes != nil {
		fmt.Printf("  Activity Type: %s\n", response.Data.Attributes.ActivityType)
		fmt.Printf("  Status: %s\n", response.Data.Attributes.Status)
		fmt.Printf("  Sub-Status: %s\n", response.Data.Attributes.SubStatus)
		if response.Data.Attributes.CreatedDateTime != nil {
			fmt.Printf("  Created: %s\n", response.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
		}
	}

	if response.Data.Links != nil {
		fmt.Printf("  Self Link: %s\n", response.Data.Links.Self)
	}

	if response.Links != nil {
		fmt.Printf("  Response Self Link: %s\n", response.Links.Self)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
