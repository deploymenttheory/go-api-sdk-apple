package configurations

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/configurations/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *Configurations {
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

func TestGetConfigurations_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldType, FieldConfiguredForPlatforms},
		Limit:  100,
	}

	result, resp, err := svc.GetV1(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 2)

	config1 := result.Data[0]
	assert.Equal(t, "configurations", config1.Type)
	assert.Equal(t, "config-12345", config1.ID)
	require.NotNil(t, config1.Attributes)
	assert.Equal(t, ConfigurationTypeCustomSetting, config1.Attributes.Type)
	assert.Equal(t, "Wi-Fi Configuration", config1.Attributes.Name)
	assert.Contains(t, config1.Attributes.ConfiguredForPlatforms, PlatformIOS)
	assert.Contains(t, config1.Attributes.ConfiguredForPlatforms, PlatformMacOS)
	// customSettingsValues is null in list response per API spec
	assert.Nil(t, config1.Attributes.CustomSettingsValues)

	config2 := result.Data[1]
	assert.Equal(t, ConfigurationTypeAirDrop, config2.Attributes.Type)
	assert.Equal(t, "Air Drop Configuration", config2.Attributes.Name)

	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Links)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetConfigurations_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
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

func TestGetConfigurations_WithLimitEnforcement(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
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

func TestGetConfigurations_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()

	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetConfigurationInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &RequestQueryOptions{
		Fields: []string{FieldName, FieldType, FieldCustomSettingsValues},
	}

	result, resp, err := svc.GetByConfigurationIDV1(ctx, "config-12345", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	config := result.Data
	assert.Equal(t, "configurations", config.Type)
	assert.Equal(t, "config-12345", config.ID)
	require.NotNil(t, config.Attributes)
	assert.Equal(t, ConfigurationTypeCustomSetting, config.Attributes.Type)
	assert.Equal(t, "Wi-Fi Configuration", config.Attributes.Name)

	// customSettingsValues is present in single resource response
	require.NotNil(t, config.Attributes.CustomSettingsValues)
	assert.NotEmpty(t, config.Attributes.CustomSettingsValues.ConfigurationProfile)
	assert.Equal(t, "WiFi.mobileconfig", config.Attributes.CustomSettingsValues.Filename)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetConfigurationInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByConfigurationIDV1(ctx, "config-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "config-12345", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetConfigurationInformation_EmptyConfigID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByConfigurationIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "configuration ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCreateConfiguration_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &ConfigurationCreateRequest{
		Data: ConfigurationCreateRequestData{
			Type: "configurations",
			Attributes: ConfigurationCreateRequestAttributes{
				Type: ConfigurationTypeCustomSetting,
				Name: "AirPlay Security Settings",
				ConfiguredForPlatforms: []string{PlatformIOS},
				CustomSettingsValues: CustomSettingsValues{
					ConfigurationProfile: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>...",
					Filename:             "Airplay.mobileconfig",
				},
			},
		},
	}

	result, resp, err := svc.CreateV1(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
	require.NotNil(t, result)

	config := result.Data
	assert.Equal(t, "configurations", config.Type)
	assert.Equal(t, "config-new-456", config.ID)
	require.NotNil(t, config.Attributes)
	assert.Equal(t, ConfigurationTypeCustomSetting, config.Attributes.Type)
	assert.Equal(t, "AirPlay Security Settings", config.Attributes.Name)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateConfiguration_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.CreateV1(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCreateConfiguration_MissingConfigurationProfile(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &ConfigurationCreateRequest{
		Data: ConfigurationCreateRequestData{
			Type: "configurations",
			Attributes: ConfigurationCreateRequestAttributes{
				Type: ConfigurationTypeCustomSetting,
				Name: "Test Config",
				// ConfigurationProfile is missing
			},
		},
	}

	result, _, err := svc.CreateV1(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "configurationProfile is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUpdateConfiguration_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &ConfigurationUpdateRequest{
		Data: ConfigurationUpdateRequestData{
			Type: "configurations",
			ID:   "config-12345",
			Attributes: ConfigurationUpdateRequestAttributes{
				Name: "Updated Wi-Fi Configuration",
				CustomSettingsValues: &CustomSettingsValues{
					ConfigurationProfile: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>...",
					Filename:             "WiFi-Updated.mobileconfig",
				},
			},
		},
	}

	result, resp, err := svc.UpdateByConfigurationIDV1(ctx, "config-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	config := result.Data
	assert.Equal(t, "config-12345", config.ID)
	require.NotNil(t, config.Attributes)
	assert.Equal(t, "Updated Wi-Fi Configuration", config.Attributes.Name)
	require.NotNil(t, config.Attributes.CustomSettingsValues)
	assert.Equal(t, "WiFi-Updated.mobileconfig", config.Attributes.CustomSettingsValues.Filename)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUpdateConfiguration_EmptyConfigID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &ConfigurationUpdateRequest{
		Data: ConfigurationUpdateRequestData{
			Type: "configurations",
			Attributes: ConfigurationUpdateRequestAttributes{
				Name: "Updated Name",
			},
		},
	}

	result, _, err := svc.UpdateByConfigurationIDV1(ctx, "", req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "configuration ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUpdateConfiguration_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.UpdateByConfigurationIDV1(ctx, "config-12345", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestDeleteConfiguration_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteByConfigurationIDV1(ctx, "config-12345")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestDeleteConfiguration_EmptyConfigID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteByConfigurationIDV1(ctx, "")

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "configuration ID is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestDeleteConfiguration_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteByConfigurationIDV1(ctx, "NONEXISTENT")

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetConfigurations_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetV1(ctx, &RequestQueryOptions{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetConfigurations_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.ConfigurationsMock{}
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

func TestConfigurationFieldConstants(t *testing.T) {
	expectedFields := []string{
		FieldType,
		FieldName,
		FieldConfiguredForPlatforms,
		FieldCustomSettingsValues,
		FieldCreatedDateTime,
		FieldUpdatedDateTime,
	}

	for _, field := range expectedFields {
		assert.NotEmpty(t, field, "Field constant should not be empty")
	}

	assert.Equal(t, "type", FieldType)
	assert.Equal(t, "name", FieldName)
	assert.Equal(t, "configuredForPlatforms", FieldConfiguredForPlatforms)
	assert.Equal(t, "customSettingsValues", FieldCustomSettingsValues)
}

func TestConfigurationTypeConstants(t *testing.T) {
	assert.Equal(t, "CUSTOM_SETTING", ConfigurationTypeCustomSetting)
	assert.Equal(t, "AIR_DROP", ConfigurationTypeAirDrop)
}

func TestPlatformConstants(t *testing.T) {
	assert.Equal(t, "PLATFORM_IOS", PlatformIOS)
	assert.Equal(t, "PLATFORM_MACOS", PlatformMacOS)
	assert.Equal(t, "PLATFORM_TVOS", PlatformTvOS)
}
