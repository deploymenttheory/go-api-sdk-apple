package firmware

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

// setupMockClient creates a firmware transport with httpmock enabled.
func setupMockClient(t *testing.T) *FirmwareService {
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
// ListAllMacFirmwareV3
// =============================================================================

const listAllMacFirmwareJSON = `{
  "devices": {
    "Mac14,3": {
      "name": "Mac mini (M2, 2023)",
      "BoardConfig": "J473AP",
      "platform": "REALBRIDGE",
      "cpid": 35200,
      "bdid": 12,
      "firmwares": [
        {
          "version": "15.3",
          "buildid": "24D60",
          "sha1sum": "abc123def456abc123def456abc123def456abc1",
          "md5sum": "fedcba9876543210fedcba9876543210",
          "size": 19734779897,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/122-12345/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE/UniversalMac_15.3_24D60_Restore.ipsw",
          "signed": true,
          "filename": "UniversalMac_15.3_24D60_Restore.ipsw"
        }
      ]
    },
    "MacBookPro18,1": {
      "name": "MacBook Pro (16-inch, 2021)",
      "BoardConfig": "J316cAP",
      "platform": "REALBRIDGE",
      "cpid": 35168,
      "bdid": 0,
      "firmwares": [
        {
          "version": "15.3",
          "buildid": "24D60",
          "sha1sum": "abc123def456abc123def456abc123def456abc1",
          "md5sum": "fedcba9876543210fedcba9876543210",
          "size": 19734779897,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/122-12345/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE/UniversalMac_15.3_24D60_Restore.ipsw",
          "signed": true,
          "filename": "UniversalMac_15.3_24D60_Restore.ipsw"
        }
      ]
    },
    "iPod9,1": {
      "name": "iPod touch (7th generation)",
      "BoardConfig": "N112AP",
      "platform": "S8000",
      "cpid": 32800,
      "bdid": 14,
      "firmwares": []
    }
  }
}`

func TestListAllMacFirmwareV3_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, listAllMacFirmwareJSON))

	result, resp, err := svc.ListAllMacFirmwareV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Len(t, result.Devices, 2, "iPod should be filtered out, leaving 2 Mac devices")
	assert.Contains(t, result.Devices, "Mac14,3")
	assert.Contains(t, result.Devices, "MacBookPro18,1")
	assert.NotContains(t, result.Devices, "iPod9,1")
}

func TestListAllMacFirmwareV3_FirmwareContents(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, listAllMacFirmwareJSON))

	result, _, err := svc.ListAllMacFirmwareV3(context.Background())

	require.NoError(t, err)
	mac := result.Devices["Mac14,3"]
	require.Len(t, mac.Firmwares, 1)
	fw := mac.Firmwares[0]
	assert.Equal(t, "15.3", fw.Version)
	assert.Equal(t, "24D60", fw.BuildID)
	assert.True(t, fw.Signed)
	assert.Equal(t, int64(19734779897), fw.Size)
}

func TestListAllMacFirmwareV3_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		httpmock.NewStringResponder(503, "Service Unavailable"))

	_, resp, err := svc.ListAllMacFirmwareV3(context.Background())

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, err.Error(), "503")
}

// =============================================================================
// ListUniqueMacFirmwareVersionsV3
// =============================================================================

const multiVersionFirmwareJSON = `{
  "devices": {
    "Mac14,3": {
      "name": "Mac mini (M2, 2023)",
      "BoardConfig": "J473AP",
      "platform": "REALBRIDGE",
      "firmwares": [
        {
          "version": "15.3",
          "buildid": "24D60",
          "sha1sum": "aaa",
          "md5sum": "bbb",
          "size": 100,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/a/UniversalMac_15.3_24D60_Restore.ipsw",
          "signed": true,
          "filename": "UniversalMac_15.3_24D60_Restore.ipsw"
        },
        {
          "version": "15.2",
          "buildid": "24C101",
          "sha1sum": "ccc",
          "md5sum": "ddd",
          "size": 99,
          "releasedate": "2024-12-11T00:00:00Z",
          "uploaddate": "2024-12-09T00:00:00Z",
          "url": "https://updates.cdn-apple.com/b/UniversalMac_15.2_24C101_Restore.ipsw",
          "signed": false,
          "filename": "UniversalMac_15.2_24C101_Restore.ipsw"
        }
      ]
    },
    "MacBookPro18,1": {
      "name": "MacBook Pro (16-inch, 2021)",
      "BoardConfig": "J316cAP",
      "platform": "REALBRIDGE",
      "firmwares": [
        {
          "version": "15.3",
          "buildid": "24D60",
          "sha1sum": "aaa",
          "md5sum": "bbb",
          "size": 100,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/a/UniversalMac_15.3_24D60_Restore.ipsw",
          "signed": true,
          "filename": "UniversalMac_15.3_24D60_Restore.ipsw"
        }
      ]
    }
  }
}`

