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
	fmt.Println("=== Apple Business Manager - Get Organization Device Activity ===")

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

	// The activity ID returned by CreateOrgDeviceActivityV1, AssignDevicesV1 or
	// UnassignDevicesV1. Apple retains activities for 30 days.
	activityID := "b1481656-b267-480d-b284-a809eed8b041"

	// Poll until the activity reaches a terminal status.
	deadline := time.Now().Add(5 * time.Minute)
	for {
		response, _, err := c.AXMAPI.DeviceManagement.GetOrgDeviceActivityByIDV1(ctx, activityID, &devicemanagement.GetOrgDeviceActivityQueryOptions{
			Fields: []string{
				devicemanagement.FieldActivityStatus,
				devicemanagement.FieldActivitySubStatus,
				devicemanagement.FieldActivityCreatedDateTime,
				devicemanagement.FieldActivityCompletedDateTime,
				devicemanagement.FieldActivityDownloadURL,
			},
		})
		if err != nil {
			log.Fatalf("Error getting org device activity: %v", err)
		}

		attrs := response.Data.Attributes
		if attrs == nil {
			log.Fatalf("Activity %s returned no attributes", activityID)
		}

		fmt.Printf("Activity %s: status=%s subStatus=%s\n", response.Data.ID, attrs.Status, attrs.SubStatus)

		terminal := attrs.Status == devicemanagement.ActivityStatusCompleted ||
			attrs.Status == devicemanagement.ActivityStatusFailed ||
			attrs.Status == devicemanagement.ActivityStatusStopped

		if terminal {
			if attrs.CompletedDateTime != nil {
				fmt.Printf("  Completed: %s\n", attrs.CompletedDateTime.Format(time.RFC3339))
			}
			// A COMPLETED activity exposes a presigned URL to the activity log in CSV
			// format, listing the per-serial-number outcome.
			if attrs.DownloadURL != "" {
				fmt.Printf("  Activity log (CSV): %s\n", attrs.DownloadURL)
			}

			jsonData, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling response: %v", err)
			}
			fmt.Println("\nFull JSON response:")
			fmt.Println(string(jsonData))
			return
		}

		if time.Now().After(deadline) {
			log.Fatalf("Activity %s did not reach a terminal status within the poll window", activityID)
		}

		time.Sleep(5 * time.Second)
	}
}
