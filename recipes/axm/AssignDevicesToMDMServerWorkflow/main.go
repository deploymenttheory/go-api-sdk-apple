package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
)

// AssignDevicesToMDMServerWorkflow demonstrates a complete workflow for:
// 1. Discovering available MDM servers
// 2. Listing organization devices
// 3. Identifying unassigned devices
// 4. Assigning those devices to a target MDM server
// 5. Verifying the assignments completed
func main() {
	fmt.Println("=== Recipe: Assign Unassigned Devices to MDM Server ===")

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

	// Step 2: List organization devices
	fmt.Println("\nStep 2: Listing organization devices...")

	devicesResponse, _, err := c.AXMAPI.Devices.GetV1(ctx, &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldStatus,
		},
		Limit: 50,
	})
	if err != nil {
		log.Fatalf("Error getting organization devices: %v", err)
	}

	fmt.Printf("Found %d devices in organization\n", len(devicesResponse.Data))

	// Step 3: Identify unassigned devices
	fmt.Println("\nStep 3: Identifying unassigned devices...")

	var unassignedDevices []devices.OrgDevice

	for _, device := range devicesResponse.Data {
		linkage, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerIDByDeviceIDV1(ctx, device.ID)
		if err != nil {
			// Error likely means no server assigned
			unassignedDevices = append(unassignedDevices, device)
			continue
		}
		if linkage.Data.ID == "" {
			unassignedDevices = append(unassignedDevices, device)
		}
	}

	fmt.Printf("Found %d unassigned devices\n", len(unassignedDevices))

	if len(unassignedDevices) == 0 {
		fmt.Println("No unassigned devices to process. Exiting.")
		return
	}

	// Step 4: Assign unassigned devices to the target server (up to 10)
	maxToAssign := 10
	if len(unassignedDevices) < maxToAssign {
		maxToAssign = len(unassignedDevices)
	}

	deviceIDsToAssign := make([]string, maxToAssign)
	for i := 0; i < maxToAssign; i++ {
		deviceIDsToAssign[i] = unassignedDevices[i].ID
	}

	fmt.Printf("\nStep 4: Assigning %d devices to %s...\n", maxToAssign, targetServerName)

	for _, device := range unassignedDevices[:maxToAssign] {
		serialNumber := ""
		if device.Attributes != nil {
			serialNumber = device.Attributes.SerialNumber
		}
		fmt.Printf("  - %s (Serial: %s)\n", device.ID, serialNumber)
	}

	assignResponse, _, err := c.AXMAPI.DeviceManagement.AssignDevicesV1(ctx, targetServer.ID, deviceIDsToAssign)
	if err != nil {
		log.Fatalf("Error assigning devices: %v", err)
	}

	fmt.Printf("Assignment activity created: %s\n", assignResponse.Data.ID)
	if assignResponse.Data.Attributes != nil {
		fmt.Printf("Status: %s / %s\n", assignResponse.Data.Attributes.Status, assignResponse.Data.Attributes.SubStatus)
	}

	// Step 5: Verify assignments after a brief wait
	fmt.Println("\nStep 5: Verifying assignments (waiting 3 seconds for processing)...")
	time.Sleep(3 * time.Second)

	assignedCount := 0
	for _, device := range unassignedDevices[:maxToAssign] {
		linkage, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerIDByDeviceIDV1(ctx, device.ID)
		if err != nil {
			fmt.Printf("  Device %s: could not verify\n", device.ID)
			continue
		}
		if linkage.Data.ID == targetServer.ID {
			assignedCount++
			fmt.Printf("  Device %s: confirmed assigned\n", device.ID)
		} else {
			fmt.Printf("  Device %s: assignment still processing (server: %s)\n", device.ID, linkage.Data.ID)
		}
	}

	fmt.Printf("\nSummary: %d/%d devices confirmed assigned to %s\n", assignedCount, maxToAssign, targetServerName)
	fmt.Println("Note: Assignments may take additional time to process in Apple's system.")
}
