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
	fmt.Println("=== Apple Business Manager - Release Devices from an Organization ===")

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

	// ---------------------------------------------------------------------------
	// WARNING: releasing devices is irreversible.
	//
	// Released devices are no longer registered to the organization, their device
	// enrollment assignments are removed, they are unenrolled from the built-in
	// device management service, and they are removed from any Blueprints.
	//
	// Apple Business API 2.4+ only. The Apple School Manager API has no equivalent.
	// ---------------------------------------------------------------------------
	deviceIDs := []string{
		"XABC123X0ABC123X0",
	}

	deviceLinkages := make([]devicemanagement.OrgDeviceActivityDeviceLinkage, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		deviceLinkages = append(deviceLinkages, devicemanagement.OrgDeviceActivityDeviceLinkage{
			Type: devicemanagement.ResourceTypeOrgDevices,
			ID:   deviceID,
		})
	}

	response, _, err := c.AXMAPI.DeviceManagement.CreateOrgDeviceActivityV1(ctx, &devicemanagement.OrgDeviceActivityCreateRequest{
		Data: devicemanagement.OrgDeviceActivityData{
			Type: devicemanagement.ResourceTypeOrgDeviceActivities,
			Attributes: devicemanagement.OrgDeviceActivityCreateAttributes{
				ActivityType: devicemanagement.ActivityTypeReleaseDevices,
			},
			Relationships: devicemanagement.OrgDeviceActivityCreateRelationships{
				Devices: &devicemanagement.OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error releasing devices: %v", err)
	}

	fmt.Printf("Release activity created:\n")
	fmt.Printf("  Activity ID: %s\n", response.Data.ID)
	fmt.Printf("  Type: %s\n", response.Data.Type)

	if response.Data.Attributes != nil {
		fmt.Printf("  Status: %s\n", response.Data.Attributes.Status)
		fmt.Printf("  Sub-Status: %s\n", response.Data.Attributes.SubStatus)
		if response.Data.Attributes.CreatedDateTime != nil {
			fmt.Printf("  Created: %s\n", response.Data.Attributes.CreatedDateTime.Format(time.RFC3339))
		}
	}

	fmt.Println("\nThe activity is asynchronous. Poll it with:")
	fmt.Printf("  c.AXMAPI.DeviceManagement.GetOrgDeviceActivityByIDV1(ctx, %q, nil)\n", response.Data.ID)

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
