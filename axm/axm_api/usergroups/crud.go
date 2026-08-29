package usergroups

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// UserGroups handles communication with the user groups
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-user-groups
type (
	UserGroups struct {
		client client.Client
	}
)

// NewService creates a new user groups service.
func NewService(c client.Client) *UserGroups {
	return &UserGroups{client: c}
}

// GetV1 retrieves a list of user groups in an organization.
// URL: GET https://api-business.apple.com/v1/userGroups
// https://developer.apple.com/documentation/applebusinessapi/get-user-groups
func (s *UserGroups) GetV1(ctx context.Context, opts *RequestQueryOptions) (*UserGroupsResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[userGroups]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allGroups []UserGroup
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointUserGroups, func(pageData []byte) error {
			var pageResponse UserGroupsResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allGroups = append(allGroups, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &UserGroupsResponse{
		Data:  allGroups,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}

// GetByUserGroupIDV1 retrieves information about a specific user group in an organization.
// URL: GET https://api-business.apple.com/v1/userGroups/{id}
// https://developer.apple.com/documentation/applebusinessapi/get-usergroup-information
func (s *UserGroups) GetByUserGroupIDV1(ctx context.Context, groupID string, opts *RequestQueryOptions) (*UserGroupResponse, *resty.Response, error) {
	if groupID == "" {
		return nil, nil, fmt.Errorf("group ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := constants.EndpointUserGroups + "/" + groupID

	params := s.client.QueryBuilder()

	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[userGroups]", opts.Fields)
	}

	var result UserGroupResponse

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

// GetUserIDsByGroupIDV1 retrieves a list of user IDs for a user group in an organization.
// URL: GET https://api-business.apple.com/v1/userGroups/{id}/relationships/users
// https://developer.apple.com/documentation/applebusinessapi/get-all-user-ids-for-a-user-group
func (s *UserGroups) GetUserIDsByGroupIDV1(ctx context.Context, groupID string, opts *RequestQueryOptions) (*UserGroupUsersLinkagesResponse, *resty.Response, error) {
	if groupID == "" {
		return nil, nil, fmt.Errorf("group ID is required")
	}

	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	endpoint := fmt.Sprintf(constants.EndpointUserGroups+"/%s/relationships/users", groupID)

	params := s.client.QueryBuilder()

	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allLinkages []UserGroupUserLinkage
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(endpoint, func(pageData []byte) error {
			var pageResponse UserGroupUsersLinkagesResponse
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

	return &UserGroupUsersLinkagesResponse{
		Data:  allLinkages,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}
