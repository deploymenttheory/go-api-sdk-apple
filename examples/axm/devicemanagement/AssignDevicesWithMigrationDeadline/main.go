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
	fmt.Println("=== Apple Business Manager - Schedule a Device Management Service Migration ===")

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
		"C39J8W1H3VG5",
		"F2KZ3N4G5HD3",
	}

	// The API rejects deadlines more than 90 days in the future.
	deadline := time.Now().UTC().Add(30 * 24 * time.Hour)

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
				ActivityType: devicemanagement.ActivityTypeAssignDevicesWithMDMMigrationDeadline,
				ActivityTypeMetadata: &devicemanagement.ActivityTypeMetadata{
					MDMMigrationDeadlineDateTime: &deadline,
				},
			},
			Relationships: devicemanagement.OrgDeviceActivityCreateRelationships{
				MDMServer: &devicemanagement.OrgDeviceActivityMDMServerRelationship{
					Data: devicemanagement.OrgDeviceActivityMDMServerLinkage{
						Type: devicemanagement.ResourceTypeMDMServers,
						ID:   mdmServerID,
					},
				},
				Devices: &devicemanagement.OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error scheduling device management service migration: %v", err)
	}

	fmt.Printf("Migration activity created:\n")
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
