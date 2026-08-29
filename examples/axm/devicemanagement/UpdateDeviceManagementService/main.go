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
	fmt.Println("=== Apple Business Manager - Update Device Management Service ===")

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
	enableDisown := true

	req := &devicemanagement.MDMServerUpdateRequest{
		Data: devicemanagement.MDMServerUpdateRequestData{
			Type: "mdmServers",
			ID:   serverID,
			Attributes: devicemanagement.MDMServerUpdateRequestAttributes{
				ServerName:             "Production MDM Updated",
				EnableMdmDisownFlag:    &enableDisown,
				DefaultProductFamilies: []string{"IPAD", "IPHONE", "MAC"},
			},
		},
	}

	response, _, err := c.AXMAPI.DeviceManagement.UpdateMDMServerByIDV1(ctx, serverID, req)
	if err != nil {
		log.Fatalf("Error updating MDM server: %v", err)
	}

	server := response.Data
	fmt.Printf("Updated MDM Server:\n")
	fmt.Printf("  ID: %s\n", server.ID)
	if server.Attributes != nil {
		fmt.Printf("  Name: %s\n", server.Attributes.ServerName)
		fmt.Printf("  Status: %s\n", server.Attributes.Status)
		fmt.Printf("  Enable MDM Disown Flag: %v\n", server.Attributes.EnableMdmDisownFlag)
		fmt.Printf("  Default Product Families: %v\n", server.Attributes.DefaultProductFamilies)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
