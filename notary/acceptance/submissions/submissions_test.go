package submissions_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary/acceptance"
	notarysubmissions "github.com/deploymenttheory/go-sdk-appleservices/notary/notary_api/submissions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetPreviousSubmissions_ReadOnly lists the team's recent submissions.
// This is safe to run against the real API because it's purely read-only.
func TestGetPreviousSubmissions_ReadOnly(t *testing.T) {
	acceptance.RequireClient(t)

	ctx, cancel := acceptance.NewContext()
	defer cancel()

	acceptance.LogTestStage(t, "List", "Fetching previous submissions")

	result, resp, err := acceptance.Client.NotaryAPI.Submissions.GetPreviousSubmissionsV2(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	acceptance.LogTestSuccess(t, "Retrieved %d previous submissions", len(result.Data))

	for _, sub := range result.Data {
		assert.NotEmpty(t, sub.ID, "submission ID should not be empty")
		assert.NotEmpty(t, sub.Type, "submission type should not be empty")
		assert.NotEmpty(t, sub.Attributes.Name, "submission name should not be empty")
		assert.NotEmpty(t, sub.Attributes.CreatedDate, "submission createdDate should not be empty")

		status := sub.Attributes.Status
		validStatuses := []string{
			notarysubmissions.SubmissionStatusAccepted,
			notarysubmissions.SubmissionStatusInProgress,
			notarysubmissions.SubmissionStatusInvalid,
			notarysubmissions.SubmissionStatusRejected,
		}
		validStatus := false
		for _, s := range validStatuses {
			if status == s {
				validStatus = true
				break
			}
		}
		assert.True(t, validStatus, "submission status %q is not one of the expected values", status)
	}
}

// TestGetSubmissionStatus_ReadOnly retrieves the status of the most recent submission.
// Skipped when no previous submissions exist.
func TestGetSubmissionStatus_ReadOnly(t *testing.T) {
	acceptance.RequireClient(t)

	ctx, cancel := acceptance.NewContext()
	defer cancel()

	acceptance.LogTestStage(t, "Setup", "Fetching previous submissions to find a candidate")

	listResult, _, err := acceptance.Client.NotaryAPI.Submissions.GetPreviousSubmissionsV2(ctx)
	require.NoError(t, err)

	if len(listResult.Data) == 0 {
		t.Skip("No previous submissions found — skipping status check")
	}

	submissionID := listResult.Data[0].ID
	acceptance.LogTestStage(t, "Status", "Fetching status for submission ID=%s", submissionID)

	ctx2, cancel2 := acceptance.NewContext()
	defer cancel2()

	result, resp, err := acceptance.Client.NotaryAPI.Submissions.GetSubmissionStatusV2(ctx2, submissionID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, submissionID, result.Data.ID)
	assert.NotEmpty(t, result.Data.Attributes.Name)
	assert.NotEmpty(t, result.Data.Attributes.Status)
	assert.NotEmpty(t, result.Data.Attributes.CreatedDate)

	acceptance.LogTestSuccess(t, "Submission ID=%s has status=%s", submissionID, result.Data.Attributes.Status)
}

// TestGetSubmissionLog_ReadOnly retrieves the log URL for a completed submission.
// Skipped when no completed submissions exist.
func TestGetSubmissionLog_ReadOnly(t *testing.T) {
	acceptance.RequireClient(t)

	ctx, cancel := acceptance.NewContext()
	defer cancel()

	acceptance.LogTestStage(t, "Setup", "Searching for a completed submission")

	listResult, _, err := acceptance.Client.NotaryAPI.Submissions.GetPreviousSubmissionsV2(ctx)
	require.NoError(t, err)

	// Find a terminal-status submission (Accepted, Invalid, or Rejected)
	var completedID string
	terminalStatuses := map[string]bool{
		notarysubmissions.SubmissionStatusAccepted: true,
		notarysubmissions.SubmissionStatusInvalid:  true,
		notarysubmissions.SubmissionStatusRejected: true,
	}
	for _, sub := range listResult.Data {
		if terminalStatuses[sub.Attributes.Status] {
			completedID = sub.ID
			break
		}
	}

	if completedID == "" {
		t.Skip("No completed submissions found — skipping log URL check")
	}

	acceptance.LogTestStage(t, "Log", "Fetching log URL for submission ID=%s", completedID)

	ctx2, cancel2 := acceptance.NewContext()
	defer cancel2()

	result, resp, err := acceptance.Client.NotaryAPI.Submissions.GetSubmissionLogV2(ctx2, completedID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)

	assert.Equal(t, completedID, result.Data.ID)
	assert.NotEmpty(t, result.Data.Attributes.DeveloperLogURL, "developer log URL should not be empty")

	acceptance.LogTestSuccess(t, "Got log URL for submission ID=%s", completedID)
}

// TestSubmissionPolling demonstrates polling a submission until it reaches terminal status.
// Only runs when NOTARY_RUN_POLLING=true is explicitly set (slow test).
func TestSubmissionPolling_IfEnabled(t *testing.T) {
	if acceptance.Config == nil {
		t.Skip("acceptance config not initialised")
	}

	// Opt-in test gate
	if !getBoolEnv("NOTARY_RUN_POLLING", false) {
		t.Skip("NOTARY_RUN_POLLING not set — skipping polling test")
	}

	acceptance.RequireClient(t)

	ctx, cancel := acceptance.NewContext()
	defer cancel()

	listResult, _, err := acceptance.Client.NotaryAPI.Submissions.GetPreviousSubmissionsV2(ctx)
	require.NoError(t, err)

	var inProgressID string
	for _, sub := range listResult.Data {
		if sub.Attributes.Status == notarysubmissions.SubmissionStatusInProgress {
			inProgressID = sub.ID
			break
		}
	}

	if inProgressID == "" {
		t.Skip("No in-progress submissions found")
	}

	acceptance.LogTestStage(t, "Poll", "Polling submission ID=%s until terminal status", inProgressID)

	terminalStatuses := map[string]bool{
		notarysubmissions.SubmissionStatusAccepted: true,
		notarysubmissions.SubmissionStatusInvalid:  true,
		notarysubmissions.SubmissionStatusRejected: true,
	}

	reached := acceptance.PollUntil(t, 10*time.Minute, 30*time.Second, func() bool {
		ctx, cancel := acceptance.NewContext()
		defer cancel()

		result, _, err := acceptance.Client.NotaryAPI.Submissions.GetSubmissionStatusV2(ctx, inProgressID)
		if err != nil {
			acceptance.LogTestWarning(t, "Poll error: %v", err)
			return false
		}

		acceptance.LogTestStage(t, "Poll", "Submission ID=%s status=%s", inProgressID, result.Data.Attributes.Status)
		return terminalStatuses[result.Data.Attributes.Status]
	})

	assert.True(t, reached, "Submission did not reach terminal status within timeout")
}

func getBoolEnv(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
