package organizationalunits

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/organizationalunits/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *OrganizationalUnits {
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

func TestGetOrganizationalUnits_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldDescription},
		Limit:  100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	unit := result.Data[0]
	assert.Equal(t, "organizationalUnits", unit.Type)
	assert.Equal(t, "OU789012", unit.ID)
	require.NotNil(t, unit.Attributes)

	assert.Equal(t, "Engineering", unit.Attributes.Name)
	assert.Equal(t, "Engineering organizational unit", unit.Attributes.Description)
	assert.NotNil(t, unit.Attributes.CreatedDateTime)
	assert.NotNil(t, unit.Attributes.UpdatedDateTime)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationalUnits_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
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

func TestGetOrganizationalUnits_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
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

func TestGetOrganizationalUnits_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationalUnitInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldDescription},
	}

	result, resp, err := svc.GetByOrganizationalUnitIDV1(ctx, "OU789012", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	unit := result.Data
	assert.Equal(t, "organizationalUnits", unit.Type)
	assert.Equal(t, "OU789012", unit.ID)
	require.NotNil(t, unit.Attributes)
	assert.Equal(t, "Engineering", unit.Attributes.Name)
	assert.Equal(t, "Engineering organizational unit", unit.Attributes.Description)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationalUnitInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByOrganizationalUnitIDV1(ctx, "OU789012", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "OU789012", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationalUnitInformation_EmptyUnitID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByOrganizationalUnitIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "organizational unit ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetOrganizationalUnitInformation_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByOrganizationalUnitIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByOrganizationalUnitID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{Limit: 100}

	result, resp, err := svc.GetUserIDsByOrganizationalUnitIDV1(ctx, "OU789012", opts)

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

func TestGetUserIDsByOrganizationalUnitID_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetUserIDsByOrganizationalUnitIDV1(ctx, "OU789012", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByOrganizationalUnitID_EmptyUnitID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetUserIDsByOrganizationalUnitIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "organizational unit ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByOrganizationalUnitID_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetUserIDsByOrganizationalUnitIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetOrganizationalUnits_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetOrganizationalUnits_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.OrganizationalUnitsMock{}
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

func TestOrganizationalUnitFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldName,
		FieldDescription,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
		FieldUsers,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "name", FieldName)
	assert.Equal(t, "description", FieldDescription)
	assert.Equal(t, "createdDateTime", FieldCreatedDateTime)
	assert.Equal(t, "updatedDateTime", FieldUpdatedDateTime)
	assert.Equal(t, "users", FieldUsers)
}
