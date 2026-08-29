package submissions

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary/notary_api/submissions/mocks"
	"github.com/deploymenttheory/go-sdk-appleservices/notary/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// setupMockClient creates a client with httpmock enabled.
func setupMockClient(t *testing.T) *Submissions {
	mockAuth := &MockAuthProvider{}

	coreClient, err := client.NewTransport(
		"test-key-id",
		"test-issuer-id",
		"dummy-key",
		client.WithAuth(mockAuth),
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(coreClient.GetHTTPClient().Client())

	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewService(coreClient)
}

// MockAuthProvider implements the AuthProvider interface for testing.
type MockAuthProvider struct{}

func (m *MockAuthProvider) ApplyAuth(req *resty.Request) error {
	return nil
}

// --- SubmitSoftwareV2 tests ---

func TestSubmitSoftwareV2_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &NewSubmissionRequest{
		SHA256:         "68d561c564ef61f718e99a81b13bcb52af11b7ac9baf538af3ea0c83326fb6a1",
		SubmissionName: "OvernightTextEditor_11.6.8.zip",
		Notifications: []NewSubmissionRequestNotification{
			{Channel: NotificationChannelWebhook, Target: "https://example.com"},
		},
	}

	result, resp, err := svc.SubmitSoftwareV2(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, "2efe2717-52ef-43a5-96dc-0797e4ca1041", result.Data.ID)
	assert.Equal(t, "submissionsPostResponse", result.Data.Type)
	assert.Equal(t, "ASIAIOSFODNN7EXAMPLE", result.Data.Attributes.AWSAccessKeyID)
	assert.Equal(t, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", result.Data.Attributes.AWSSecretAccessKey)
	assert.Equal(t, "AQoDYXdzEJr...", result.Data.Attributes.AWSSessionToken)
	assert.Equal(t, "EXAMPLE-BUCKET", result.Data.Attributes.Bucket)
	assert.Equal(t, "EXAMPLE-KEY-NAME", result.Data.Attributes.Object)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestSubmitSoftwareV2_NilRequest(t *testing.T) {
	svc := setupMockClient(t)

	ctx := context.Background()
	result, _, err := svc.SubmitSoftwareV2(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "request is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSubmitSoftwareV2_EmptySubmissionName(t *testing.T) {
	svc := setupMockClient(t)

	ctx := context.Background()
	req := &NewSubmissionRequest{
		SHA256:         "68d561c564ef61f718e99a81b13bcb52af11b7ac9baf538af3ea0c83326fb6a1",
		SubmissionName: "",
	}

	result, _, err := svc.SubmitSoftwareV2(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "submissionName is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestSubmitSoftwareV2_InvalidSHA256(t *testing.T) {
	svc := setupMockClient(t)

	tests := []struct {
		name string
		sha  string
	}{
		{"Too short", "68d561c5"},
		{"Too long", "68d561c564ef61f718e99a81b13bcb52af11b7ac9baf538af3ea0c83326fb6a1XX"},
		{"Non-hex characters", "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"},
		{"Empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := &NewSubmissionRequest{
				SHA256:         tt.sha,
				SubmissionName: "test.zip",
			}

			result, _, err := svc.SubmitSoftwareV2(ctx, req)

			require.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "sha256")
		})
	}
}

func TestSubmitSoftwareV2_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	req := &NewSubmissionRequest{
		SHA256:         "68d561c564ef61f718e99a81b13bcb52af11b7ac9baf538af3ea0c83326fb6a1",
		SubmissionName: "OvernightTextEditor_11.6.8.zip",
	}

	result, _, err := svc.SubmitSoftwareV2(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestSubmitSoftwareV2_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &NewSubmissionRequest{
		SHA256:         "68d561c564ef61f718e99a81b13bcb52af11b7ac9baf538af3ea0c83326fb6a1",
		SubmissionName: "OvernightTextEditor_11.6.8.zip",
	}

	result, _, err := svc.SubmitSoftwareV2(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

// --- GetPreviousSubmissionsV2 tests ---

func TestGetPreviousSubmissionsV2_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, resp, err := svc.GetPreviousSubmissionsV2(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	require.Len(t, result.Data, 3)

	first := result.Data[0]
	assert.Equal(t, "2efe2717-52ef-43a5-96dc-0797e4ca1041", first.ID)
	assert.Equal(t, "submissions", first.Type)
	assert.Equal(t, "OvernightTextEditor_11.6.8.zip", first.Attributes.Name)
	assert.Equal(t, SubmissionStatusAccepted, first.Attributes.Status)
	assert.Equal(t, "2021-04-29T01:38:09.498Z", first.Attributes.CreatedDate)

	third := result.Data[2]
	assert.Equal(t, SubmissionStatusInvalid, third.Attributes.Status)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPreviousSubmissionsV2_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := svc.GetPreviousSubmissionsV2(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetPreviousSubmissionsV2_ContextTimeout(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Millisecond)

	result, _, err := svc.GetPreviousSubmissionsV2(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// --- GetSubmissionStatusV2 tests ---

func TestGetSubmissionStatusV2_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, resp, err := svc.GetSubmissionStatusV2(ctx, "2efe2717-52ef-43a5-96dc-0797e4ca1041")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, "2efe2717-52ef-43a5-96dc-0797e4ca1041", result.Data.ID)
	assert.Equal(t, "submissions", result.Data.Type)
	assert.Equal(t, "OvernightTextEditor_11.6.8.zip", result.Data.Attributes.Name)
	assert.Equal(t, SubmissionStatusAccepted, result.Data.Attributes.Status)
	assert.Equal(t, "2022-06-08T01:38:09.498Z", result.Data.Attributes.CreatedDate)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetSubmissionStatusV2_EmptyID(t *testing.T) {
	svc := setupMockClient(t)

	ctx := context.Background()
	result, _, err := svc.GetSubmissionStatusV2(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "submissionID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetSubmissionStatusV2_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := svc.GetSubmissionStatusV2(ctx, "00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetSubmissionStatusV2_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _, err := svc.GetSubmissionStatusV2(ctx, "2efe2717-52ef-43a5-96dc-0797e4ca1041")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

// --- GetSubmissionLogV2 tests ---

func TestGetSubmissionLogV2_Success(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, resp, err := svc.GetSubmissionLogV2(ctx, "2efe2717-52ef-43a5-96dc-0797e4ca1041")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, "2efe2717-52ef-43a5-96dc-0797e4ca1041", result.Data.ID)
	assert.Equal(t, "submissionsLog", result.Data.Type)
	assert.NotEmpty(t, result.Data.Attributes.DeveloperLogURL)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetSubmissionLogV2_EmptyID(t *testing.T) {
	svc := setupMockClient(t)

	ctx := context.Background()
	result, _, err := svc.GetSubmissionLogV2(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "submissionID is required")
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestGetSubmissionLogV2_NotFound(t *testing.T) {
	svc := setupMockClient(t)
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := svc.GetSubmissionLogV2(ctx, "00000000-0000-0000-0000-000000000000")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetSubmissionLogV2_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.Reset()
	mockHandler := &mocks.SubmissionsMock{}
	mockHandler.RegisterErrorMocks()
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := svc.GetSubmissionLogV2(ctx, "2efe2717-52ef-43a5-96dc-0797e4ca1041")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// --- Constants tests ---

func TestSubmissionStatusConstants(t *testing.T) {
	assert.Equal(t, "Accepted", SubmissionStatusAccepted)
	assert.Equal(t, "In Progress", SubmissionStatusInProgress)
	assert.Equal(t, "Invalid", SubmissionStatusInvalid)
	assert.Equal(t, "Rejected", SubmissionStatusRejected)
}

func TestNotificationChannelConstant(t *testing.T) {
	assert.Equal(t, "webhook", NotificationChannelWebhook)
}

func TestStatusConstants_AllPresent(t *testing.T) {
	statuses := []string{
		SubmissionStatusAccepted,
		SubmissionStatusInProgress,
		SubmissionStatusInvalid,
		SubmissionStatusRejected,
	}

	for _, s := range statuses {
		assert.NotEmpty(t, s, "Status constant should not be empty")
	}
}
