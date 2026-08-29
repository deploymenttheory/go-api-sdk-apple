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
	fmt.Println("=== Apple Business Manager - Get AppleCare Information by Device ID ===")

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
			devices.FieldAppleCareStatus,
			devices.FieldAppleCarePaymentType,
			devices.FieldAppleCareDescription,
			devices.FieldAppleCareAgreementNumber,
			devices.FieldAppleCareStartDateTime,
			devices.FieldAppleCareEndDateTime,
			devices.FieldAppleCareIsRenewable,
			devices.FieldAppleCareIsCanceled,
			devices.FieldAppleCareContractCancelDateTime,
		},
		Limit: 100,
	}

	response, _, err := c.AXMAPI.Devices.GetAppleCareByDeviceIDV1(ctx, deviceID, opts)
	if err != nil {
		log.Fatalf("Error getting AppleCare information for device %s: %v", deviceID, err)
	}

	fmt.Printf("Found %d AppleCare coverage plan(s) for device %s\n\n", len(response.Data), deviceID)

	for i, coverage := range response.Data {
		fmt.Printf("Coverage Plan %d:\n", i+1)
		fmt.Printf("  ID: %s\n", coverage.ID)
		fmt.Printf("  Type: %s\n", coverage.Type)

		if coverage.Attributes != nil {
			fmt.Printf("  Description: %s\n", coverage.Attributes.Description)
			fmt.Printf("  Status: %s\n", coverage.Attributes.Status)
			fmt.Printf("  Payment Type: %s\n", coverage.Attributes.PaymentType)
			fmt.Printf("  Is Renewable: %v\n", coverage.Attributes.IsRenewable)
			fmt.Printf("  Is Canceled: %v\n", coverage.Attributes.IsCanceled)

			if coverage.Attributes.AgreementNumber != nil {
				fmt.Printf("  Agreement Number: %s\n", *coverage.Attributes.AgreementNumber)
			}
			if coverage.Attributes.StartDateTime != nil {
				fmt.Printf("  Start Date: %s\n", coverage.Attributes.StartDateTime.Format(time.RFC3339))
			}
			if coverage.Attributes.EndDateTime != nil {
				fmt.Printf("  End Date: %s\n", coverage.Attributes.EndDateTime.Format(time.RFC3339))
			}
			if coverage.Attributes.ContractCancelDateTime != nil {
				fmt.Printf("  Contract Cancel Date: %s\n", coverage.Attributes.ContractCancelDateTime.Format(time.RFC3339))
			}
		}
		fmt.Println()
	}

	if response.Links != nil && response.Links.Next != "" {
		fmt.Printf("Next page: %s\n", response.Links.Next)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %v", err)
	}
	fmt.Println("Full JSON response:")
	fmt.Println(string(jsonData))
}
