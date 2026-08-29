package blueprints

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/blueprints/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *Blueprints {
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

// --- Create ---

func TestCreateBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintCreateRequest{
		Data: BlueprintCreateRequestData{
			Type: "blueprints",
			Attributes: BlueprintCreateRequestAttributes{
				Name:        "Marketing Team Blueprint",
				Description: "Apps and settings for marketing team members",
			},
			Relationships: &BlueprintRequestRelationships{
				Apps: &BlueprintLinkageData{
					Data: []BlueprintLinkage{{Type: "apps", ID: "361309726"}},
				},
			},
		},
	}

	result, resp, err := svc.CreateV1(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 201, resp.StatusCode())
	require.NotNil(t, result)

	bp := result.Data
	assert.Equal(t, "blueprints", bp.Type)
	assert.Equal(t, "blueprint-new-123", bp.ID)
	require.NotNil(t, bp.Attributes)
	assert.Equal(t, "Marketing Team Blueprint", bp.Attributes.Name)
	assert.Equal(t, "Apps and settings for marketing team members", bp.Attributes.Description)
	assert.Equal(t, BlueprintStatusActive, bp.Attributes.Status)
	assert.False(t, bp.Attributes.AppLicenseDeficient)
	require.NotNil(t, bp.Attributes.CreatedDateTime)
	require.NotNil(t, bp.Attributes.UpdatedDateTime)

	require.NotNil(t, bp.Relationships)
	assert.NotNil(t, bp.Relationships.Apps)
	assert.NotNil(t, bp.Relationships.Users)
	assert.NotNil(t, bp.Relationships.Configurations)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateBlueprint_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.CreateV1(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestCreateBlueprint_EmptyName(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintCreateRequest{
		Data: BlueprintCreateRequestData{
			Type:       "blueprints",
			Attributes: BlueprintCreateRequestAttributes{},
		},
	}

	result, _, err := svc.CreateV1(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint name is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

// --- Get by ID ---

func TestGetBlueprintInformation_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	opts := &GetBlueprintQueryOptions{
		Fields:  []string{FieldName, FieldDescription, FieldStatus},
		Include: []string{IncludeApps, IncludeUsers},
	}

	result, resp, err := svc.GetByBlueprintIDV1(ctx, "blueprint-12345", opts)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	bp := result.Data
	assert.Equal(t, "blueprints", bp.Type)
	assert.Equal(t, "blueprint-12345", bp.ID)
	require.NotNil(t, bp.Attributes)
	assert.Equal(t, "Engineering Onboarding", bp.Attributes.Name)
	assert.Equal(t, "Standard apps for engineering team", bp.Attributes.Description)
	assert.Equal(t, BlueprintStatusActive, bp.Attributes.Status)

	require.NotNil(t, bp.Relationships)
	assert.NotNil(t, bp.Relationships.Apps)
	assert.NotNil(t, bp.Relationships.Users)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetBlueprintInformation_WithNilOptions(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "blueprint-12345", result.Data.ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetBlueprintInformation_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetBlueprintInformation_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetByBlueprintIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetBlueprintInformation_HTTPError(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetBlueprintInformation_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGetBlueprintInformation_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond)

	result, _, err := svc.GetByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// --- Update ---

func TestUpdateBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintUpdateRequest{
		Data: BlueprintUpdateRequestData{
			Type: "blueprints",
			ID:   "blueprint-12345",
			Attributes: BlueprintUpdateRequestAttributes{
				Description: "Updated description for engineering team",
			},
		},
	}

	result, resp, err := svc.UpdateByBlueprintIDV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	bp := result.Data
	assert.Equal(t, "blueprint-12345", bp.ID)
	require.NotNil(t, bp.Attributes)
	assert.Equal(t, "Updated description for engineering team", bp.Attributes.Description)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUpdateBlueprint_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintUpdateRequest{
		Data: BlueprintUpdateRequestData{
			Type:       "blueprints",
			Attributes: BlueprintUpdateRequestAttributes{Description: "Updated"},
		},
	}

	result, _, err := svc.UpdateByBlueprintIDV1(ctx, "", req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUpdateBlueprint_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.UpdateByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

// --- Delete ---

func TestDeleteBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteByBlueprintIDV1(ctx, "blueprint-12345")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestDeleteBlueprint_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteByBlueprintIDV1(ctx, "")

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestDeleteBlueprint_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.DeleteByBlueprintIDV1(ctx, "NONEXISTENT")

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Relationship: Apps ---

func TestGetAppIDsByBlueprintID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetAppIDsByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "apps", result.Data[0].Type)
	assert.Equal(t, "361309726", result.Data[0].ID)
	assert.Equal(t, "409201541", result.Data[1].ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetAppIDsByBlueprintID_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetAppIDsByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetAppIDsByBlueprintID_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetAppIDsByBlueprintIDV1(ctx, "NONEXISTENT", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestAddAppsToBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintAppsLinkagesRequest{
		Data: []BlueprintLinkage{
			{Type: "apps", ID: "361309726"},
			{Type: "apps", ID: "409201541"},
		},
	}

	resp, err := svc.AddAppsToBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestAddAppsToBlueprint_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintAppsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "apps", ID: "361309726"}},
	}

	resp, err := svc.AddAppsToBlueprintV1(ctx, "", req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAddAppsToBlueprint_NilRequest(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	resp, err := svc.AddAppsToBlueprintV1(ctx, "blueprint-12345", nil)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "request is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestRemoveAppsFromBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintAppsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "apps", ID: "361309726"}},
	}

	resp, err := svc.RemoveAppsFromBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRemoveAppsFromBlueprint_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintAppsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "apps", ID: "361309726"}},
	}

	resp, err := svc.RemoveAppsFromBlueprintV1(ctx, "", req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

// --- Relationship: Configurations ---

func TestGetConfigurationIDsByBlueprintID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetConfigurationIDsByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "configurations", result.Data[0].Type)
	assert.Equal(t, "config-12345", result.Data[0].ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetConfigurationIDsByBlueprintID_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetConfigurationIDsByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAddConfigurationsToBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintConfigurationsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "configurations", ID: "config-12345"}},
	}

	resp, err := svc.AddConfigurationsToBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRemoveConfigurationsFromBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintConfigurationsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "configurations", ID: "config-12345"}},
	}

	resp, err := svc.RemoveConfigurationsFromBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Relationship: Packages ---

