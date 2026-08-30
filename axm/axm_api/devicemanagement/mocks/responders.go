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

// splitPath splits a URL path into non-empty segments.
func splitPath(path string) []string {
	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

var mockState struct {
	sync.Mutex
	mdmServers       map[string]map[string]any
	deviceActivities map[string]map[string]any
}

func init() {
	mockState.mdmServers = make(map[string]map[string]any)
	mockState.deviceActivities = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`))
}

// loadMockResponse loads JSON response from the mocks folder
func loadMockResponse(filename string) ([]byte, error) {
	mockPath := filepath.Join("mocks", filename)
	return os.ReadFile(mockPath)
}

type DeviceManagementMock struct{}

// RegisterMocks registers all the HTTP mock responders for device management services
func (m *DeviceManagementMock) RegisterMocks() {
	mockState.Lock()
	mockState.mdmServers = make(map[string]map[string]any)
	mockState.deviceActivities = make(map[string]map[string]any)
	mockState.Unlock()

	// Seed with default test MDM server
	m.seedTestMDMServer()

	// Seed with a COMPLETED org device activity
	m.seedTestOrgDeviceActivity()

	// GET /mdmServers - List device management services
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/mdmServers", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_device_management_services.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// If we have MDM servers in mock state, return them
		if len(mockState.mdmServers) > 0 {
			serverList := make([]map[string]any, 0, len(mockState.mdmServers))
			for _, server := range mockState.mdmServers {
				serverCopy := make(map[string]any)
				for k, v := range server {
					serverCopy[k] = v
				}
				serverList = append(serverList, serverCopy)
			}
			responseObj["data"] = serverList
		}

		// Handle query parameters
		query := req.URL.Query()

		// Handle fields parameter
		if fields := query.Get("fields[mdmServers]"); fields != "" {
			// In a real implementation, we would filter fields here
		}

		// Handle limit parameter
		if limit := query.Get("limit"); limit != "" {
			// In a real implementation, we would apply pagination here
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /mdmServers/{id} - Get a specific MDM server by ID
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/mdmServers/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := splitPath(req.URL.Path)
		serverID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.mdmServers[serverID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_get_mdm_server.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /mdmServers - Create a new MDM server
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/mdmServers", func(req *http.Request) (*http.Response, error) {
		mockData, err := loadMockResponse("validate_create_mdm_server_response.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// PATCH /mdmServers/{id} - Update an MDM server
	httpmock.RegisterResponder("PATCH", `=~^https://api-business\.apple\.com/v1/mdmServers/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := splitPath(req.URL.Path)
		serverID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.mdmServers[serverID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
		}

		mockData, err := loadMockResponse("validate_update_mdm_server_response.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// DELETE /mdmServers/{id} - Delete an MDM server
	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/mdmServers/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := splitPath(req.URL.Path)
		serverID := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.mdmServers[serverID]
		if exists {
			delete(mockState.mdmServers, serverID)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// GET /mdmServers/{id}/relationships/devices - Get MDM server device linkages
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices", func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_mdm_server_device_linkages.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /orgDevices/{id}/relationships/assignedServer - Get assigned server linkage
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/relationships/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_assigned_server_linkage.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /orgDevices/{id}/assignedServer - Get assigned server information
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		// Load the mock response template
		mockData, err := loadMockResponse("validate_get_assigned_server_info.json")
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
		if fields := query.Get("fields[mdmServers]"); fields != "" {
			// In a real implementation, we would filter fields here
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /orgDeviceActivities - Assign/Unassign devices
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/orgDeviceActivities", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Invalid request body"}]}`), nil
		}

		// Extract activity type from request
		data, ok := requestBody["data"].(map[string]any)
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Invalid request structure"}]}`), nil
		}

		attributes, ok := data["attributes"].(map[string]any)
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Missing attributes"}]}`), nil
		}

		activityType, ok := attributes["activityType"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Missing activityType"}]}`), nil
		}

		// Choose appropriate response based on activity type
		mockFile, ok := activityResponseFixtures[activityType]
		if !ok {
			return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Invalid activityType"}]}`), nil
		}

		// Load the mock response template
		mockData, err := loadMockResponse(mockFile)
		if err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to load mock data"}]}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(mockData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Failed to parse mock data"}]}`), nil
		}

		// Record the activity so GET /orgDeviceActivities/{id} can read it back
		if activity, ok := responseObj["data"].(map[string]any); ok {
			if id, ok := activity["id"].(string); ok {
				mockState.Lock()
				mockState.deviceActivities[id] = activity
				mockState.Unlock()
			}
		}

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// GET /orgDeviceActivities/{id} - Get organization device activity information
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDeviceActivities/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := splitPath(req.URL.Path)
		activityID := parts[len(parts)-1]

		mockState.Lock()
		activity, exists := mockState.deviceActivities[activityID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Resource Not Found","detail":"The requested resource was not found"}]}`), nil
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"data": activity,
			"links": map[string]any{
				"self": "https://api-business.apple.com/v1/orgDeviceActivities/" + activityID,
			},
		})
	})
}

