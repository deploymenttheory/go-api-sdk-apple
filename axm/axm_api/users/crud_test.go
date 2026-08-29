package users

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/users/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *Users {
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

func TestGetUsers_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldFirstName, FieldLastName, FieldStatus},
		Limit:  100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	user := result.Data[0]
	assert.Equal(t, "users", user.Type)
	assert.Equal(t, "1234567890", user.ID)
	require.NotNil(t, user.Attributes)

	assert.Equal(t, "John", user.Attributes.FirstName)
	assert.Equal(t, "Doe", user.Attributes.LastName)
	assert.Equal(t, "A", user.Attributes.MiddleName)
	assert.Equal(t, "ACTIVE", user.Attributes.Status)
	assert.Equal(t, "john.doe@appleid.example.com", user.Attributes.ManagedAppleAccount)
	assert.False(t, user.Attributes.IsExternalUser)
	assert.Equal(t, "john.doe@example.com", user.Attributes.Email)
	assert.Equal(t, "EMP001", user.Attributes.EmployeeNumber)
	assert.Equal(t, "CC100", user.Attributes.CostCenter)
	assert.Equal(t, "Engineering", user.Attributes.Division)
	assert.Equal(t, "IT", user.Attributes.Department)
	assert.Equal(t, "Software Engineer", user.Attributes.JobTitle)

	require.Len(t, user.Attributes.RoleOuList, 1)
	assert.Equal(t, "Administrator", user.Attributes.RoleOuList[0].RoleName)
	assert.Equal(t, "OU123456", user.Attributes.RoleOuList[0].OuId)

	require.Len(t, user.Attributes.PhoneNumbers, 1)
	assert.Equal(t, "+1-555-123-4567", user.Attributes.PhoneNumbers[0].PhoneNumber)
	assert.Equal(t, "WORK", user.Attributes.PhoneNumbers[0].Type)

	assert.NotNil(t, user.Attributes.StartDateTime)
	assert.NotNil(t, user.Attributes.CreatedDateTime)
	assert.NotNil(t, user.Attributes.UpdatedDateTime)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUsers_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
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

func TestGetUsers_WithFieldSelection(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldFirstName, FieldLastName, FieldEmail, FieldStatus},
		Limit:  50,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUsers_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Limit: 1500,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUsers_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{}

	result, _, err := svc.GetV1(ctx, opts)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	userID := "1234567890"
	opts := &RequestQueryOptions{
		Fields: []string{FieldFirstName, FieldLastName, FieldStatus},
	}

	result, resp, err := svc.GetByUserIDV1(ctx, userID, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	user := result.Data
	assert.Equal(t, "users", user.Type)
	assert.Equal(t, "1234567890", user.ID)
	require.NotNil(t, user.Attributes)

	assert.Equal(t, "John", user.Attributes.FirstName)
	assert.Equal(t, "Doe", user.Attributes.LastName)
	assert.Equal(t, "ACTIVE", user.Attributes.Status)
	assert.Equal(t, "john.doe@example.com", user.Attributes.Email)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByUserIDV1(ctx, "1234567890", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "1234567890", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserInformation_EmptyUserID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByUserIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetUserInformation_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByUserIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserInformation_HTTPError(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByUserIDV1(ctx, "1234567890", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserInformation_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetByUserIDV1(ctx, "1234567890", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetUsers_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.UsersMock{}
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

func TestUserFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldFirstName,
		FieldLastName,
		FieldMiddleName,
		FieldStatus,
		FieldManagedAppleAccount,
		FieldIsExternalUser,
		FieldRoleOuList,
		FieldEmail,
		FieldEmployeeNumber,
		FieldCostCenter,
		FieldDivision,
		FieldDepartment,
		FieldJobTitle,
		FieldStartDateTime,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
		FieldPhoneNumbers,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "firstName", FieldFirstName)
	assert.Equal(t, "lastName", FieldLastName)
	assert.Equal(t, "status", FieldStatus)
	assert.Equal(t, "email", FieldEmail)
}

func TestUserStatusConstants(t *testing.T) {
	assert.Equal(t, "ACTIVE", UserStatusActive)
	assert.Equal(t, "INACTIVE", UserStatusInactive)
}
