package devicemanagement

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// DeviceManagement handles communication with the device management
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessmanagerapi/
type (
	DeviceManagement struct {
		client client.Client
	}
)

// NewService creates a new device management service.
func NewService(c client.Client) *DeviceManagement {
	return &DeviceManagement{client: c}
}

// GetV1 retrieves a list of device management services (MDM servers) in an organization.
// URL: GET https://api-business.apple.com/v1/mdmServers
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-device-management-services
func (s *DeviceManagement) GetV1(ctx context.Context, opts *RequestQueryOptions) (*ResponseMDMServers, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		params.AddInt("limit", opts.Limit)
	}

	var allServers []MDMServer
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointMDMServers, func(pageData []byte) error {
			var pageResponse ResponseMDMServers
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allServers = append(allServers, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &ResponseMDMServers{
		Data:  allServers,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByMDMServerIDV1 retrieves information about a specific device management service.
// URL: GET https://api-business.apple.com/v1/mdmServers/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-device-management-service-information
func (s *DeviceManagement) GetByMDMServerIDV1(ctx context.Context, serverID string, opts *RequestQueryOptions) (*MDMServerResponse, *resty.Response, error) {
	if serverID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointMDMServers+"/%s", serverID)

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	var result MDMServerResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// CreateMDMServerV1 creates a new device management service in an organization.
// URL: POST https://api-business.apple.com/v1/mdmServers
// https://developer.apple.com/documentation/applebusinessmanagerapi/create-a-device-management-service
// Note: serverName and serverCertificate are required.
func (s *DeviceManagement) CreateMDMServerV1(ctx context.Context, req *MDMServerCreateRequest) (*MDMServerResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	if req.Data.Attributes.ServerName == "" {
		return nil, nil, fmt.Errorf("serverName is required")
	}
	if req.Data.Attributes.ServerCertificate.Name == "" {
		return nil, nil, fmt.Errorf("serverCertificate.name is required")
	}
	if req.Data.Attributes.ServerCertificate.Data == "" {
		return nil, nil, fmt.Errorf("serverCertificate.data is required")
	}

	var result MDMServerResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		SetResult(&result).
		Post(constants.EndpointMDMServers)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UpdateMDMServerByIDV1 updates an existing device management service in an organization.
// URL: PATCH https://api-business.apple.com/v1/mdmServers/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/update-a-device-management-service
func (s *DeviceManagement) UpdateMDMServerByIDV1(ctx context.Context, serverID string, req *MDMServerUpdateRequest) (*MDMServerResponse, *resty.Response, error) {
	if serverID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointMDMServers+"/%s", serverID)

	var result MDMServerResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		SetResult(&result).
		Patch(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DeleteMDMServerByIDV1 deletes a device management service from an organization.
// URL: DELETE https://api-business.apple.com/v1/mdmServers/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/delete-a-device-management-service
// Note: A server with devices assigned cannot be deleted. Returns 204 No Content on success.
func (s *DeviceManagement) DeleteMDMServerByIDV1(ctx context.Context, serverID string) (*resty.Response, error) {
	if serverID == "" {
		return nil, fmt.Errorf("MDM server ID is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointMDMServers+"/%s", serverID)

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetDeviceSerialNumbersByServerIDV1 retrieves a list of device IDs assigned to an MDM server.
// URL: GET https://api-business.apple.com/v1/mdmServers/{id}/relationships/devices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-device-ids-for-a-device-management-service
func (s *DeviceManagement) GetDeviceSerialNumbersByServerIDV1(ctx context.Context, mdmServerID string, opts *RequestQueryOptions) (*ResponseMDMServerDevicesLinkages, *resty.Response, error) {
	if mdmServerID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointMDMServers+"/%s/relationships/devices", mdmServerID)

	params := s.client.QueryBuilder()

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		params.AddInt("limit", opts.Limit)
	}

	var allDevices []MDMServerDeviceLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse ResponseMDMServerDevicesLinkages
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allDevices = append(allDevices, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &ResponseMDMServerDevicesLinkages{
		Data:  allDevices,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetAssignedServerIDByDeviceIDV1 retrieves the assigned device management service ID linkage for a device.
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/relationships/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-id-for-an-orgdevice
func (s *DeviceManagement) GetAssignedServerIDByDeviceIDV1(ctx context.Context, deviceID string) (*ResponseOrgDeviceAssignedServerLinkage, *resty.Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device ID is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointOrgDevices+"/%s/relationships/assignedServer", deviceID)

	var result ResponseOrgDeviceAssignedServerLinkage

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetAssignedServerInfoByDeviceIDV1 retrieves the assigned device management service information for a device.
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/assignedServer
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-the-assigned-device-management-service-information-for-an-orgdevice
func (s *DeviceManagement) GetAssignedServerInfoByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*MDMServerResponse, *resty.Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointOrgDevices+"/%s/assignedServer", deviceID)

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[mdmServers]", opts.Fields)
	}

	var result MDMServerResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// deviceLinkages converts a slice of device IDs into orgDevices resource linkages.
func deviceLinkages(deviceIDs []string) []OrgDeviceActivityDeviceLinkage {
	linkages := make([]OrgDeviceActivityDeviceLinkage, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		linkages = append(linkages, OrgDeviceActivityDeviceLinkage{
			Type: ResourceTypeOrgDevices,
			ID:   deviceID,
		})
	}
	return linkages
}

// CreateOrgDeviceActivityV1 creates an organization device activity: assign or unassign
// devices to a device management service, schedule, update or cancel a device management
// service migration, or release devices from the organization.
//
// The activity is asynchronous. The response describes an IN_PROGRESS activity; poll
// GetOrgDeviceActivityByIDV1 with the returned ID for the outcome.
//
// Activity-type semantics are enforced by the API, not by this method: mdmServer is
// required for ASSIGN_DEVICES, UNASSIGN_DEVICES and
// ASSIGN_DEVICES_WITH_MDM_MIGRATION_DEADLINE, activityTypeMetadata is required for the
// two deadline activity types, and a migration deadline cannot be more than 90 days in
// the future. The API returns 409 Conflict when any of these is violated.
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessapi/create-an-orgdeviceactivity
func (s *DeviceManagement) CreateOrgDeviceActivityV1(ctx context.Context, req *OrgDeviceActivityCreateRequest) (*ResponseOrgDeviceActivity, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	if req.Data.Attributes.ActivityType == "" {
		return nil, nil, fmt.Errorf("activity type is required")
	}
	if req.Data.Relationships.Devices == nil || len(req.Data.Relationships.Devices.Data) == 0 {
		return nil, nil, fmt.Errorf("at least one device ID is required")
	}

	var result ResponseOrgDeviceActivity

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		SetResult(&result).
		Post(constants.EndpointOrgDeviceActivities)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetOrgDeviceActivityByIDV1 retrieves the information for an organization device activity.
//
// The API retains activities for the past 30 days. completedDateTime and downloadUrl are
// present only once the activity reaches a COMPLETED status.
// URL: GET https://api-business.apple.com/v1/orgDeviceActivities/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-orgdeviceactivity-information
func (s *DeviceManagement) GetOrgDeviceActivityByIDV1(ctx context.Context, activityID string, opts *GetOrgDeviceActivityQueryOptions) (*ResponseOrgDeviceActivity, *resty.Response, error) {
	if activityID == "" {
		return nil, nil, fmt.Errorf("activity ID is required")
	}
	if opts == nil {
		opts = &GetOrgDeviceActivityQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointOrgDeviceActivities+"/%s", activityID)

	params := s.client.QueryBuilder()
	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[orgDeviceActivities]", opts.Fields)
	}

	var result ResponseOrgDeviceActivity

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// AssignDevicesV1 assigns devices to an MDM server.
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessapi/create-an-orgdeviceactivity
func (s *DeviceManagement) AssignDevicesV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, *resty.Response, error) {
	if mdmServerID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}

	return s.CreateOrgDeviceActivityV1(ctx, &OrgDeviceActivityCreateRequest{
		Data: OrgDeviceActivityData{
			Type: ResourceTypeOrgDeviceActivities,
			Attributes: OrgDeviceActivityCreateAttributes{
				ActivityType: ActivityTypeAssignDevices,
			},
			Relationships: OrgDeviceActivityCreateRelationships{
				MDMServer: &OrgDeviceActivityMDMServerRelationship{
					Data: OrgDeviceActivityMDMServerLinkage{
						Type: ResourceTypeMDMServers,
						ID:   mdmServerID,
					},
				},
				Devices: &OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages(deviceIDs),
				},
			},
		},
	})
}

// UnassignDevicesV1 unassigns devices from an MDM server.
// URL: POST https://api-business.apple.com/v1/orgDeviceActivities
// https://developer.apple.com/documentation/applebusinessapi/create-an-orgdeviceactivity
func (s *DeviceManagement) UnassignDevicesV1(ctx context.Context, mdmServerID string, deviceIDs []string) (*ResponseOrgDeviceActivity, *resty.Response, error) {
	if mdmServerID == "" {
		return nil, nil, fmt.Errorf("MDM server ID is required")
	}

	return s.CreateOrgDeviceActivityV1(ctx, &OrgDeviceActivityCreateRequest{
		Data: OrgDeviceActivityData{
			Type: ResourceTypeOrgDeviceActivities,
			Attributes: OrgDeviceActivityCreateAttributes{
				ActivityType: ActivityTypeUnassignDevices,
			},
			Relationships: OrgDeviceActivityCreateRelationships{
				MDMServer: &OrgDeviceActivityMDMServerRelationship{
					Data: OrgDeviceActivityMDMServerLinkage{
						Type: ResourceTypeMDMServers,
						ID:   mdmServerID,
					},
				},
				Devices: &OrgDeviceActivityDevicesRelationship{
					Data: deviceLinkages(deviceIDs),
				},
			},
		},
	})
}
