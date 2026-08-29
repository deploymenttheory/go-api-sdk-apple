package organizationalunits

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// OrganizationalUnits handles communication with the organizational units
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-organizational-units
type (
	OrganizationalUnits struct {
		client client.Client
	}
)

// NewService creates a new organizational units service.
func NewService(c client.Client) *OrganizationalUnits {
	return &OrganizationalUnits{client: c}
}

// GetV1 retrieves a list of organizational units in an organization.
// URL: GET https://api-business.apple.com/v1/organizationalUnits
// https://developer.apple.com/documentation/applebusinessapi/get-organizational-units
func (s *OrganizationalUnits) GetV1(ctx context.Context, opts *RequestQueryOptions) (*OrganizationalUnitsResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[organizationalUnits]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allUnits []OrganizationalUnit
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointOrganizationalUnits, func(pageData []byte) error {
			var pageResponse OrganizationalUnitsResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allUnits = append(allUnits, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &OrganizationalUnitsResponse{
		Data:  allUnits,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByOrganizationalUnitIDV1 retrieves information about a specific organizational unit in an organization.
// URL: GET https://api-business.apple.com/v1/organizationalUnits/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-organizational-unit-information
func (s *OrganizationalUnits) GetByOrganizationalUnitIDV1(ctx context.Context, unitID string, opts *RequestQueryOptions) (*OrganizationalUnitResponse, *resty.Response, error) {
	if unitID == "" {
		return nil, nil, fmt.Errorf("organizational unit ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointOrganizationalUnits + "/" + unitID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[organizationalUnits]", opts.Fields)
	}

	var result OrganizationalUnitResponse

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

// GetUserIDsByOrganizationalUnitIDV1 retrieves a list of user IDs for an organizational unit in an organization.
// URL: GET https://api-business.apple.com/v1/organizationalUnits/{id}/relationships/users
// https://developer.apple.com/documentation/applebusinessapi/get-user-ids-for-an-organizational-unit
func (s *OrganizationalUnits) GetUserIDsByOrganizationalUnitIDV1(ctx context.Context, unitID string, opts *RequestQueryOptions) (*OrganizationalUnitUsersLinkagesResponse, *resty.Response, error) {
	if unitID == "" {
		return nil, nil, fmt.Errorf("organizational unit ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointOrganizationalUnits+"/%s/relationships/users", unitID)

	params := s.client.QueryBuilder()

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []OrganizationalUnitUserLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse OrganizationalUnitUsersLinkagesResponse
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

	return &OrganizationalUnitUsersLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}
