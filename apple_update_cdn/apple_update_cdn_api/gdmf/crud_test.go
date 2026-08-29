package gdmf

import (
	"context"
	"net/http"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a GDMF transport with httpmock enabled.
func setupMockClient(t *testing.T) *GDMFService {
	t.Helper()

	transport, err := client.NewTransport(
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewService(transport)
}

// jsonResponder wraps a body string in an HTTP JSON response.
func jsonResponder(status int, body string) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(status, body)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	}
}

// =============================================================================
// GetPublicVersionsV2
// =============================================================================

const gdmfResponseJSON = `{
  "PublicAssetSets": {
    "iOS": [
      {
        "ProductVersion": "18.3",
        "Build": "22D60",
        "PostingDate": "2025-01-27",
        "ExpirationDate": "2025-05-01",
        "SupportedDevices": ["iPhone15,2", "iPhone15,3"]
      }
    ],
    "macOS": [
      {
        "ProductVersion": "15.3",
        "Build": "24D60",
        "PostingDate": "2025-01-27",
        "ExpirationDate": "2025-05-01",
        "SupportedDevices": ["J473AP", "J316cAP", "Mac-1E7E29AD0135F9BC"]
      }
    ],
    "visionOS": [
      {
        "ProductVersion": "2.3",
        "Build": "22N330",
        "PostingDate": "2025-01-27",
        "ExpirationDate": "2025-05-01",
        "SupportedDevices": ["RealityDevice14,1"]
      }
    ]
  },
  "AssetSets": {
    "iOS": [],
    "macOS": [
      {
        "ProductVersion": "15.3.1",
        "Build": "24D70",
        "PostingDate": "2025-02-05",
        "ExpirationDate": "2025-05-15",
        "SupportedDevices": ["J473AP"]
      }
    ],
    "visionOS": []
  },
  "PublicBackgroundSecurityImprovements": {
    "iOS": [],
    "macOS": [],
    "visionOS": []
  }
}`

func TestGetPublicVersionsV2_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://gdmf.apple.com/v2/pmv",
		jsonResponder(200, gdmfResponseJSON))

	result, resp, err := svc.GetPublicVersionsV2(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
}

func TestGetPublicVersionsV2_PublicAssetSets(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://gdmf.apple.com/v2/pmv",
		jsonResponder(200, gdmfResponseJSON))

	result, _, err := svc.GetPublicVersionsV2(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result.PublicAssetSets)

	macos := result.PublicAssetSets.MacOS
	require.Len(t, macos, 1)
	assert.Equal(t, "15.3", macos[0].ProductVersion)
	assert.Equal(t, "24D60", macos[0].Build)
	assert.Equal(t, "2025-01-27", macos[0].PostingDate)
	assert.Contains(t, macos[0].SupportedDevices, "J473AP")

	ios := result.PublicAssetSets.IOS
	require.Len(t, ios, 1)
	assert.Equal(t, "18.3", ios[0].ProductVersion)

	visionos := result.PublicAssetSets.VisionOS
	require.Len(t, visionos, 1)
	assert.Equal(t, "2.3", visionos[0].ProductVersion)
}

func TestGetPublicVersionsV2_AssetSets(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://gdmf.apple.com/v2/pmv",
		jsonResponder(200, gdmfResponseJSON))

	result, _, err := svc.GetPublicVersionsV2(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result.AssetSets)
	require.Len(t, result.AssetSets.MacOS, 1)
	assert.Equal(t, "15.3.1", result.AssetSets.MacOS[0].ProductVersion)
}

func TestGetPublicVersionsV2_PublicBackgroundSecurityImprovements(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://gdmf.apple.com/v2/pmv",
		jsonResponder(200, gdmfResponseJSON))

	result, _, err := svc.GetPublicVersionsV2(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result.PublicBackgroundSecurityImprovements)
	assert.Empty(t, result.PublicBackgroundSecurityImprovements.MacOS)
}

func TestGetPublicVersionsV2_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://gdmf.apple.com/v2/pmv",
		httpmock.NewStringResponder(503, "Service Unavailable"))

	_, resp, err := svc.GetPublicVersionsV2(context.Background())

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, err.Error(), "503")
}

func TestGetPublicVersionsV2_PlatformConstants(t *testing.T) {
	assert.Equal(t, "macOS", PlatformMacOS)
	assert.Equal(t, "iOS", PlatformIOS)
	assert.Equal(t, "visionOS", PlatformVisionOS)
}
