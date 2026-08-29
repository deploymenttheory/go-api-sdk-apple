package packages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Packages handles communication with the packages
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-packages
type (
	Packages struct {
		client client.Client
	}
)

// NewService creates a new packages service.
func NewService(c client.Client) *Packages {
	return &Packages{client: c}
}

// GetV1 retrieves a list of packages available in an organization.
// URL: GET https://api-business.apple.com/v1/packages
// https://developer.apple.com/documentation/applebusinessapi/get-packages
func (s *Packages) GetV1(ctx context.Context, opts *RequestQueryOptions) (*PackagesResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[packages]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allPackages []Package
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointPackages, func(pageData []byte) error {
			var pageResponse PackagesResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allPackages = append(allPackages, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &PackagesResponse{
		Data:  allPackages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByPackageIDV1 retrieves information about a specific package in an organization.
// URL: GET https://api-business.apple.com/v1/packages/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-package-information
func (s *Packages) GetByPackageIDV1(ctx context.Context, packageID string, opts *RequestQueryOptions) (*PackageResponse, *resty.Response, error) {
	if packageID == "" {
		return nil, nil, fmt.Errorf("package ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointPackages + "/" + packageID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[packages]", opts.Fields)
	}

	var result PackageResponse

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