// activityResponseFixtures maps an orgDeviceActivities activityType to the fixture
// used for its create response.
var activityResponseFixtures = map[string]string{
	"ASSIGN_DEVICES":   "validate_assign_devices_response.json",
	"UNASSIGN_DEVICES": "validate_unassign_devices_response.json",
	"ASSIGN_DEVICES_WITH_MDM_MIGRATION_DEADLINE": "validate_create_mdm_migration_activity_response.json",
	"UPDATE_MDM_MIGRATION_DEADLINE":              "validate_create_mdm_migration_activity_response.json",
	"CANCEL_MDM_MIGRATION":                       "validate_create_mdm_migration_activity_response.json",
	"RELEASE_DEVICES":                            "validate_release_devices_response.json",
}

// SeededOrgDeviceActivityID is the ID of the COMPLETED activity seeded by RegisterMocks,
// for tests that need to read back a finished activity.
const SeededOrgDeviceActivityID = "a0000000-1111-2222-3333-444444444444"

// seedTestOrgDeviceActivity seeds a COMPLETED org device activity into the mock state.
func (m *DeviceManagementMock) seedTestOrgDeviceActivity() {
	mockData, err := loadMockResponse("validate_get_org_device_activity_information.json")
	if err != nil {
		return
	}

	var responseObj map[string]any
	if err := json.Unmarshal(mockData, &responseObj); err != nil {
		return
	}

	activity, ok := responseObj["data"].(map[string]any)
	if !ok {
		return
	}

	mockState.Lock()
	defer mockState.Unlock()
	mockState.deviceActivities[SeededOrgDeviceActivityID] = activity
}

// RegisterErrorMocks registers mock responders that return error responses
func (m *DeviceManagementMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.mdmServers = make(map[string]map[string]any)
	mockState.deviceActivities = make(map[string]map[string]any)
	mockState.Unlock()

	// GET /mdmServers - Return error
	httpmock.RegisterResponder("GET", "https://api-business.apple.com/v1/mdmServers", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(500, `{"errors":[{"status":"500","code":"INTERNAL_ERROR","title":"Internal Server Error","detail":"Mock error for testing"}]}`), nil
	})

	// GET /mdmServers/{id} - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/mdmServers/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
	})

	// POST /mdmServers - Return error
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/mdmServers", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Mock error for testing"}]}`), nil
	})

	// PATCH /mdmServers/{id} - Return not found error
	httpmock.RegisterResponder("PATCH", `=~^https://api-business\.apple\.com/v1/mdmServers/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
	})

	// DELETE /mdmServers/{id} - Return not found error
	httpmock.RegisterResponder("DELETE", `=~^https://api-business\.apple\.com/v1/mdmServers/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
	})

	// GET /mdmServers/{id}/relationships/devices - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/mdmServers/.+/relationships/devices$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"MDM Server Not Found","detail":"The requested MDM server was not found"}]}`), nil
	})

	// GET /orgDevices/{id}/relationships/assignedServer - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/relationships/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
	})

	// GET /orgDevices/{id}/assignedServer - Return not found error
	httpmock.RegisterResponder("GET", `=~^https://api-business\.apple\.com/v1/orgDevices/.+/assignedServer$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"errors":[{"status":"404","code":"RESOURCE_NOT_FOUND","title":"Device Not Found","detail":"The requested device was not found"}]}`), nil
	})

	// POST /orgDeviceActivities - Return error
	httpmock.RegisterResponder("POST", "https://api-business.apple.com/v1/orgDeviceActivities", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"errors":[{"status":"400","code":"BAD_REQUEST","title":"Bad Request","detail":"Mock error for testing"}]}`), nil
	})
}

// CleanupMockState clears all mock state data
func (m *DeviceManagementMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.mdmServers {
		delete(mockState.mdmServers, id)
	}
	for id := range mockState.deviceActivities {
		delete(mockState.deviceActivities, id)
	}
}

// AddMockMDMServer adds an MDM server to the mock state
func (m *DeviceManagementMock) AddMockMDMServer(serverID string, server map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.mdmServers[serverID] = server
}

// GetMockMDMServer retrieves an MDM server from the mock state
func (m *DeviceManagementMock) GetMockMDMServer(serverID string) (map[string]any, bool) {
	mockState.Lock()
	defer mockState.Unlock()
	server, exists := mockState.mdmServers[serverID]
	return server, exists
}

// seedTestMDMServer adds a default test MDM server to the mock state
func (m *DeviceManagementMock) seedTestMDMServer() {
	testServer := map[string]any{
		"type": "mdmServers",
		"id":   "1F97349736CF4614A94F624E705841AD",
		"attributes": map[string]any{
			"serverName":      "Test Device Management Service",
			"serverType":      "MDM",
			"createdDateTime": "2025-05-01T03:21:44.685Z",
			"updatedDateTime": "2025-05-01T03:21:46.284Z",
		},
		"relationships": map[string]any{
			"devices": map[string]any{
				"links": map[string]any{
					"self": "https://api-business.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices",
				},
			},
		},
	}

	m.AddMockMDMServer("1F97349736CF4614A94F624E705841AD", testServer)
}
