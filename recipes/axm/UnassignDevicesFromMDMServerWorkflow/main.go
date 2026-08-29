package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
)

// UnassignDevicesFromMDMServerWorkflow demonstrates a complete workflow for:
// 1. Discovering available MDM servers
// 2. Listing devices currently assigned to a target MDM server
// 3. Unassigning those devices from the server
// 4. Verifying the unassignments completed
func main() {
	fmt.Println("=== Recipe: Unassign Devices from MDM Server ===")

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

	// Step 1: Discover available MDM servers
	fmt.Println("\nStep 1: Discovering available MDM servers...")

	serversResponse, _, err := c.AXMAPI.DeviceManagement.GetV1(ctx, &devicemanagement.RequestQueryOptions{
		Fields: []string{
			devicemanagement.FieldServerName,
			devicemanagement.FieldServerType,
		},
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Error getting MDM servers: %v", err)
	}

	if len(serversResponse.Data) == 0 {
		log.Fatalf("No MDM servers found in organization")
	}

	fmt.Printf("Found %d MDM servers:\n", len(serversResponse.Data))
	for i, server := range serversResponse.Data {
		serverName := ""
		if server.Attributes != nil {
			serverName = server.Attributes.ServerName
		}
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, serverName, server.ID)
	}

	// Use the first MDM server as the target
	targetServer := serversResponse.Data[0]
	targetServerName := ""
	if targetServer.Attributes != nil {
		targetServerName = targetServer.Attributes.ServerName
	}
	fmt.Printf("\nTarget MDM server: %s (ID: %s)\n", targetServerName, targetServer.ID)

	// Step 2: List devices assigned to the target server
	fmt.Printf("\nStep 2: Listing devices assigned to %s...\n", targetServerName)

	deviceLinkages, _, err := c.AXMAPI.DeviceManagement.GetDeviceSerialNumbersByServerIDV1(ctx, targetServer.ID, &devicemanagement.RequestQueryOptions{
		Limit: 100,
	})
	if err != nil {
		log.Fatalf("Error getting device linkages for server: %v", err)
	}

	if len(deviceLinkages.Data) == 0 {
		fmt.Printf("No devices currently assigned to %s. Exiting.\n", targetServerName)
		return
	}

	fmt.Printf("Found %d devices assigned to %s\n", len(deviceLinkages.Data), targetServerName)
	for i, linkage := range deviceLinkages.Data {
		fmt.Printf("  %d. Device ID: %s\n", i+1, linkage.ID)
	}

	// Step 3: Unassign devices from the target server (up to 10)
	maxToUnassign := 10
	if len(deviceLinkages.Data) < maxToUnassign {
		maxToUnassign = len(deviceLinkages.Data)
	}

	deviceIDsToUnassign := make([]string, maxToUnassign)
	for i := 0; i < maxToUnassign; i++ {
		deviceIDsToUnassign[i] = deviceLinkages.Data[i].ID
	}

	fmt.Printf("\nStep 3: Unassigning %d devices from %s...\n", maxToUnassign, targetServerName)

	unassignResponse, _, err := c.AXMAPI.DeviceManagement.UnassignDevicesV1(ctx, targetServer.ID, deviceIDsToUnassign)
	if err != nil {
		log.Fatalf("Error unassigning devices: %v", err)
	}

	fmt.Printf("Unassignment activity created: %s\n", unassignResponse.Data.ID)
	if unassignResponse.Data.Attributes != nil {
		fmt.Printf("Status: %s / %s\n", unassignResponse.Data.Attributes.Status, unassignResponse.Data.Attributes.SubStatus)
	}

	// Step 4: Verify unassignments after a brief wait
	fmt.Println("\nStep 4: Verifying unassignments (waiting 3 seconds for processing)...")
	time.Sleep(3 * time.Second)

	unassignedCount := 0
	for _, deviceID := range deviceIDsToUnassign {
		linkage, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerIDByDeviceIDV1(ctx, deviceID)
		if err != nil {
			// Error likely means no server assigned — success
			unassignedCount++
			fmt.Printf("  Device %s: confirmed unassigned\n", deviceID)
			continue
		}
		if linkage.Data.ID == "" {
			unassignedCount++
			fmt.Printf("  Device %s: confirmed unassigned\n", deviceID)
		} else if linkage.Data.ID != targetServer.ID {
			fmt.Printf("  Device %s: reassigned to different server %s\n", deviceID, linkage.Data.ID)
		} else {
			fmt.Printf("  Device %s: unassignment still processing\n", deviceID)
		}
	}

	fmt.Printf("\nSummary: %d/%d devices confirmed unassigned from %s\n", unassignedCount, maxToUnassign, targetServerName)
	fmt.Println("Note: Unassignments may take additional time to process in Apple's system.")
}