func TestListUniqueMacFirmwareVersionsV3_Deduplicated(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, multiVersionFirmwareJSON))

	result, resp, err := svc.ListUniqueMacFirmwareVersionsV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	// 24D60 appears in both Mac14,3 and MacBookPro18,1 — should deduplicate to 1 entry.
	// 24C101 appears only in Mac14,3.
	assert.Len(t, result, 2, "should have 2 unique build IDs")
}

func TestListUniqueMacFirmwareVersionsV3_SortedNewestFirst(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, multiVersionFirmwareJSON))

	result, _, err := svc.ListUniqueMacFirmwareVersionsV3(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "15.3", result[0].Version, "newest version should be first")
	assert.Equal(t, "15.2", result[1].Version)
}

func TestListUniqueMacFirmwareVersionsV3_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		httpmock.NewStringResponder(429, "Too Many Requests"))

	_, _, err := svc.ListUniqueMacFirmwareVersionsV3(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "429")
}

// =============================================================================
// GetByDeviceV4
// =============================================================================

const deviceV4JSON = `{
  "name": "Mac mini (M2, 2023)",
  "identifier": "Mac14,3",
  "boardconfig": "J473AP",
  "platform": "REALBRIDGE",
  "cpid": 35200,
  "bdid": 12,
  "boards": ["J473AP"],
  "firmwares": [
    {
      "identifier": "Mac14,3",
      "version": "15.3",
      "buildid": "24D60",
      "sha1sum": "abc123def456abc123def456abc123def456abc1",
      "md5sum": "fedcba9876543210fedcba9876543210",
      "sha256sum": "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
      "filesize": 19734779897,
      "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/122-12345/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE/UniversalMac_15.3_24D60_Restore.ipsw",
      "releasedate": "2025-01-27T00:00:00Z",
      "uploaddate": "2025-01-25T00:00:00Z",
      "signed": true
    }
  ]
}`

func TestGetByDeviceV4_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v4/device/Mac14,3",
		jsonResponder(200, deviceV4JSON))

	result, resp, err := svc.GetByDeviceV4(context.Background(), "Mac14,3")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, "Mac14,3", result.Identifier)
	assert.Equal(t, "Mac mini (M2, 2023)", result.Name)
	require.Len(t, result.Firmwares, 1)
	fw := result.Firmwares[0]
	assert.Equal(t, "15.3", fw.Version)
	assert.Equal(t, "24D60", fw.BuildID)
	assert.Equal(t, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab", fw.SHA256Sum)
	assert.Equal(t, int64(19734779897), fw.FileSize)
	assert.True(t, fw.Signed)
}

func TestGetByDeviceV4_QueryParamTypeIPSW(t *testing.T) {
	svc := setupMockClient(t)

	var capturedType string
	httpmock.RegisterResponder("GET", `=~^https://api\.ipsw\.me/v4/device/`,
		func(req *http.Request) (*http.Response, error) {
			capturedType = req.URL.Query().Get("type")
			resp := httpmock.NewStringResponse(200, deviceV4JSON)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	_, _, err := svc.GetByDeviceV4(context.Background(), "Mac14,3")

	require.NoError(t, err)
	assert.Equal(t, FirmwareTypeIPSW, capturedType, "type=ipsw query param must be sent")
}

func TestGetByDeviceV4_EmptyIdentifierError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.GetByDeviceV4(context.Background(), "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "device identifier is required")
}

func TestGetByDeviceV4_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v4/device/Mac99,99",
		httpmock.NewStringResponder(404, "Not Found"))

	_, resp, err := svc.GetByDeviceV4(context.Background(), "Mac99,99")

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, err.Error(), "404")
}

// =============================================================================
// ListAllFirmwareV3 (unfiltered)
// =============================================================================

const mixedDevicesJSON = `{
  "devices": {
    "Mac14,3": {
      "name": "Mac mini (M2, 2023)",
      "BoardConfig": "J473AP",
      "platform": "REALBRIDGE",
      "firmwares": []
    },
    "iPhone15,2": {
      "name": "iPhone 14 Pro",
      "BoardConfig": "D73AP",
      "platform": "T8120",
      "firmwares": []
    },
    "iPad14,4": {
      "name": "iPad mini (6th generation)",
      "BoardConfig": "J310AP",
      "platform": "T8101",
      "firmwares": []
    },
    "iPod9,1": {
      "name": "iPod touch (7th generation)",
      "BoardConfig": "N112AP",
      "platform": "S8000",
      "firmwares": []
    }
  }
}`

