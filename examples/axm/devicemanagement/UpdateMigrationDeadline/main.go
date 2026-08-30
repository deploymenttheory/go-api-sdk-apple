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
	fmt.Println("=== Apple Business Manager - Update a Device Management Service Migration Deadline ===")

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

	deviceIDs := []string{
		"C39J8W1H3VG5",
		"F2KZ3N4G5HD3",
	}

	// A deadline earlier than the existing one — or in the past — is honoured by the
	// device immediately, without offering the user the option to delay.
	newDeadline := time.Now().UTC().Add(14 * 24 * time.Hour)

	deviceLinkages := make([]devicemanagement.OrgDeviceActivityDeviceLinkage, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		deviceLinkages = append(deviceLinkages, devicemanagement.OrgDeviceActivityDeviceLinkage{
			Type: devicemanagement.ResourceTypeOrgDevices,
			ID:   deviceID,
		})
	}

	// UPDATE_MDM_MIGRATION_DEADLINE takes no mdmServer relationship.
	response, _, err := c.AXMAPI.DeviceManagement.CreateOrgDeviceActivityV1(ctx, &devicemanagement.OrgDeviceActivityCreateRequest{
		Data: devicemanagement.OrgDeviceActivityData{
			Type: devicemanagement.ResourceTypeOrgDeviceActivities,
			Attributes: devicemanagement.OrgDeviceActivityCreateAttributes{
				ActivityType: devicemanagement.ActivityTypeUpdateMDMMigrationDeadline,
				ActivityTypeMetadata: &devicemanagement.ActivityTypeMetadata{
					MDMMigrationDeadlineDateTime: &newDeadline,
				},
			},
			Relationships: devicemanagement.OrgDeviceActivityCreateRelationships{
				Devices: &devicemanagement.OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error updating device management service migration deadline: %v", err)
	}

	fmt.Printf("Deadline update activity created:\n")
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
