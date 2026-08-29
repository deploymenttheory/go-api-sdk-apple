package submissions

import (
	"context"
	"fmt"
	"regexp"

	"github.com/deploymenttheory/go-sdk-appleservices/notary/client"
	"github.com/deploymenttheory/go-sdk-appleservices/notary/constants"
	"resty.dev/v3"
)

var sha256Re = regexp.MustCompile(`^[A-Fa-f0-9]{64}$`)

// Submissions handles communication with the Apple Notary API submissions endpoints.
//
// Apple Notary API docs: https://developer.apple.com/documentation/notaryapi
type Submissions struct {
	client client.Client
}

// NewService creates a new submissions service.
func NewService(c client.Client) *Submissions {
	return &Submissions{client: c}
}

// SubmitSoftwareV2 starts the notarization process for a new version of software.
// URL: POST https://appstoreconnect.apple.com/notary/v2/submissions
//
// The response contains temporary AWS credentials to upload the software to S3
// and a submission ID for tracking the notarization outcome.
func (s *Submissions) SubmitSoftwareV2(ctx context.Context, req *NewSubmissionRequest) (*NewSubmissionResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	if req.SubmissionName == "" {
		return nil, nil, fmt.Errorf("submissionName is required")
	}
	if !sha256Re.MatchString(req.SHA256) {
		return nil, nil, fmt.Errorf("sha256 must be a 64-character hexadecimal string")
	}

	var result NewSubmissionResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(req).
		SetResult(&result).
		Post(constants.EndpointSubmissions)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetPreviousSubmissionsV2 retrieves the list of the team's most recent notarization
// submissions (up to 100 most recent).
// URL: GET https://appstoreconnect.apple.com/notary/v2/submissions
func (s *Submissions) GetPreviousSubmissionsV2(ctx context.Context) (*SubmissionListResponse, *resty.Response, error) {
	var result SubmissionListResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointSubmissions)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetSubmissionStatusV2 fetches the current status of a notarization submission.
// URL: GET https://appstoreconnect.apple.com/notary/v2/submissions/{submissionId}
func (s *Submissions) GetSubmissionStatusV2(ctx context.Context, submissionID string) (*SubmissionResponse, *resty.Response, error) {
	if submissionID == "" {
		return nil, nil, fmt.Errorf("submissionID is required")
	}

	endpoint := constants.EndpointSubmissions + "/" + submissionID

	var result SubmissionResponse
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

// GetSubmissionLogV2 retrieves a temporary URL for downloading the notarization log.
// URL: GET https://appstoreconnect.apple.com/notary/v2/submissions/{submissionId}/logs
//
// The returned URL is valid for only a few hours. If you need the log again later,
// call this endpoint again to get a fresh URL.
func (s *Submissions) GetSubmissionLogV2(ctx context.Context, submissionID string) (*SubmissionLogURLResponse, *resty.Response, error) {
	if submissionID == "" {
		return nil, nil, fmt.Errorf("submissionID is required")
	}

	endpoint := constants.EndpointSubmissions + "/" + submissionID + constants.EndpointLogs

	var result SubmissionLogURLResponse
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