func TestListAllFirmwareV3_ReturnsAllDevices(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, mixedDevicesJSON))

	result, resp, err := svc.ListAllFirmwareV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, result.Devices, 4, "all device types should be present")
	assert.Contains(t, result.Devices, "Mac14,3")
	assert.Contains(t, result.Devices, "iPhone15,2")
	assert.Contains(t, result.Devices, "iPad14,4")
	assert.Contains(t, result.Devices, "iPod9,1")
}

// =============================================================================
// ListAllIOSFirmwareV3
// =============================================================================

const iosMixedJSON = `{
  "devices": {
    "iPhone15,2": {
      "name": "iPhone 14 Pro",
      "BoardConfig": "D73AP",
      "platform": "T8120",
      "firmwares": [
        {
          "version": "18.3",
          "buildid": "22D60",
          "sha1sum": "aaa",
          "md5sum": "bbb",
          "size": 7000000000,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/001-11111/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE/iPhone15,2_18.3_22D60_Restore.ipsw",
          "signed": true,
          "filename": "iPhone15,2_18.3_22D60_Restore.ipsw"
        }
      ]
    },
    "iPhone16,1": {
      "name": "iPhone 15 Pro",
      "BoardConfig": "D83AP",
      "platform": "T8130",
      "firmwares": [
        {
          "version": "18.3",
          "buildid": "22D60",
          "sha1sum": "aaa",
          "md5sum": "bbb",
          "size": 7100000000,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/001-22222/BBBBBBBB-CCCC-DDDD-EEEE-FFFFFFFFFFFF/iPhone16,1_18.3_22D60_Restore.ipsw",
          "signed": true,
          "filename": "iPhone16,1_18.3_22D60_Restore.ipsw"
        }
      ]
    },
    "Mac14,3": {
      "name": "Mac mini (M2, 2023)",
      "BoardConfig": "J473AP",
      "platform": "REALBRIDGE",
      "firmwares": []
    },
    "iPad14,4": {
      "name": "iPad mini (6th generation)",
      "BoardConfig": "J310AP",
      "platform": "T8101",
      "firmwares": []
    }
  }
}`

func TestListAllIOSFirmwareV3_FiltersToiPhone(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, iosMixedJSON))

	result, resp, err := svc.ListAllIOSFirmwareV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, result.Devices, 2, "only iPhone devices should remain")
	assert.Contains(t, result.Devices, "iPhone15,2")
	assert.Contains(t, result.Devices, "iPhone16,1")
	assert.NotContains(t, result.Devices, "Mac14,3")
	assert.NotContains(t, result.Devices, "iPad14,4")
}

func TestListAllIOSFirmwareV3_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		httpmock.NewStringResponder(503, "Service Unavailable"))

	_, _, err := svc.ListAllIOSFirmwareV3(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "503")
}

// =============================================================================
// ListAllIPadOSFirmwareV3
// =============================================================================

const ipadMixedJSON = `{
  "devices": {
    "iPad14,4": {
      "name": "iPad mini (6th generation)",
      "BoardConfig": "J310AP",
      "platform": "T8101",
      "firmwares": [
        {
          "version": "18.3",
          "buildid": "22D60",
          "sha1sum": "ccc",
          "md5sum": "ddd",
          "size": 8000000000,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/002-11111/CCCCCCCC-DDDD-EEEE-FFFF-111111111111/iPad14,4_18.3_22D60_Restore.ipsw",
          "signed": true,
          "filename": "iPad14,4_18.3_22D60_Restore.ipsw"
        }
      ]
    },
    "iPad13,16": {
      "name": "iPad Air (5th generation)",
      "BoardConfig": "J407AP",
      "platform": "T8110",
      "firmwares": [
        {
          "version": "18.3",
          "buildid": "22D60",
          "sha1sum": "eee",
          "md5sum": "fff",
          "size": 8100000000,
          "releasedate": "2025-01-27T00:00:00Z",
          "uploaddate": "2025-01-25T00:00:00Z",
          "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/002-22222/DDDDDDDD-EEEE-FFFF-0000-222222222222/iPad13,16_18.3_22D60_Restore.ipsw",
          "signed": true,
          "filename": "iPad13,16_18.3_22D60_Restore.ipsw"
        }
      ]
    },
    "iPhone15,2": {
      "name": "iPhone 14 Pro",
      "BoardConfig": "D73AP",
      "platform": "T8120",
      "firmwares": []
    },
    "Mac14,3": {
      "name": "Mac mini (M2, 2023)",
      "BoardConfig": "J473AP",
      "platform": "REALBRIDGE",
      "firmwares": []
    }
  }
}`

func TestListAllIPadOSFirmwareV3_FiltersToiPad(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, ipadMixedJSON))

	result, resp, err := svc.ListAllIPadOSFirmwareV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, result.Devices, 2, "only iPad devices should remain")
	assert.Contains(t, result.Devices, "iPad14,4")
	assert.Contains(t, result.Devices, "iPad13,16")
	assert.NotContains(t, result.Devices, "iPhone15,2")
	assert.NotContains(t, result.Devices, "Mac14,3")
}

