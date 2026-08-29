package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
)

func main() {
	fmt.Println("=== Apple Business Manager - Get Device Information by Device ID ===")

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

	deviceID := "XABC123X0ABC123X0"

	opts := &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldDeviceModel,
			devices.FieldProductFamily,
			devices.FieldProductType,
			devices.FieldDeviceCapacity,
			devices.FieldPartNumber,
			devices.FieldOrderNumber,
			devices.FieldColor,
			devices.FieldStatus,
			devices.FieldAddedToOrgDateTime,
			devices.FieldUpdatedDateTime,
			devices.FieldOrderDateTime,
			devices.FieldIMEI,
			devices.FieldMEID,
			devices.FieldEID,
			devices.FieldWiFiMACAddress,
			devices.FieldBluetoothMACAddress,
			devices.FieldPurchaseSourceId,
			devices.FieldPurchaseSourceType,
		},
	}

	response, _, err := c.AXMAPI.Devices.GetByDeviceIDV1(ctx, deviceID, opts)
	if err != nil {
		log.Fatalf("Error getting device information for %s: %v", deviceID, err)
	}

	device := response.Data
	fmt.Printf("Device Information:\n")
	fmt.Printf("  ID: %s\n", device.ID)
	fmt.Printf("  Type: %s\n", device.Type)

	if device.Attributes != nil {
		fmt.Printf("  Serial Number: %s\n", device.Attributes.SerialNumber)
		fmt.Printf("  Model: %s\n", device.Attributes.DeviceModel)
		fmt.Printf("  Product Family: %s\n", device.Attributes.ProductFamily)
		fmt.Printf("  Product Type: %s\n", device.Attributes.ProductType)
		fmt.Printf("  Capacity: %s\n", device.Attributes.DeviceCapacity)
		fmt.Printf("  Part Number: %s\n", device.Attributes.PartNumber)
		fmt.Printf("  Order Number: %s\n", device.Attributes.OrderNumber)
		fmt.Printf("  Color: %s\n", device.Attributes.Color)
		fmt.Printf("  Status: %s\n", device.Attributes.Status)

		if device.Attributes.AddedToOrgDateTime != nil {
			fmt.Printf("  Added to Org: %s\n", device.Attributes.AddedToOrgDateTime.Format(time.RFC3339))
		}
		if device.Attributes.UpdatedDateTime != nil {
			fmt.Printf("  Updated: %s\n", device.Attributes.UpdatedDateTime.Format(time.RFC3339))
		}
		if device.Attributes.OrderDateTime != nil {
			fmt.Printf("  Ordered: %s\n", device.Attributes.OrderDateTime.Format(time.RFC3339))
		}

		if len(device.Attributes.IMEI) > 0 {
			fmt.Printf("  IMEI: %v\n", device.Attributes.IMEI)
		}
		if len(device.Attributes.MEID) > 0 {
			fmt.Printf("  MEID: %v\n", device.Attributes.MEID)
		}
		if device.Attributes.EID != "" {
			fmt.Printf("  EID: %s\n", device.Attributes.EID)
		}
		if device.Attributes.WiFiMACAddress != "" {
			fmt.Printf("  WiFi MAC: %s\n", device.Attributes.WiFiMACAddress)
		}
		if device.Attributes.BluetoothMACAddress != "" {
			fmt.Printf("  Bluetooth MAC: %s\n", device.Attributes.BluetoothMACAddress)
		}
		if device.Attributes.PurchaseSourceId != "" {
			fmt.Printf("  Purchase Source ID: %s\n", device.Attributes.PurchaseSourceId)
		}
		if device.Attributes.PurchaseSourceType != "" {
			fmt.Printf("  Purchase Source Type: %s\n", device.Attributes.PurchaseSourceType)
		}
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("\nFull JSON response:")
	fmt.Println(string(jsonData))
}
