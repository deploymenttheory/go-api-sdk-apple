package usergroups

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/usergroups/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *UserGroups {
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

func TestGetUserGroups_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldStatus, FieldTotalMemberCount},
		Limit:  100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	group := result.Data[0]
	assert.Equal(t, "userGroups", group.Type)
	assert.Equal(t, "UG123456", group.ID)
	require.NotNil(t, group.Attributes)

	assert.Equal(t, "OU789012", group.Attributes.OuId)
	assert.Equal(t, "Engineering Team", group.Attributes.Name)
	assert.Equal(t, "STANDARD", group.Attributes.Type)
	assert.Equal(t, 25, group.Attributes.TotalMemberCount)
	assert.Equal(t, "ACTIVE", group.Attributes.Status)
	assert.NotNil(t, group.Attributes.CreatedDateTime)
	assert.NotNil(t, group.Attributes.UpdatedDateTime)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroups_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetV1(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroups_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{Limit: 1500}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroups_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroupInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldStatus},
	}

	result, resp, err := svc.GetByUserGroupIDV1(ctx, "UG123456", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	group := result.Data
	assert.Equal(t, "userGroups", group.Type)
	assert.Equal(t, "UG123456", group.ID)
	require.NotNil(t, group.Attributes)
	assert.Equal(t, "Engineering Team", group.Attributes.Name)
	assert.Equal(t, "ACTIVE", group.Attributes.Status)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroupInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByUserGroupIDV1(ctx, "UG123456", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "UG123456", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroupInformation_EmptyGroupID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByUserGroupIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "group ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetUserGroupInformation_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByUserGroupIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByGroupID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{Limit: 100}

	result, resp, err := svc.GetUserIDsByGroupIDV1(ctx, "UG123456", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Len(t, result.Data, 2)
	assert.Equal(t, "users", result.Data[0].Type)
	assert.Equal(t, "1234567890", result.Data[0].ID)
	assert.Equal(t, "users", result.Data[1].Type)
	assert.Equal(t, "0987654321", result.Data[1].ID)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByGroupID_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetUserIDsByGroupIDV1(ctx, "UG123456", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByGroupID_EmptyGroupID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetUserIDsByGroupIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "group ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByGroupID_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetUserIDsByGroupIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroups_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetUserGroups_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UserGroupsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Millisecond)

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestUserGroupFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldOuId,
		FieldName,
		FieldType,
		FieldTotalMemberCount,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
		FieldStatus,
		FieldUsers,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "ouId", FieldOuId)
	assert.Equal(t, "name", FieldName)
	assert.Equal(t, "status", FieldStatus)
	assert.Equal(t, "totalMemberCount", FieldTotalMemberCount)
}

func TestUserGroupStatusConstants(t *testing.T) {
	assert.Equal(t, "ACTIVE", UserGroupStatusActive)
	assert.Equal(t, "INACTIVE", UserGroupStatusInactive)
	assert.Equal(t, "STANDARD", UserGroupTypeStandard)
}
