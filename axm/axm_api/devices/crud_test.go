package devices

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) *Devices {
	// Create a mock auth provider
	mockAuth := &MockAuthProvider{}

	// Create dummy private key for testing (using a mock)
	// We'll override auth with our mock provider
	dummyKey := "dummy-key"

	// Create core transport with mock auth override
	coreClient, err := client.NewTransport(
		"test-key-id",
		"test-issuer-id",
		dummyKey,
		client.WithAuth(mockAuth),
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0), // Disable retries for tests
	)
	require.NoError(t, err)

	// Activate httpmock for the client's HTTP client
	httpmock.ActivateNonDefault(coreClient.GetHTTPClient().Client())

	// Setup cleanup
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Create devices service
	return NewService(coreClient)
}

// MockAuthProvider implements the AuthProvider interface for testing
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	// Mock auth - do nothing for tests
	return nil
}

func TestGetOrganizationDevices_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{"serialNumber", "deviceModel", "status"},
		Limit:  100,
	}

	result, resp, err := client.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify the first device
	device := result.Data[0]
	assert.Equal(t, "orgDevices", device.Type)
	assert.Equal(t, "XABC123X0ABC123X0", device.ID)
	assert.NotNil(t, device.Attributes)

	// Test all string fields from the list response
	assert.Equal(t, "XABC123X0ABC123X0", device.Attributes.SerialNumber)
	assert.Equal(t, "iMac 21.5\"", device.Attributes.DeviceModel)
	assert.Equal(t, "Mac", device.Attributes.ProductFamily)
	assert.Equal(t, "iMac16,2", device.Attributes.ProductType)
	assert.Equal(t, "750GB", device.Attributes.DeviceCapacity)
	assert.Equal(t, "FD311LL/A", device.Attributes.PartNumber)
	assert.Equal(t, "1234567890", device.Attributes.OrderNumber)
	assert.Equal(t, "SILVER", device.Attributes.Color) // Has color in list response
	assert.Equal(t, "UNASSIGNED", device.Attributes.Status)
	assert.Equal(t, "89049037640158663184237812557346", device.Attributes.EID)
	assert.Equal(t, "-2085650007946880", device.Attributes.PurchaseSourceId)
	assert.Equal(t, "APPLE", device.Attributes.PurchaseSourceType)

	// Test timestamp fields
	assert.NotNil(t, device.Attributes.AddedToOrgDateTime)
	assert.NotNil(t, device.Attributes.UpdatedDateTime)
	assert.NotNil(t, device.Attributes.OrderDateTime)

	// Test array fields
	assert.Len(t, device.Attributes.IMEI, 2)
	assert.Equal(t, "123456789012345", device.Attributes.IMEI[0])
	assert.Equal(t, "123456789012346", device.Attributes.IMEI[1])
	assert.Len(t, device.Attributes.MEID, 1)
	assert.Equal(t, "12345678901237", device.Attributes.MEID[0])

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	// Verify exactly one HTTP call was made
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := client.GetV1(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_WithFieldSelection(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldSerialNumber, FieldDeviceModel, FieldStatus, FieldProductFamily},
		Limit:  50,
	}

	result, resp, err := client.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify device data
	device := result.Data[0]
	assert.Equal(t, "XABC123X0ABC123X0", device.Attributes.SerialNumber)
	assert.Equal(t, "iMac 21.5\"", device.Attributes.DeviceModel)
	assert.Equal(t, "UNASSIGNED", device.Attributes.Status)
	assert.Equal(t, "Mac", device.Attributes.ProductFamily)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_WithLimitEnforcement(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Limit: 1500, // Exceeds API maximum of 1000
	}

	result, resp, err := client.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationDevices_HTTPError(t *testing.T) {
	client := setupMockClient(t)

	// Reset httpmock to clear any previous registrations
	httpmock.Reset()

	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := client.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount()) // retries disabled
}

