package devicemanagement

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) *DeviceManagement {
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

	// Create device management service
	return NewService(coreClient)
}

// MockAuthProvider implements the AuthProvider interface for testing
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	// Mock auth - do nothing for tests
	return nil
}

func TestGetDeviceManagementServices_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldServerName, FieldServerType, FieldCreatedDateTime},
		Limit:  100,
	}

	result, resp, err := client.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify the first MDM server
	server := result.Data[0]
	assert.Equal(t, "mdmServers", server.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", server.ID)
	assert.NotNil(t, server.Attributes)

	// Test all server attributes
	assert.Equal(t, "Test Device Management Service", server.Attributes.ServerName)
	assert.Equal(t, "MDM", server.Attributes.ServerType)
	assert.NotNil(t, server.Attributes.CreatedDateTime)
	assert.NotNil(t, server.Attributes.UpdatedDateTime)

	// Verify relationships
	assert.NotNil(t, server.Relationships)
	assert.NotNil(t, server.Relationships.Devices)
	assert.NotNil(t, server.Relationships.Devices.Links)
	assert.Contains(t, server.Relationships.Devices.Links.Self, "relationships/devices")

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServices_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
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

func TestGetDeviceManagementServices_WithLimitEnforcement(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
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

func TestGetDeviceManagementServices_HTTPError(t *testing.T) {
	client := setupMockClient(t)

	// Reset httpmock to clear any previous registrations
	httpmock.Reset()

	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := client.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount()) // retries disabled
}

func TestGetDeviceSerialNumbersForDeviceManagementService_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"
	opts := &RequestQueryOptions{
		Limit: 100,
	}

	result, resp, err := client.GetDeviceSerialNumbersByServerIDV1(ctx, serverID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	// Verify device linkage
	linkage := result.Data[0]
	assert.Equal(t, "orgDevices", linkage.Type)
	assert.Equal(t, "XABC123X0ABC123X0", linkage.ID)

	// Verify pagination metadata
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceSerialNumbersForDeviceManagementService_EmptyServerID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := client.GetDeviceSerialNumbersByServerIDV1(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetDeviceSerialNumbersForDeviceManagementService_WithNilOptions(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"

	result, resp, err := client.GetDeviceSerialNumbersByServerIDV1(ctx, serverID, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceIDForADevice_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "DVVS36G1YD3JKQNI"

	result, resp, err := client.GetAssignedServerIDByDeviceIDV1(ctx, deviceID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	// Verify linkage data
	assert.Equal(t, "mdmServers", result.Data.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", result.Data.ID)

	// Verify links
	assert.NotNil(t, result.Links)
	assert.Contains(t, result.Links.Self, "relationships/assignedServer")
	assert.Contains(t, result.Links.Related, "assignedServer")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceIDForADevice_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := client.GetAssignedServerIDByDeviceIDV1(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceInformationByDeviceID_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceID := "DVVS36G1YD3JKQNI"
	opts := &RequestQueryOptions{
		Fields: []string{FieldServerName, FieldServerType},
	}

	result, resp, err := client.GetAssignedServerInfoByDeviceIDV1(ctx, deviceID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	// Verify server data
	server := result.Data
	assert.Equal(t, "mdmServers", server.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", server.ID)
	assert.NotNil(t, server.Attributes)

	// Test all server attributes
	assert.Equal(t, "Test Device Management Service", server.Attributes.ServerName)
	assert.Equal(t, "APPLE_CONFIGURATOR", server.Attributes.ServerType)
	assert.NotNil(t, server.Attributes.CreatedDateTime)
	assert.NotNil(t, server.Attributes.UpdatedDateTime)

	// Verify relationships
	assert.NotNil(t, server.Relationships)
	assert.NotNil(t, server.Relationships.Devices)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAssignedDeviceManagementServiceInformationByDeviceID_EmptyDeviceID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := client.GetAssignedServerInfoByDeviceIDV1(ctx, "", opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAssignDevicesToServer_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"
	deviceIDs := []string{"XABC123X0ABC123X0", "YDEF456Y1DEF456Y1"}

	result, resp, err := client.AssignDevicesV1(ctx, serverID, deviceIDs)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
	require.NotNil(t, result)

	// Verify activity data
	activity := result.Data
	assert.Equal(t, "orgDeviceActivities", activity.Type)
	assert.NotEmpty(t, activity.ID)
	assert.NotNil(t, activity.Attributes)

	// Test activity attributes
	assert.Equal(t, ActivityStatusInProgress, activity.Attributes.Status)
	assert.Equal(t, ActivitySubStatusSubmitted, activity.Attributes.SubStatus)
	assert.NotNil(t, activity.Attributes.CreatedDateTime)

	// Verify links
	assert.NotNil(t, activity.Links)
	assert.Contains(t, activity.Links.Self, "orgDeviceActivities")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestAssignDevicesToServer_EmptyServerID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceIDs := []string{"XABC123X0ABC123X0"}

	result, _, err := client.AssignDevicesV1(ctx, "", deviceIDs)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAssignDevicesToServer_EmptyDeviceIDs(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"

	result, _, err := client.AssignDevicesV1(ctx, serverID, []string{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least one device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUnassignDevicesFromServer_Success(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"
	deviceIDs := []string{"XABC123X0ABC123X0"}

	result, resp, err := client.UnassignDevicesV1(ctx, serverID, deviceIDs)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
	require.NotNil(t, result)

	// Verify activity data
	activity := result.Data
	assert.Equal(t, "orgDeviceActivities", activity.Type)
	assert.NotEmpty(t, activity.ID)
	assert.NotNil(t, activity.Attributes)

	// Test activity attributes
	assert.Equal(t, ActivityStatusInProgress, activity.Attributes.Status)
	assert.Equal(t, ActivitySubStatusSubmitted, activity.Attributes.SubStatus)
	assert.NotNil(t, activity.Attributes.CreatedDateTime)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUnassignDevicesFromServer_EmptyServerID(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceIDs := []string{"XABC123X0ABC123X0"}

	result, _, err := client.UnassignDevicesV1(ctx, "", deviceIDs)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUnassignDevicesFromServer_EmptyDeviceIDs(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	serverID := "1F97349736CF4614A94F624E705841AD"

	result, _, err := client.UnassignDevicesV1(ctx, serverID, []string{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least one device ID is required")

	// No HTTP call should be made
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestContextCancellation(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	opts := &RequestQueryOptions{}

	result, _, err := client.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestContextTimeout(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
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

func TestFieldConstants(t *testing.T) {
	// Test that field constants are properly defined
	expectedFields := []string{
		FieldServerName,
		FieldServerType,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
		FieldDevices,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	// Test specific field values
	assert.Equal(t, "serverName", FieldServerName)
	assert.Equal(t, "serverType", FieldServerType)
	assert.Equal(t, "createdDateTime", FieldCreatedDateTime)
	assert.Equal(t, "updatedDateTime", FieldUpdatedDateTime)
	assert.Equal(t, "devices", FieldDevices)
}

func TestActivityConstants(t *testing.T) {
	// Test activity type constants
	assert.Equal(t, "ASSIGN_DEVICES", ActivityTypeAssignDevices)
	assert.Equal(t, "UNASSIGN_DEVICES", ActivityTypeUnassignDevices)

	// Test activity status constants
	assert.Equal(t, "IN_PROGRESS", ActivityStatusInProgress)
	assert.Equal(t, "COMPLETED", ActivityStatusCompleted)
	assert.Equal(t, "FAILED", ActivityStatusFailed)

	// Test activity sub-status constants
	assert.Equal(t, "SUBMITTED", ActivitySubStatusSubmitted)
	assert.Equal(t, "PROCESSING", ActivitySubStatusProcessing)
}

func TestOptionsStructures(t *testing.T) {
	// Test RequestQueryOptions
	opts1 := &RequestQueryOptions{
		Fields: []string{FieldServerName, FieldServerType},
		Limit:  100,
	}
	assert.Len(t, opts1.Fields, 2)
	assert.Equal(t, 100, opts1.Limit)

	// Test RequestQueryOptions
	opts2 := &RequestQueryOptions{
		Limit: 50,
	}
	assert.Equal(t, 50, opts2.Limit)

	// Test RequestQueryOptions
	opts3 := &RequestQueryOptions{
		Fields: []string{FieldServerName, FieldCreatedDateTime},
	}
	assert.Len(t, opts3.Fields, 2)
}

func TestComprehensiveFieldCoverage(t *testing.T) {
	client := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	opts := &RequestQueryOptions{
		Fields: []string{
			FieldServerName,
			FieldServerType,
			FieldEnableMdmDisownFlag,
			FieldDefaultProductFamilies,
			FieldStatus,
			FieldDeviceCount,
			FieldLastConnectedDateTime,
			FieldLastConnectedIp,
			FieldCreatedDateTime,
			FieldUpdatedDateTime,
			FieldDevices,
		},
	}

	result, resp, err := client.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	if len(result.Data) > 0 {
		attrs := result.Data[0].Attributes
		_ = attrs.ServerName
		_ = attrs.ServerType
		_ = attrs.EnableMdmDisownFlag
		_ = attrs.DefaultProductFamilies
		_ = attrs.Status
		_ = attrs.DeviceCount
		_ = attrs.LastConnectedDateTime
		_ = attrs.LastConnectedIp
		_ = attrs.CreatedDateTime
		_ = attrs.UpdatedDateTime
		_ = attrs.Devices
	}
}

// ====== GetByMDMServerIDV1 tests ======

func TestGetDeviceManagementServiceByID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldServerName, FieldServerType, FieldStatus},
	}

	result, resp, err := svc.GetByMDMServerIDV1(ctx, "1F97349736CF4614A94F624E705841AD", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	server := result.Data
	assert.Equal(t, "mdmServers", server.Type)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", server.ID)
	require.NotNil(t, server.Attributes)
	assert.Equal(t, "Production MDM", server.Attributes.ServerName)
	assert.Equal(t, "MDM", server.Attributes.ServerType)
	assert.Equal(t, MDMServerStatusActive, server.Attributes.Status)
	assert.Equal(t, 128, server.Attributes.DeviceCount)
	assert.False(t, server.Attributes.EnableMdmDisownFlag)
	assert.Contains(t, server.Attributes.DefaultProductFamilies, "MAC")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServiceByID_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByMDMServerIDV1(ctx, "1F97349736CF4614A94F624E705841AD", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServiceByID_EmptyServerID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByMDMServerIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetDeviceManagementServiceByID_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	httpmock.Reset()
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByMDMServerIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// ====== CreateMDMServerV1 tests ======

func TestCreateDeviceManagementService_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &MDMServerCreateRequest{
		Data: MDMServerCreateRequestData{
			Type: "mdmServers",
			Attributes: MDMServerCreateRequestAttributes{
				ServerName: "Marketing Team MDM",
				ServerCertificate: MDMServerCertificate{
					Name: "marketing-mdm.cer",
					Data: "MIIDXTCCAkWgAwIBAgIJALxxxxxxx...",
				},
				EnableMdmDisownFlag: true,
			},
		},
	}

	result, resp, err := svc.CreateMDMServerV1(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
	require.NotNil(t, result)

	server := result.Data
	assert.Equal(t, "mdmServers", server.Type)
	assert.Equal(t, "2A87349736CF4614A94F624E705841BE", server.ID)
	require.NotNil(t, server.Attributes)
	assert.Equal(t, "Marketing Team MDM", server.Attributes.ServerName)
	assert.True(t, server.Attributes.EnableMdmDisownFlag)
	assert.Equal(t, MDMServerStatusActive, server.Attributes.Status)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateDeviceManagementService_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.CreateMDMServerV1(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCreateDeviceManagementService_MissingServerName(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &MDMServerCreateRequest{
		Data: MDMServerCreateRequestData{
			Type: "mdmServers",
			Attributes: MDMServerCreateRequestAttributes{
				ServerCertificate: MDMServerCertificate{
					Name: "cert.cer",
					Data: "MIIDXTCCAkWg...",
				},
			},
		},
	}

	result, _, err := svc.CreateMDMServerV1(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "serverName is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCreateDeviceManagementService_MissingCertificateName(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &MDMServerCreateRequest{
		Data: MDMServerCreateRequestData{
			Type: "mdmServers",
			Attributes: MDMServerCreateRequestAttributes{
				ServerName: "Test MDM",
				ServerCertificate: MDMServerCertificate{
					Data: "MIIDXTCCAkWg...",
				},
			},
		},
	}

	result, _, err := svc.CreateMDMServerV1(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "serverCertificate.name is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCreateDeviceManagementService_MissingCertificateData(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &MDMServerCreateRequest{
		Data: MDMServerCreateRequestData{
			Type: "mdmServers",
			Attributes: MDMServerCreateRequestAttributes{
				ServerName: "Test MDM",
				ServerCertificate: MDMServerCertificate{
					Name: "cert.cer",
				},
			},
		},
	}

	result, _, err := svc.CreateMDMServerV1(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "serverCertificate.data is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

// ====== UpdateMDMServerByIDV1 tests ======

func TestUpdateDeviceManagementService_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	enableDisown := true
	req := &MDMServerUpdateRequest{
		Data: MDMServerUpdateRequestData{
			Type: "mdmServers",
			ID:   "1F97349736CF4614A94F624E705841AD",
			Attributes: MDMServerUpdateRequestAttributes{
				ServerName:          "Production MDM Updated",
				EnableMdmDisownFlag: &enableDisown,
			},
		},
	}

	result, resp, err := svc.UpdateMDMServerByIDV1(ctx, "1F97349736CF4614A94F624E705841AD", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	server := result.Data
	assert.Equal(t, "1F97349736CF4614A94F624E705841AD", server.ID)
	require.NotNil(t, server.Attributes)
	assert.Equal(t, "Production MDM Updated", server.Attributes.ServerName)
	assert.True(t, server.Attributes.EnableMdmDisownFlag)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUpdateDeviceManagementService_EmptyServerID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &MDMServerUpdateRequest{
		Data: MDMServerUpdateRequestData{
			Type: "mdmServers",
			Attributes: MDMServerUpdateRequestAttributes{
				ServerName: "Updated Name",
			},
		},
	}

	result, _, err := svc.UpdateMDMServerByIDV1(ctx, "", req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUpdateDeviceManagementService_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.UpdateMDMServerByIDV1(ctx, "1F97349736CF4614A94F624E705841AD", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUpdateDeviceManagementService_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	httpmock.Reset()
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &MDMServerUpdateRequest{
		Data: MDMServerUpdateRequestData{
			Type: "mdmServers",
			ID:   "NONEXISTENT",
			Attributes: MDMServerUpdateRequestAttributes{
				ServerName: "Updated Name",
			},
		},
	}

	result, resp, err := svc.UpdateMDMServerByIDV1(ctx, "NONEXISTENT", req)

	require.Error(t, err)
	assert.Nil(t, result)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// ====== DeleteMDMServerByIDV1 tests ======

func TestDeleteDeviceManagementService_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteMDMServerByIDV1(ctx, "1F97349736CF4614A94F624E705841AD")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestDeleteDeviceManagementService_EmptyServerID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteMDMServerByIDV1(ctx, "")

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "MDM server ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestDeleteDeviceManagementService_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	httpmock.Reset()
	mockHandler := &mocks.DeviceManagementMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteMDMServerByIDV1(ctx, "NONEXISTENT")

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// ====== New field and status constants tests ======

func TestMDMServerFieldConstants(t *testing.T) {
	assert.Equal(t, "serverName", FieldServerName)
	assert.Equal(t, "serverType", FieldServerType)
	assert.Equal(t, "enableMdmDisownFlag", FieldEnableMdmDisownFlag)
	assert.Equal(t, "defaultProductFamilies", FieldDefaultProductFamilies)
	assert.Equal(t, "status", FieldStatus)
	assert.Equal(t, "deviceCount", FieldDeviceCount)
	assert.Equal(t, "lastConnectedDateTime", FieldLastConnectedDateTime)
	assert.Equal(t, "lastConnectedIp", FieldLastConnectedIp)
	assert.Equal(t, "createdDateTime", FieldCreatedDateTime)
	assert.Equal(t, "updatedDateTime", FieldUpdatedDateTime)
	assert.Equal(t, "devices", FieldDevices)
}

func TestMDMServerStatusConstants(t *testing.T) {
	assert.Equal(t, "ACTIVE", MDMServerStatusActive)
	assert.Equal(t, "INACTIVE", MDMServerStatusInactive)
}
