package main

import (
	"context"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
)

func main() {
	fmt.Println("=== Apple Business Manager - Delete Device Management Service ===")

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

	fmt.Printf("Deleting MDM server with ID: %s\n", serverID)
	fmt.Println("Note: A server with devices assigned cannot be deleted.")
	fmt.Println("Unassign all devices before deleting.")

	resp, err := c.AXMAPI.DeviceManagement.DeleteMDMServerByIDV1(ctx, serverID)
	if err != nil {
		log.Fatalf("Error deleting MDM server: %v", err)
	}

	fmt.Printf("MDM server deleted successfully (HTTP %d)\n", resp.StatusCode())
}
