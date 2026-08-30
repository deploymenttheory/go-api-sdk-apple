package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	orgDevices map[string]map[string]any
}

func init() {
	mockState.orgDevices = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

type OrgDevicesMock struct{}

// RegisterMocks registers all the HTTP mock responders for organization devices
func (m *OrgDevicesMock) RegisterMocks() {
	mockState.Lock()
	mockState.orgDevices = make(map[string]map[string]any)
	mockState.Unlock()

	// Seed with default test device
	m.seedTestDevice()

	// GET /orgDevices - List organization devices
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/orgDevices", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_organization_devices.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// If we have devices in mock state, return them
		if len(mockState.orgDevices) > 0 {
			deviceList := make([]map[string]any, 0, len(mockState.orgDevices))
			for _, device := range mockState.orgDevices {
				deviceCopy := make(map[string]any)
				for k, v := range device {
					deviceCopy[k] = v
				}
				deviceList = append(deviceList, deviceCopy)
			}
			responseObj["data"] = deviceList
		}

		// Handle query parameters
		query := req.URL.Query()

		// Handle fields parameter
		if fields := query.Get("fields[orgDevices]"); fields != "" {
			// In a real implementation, we would filter fields here
			// For mock purposes, we'll just validate the parameter exists
		}

		// Handle limit parameter
		if limit := query.Get("limit"); limit != "" {
			// In a real implementation, we would apply pagination here
			// For mock purposes, we'll just validate the parameter exists
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /orgDevices/{id} - Get specific device information
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		deviceID := parts[len(parts)-1]

		mockState.Lock()
		device, exists := mockState.orgDevices[deviceID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
		}

		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_device_information.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// Override template values with actual device values
		if data, ok := responseObj["data"].(map[string]any); ok {
			for k, v := range device {
				data[k] = v
			}
		}

		// Handle query parameters
		query := req.URL.Query()

		// Handle fields parameter
		if fields := query.Get("fields[orgDevices]"); fields != "" {
			// In a real implementation, we would filter fields here
			// For mock purposes, we'll just validate the parameter exists
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /orgDevices/{id}/appleCareCoverage - Get AppleCare coverage for device
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/[^/]+/appleCareCoverage$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		deviceID := parts[len(parts)-2]

		mockState.Lock()
		_, exists := mockState.orgDevices[deviceID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
		}

		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_applecare_coverage.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// Handle query parameters
		query := req.URL.Query()

		// Handle fields parameter
		if fields := query.Get("fields[appleCareCoverage]"); fields != "" {
			// In a real implementation, we would filter fields here
			// For mock purposes, we'll just validate the parameter exists
		}

		// Handle limit parameter
		if limit := query.Get("limit"); limit != "" {
			// In a real implementation, we would apply pagination here
			// For mock purposes, we'll just validate the parameter exists
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

// RegisterErrorMocks registers mock responders that return error responses
func (m *OrgDevicesMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.orgDevices = make(map[string]map[string]any)
	mockState.Unlock()

	// GET /orgDevices - Return error
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/orgDevices", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	// GET /orgDevices/{id} - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
	})

	// GET /orgDevices/{id}/appleCareCoverage - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/[^/]+/appleCareCoverage$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
	})
}

// CleanupMockState clears all mock state data
func (m *OrgDevicesMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.orgDevices {
		delete(mockState.orgDevices, id)
	}
}

// AddMockDevice adds a device to the mock state
func (m *OrgDevicesMock) AddMockDevice(deviceID string, device map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.orgDevices[deviceID] = device
}

// GetMockDevice retrieves a device from the mock state
func (m *OrgDevicesMock) GetMockDevice(deviceID string) (map[string]any, bool) {
	mockState.Lock()
	defer mockState.Unlock()
	device, exists := mockState.orgDevices[deviceID]
	return device, exists
}

// seedTestDevice adds a default test device to the mock state
func (m *OrgDevicesMock) seedTestDevice() {
	testDevice := map[string]any{
		"type": "orgDevices",
		"id":   "XABC123X0ABC123X0",
		"attributes": map[string]any{
			"serialNumber":       "XABC123X0ABC123X0",
			"addedToOrgDateTime": "2025-04-30T22:05:14.192Z",
			"updatedDateTime":    "2025-05-01T15:33:54.164Z",
			"deviceModel":        "iMac 21.5\"",
			"productFamily":      "Mac",
			"productType":        "iMac16,2",
			"deviceCapacity":     "750GB",
			"partNumber":         "FD311LL/A",
			"orderNumber":        "1234567890",
			"color":              "SILVER",
			"status":             "UNASSIGNED",
			"orderDateTime":      "2011-08-15T07:00:00Z",
			"imei":               []string{"123456789012345", "123456789012346"},
			"meid":               []string{"12345678901237"},
			"eid":                "89049037640158663184237812557346",
			"purchaseSourceId":   "-2085650007946880",
			"purchaseSourceType": "APPLE",

			// Device management service migration attributes (ASM 1.6 / ABM 2.3)
			"isMdmMigrationCapable":        true,
			"mdmMigrationStatus":           "REQUESTED",
			"mdmMigrationDeadlineDateTime": "2026-03-15T17:00:00.000Z",
		},
		"relationships": map[string]any{
			"assignedServer": map[string]any{
				"links": map[string]any{
					"self":    "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0/relationships/assignedServer",
					"related": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0/assignedServer",
				},
			},
		},
		"links": map[string]any{
			"self": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0",
		},
	}

	m.AddMockDevice("XABC123X0ABC123X0", testDevice)
}