func TestListAllIPadOSFirmwareV3_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		httpmock.NewStringResponder(503, "Service Unavailable"))

	_, _, err := svc.ListAllIPadOSFirmwareV3(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "503")
}

// =============================================================================
// ListUniqueIOSFirmwareVersionsV3
// =============================================================================

func TestListUniqueIOSFirmwareVersionsV3_Deduplicated(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, iosMixedJSON))

	result, resp, err := svc.ListUniqueIOSFirmwareVersionsV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	// Both iPhone15,2 and iPhone16,1 share build 22D60 — should deduplicate to 1 entry.
	assert.Len(t, result, 1, "should have 1 unique build ID")
	assert.Equal(t, "18.3", result[0].Version)
}

// =============================================================================
// ListUniqueIPadOSFirmwareVersionsV3
// =============================================================================

func TestListUniqueIPadOSFirmwareVersionsV3_Deduplicated(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v3/firmwares.json/condensed",
		jsonResponder(200, ipadMixedJSON))

	result, resp, err := svc.ListUniqueIPadOSFirmwareVersionsV3(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	// Both iPad14,4 and iPad13,16 share build 22D60 — should deduplicate to 1 entry.
	assert.Len(t, result, 1, "should have 1 unique build ID")
	assert.Equal(t, "18.3", result[0].Version)
}

// =============================================================================
// GetByDeviceV4 — iPhone/iPad
// =============================================================================

const iPhoneDeviceV4JSON = `{
  "name": "iPhone 14 Pro",
  "identifier": "iPhone15,2",
  "boardconfig": "D73AP",
  "platform": "T8120",
  "cpid": 33040,
  "bdid": 6,
  "boards": ["D73AP"],
  "firmwares": [
    {
      "identifier": "iPhone15,2",
      "version": "18.3",
      "buildid": "22D60",
      "sha1sum": "aaabbbccc",
      "md5sum": "dddeeefff",
      "sha256sum": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12",
      "filesize": 7000000000,
      "url": "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/001-11111/AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE/iPhone15,2_18.3_22D60_Restore.ipsw",
      "releasedate": "2025-01-27T00:00:00Z",
      "uploaddate": "2025-01-25T00:00:00Z",
      "signed": true
    }
  ]
}`

func TestGetByDeviceV4_iPhone(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://api.ipsw.me/v4/device/iPhone15,2",
		jsonResponder(200, iPhoneDeviceV4JSON))

	result, resp, err := svc.GetByDeviceV4(context.Background(), "iPhone15,2")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "iPhone15,2", result.Identifier)
	assert.Equal(t, "iPhone 14 Pro", result.Name)
	require.Len(t, result.Firmwares, 1)
	assert.Equal(t, "18.3", result.Firmwares[0].Version)
	assert.Equal(t, "22D60", result.Firmwares[0].BuildID)
}

// =============================================================================
// Identifier helper functions
// =============================================================================

func TestIsMacIdentifier(t *testing.T) {
	cases := []struct {
		identifier string
		want       bool
	}{
		{"Mac14,3", true},
		{"MacBook8,1", true},
		{"MacBookPro18,1", true},
		{"MacBookAir10,1", true},
		{"iMac21,1", true},
		{"MacPro7,1", true},
		{"Macmini9,1", true},
		{"VirtualMac2,1", true},
		{"iPhone15,2", false},
		{"iPod9,1", false},
		{"iPad14,4", false},
		{"AppleTV6,2", false},
		{"", false},
	}

	for _, tc := range cases {
		got := isMacIdentifier(tc.identifier)
		assert.Equal(t, tc.want, got, "isMacIdentifier(%q)", tc.identifier)
	}
}

func TestIsIOSIdentifier(t *testing.T) {
	cases := []struct {
		identifier string
		want       bool
	}{
		{"iPhone15,2", true},
		{"iPhone16,1", true},
		{"iPhone14,3", true},
		{"iPad14,4", false},
		{"Mac14,3", false},
		{"iPod9,1", false},
		{"", false},
	}

	for _, tc := range cases {
		got := isIOSIdentifier(tc.identifier)
		assert.Equal(t, tc.want, got, "isIOSIdentifier(%q)", tc.identifier)
	}
}

func TestIsIPadIdentifier(t *testing.T) {
	cases := []struct {
		identifier string
		want       bool
	}{
		{"iPad14,4", true},
		{"iPad13,16", true},
		{"iPad8,1", true},
		{"iPhone15,2", false},
		{"Mac14,3", false},
		{"iPod9,1", false},
		{"", false},
	}

	for _, tc := range cases {
		got := isIPadIdentifier(tc.identifier)
		assert.Equal(t, tc.want, got, "isIPadIdentifier(%q)", tc.identifier)
	}
}
