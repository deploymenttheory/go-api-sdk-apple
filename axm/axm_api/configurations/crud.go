package configurations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Configurations handles communication with the configurations
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-configurations
type (
	Configurations struct {
		client client.Client
	}
)

// NewService creates a new configurations service.
func NewService(c client.Client) *Configurations {
	return &Configurations{client: c}
}

// GetV1 retrieves a list of configurations in an organization.
// URL: GET https://api-business.apple.com/v1/configurations
// https://developer.apple.com/documentation/applebusinessapi/get-configurations
//
// Note: customSettingsValues is always null in list responses. Use GetByConfigurationIDV1 to
// retrieve customSettingsValues for CUSTOM_SETTING configurations.
func (s *Configurations) GetV1(ctx context.Context, opts *RequestQueryOptions) (*ConfigurationsResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[configurations]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allConfigurations []Configuration
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointConfigurations, func(pageData []byte) error {
			var pageResponse ConfigurationsResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allConfigurations = append(allConfigurations, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &ConfigurationsResponse{
		Data:  allConfigurations,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByConfigurationIDV1 retrieves information about a specific configuration in an organization.
// URL: GET https://api-business.apple.com/v1/configurations/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-configuration-information
//
// Use this endpoint to retrieve customSettingsValues for CUSTOM_SETTING configurations.
func (s *Configurations) GetByConfigurationIDV1(ctx context.Context, configID string, opts *RequestQueryOptions) (*ConfigurationResponse, *resty.Response, error) {
	if configID == "" {
		return nil, nil, fmt.Errorf("configuration ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointConfigurations + "/" + configID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[configurations]", opts.Fields)
	}

	var result ConfigurationResponse

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

// CreateV1 creates a new custom configuration in an organization.
// URL: POST https://api-business.apple.com/v1/configurations
// https://developer.apple.com/documentation/applebusinessapi/create-a-configuration
//
// Only configurations with type CUSTOM_SETTING can be created via the API.
// configurationProfile is required. filename and configuredForPlatforms are optional.
func (s *Configurations) CreateV1(ctx context.Context, req *ConfigurationCreateRequest) (*ConfigurationResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	if req.Data.Attributes.CustomSettingsValues.ConfigurationProfile == "" {
		return nil, nil, fmt.Errorf("configurationProfile is required")
	}

	var result ConfigurationResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		SetResult(&result).
		Post(constants.EndpointConfigurations)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UpdateByConfigurationIDV1 updates an existing custom configuration in an organization.
// URL: PATCH https://api-business.apple.com/v1/configurations/{id}
// https://developer.apple.com/documentation/applebusinessapi/update-a-configuration
//
// Only CUSTOM_SETTING configurations can be updated. Only provided fields are changed.
// If filename is provided it must end in .mobileconfig.
func (s *Configurations) UpdateByConfigurationIDV1(ctx context.Context, configID string, req *ConfigurationUpdateRequest) (*ConfigurationResponse, *resty.Response, error) {
	if configID == "" {
		return nil, nil, fmt.Errorf("configuration ID is required")
	}
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}

	endpoint := constants.EndpointConfigurations + "/" + configID

	var result ConfigurationResponse

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

// DeleteByConfigurationIDV1 deletes a configuration from an organization.
// URL: DELETE https://api-business.apple.com/v1/configurations/{id}
// https://developer.apple.com/documentation/applebusinessapi/delete-a-configuration
//
// Returns 204 No Content on success.
func (s *Configurations) DeleteByConfigurationIDV1(ctx context.Context, configID string) (*resty.Response, error) {
	if configID == "" {
		return nil, fmt.Errorf("configuration ID is required")
	}

	endpoint := constants.EndpointConfigurations + "/" + configID

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		Delete(endpoint)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