func TestGetPackageIDsByBlueprintID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetPackageIDsByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "packages", result.Data[0].Type)
	assert.Equal(t, "pkg-12345", result.Data[0].ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPackageIDsByBlueprintID_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetPackageIDsByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAddPackagesToBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintPackagesLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "packages", ID: "pkg-12345"}},
	}

	resp, err := svc.AddPackagesToBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRemovePackagesFromBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintPackagesLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "packages", ID: "pkg-12345"}},
	}

	resp, err := svc.RemovePackagesFromBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Relationship: OrgDevices ---

func TestGetDeviceIDsByBlueprintID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetDeviceIDsByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "orgDevices", result.Data[0].Type)
	assert.Equal(t, "ABC123DEF456", result.Data[0].ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetDeviceIDsByBlueprintID_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetDeviceIDsByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAddDevicesToBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintOrgDevicesLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "orgDevices", ID: "ABC123DEF456"}},
	}

	resp, err := svc.AddDevicesToBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRemoveDevicesFromBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintOrgDevicesLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "orgDevices", ID: "ABC123DEF456"}},
	}

	resp, err := svc.RemoveDevicesFromBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Relationship: Users ---

func TestGetUserIDsByBlueprintID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetUserIDsByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "users", result.Data[0].Type)
	assert.Equal(t, "1234567890", result.Data[0].ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserIDsByBlueprintID_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetUserIDsByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAddUsersToBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintUsersLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "users", ID: "1234567890"}},
	}

	resp, err := svc.AddUsersToBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRemoveUsersFromBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintUsersLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "users", ID: "1234567890"}},
	}

	resp, err := svc.RemoveUsersFromBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Relationship: UserGroups ---

func TestGetUserGroupIDsByBlueprintID_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, resp, err := svc.GetUserGroupIDsByBlueprintIDV1(ctx, "blueprint-12345", nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "userGroups", result.Data[0].Type)
	assert.Equal(t, "e0484524-fff2-4132-ad29-fe7c6258ce53", result.Data[0].ID)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetUserGroupIDsByBlueprintID_EmptyBlueprintID(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := svc.GetUserGroupIDsByBlueprintIDV1(ctx, "", nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blueprint ID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestAddUserGroupsToBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintUserGroupsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "userGroups", ID: "e0484524-fff2-4132-ad29-fe7c6258ce53"}},
	}

	resp, err := svc.AddUserGroupsToBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRemoveUserGroupsFromBlueprint_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.BlueprintsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &BlueprintUserGroupsLinkagesRequest{
		Data: []BlueprintLinkage{{Type: "userGroups", ID: "e0484524-fff2-4132-ad29-fe7c6258ce53"}},
	}

	resp, err := svc.RemoveUserGroupsFromBlueprintV1(ctx, "blueprint-12345", req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 204, resp.StatusCode())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Constants ---

func TestBlueprintFieldConstants(t *testing.T) {
	assert.Equal(t, "name", FieldName)
	assert.Equal(t, "description", FieldDescription)
	assert.Equal(t, "status", FieldStatus)
	assert.Equal(t, "createdDateTime", FieldCreatedDateTime)
	assert.Equal(t, "updatedDateTime", FieldUpdatedDateTime)
	assert.Equal(t, "appLicenseDeficient", FieldAppLicenseDeficient)
	assert.Equal(t, "apps", FieldApps)
	assert.Equal(t, "packages", FieldPackages)
	assert.Equal(t, "configurations", FieldConfigurations)
	assert.Equal(t, "orgDevices", FieldOrgDevices)
	assert.Equal(t, "users", FieldUsers)
	assert.Equal(t, "userGroups", FieldUserGroups)
}

func TestBlueprintIncludeConstants(t *testing.T) {
	assert.Equal(t, "apps", IncludeApps)
	assert.Equal(t, "packages", IncludePackages)
	assert.Equal(t, "configurations", IncludeConfigurations)
	assert.Equal(t, "orgDevices", IncludeOrgDevices)
	assert.Equal(t, "users", IncludeUsers)
	assert.Equal(t, "userGroups", IncludeUserGroups)
}

func TestBlueprintStatusConstants(t *testing.T) {
	assert.Equal(t, "ACTIVE", BlueprintStatusActive)
	assert.Equal(t, "INACTIVE", BlueprintStatusInactive)
}
