package blueprints

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Blueprints handles communication with the blueprints
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/blueprints
type (
	Blueprints struct {
		client client.Client
	}
)

// NewService creates a new blueprints service.
func NewService(c client.Client) *Blueprints {
	return &Blueprints{client: c}
}

// CreateV1 creates a new Blueprint in an organization.
// URL: POST https://api-business.apple.com/v1/blueprints
// https://developer.apple.com/documentation/applebusinessapi/create-a-blueprint
//
// name is required. Any invalid resource IDs in relationships are silently dropped by the API.
func (s *Blueprints) CreateV1(ctx context.Context, req *BlueprintCreateRequest) (*BlueprintResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	if req.Data.Attributes.Name == "" {
		return nil, nil, fmt.Errorf("blueprint name is required")
	}

	var result BlueprintResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		SetResult(&result).
		Post(constants.EndpointBlueprints)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetByBlueprintIDV1 retrieves information about a specific Blueprint in an organization.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-blueprint-information
func (s *Blueprints) GetByBlueprintIDV1(ctx context.Context, blueprintID string, opts *GetBlueprintQueryOptions) (*BlueprintResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &GetBlueprintQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[blueprints]", opts.Fields)
	}
	if len(opts.Include) > 0 {
		params.AddStringSlice("include", opts.Include)
	}
	if opts.LimitApps > 0 {
		if opts.LimitApps > 1000 {
			opts.LimitApps = 1000
		}
		params.AddInt("limit[apps]", opts.LimitApps)
	}
	if opts.LimitConfigurations > 0 {
		if opts.LimitConfigurations > 1000 {
			opts.LimitConfigurations = 1000
		}
		params.AddInt("limit[configurations]", opts.LimitConfigurations)
	}
	if opts.LimitPackages > 0 {
		if opts.LimitPackages > 1000 {
			opts.LimitPackages = 1000
		}
		params.AddInt("limit[packages]", opts.LimitPackages)
	}
	if opts.LimitOrgDevices > 0 {
		if opts.LimitOrgDevices > 1000 {
			opts.LimitOrgDevices = 1000
		}
		params.AddInt("limit[orgDevices]", opts.LimitOrgDevices)
	}
	if opts.LimitUsers > 0 {
		if opts.LimitUsers > 1000 {
			opts.LimitUsers = 1000
		}
		params.AddInt("limit[users]", opts.LimitUsers)
	}
	if opts.LimitUserGroups > 0 {
		if opts.LimitUserGroups > 1000 {
			opts.LimitUserGroups = 1000
		}
		params.AddInt("limit[userGroups]", opts.LimitUserGroups)
	}

	var result BlueprintResponse

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

// UpdateByBlueprintIDV1 updates an existing Blueprint in an organization.
// URL: PATCH https://api-business.apple.com/v1/blueprints/{id}
// https://developer.apple.com/documentation/applebusinessapi/update-a-blueprint
//
// Only provided fields are updated. Any invalid resource IDs in relationships are silently dropped.
func (s *Blueprints) UpdateByBlueprintIDV1(ctx context.Context, blueprintID string, req *BlueprintUpdateRequest) (*BlueprintResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID

	var result BlueprintResponse

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

// DeleteByBlueprintIDV1 deletes a Blueprint from an organization.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}
// https://developer.apple.com/documentation/applebusinessapi/delete-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) DeleteByBlueprintIDV1(ctx context.Context, blueprintID string) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetAppIDsByBlueprintIDV1 retrieves a list of app IDs associated with a Blueprint.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}/relationships/apps
// https://developer.apple.com/documentation/applebusinessapi/get-app-i-ds-for-a-blueprint
func (s *Blueprints) GetAppIDsByBlueprintIDV1(ctx context.Context, blueprintID string, opts *RequestQueryOptions) (*BlueprintAppsLinkagesResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/apps"

	params := s.client.QueryBuilder()
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []BlueprintLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse BlueprintAppsLinkagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allLinkages = append(allLinkages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &BlueprintAppsLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// AddAppsToBlueprintV1 adds apps to a Blueprint.
// URL: POST https://api-business.apple.com/v1/blueprints/{id}/relationships/apps
// https://developer.apple.com/documentation/applebusinessapi/add-apps-to-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) AddAppsToBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintAppsLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/apps"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Post(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveAppsFromBlueprintV1 removes apps from a Blueprint.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}/relationships/apps
// https://developer.apple.com/documentation/applebusinessapi/remove-apps-from-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) RemoveAppsFromBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintAppsLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/apps"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetConfigurationIDsByBlueprintIDV1 retrieves a list of configuration IDs associated with a Blueprint.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}/relationships/configurations
// https://developer.apple.com/documentation/applebusinessapi/get-configuration-i-ds-for-a-blueprint
func (s *Blueprints) GetConfigurationIDsByBlueprintIDV1(ctx context.Context, blueprintID string, opts *RequestQueryOptions) (*BlueprintConfigurationsLinkagesResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/configurations"

	params := s.client.QueryBuilder()
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []BlueprintLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse BlueprintConfigurationsLinkagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allLinkages = append(allLinkages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &BlueprintConfigurationsLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// AddConfigurationsToBlueprintV1 adds configurations to a Blueprint.
// URL: POST https://api-business.apple.com/v1/blueprints/{id}/relationships/configurations
// https://developer.apple.com/documentation/applebusinessapi/add-configurations-to-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) AddConfigurationsToBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintConfigurationsLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/configurations"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Post(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveConfigurationsFromBlueprintV1 removes configurations from a Blueprint.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}/relationships/configurations
// https://developer.apple.com/documentation/applebusinessapi/remove-configurations-from-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) RemoveConfigurationsFromBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintConfigurationsLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/configurations"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetPackageIDsByBlueprintIDV1 retrieves a list of package IDs associated with a Blueprint.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}/relationships/packages
// https://developer.apple.com/documentation/applebusinessapi/get-package-i-ds-for-a-blueprint
func (s *Blueprints) GetPackageIDsByBlueprintIDV1(ctx context.Context, blueprintID string, opts *RequestQueryOptions) (*BlueprintPackagesLinkagesResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/packages"

	params := s.client.QueryBuilder()
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []BlueprintLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse BlueprintPackagesLinkagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allLinkages = append(allLinkages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &BlueprintPackagesLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// AddPackagesToBlueprintV1 adds packages to a Blueprint.
// URL: POST https://api-business.apple.com/v1/blueprints/{id}/relationships/packages
// https://developer.apple.com/documentation/applebusinessapi/add-packages-to-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) AddPackagesToBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintPackagesLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/packages"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Post(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemovePackagesFromBlueprintV1 removes packages from a Blueprint.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}/relationships/packages
// https://developer.apple.com/documentation/applebusinessapi/remove-packages-from-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) RemovePackagesFromBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintPackagesLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/packages"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetDeviceIDsByBlueprintIDV1 retrieves a list of device IDs associated with a Blueprint.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}/relationships/orgDevices
// https://developer.apple.com/documentation/applebusinessapi/get-device-i-ds-for-a-blueprint
func (s *Blueprints) GetDeviceIDsByBlueprintIDV1(ctx context.Context, blueprintID string, opts *RequestQueryOptions) (*BlueprintOrgDevicesLinkagesResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/orgDevices"

	params := s.client.QueryBuilder()
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []BlueprintLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse BlueprintOrgDevicesLinkagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allLinkages = append(allLinkages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &BlueprintOrgDevicesLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// AddDevicesToBlueprintV1 adds devices to a Blueprint.
// URL: POST https://api-business.apple.com/v1/blueprints/{id}/relationships/orgDevices
// https://developer.apple.com/documentation/applebusinessapi/add-devices-to-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) AddDevicesToBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintOrgDevicesLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/orgDevices"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Post(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveDevicesFromBlueprintV1 removes devices from a Blueprint.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}/relationships/orgDevices
// https://developer.apple.com/documentation/applebusinessapi/remove-devices-from-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) RemoveDevicesFromBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintOrgDevicesLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/orgDevices"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetUserIDsByBlueprintIDV1 retrieves a list of user IDs associated with a Blueprint.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}/relationships/users
// https://developer.apple.com/documentation/applebusinessapi/get-user-i-ds-for-a-blueprint
func (s *Blueprints) GetUserIDsByBlueprintIDV1(ctx context.Context, blueprintID string, opts *RequestQueryOptions) (*BlueprintUsersLinkagesResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/users"

	params := s.client.QueryBuilder()
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []BlueprintLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse BlueprintUsersLinkagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allLinkages = append(allLinkages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &BlueprintUsersLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// AddUsersToBlueprintV1 adds users to a Blueprint.
// URL: POST https://api-business.apple.com/v1/blueprints/{id}/relationships/users
// https://developer.apple.com/documentation/applebusinessapi/add-users-to-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) AddUsersToBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintUsersLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/users"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Post(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveUsersFromBlueprintV1 removes users from a Blueprint.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}/relationships/users
// https://developer.apple.com/documentation/applebusinessapi/remove-users-from-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) RemoveUsersFromBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintUsersLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/users"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetUserGroupIDsByBlueprintIDV1 retrieves a list of user group IDs associated with a Blueprint.
// URL: GET https://api-business.apple.com/v1/blueprints/{id}/relationships/userGroups
// https://developer.apple.com/documentation/applebusinessapi/get-user-group-i-ds-for-a-blueprint
func (s *Blueprints) GetUserGroupIDsByBlueprintIDV1(ctx context.Context, blueprintID string, opts *RequestQueryOptions) (*BlueprintUserGroupsLinkagesResponse, *resty.Response, error) {
	if blueprintID == "" {
		return nil, nil, fmt.Errorf("blueprint ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/userGroups"

	params := s.client.QueryBuilder()
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []BlueprintLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse BlueprintUserGroupsLinkagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allLinkages = append(allLinkages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &BlueprintUserGroupsLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// AddUserGroupsToBlueprintV1 adds user groups to a Blueprint.
// URL: POST https://api-business.apple.com/v1/blueprints/{id}/relationships/userGroups
// https://developer.apple.com/documentation/applebusinessapi/add-user-groups-to-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) AddUserGroupsToBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintUserGroupsLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/userGroups"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Post(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveUserGroupsFromBlueprintV1 removes user groups from a Blueprint.
// URL: DELETE https://api-business.apple.com/v1/blueprints/{id}/relationships/userGroups
// https://developer.apple.com/documentation/applebusinessapi/remove-user-groups-from-a-blueprint
//
// Returns 204 No Content on success.
func (s *Blueprints) RemoveUserGroupsFromBlueprintV1(ctx context.Context, blueprintID string, req *BlueprintUserGroupsLinkagesRequest) (*resty.Response, error) {
	if blueprintID == "" {
		return nil, fmt.Errorf("blueprint ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointBlueprints + "/" + blueprintID + "/relationships/userGroups"

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