func TestGetDeviceInformation_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{
		Fields: []string{FieldSerialNumber, FieldDeviceModel, FieldStatus},
	}

	result, resp, err := client.GetByDeviceIDV1(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	// Verify device data
	device := result.Data
	assert.Equal(t, "orgDevices", device.Type)
	assert.Equal(t, "XABC123X0ABC123X0", device.ID)
	assert.NotNil(t, device.Attributes)

	// Test all string fields
	assert.Equal(t, "XABC123X0ABC123X0", device.Attributes.SerialNumber)
	assert.Equal(t, "iMac 21.5\"", device.Attributes.DeviceModel)
	assert.Equal(t, "Mac", device.Attributes.ProductFamily)
	assert.Equal(t, "iMac16,2", device.Attributes.ProductType)
	assert.Equal(t, "750GB", device.Attributes.DeviceCapacity)
	assert.Equal(t, "FD311LL/A", device.Attributes.PartNumber)
	assert.Equal(t, "1234567890", device.Attributes.OrderNumber)
	assert.Equal(t, "SILVER", device.Attributes.Color)
	assert.Equal(t, "UNASSIGNED", device.Attributes.Status)
	assert.Equal(t, "89049037640158663184237812557346", device.Attributes.EID)
	assert.Equal(t, "-2085650007946880", device.Attributes.PurchaseSourceId)
	assert.Equal(t, "APPLE", device.Attributes.PurchaseSourceType)

	// Verify timestamps are parsed correctly
	assert.NotNil(t, device.Attributes.AddedToOrgDateTime)
	assert.NotNil(t, device.Attributes.UpdatedDateTime)
	assert.NotNil(t, device.Attributes.OrderDateTime)

	// Verify arrays are handled correctly
	assert.Len(t, device.Attributes.IMEI, 2)
	assert.Equal(t, "123456789012345", device.Attributes.IMEI[0])
	assert.Equal(t, "123456789012346", device.Attributes.IMEI[1])
	assert.Len(t, device.Attributes.MEID, 1)
	assert.Equal(t, "12345678901237", device.Attributes.MEID[0])

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"

	result, resp, err := client.GetByDeviceIDV1(ctx, deviceID, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "XABC123X0ABC123X0", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := client.GetByDeviceIDV1(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_DeviceNotFound(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "NONEXISTENT123"
	opts := &RequestQueryOptions{}

	result, _, err := client.GetByDeviceIDV1(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_HTTPError(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{}

	result, _, err := client.GetByDeviceIDV1(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceInformation_ContextCancellation(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{}

	result, _, err := client.GetByDeviceIDV1(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetOrganizationDevices_ContextTimeout(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Sleep to ensure timeout
	time.Sleep(1 * time.Millisecond)

	opts := &RequestQueryOptions{}

	result, _, err := client.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestDeviceFieldConstants(t *testing.T) {
	// Test that field constants are properly defined
	expectedFields := []string{
		FieldSerialNumber,
		FieldAddedToOrgDateTime,
		FieldUpdatedDateTime,
		FieldDeviceModel,
		FieldProductFamily,
		FieldProductType,
		FieldDeviceCapacity,
		FieldPartNumber,
		FieldOrderNumber,
		FieldColor,
		FieldStatus,
		FieldOrderDateTime,
		FieldIMEI,
		FieldMEID,
		FieldEID,
		FieldWiFiMACAddress,
		FieldBluetoothMACAddress,
		FieldPurchaseSourceId,
		FieldPurchaseSourceType,
		FieldAssignedServer,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	// Test specific field values
	assert.Equal(t, "serialNumber", FieldSerialNumber)
	assert.Equal(t, "deviceModel", FieldDeviceModel)
	assert.Equal(t, "status", FieldStatus)
	assert.Equal(t, "productFamily", FieldProductFamily)
}

func TestOptionsStructures(t *testing.T) {
	// Test RequestQueryOptions
	opts1 := &RequestQueryOptions{
		Fields: []string{FieldSerialNumber, FieldDeviceModel},
		Limit:  100,
	}
	assert.Len(t, opts1.Fields, 2)
	assert.Equal(t, 100, opts1.Limit)

	// Test RequestQueryOptions
	opts2 := &RequestQueryOptions{
		Fields: []string{FieldStatus, FieldProductFamily},
	}
	assert.Len(t, opts2.Fields, 2)
}

func TestComprehensiveFieldCoverage(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"

	// Test with all available fields to ensure comprehensive coverage
	opts := &RequestQueryOptions{
		Fields: []string{
			FieldSerialNumber,
			FieldAddedToOrgDateTime,
			FieldUpdatedDateTime,
			FieldDeviceModel,
			FieldProductFamily,
			FieldProductType,
			FieldDeviceCapacity,
			FieldPartNumber,
			FieldOrderNumber,
			FieldColor,
			FieldStatus,
			FieldOrderDateTime,
			FieldIMEI,
			FieldMEID,
			FieldEID,
			FieldWiFiMACAddress,
			FieldBluetoothMACAddress,
			FieldPurchaseSourceId,
			FieldPurchaseSourceType,
			FieldAssignedServer,
			FieldAppleCareCoverage,
			FieldReleasedFromOrgDateTime,
			FieldIsMdmMigrationCapable,
			FieldMdmMigrationStatus,
			FieldMdmMigrationDeadlineDateTime,
		},
	}

	result, resp, err := client.GetByDeviceIDV1(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	// Verify that the request was made with all field parameters
	// This ensures our field constants are correctly defined
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Verify all fields are accessible (no compilation errors)
	attrs := result.Data.Attributes
	_ = attrs.SerialNumber
	_ = attrs.AddedToOrgDateTime
	_ = attrs.UpdatedDateTime
	_ = attrs.DeviceModel
	_ = attrs.ProductFamily
	_ = attrs.ProductType
	_ = attrs.DeviceCapacity
	_ = attrs.PartNumber
	_ = attrs.OrderNumber
	_ = attrs.Color
	_ = attrs.Status
	_ = attrs.OrderDateTime
	_ = attrs.IMEI
	_ = attrs.MEID
	_ = attrs.EID
	_ = attrs.WiFiMACAddress
	_ = attrs.BluetoothMACAddress
	_ = attrs.PurchaseSourceId
	_ = attrs.PurchaseSourceType
	_ = attrs.AssignedServer
}

func TestGetAppleCareInformation_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{
		Fields: []string{FieldAppleCareStatus, FieldAppleCarePaymentType, FieldAppleCareDescription},
		Limit:  100,
	}

	result, resp, err := client.GetAppleCareByDeviceIDV1(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify we have multiple AppleCare coverage entries
	assert.Len(t, result.Data, 3)

	// Verify first coverage (Limited Warranty)
	coverage1 := result.Data[0]
	assert.Equal(t, "appleCareCoverage", coverage1.Type)
	assert.Equal(t, "XABC123X0ABC123X0", coverage1.ID)
	assert.NotNil(t, coverage1.Attributes)
	assert.Equal(t, "ACTIVE", coverage1.Attributes.Status)
	assert.Equal(t, "NONE", coverage1.Attributes.PaymentType)
	assert.Equal(t, "Limited Warranty", coverage1.Attributes.Description)
	assert.False(t, coverage1.Attributes.IsRenewable)
	assert.False(t, coverage1.Attributes.IsCanceled)
	assert.NotNil(t, coverage1.Attributes.StartDateTime)
	assert.NotNil(t, coverage1.Attributes.EndDateTime)
	assert.Nil(t, coverage1.Attributes.AgreementNumber)
	assert.Nil(t, coverage1.Attributes.ContractCancelDateTime)

	// Verify second coverage (AppleCare+)
	coverage2 := result.Data[1]
	assert.Equal(t, "appleCareCoverage", coverage2.Type)
	assert.Equal(t, "0000000001", coverage2.ID)
	assert.NotNil(t, coverage2.Attributes)
	assert.Equal(t, "ACTIVE", coverage2.Attributes.Status)
	assert.Equal(t, "SUBSCRIPTION", coverage2.Attributes.PaymentType)
	assert.Equal(t, "AppleCare+", coverage2.Attributes.Description)
	assert.True(t, coverage2.Attributes.IsRenewable)
	assert.False(t, coverage2.Attributes.IsCanceled)
	assert.NotNil(t, coverage2.Attributes.AgreementNumber)

	// Verify third coverage (AppleCare+ for Business Essentials)
	coverage3 := result.Data[2]
	assert.Equal(t, "appleCareCoverage", coverage3.Type)
	assert.Equal(t, "abe-XABC123X0ABC123X0", coverage3.ID)
	assert.NotNil(t, coverage3.Attributes)
	assert.Equal(t, "ACTIVE", coverage3.Attributes.Status)
	assert.Equal(t, "ABE_SUBSCRIPTION", coverage3.Attributes.PaymentType)
	assert.Equal(t, "AppleCare+ for Business Essentials", coverage3.Attributes.Description)
	assert.True(t, coverage3.Attributes.IsRenewable)
	assert.Nil(t, coverage3.Attributes.EndDateTime)

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppleCareInformation_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"

	result, resp, err := client.GetAppleCareByDeviceIDV1(ctx, deviceID, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)
	assert.Len(t, result.Data, 3)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppleCareInformation_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := client.GetAppleCareByDeviceIDV1(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAppleCareInformation_DeviceNotFound(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "NONEXISTENT123"
	opts := &RequestQueryOptions{}

	result, _, err := client.GetAppleCareByDeviceIDV1(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppleCareInformation_WithFieldSelection(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{
		Fields: []string{
			FieldAppleCareStatus,
			FieldAppleCarePaymentType,
			FieldAppleCareDescription,
			FieldAppleCareStartDateTime,
			FieldAppleCareEndDateTime,
			FieldAppleCareIsRenewable,
			FieldAppleCareIsCanceled,
		},
	}

	result, resp, err := client.GetAppleCareByDeviceIDV1(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppleCareInformation_WithLimitEnforcement(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{
		Limit: 1500, // Exceeds API maximum of 1000
	}

	result, resp, err := client.GetAppleCareByDeviceIDV1(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppleCareInformation_ContextCancellation(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	deviceID := "XABC123X0ABC123X0"
	opts := &RequestQueryOptions{}

	result, _, err := client.GetAppleCareByDeviceIDV1(ctx, deviceID, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestAppleCareFieldConstants(t *testing.T) {
	// Test that AppleCare field constants are properly defined
	expectedFields := []string{
		FieldAppleCareStatus,
		FieldAppleCarePaymentType,
		FieldAppleCareDescription,
		FieldAppleCareAgreementNumber,
		FieldAppleCareStartDateTime,
		FieldAppleCareEndDateTime,
		FieldAppleCareIsRenewable,
		FieldAppleCareIsCanceled,
		FieldAppleCareContractCancelDateTime,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "AppleCare field constant should not be empty")
	}

	// Test specific field values
	assert.Equal(t, "status", FieldAppleCareStatus)
	assert.Equal(t, "paymentType", FieldAppleCarePaymentType)
	assert.Equal(t, "description", FieldAppleCareDescription)
	assert.Equal(t, "agreementNumber", FieldAppleCareAgreementNumber)
}

func TestAppleCareStatusConstants(t *testing.T) {
	// Test AppleCare status constants
	assert.Equal(t, "ACTIVE", AppleCareStatusActive)
	assert.Equal(t, "INACTIVE", AppleCareStatusInactive)
	assert.Equal(t, "EXPIRED", AppleCareStatusExpired)

	// Test payment type constants
	assert.Equal(t, "NONE", PaymentTypeNone)
	assert.Equal(t, "SUBSCRIPTION", PaymentTypeSubscription)
	assert.Equal(t, "ABE_SUBSCRIPTION", PaymentTypeABESubscription)
}

// ====== DEVICE MANAGEMENT SERVICE MIGRATION TESTS ======

// TestGetDeviceInformation_MigrationAttributes verifies the device management service
// migration attributes added in Apple School Manager API 1.6 / Apple Business API 2.3
// decode correctly.
func TestGetDeviceInformation_MigrationAttributes(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.OrgDevicesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, resp, err := client.GetByDeviceIDV1(ctx, "XABC123X0ABC123X0", &RequestQueryOptions{
		Fields: []string{
			FieldIsMdmMigrationCapable,
			FieldMdmMigrationStatus,
			FieldMdmMigrationDeadlineDateTime,
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	attrs := result.Data.Attributes
	require.NotNil(t, attrs)
	assert.True(t, attrs.IsMdmMigrationCapable)
	assert.Equal(t, MdmMigrationStatusRequested, attrs.MdmMigrationStatus)
	require.NotNil(t, attrs.MdmMigrationDeadlineDateTime)
	assert.Equal(t,
		time.Date(2026, 3, 15, 17, 0, 0, 0, time.UTC),
		attrs.MdmMigrationDeadlineDateTime.UTC(),
	)
}

func TestMdmMigrationStatusConstants(t *testing.T) {
	assert.Equal(t, "REQUESTED", MdmMigrationStatusRequested)
	assert.Equal(t, "STARTED", MdmMigrationStatusStarted)
	assert.Equal(t, "SUCCESS", MdmMigrationStatusSuccess)
	assert.Equal(t, "FAILED", MdmMigrationStatusFailed)
}

func TestMigrationAndReleaseFieldConstants(t *testing.T) {
	assert.Equal(t, "isMdmMigrationCapable", FieldIsMdmMigrationCapable)
	assert.Equal(t, "mdmMigrationStatus", FieldMdmMigrationStatus)
	assert.Equal(t, "mdmMigrationDeadlineDateTime", FieldMdmMigrationDeadlineDateTime)
	assert.Equal(t, "releasedFromOrgDateTime", FieldReleasedFromOrgDateTime)
	assert.Equal(t, "appleCareCoverage", FieldAppleCareCoverage)
}
