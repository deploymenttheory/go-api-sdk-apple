package users

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Users handles communication with the users
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-users
type (
	Users struct {
		client client.Client
	}
)

// NewService creates a new users service.
func NewService(c client.Client) *Users {
	return &Users{client: c}
}

// GetV1 retrieves a list of users in an organization.
// URL: GET https://api-business.apple.com/v1/users
// https://developer.apple.com/documentation/applebusinessapi/get-users
func (s *Users) GetV1(ctx context.Context, opts *RequestQueryOptions) (*UsersResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[users]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allUsers []User
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointUsers, func(pageData []byte) error {
			var pageResponse UsersResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allUsers = append(allUsers, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &UsersResponse{
		Data:  allUsers,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByUserIDV1 retrieves information about a specific user in an organization.
// URL: GET https://api-business.apple.com/v1/users/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-user-information
func (s *Users) GetByUserIDV1(ctx context.Context, userID string, opts *RequestQueryOptions) (*UserResponse, *resty.Response, error) {
	if userID == "" {
		return nil, nil, fmt.Errorf("user ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointUsers + "/" + userID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[users]", opts.Fields)
	}

	var result UserResponse

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
