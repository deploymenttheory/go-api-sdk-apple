package auditevents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// AuditEvents handles communication with the audit events
// related methods of the Apple Business Manager API.
//
// Apple Business Manager API docs: https://developer.apple.com/documentation/applebusinessapi/get-audit-events
type (
	AuditEvents struct {
		client client.Client
	}
)

// NewService creates a new audit events service.
func NewService(c client.Client) *AuditEvents {
	return &AuditEvents{client: c}
}

// GetV1 retrieves a list of audit events in an organization that satisfies the query criteria.
// URL: GET https://api-business.apple.com/v1/auditEvents
// https://developer.apple.com/documentation/applebusinessapi/get-audit-events
//
// filter[startTimestamp] and filter[endTimestamp] are required.
func (s *AuditEvents) GetV1(ctx context.Context, opts *RequestQueryOptions) (*AuditEventsResponse, *resty.Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("opts is required: filter[startTimestamp] and filter[endTimestamp] must be provided")
	}
	if opts.FilterStartTimestamp == "" {
		return nil, nil, fmt.Errorf("filter[startTimestamp] is required")
	}
	if opts.FilterEndTimestamp == "" {
		return nil, nil, fmt.Errorf("filter[endTimestamp] is required")
	}

	params := s.client.QueryBuilder()
	params.AddString("filter[startTimestamp]", opts.FilterStartTimestamp)
	params.AddString("filter[endTimestamp]", opts.FilterEndTimestamp)

	if opts.FilterActorID != "" {
		params.AddString("filter[actorId]", opts.FilterActorID)
	}
	if opts.FilterSubjectID != "" {
		params.AddString("filter[subjectId]", opts.FilterSubjectID)
	}
	if opts.FilterType != "" {
		params.AddString("filter[type]", opts.FilterType)
	}
	if len(opts.Fields) > 0 {
		params.AddStringSlice("fields[auditEvents]", opts.Fields)
	}
	if opts.Limit > 0 {
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
		params.AddInt("limit", opts.Limit)
	}

	var allEvents []AuditEvent
	var lastMeta *Meta
	var lastLinks *Links

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		GetPaginated(constants.EndpointAuditEvents, func(pageData []byte) error {
			var pageResponse AuditEventsResponse
			if err := json.Unmarshal(pageData, &pageResponse); err != nil {
				return fmt.Errorf("failed to unmarshal page: %w", err)
			}
			allEvents = append(allEvents, pageResponse.Data...)
			lastMeta = pageResponse.Meta
			lastLinks = pageResponse.Links
			return nil
		})

	if err != nil {
		return nil, resp, err
	}

	return &AuditEventsResponse{
		Data:  allEvents,
		Meta:  lastMeta,
		Links: lastLinks,
	}, resp, nil
}
