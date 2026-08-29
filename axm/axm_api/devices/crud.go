package devices

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Devices handles communication with the device
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessmanagerapi/
type (
	Devices struct {
		client client.Client
	}
)

// NewService creates a new devices service.
func NewService(c client.Client) *Devices {
	return &Devices{client: c}
}

// GetV1 retrieves a list of devices in an organization that enroll using Automated Device Enrollment.
// URL: GET https://api-business.apple.com/v1/orgDevices
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-org-devices
func (s *Devices) GetV1(ctx context.Context, opts *RequestQueryOptions) (*OrgDevicesResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[orgDevices]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		params.AddInt("limit", opts.Limit)
	}

	var allDevices []OrgDevice
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointOrgDevices, func(pageData []byte) error {
			var pageResponse OrgDevicesResponse
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

	return &OrgDevicesResponse{
		Data:  allDevices,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByDeviceIDV1 retrieves information about a specific device in an organization.
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-orgdevice-information
func (s *Devices) GetByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*OrgDeviceResponse, *resty.Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointOrgDevices + "/" + deviceID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[orgDevices]", opts.Fields)
	}

	var result OrgDeviceResponse

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

// GetAppleCareByDeviceIDV1 retrieves AppleCare coverage information for a specific device.
// URL: GET https://api-business.apple.com/v1/orgDevices/{id}/appleCareCoverage
// https://developer.apple.com/documentation/applebusinessmanagerapi/get-all-apple-care-coverage-for-an-orgdevice
func (s *Devices) GetAppleCareByDeviceIDV1(ctx context.Context, deviceID string, opts *RequestQueryOptions) (*AppleCareCoverageResponse, *resty.Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointOrgDevices + "/" + deviceID + "/appleCareCoverage"

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[appleCareCoverage]", opts.Fields)
	}

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000 // Enforce API maximum
		}
		params.AddInt("limit", opts.Limit)
	}

	var allCoverage []AppleCareCoverage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse AppleCareCoverageResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allCoverage = append(allCoverage, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &AppleCareCoverageResponse{
		Data:  allCoverage,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}
