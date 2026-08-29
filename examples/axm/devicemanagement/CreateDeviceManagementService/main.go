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
	fmt.Println("=== Apple Business Manager - Create Device Management Service ===")

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

	req := &devicemanagement.MDMServerCreateRequest{
		Data: devicemanagement.MDMServerCreateRequestData{
			Type: "mdmServers",
			Attributes: devicemanagement.MDMServerCreateRequestAttributes{
				ServerName: "Marketing Team MDM",
				ServerCertificate: devicemanagement.MDMServerCertificate{
					Name: "marketing-mdm.cer",
					Data: "MIIDXTCCAkWgAwIBAgIJALxxxxxxx...",
				},
				EnableMdmDisownFlag: true,
			},
		},
	}

	response, _, err := c.AXMAPI.DeviceManagement.CreateMDMServerV1(ctx, req)
	if err != nil {
		log.Fatalf("Error creating MDM server: %v", err)
	}

	server := response.Data
	fmt.Printf("Created MDM Server:\n")
	fmt.Printf("  ID: %s\n", server.ID)
	fmt.Printf("  Type: %s\n", server.Type)
	if server.Attributes != nil {
		fmt.Printf("  Name: %s\n", server.Attributes.ServerName)
		fmt.Printf("  Server Type: %s\n", server.Attributes.ServerType)
		fmt.Printf("  Status: %s\n", server.Attributes.Status)
		fmt.Printf("  Enable MDM Disown Flag: %v\n", server.Attributes.EnableMdmDisownFlag)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
