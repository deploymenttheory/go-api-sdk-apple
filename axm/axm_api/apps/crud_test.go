package apps

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/apps/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *Apps {
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

func TestGetApps_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldBundleId, FieldVersion},
		Limit:  100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Data)

	app := result.Data[0]
	assert.Equal(t, "apps", app.Type)
	assert.Equal(t, "361309726", app.ID)
	require.NotNil(t, app.Attributes)

	assert.Equal(t, "Pages", app.Attributes.Name)
	assert.Equal(t, "com.apple.Pages", app.Attributes.BundleId)
	assert.Equal(t, "https://www.apple.com/pages/", app.Attributes.WebsiteUrl)
	assert.Equal(t, "14.0", app.Attributes.Version)
	assert.False(t, app.Attributes.IsCustomApp)
	assert.Equal(t, "https://apps.apple.com/app/pages/id361309726", app.Attributes.AppStoreUrl)

	require.Len(t, app.Attributes.SupportedOS, 2)
	assert.Contains(t, app.Attributes.SupportedOS, SupportedOSiOS)
	assert.Contains(t, app.Attributes.SupportedOS, SupportedOSmacOS)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetApps_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
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

func TestGetApps_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
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

func TestGetApps_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldBundleId},
	}

	result, resp, err := svc.GetByAppIDV1(ctx, "361309726", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	app := result.Data
	assert.Equal(t, "apps", app.Type)
	assert.Equal(t, "361309726", app.ID)
	require.NotNil(t, app.Attributes)
	assert.Equal(t, "Pages", app.Attributes.Name)
	assert.Equal(t, "com.apple.Pages", app.Attributes.BundleId)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByAppIDV1(ctx, "361309726", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "361309726", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppInformation_EmptyAppID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByAppIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "app ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAppInformation_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByAppIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppInformation_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetByAppIDV1(ctx, "361309726", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetApps_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.AppsMock{}
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

func TestAppFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldName,
		FieldBundleId,
		FieldWebsiteUrl,
		FieldVersion,
		FieldSupportedOS,
		FieldIsCustomApp,
		FieldAppStoreUrl,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "name", FieldName)
	assert.Equal(t, "bundleId", FieldBundleId)
	assert.Equal(t, "version", FieldVersion)
	assert.Equal(t, "supportedOS", FieldSupportedOS)
}

func TestSupportedOSConstants(t *testing.T) {
	assert.Equal(t, "SUPPORTED_OS_IOS", SupportedOSiOS)
	assert.Equal(t, "SUPPORTED_OS_MACOS", SupportedOSmacOS)
	assert.Equal(t, "SUPPORTED_OS_TVOS", SupportedOStvOS)
	assert.Equal(t, "SUPPORTED_OS_WATCHOS", SupportedOSwatchOS)
}
