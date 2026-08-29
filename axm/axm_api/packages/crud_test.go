package packages

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/packages/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *Packages {
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

func TestGetPackages_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldVersion, FieldDescription},
		Limit:  100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	pkg := result.Data[0]
	assert.Equal(t, "packages", pkg.Type)
	assert.Equal(t, "pkg-12345", pkg.ID)
	require.NotNil(t, pkg.Attributes)

	assert.Equal(t, "Enterprise Software Suite", pkg.Attributes.Name)
	assert.Equal(t, "https://example.com/packages/enterprise-suite.pkg", pkg.Attributes.URL)
	assert.Equal(t, "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6", pkg.Attributes.Hash)
	assert.Equal(t, "Complete enterprise productivity software suite", pkg.Attributes.Description)
	assert.Equal(t, "2.5.1", pkg.Attributes.Version)
	assert.NotNil(t, pkg.Attributes.CreatedDateTime)
	assert.NotNil(t, pkg.Attributes.UpdatedDateTime)

	require.Len(t, pkg.Attributes.BundleIds, 1)
	assert.Equal(t, "com.example.enterpriseapp", pkg.Attributes.BundleIds[0])

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPackages_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
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

func TestGetPackages_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
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

func TestGetPackages_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPackageInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldVersion},
	}

	result, resp, err := svc.GetByPackageIDV1(ctx, "pkg-12345", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	pkg := result.Data
	assert.Equal(t, "packages", pkg.Type)
	assert.Equal(t, "pkg-12345", pkg.ID)
	require.NotNil(t, pkg.Attributes)
	assert.Equal(t, "Enterprise Software Suite", pkg.Attributes.Name)
	assert.Equal(t, "2.5.1", pkg.Attributes.Version)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPackageInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByPackageIDV1(ctx, "pkg-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "pkg-12345", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPackageInformation_EmptyPackageID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByPackageIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "package ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetPackageInformation_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByPackageIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPackageInformation_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetByPackageIDV1(ctx, "pkg-12345", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetPackages_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.PackagesMock{}
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

func TestPackageFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldName,
		FieldURL,
		FieldHash,
		FieldBundleIds,
		FieldDescription,
		FieldVersion,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "name", FieldName)
	assert.Equal(t, "url", FieldURL)
	assert.Equal(t, "hash", FieldHash)
	assert.Equal(t, "bundleIds", FieldBundleIds)
	assert.Equal(t, "version", FieldVersion)
}
