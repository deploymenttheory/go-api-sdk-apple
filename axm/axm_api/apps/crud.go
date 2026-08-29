package apps

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Apps handles communication with the apps
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-apps
type (
	Apps struct {
		client client.Client
	}
)

// NewService creates a new apps service.
func NewService(c client.Client) *Apps {
	return &Apps{client: c}
}

// GetV1 retrieves a list of licensed apps in an organization.
// URL: GET https://api-business.apple.com/v1/apps
// https://developer.apple.com/documentation/applebusinessapi/get-apps
func (s *Apps) GetV1(ctx context.Context, opts *RequestQueryOptions) (*AppsResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[apps]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allApps []App
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointApps, func(pageData []byte) error {
			var pageResponse AppsResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allApps = append(allApps, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &AppsResponse{
		Data:  allApps,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByAppIDV1 retrieves information about a specific app in an organization.
// URL: GET https://api-business.apple.com/v1/apps/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-app-information
func (s *Apps) GetByAppIDV1(ctx context.Context, appID string, opts *RequestQueryOptions) (*AppResponse, *resty.Response, error) {
	if appID == "" {
		return nil, nil, fmt.Errorf("app ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointApps + "/" + appID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[apps]", opts.Fields)
	}

	var result AppResponse

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
