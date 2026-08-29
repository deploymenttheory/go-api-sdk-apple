package auditevents

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/auditevents/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *AuditEvents {
	mockAuth := &MockAuthProvider{}

	coreClient, err := client.NewTransport(
		"test-key-id",
		"test-issuer-id",
		"dummy-key",
		client.WithAuth(mockAuth),
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(coreClient.GetHTTPClient().Client())

	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewService(coreClient)
}

// MockAuthProvider implements the AuthProvider interface for testing.
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	return nil
}

func TestGetAuditEvents_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
		FilterEndTimestamp:   "2026-02-14T23:59:59Z",
		Limit:                100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	event := result.Data[0]
	assert.Equal(t, "auditEvents", event.Type)
	assert.Equal(t, "event-12345", event.ID)
	require.NotNil(t, event.Attributes)

	assert.Equal(t, "DEVICE_ADDED_TO_ORG", event.Attributes.Type)
	assert.Equal(t, "DEVICE_INVENTORY", event.Attributes.Category)
	assert.Equal(t, "USER", event.Attributes.ActorType)
	assert.Equal(t, "user-abc123", event.Attributes.ActorId)
	assert.Equal(t, "elana.landot@melardclothing.com", event.Attributes.ActorName)
	assert.Equal(t, "DEVICE", event.Attributes.SubjectType)
	assert.Equal(t, "device-xyz789", event.Attributes.SubjectId)
	assert.Equal(t, "MacBook Pro", event.Attributes.SubjectName)
	assert.Equal(t, "SUCCESS", event.Attributes.Outcome)
	assert.Equal(t, "group-001", event.Attributes.GroupId)
	assert.Equal(t, "eventDataDeviceAddedToOrg", event.Attributes.EventDataPropertyKey)

	require.NotNil(t, event.Attributes.EventDataDeviceAddedToOrg)
	assert.Equal(t, "C02X1234ABCD", event.Attributes.EventDataDeviceAddedToOrg.SerialNumber)
	assert.Equal(t, "APPLE", event.Attributes.EventDataDeviceAddedToOrg.PurchaseSourceType)
	assert.Equal(t, "order-56789", event.Attributes.EventDataDeviceAddedToOrg.PurchaseSourceId)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetV1(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "opts is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_MissingStartTimestamp(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		FilterEndTimestamp: "2026-02-14T23:59:59Z",
	}

	result, _, err := svc.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "filter[startTimestamp] is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_MissingEndTimestamp(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
	}

	result, _, err := svc.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "filter[endTimestamp] is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_WithAllFilters(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
		FilterEndTimestamp:   "2026-02-14T23:59:59Z",
		FilterActorID:        "user-abc123",
		FilterSubjectID:      "device-xyz789",
		FilterType:           AuditEventTypeDeviceAddedToOrg,
		Fields: []string{
			FieldEventDateTime,
			FieldType,
			FieldCategory,
			FieldActorType,
			FieldActorId,
			FieldOutcome,
		},
		Limit: 50,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
		FilterEndTimestamp:   "2026-02-14T23:59:59Z",
		Limit:                1500,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
		FilterEndTimestamp:   "2026-02-14T23:59:59Z",
	}

	result, _, err := svc.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAuditEvents_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
		FilterEndTimestamp:   "2026-02-14T23:59:59Z",
	}

	result, _, err := svc.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetAuditEvents_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AuditEventsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Millisecond)

	opts := &RequestQueryOptions{
		FilterStartTimestamp: "2026-02-14T00:00:00Z",
		FilterEndTimestamp:   "2026-02-14T23:59:59Z",
	}

	result, _, err := svc.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestAuditEventTypeConstants(t *testing.T) {
	assert.Equal(t, "DEVICE_ADDED_TO_ORG", AuditEventTypeDeviceAddedToOrg)
	assert.Equal(t, "DEVICE_REMOVED_FROM_ORG", AuditEventTypeDeviceRemovedFromOrg)
	assert.Equal(t, "DEVICE_ASSIGNED_TO_SERVER", AuditEventTypeDeviceAssignedToServer)
	assert.Equal(t, "DEVICE_UNASSIGNED_FROM_SERVER", AuditEventTypeDeviceUnassignedFromServer)
	assert.Equal(t, "DEVICE_IS_ERASED", AuditEventTypeDeviceIsErased)
	assert.Equal(t, "CONFIG_SETTINGS_CREATED", AuditEventTypeConfigSettingsCreated)
	assert.Equal(t, "CONFIG_SETTINGS_UPDATED", AuditEventTypeConfigSettingsUpdated)
	assert.Equal(t, "CONFIG_SETTINGS_DELETED", AuditEventTypeConfigSettingsDeleted)
	assert.Equal(t, "ACCOUNT_ADDED", AuditEventTypeAccountAdded)
	assert.Equal(t, "ACCOUNT_DELETED", AuditEventTypeAccountDeleted)
	assert.Equal(t, "API_ACCOUNT_CREATED_WITH_KEY", AuditEventTypeAPIAccountCreatedWithKey)
	assert.Equal(t, "API_ACCOUNT_DELETED", AuditEventTypeAPIAccountDeleted)
	assert.Equal(t, "SUCCESS", AuditEventOutcomeSuccess)
	assert.Equal(t, "FAILURE", AuditEventOutcomeFailure)
}

func TestAuditEventFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldEventDateTime,
		FieldType,
		FieldCategory,
		FieldActorType,
		FieldActorId,
		FieldActorName,
		FieldSubjectType,
		FieldSubjectId,
		FieldSubjectName,
		FieldOutcome,
		FieldGroupId,
		FieldEventDataPropertyKey,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "eventDateTime", FieldEventDateTime)
	assert.Equal(t, "type", FieldType)
	assert.Equal(t, "category", FieldCategory)
	assert.Equal(t, "outcome", FieldOutcome)
}
