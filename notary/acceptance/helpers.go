package acceptance

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary"
	"github.com/stretchr/testify/require"
)

// SkipIfNotConfigured skips the test when Notary credentials are not set.
func SkipIfNotConfigured(t *testing.T) {
	t.Helper()
	if !IsConfigured() {
		t.Skip("APPLE_KEY_ID, APPLE_ISSUER_ID, or APPLE_PRIVATE_KEY_PEM/APPLE_PRIVATE_KEY_PATH not set — skipping acceptance test")
	}
}

// RequireClient ensures the shared client is initialised, skipping the test
// when credentials are absent or initialisation fails.
func RequireClient(t *testing.T) {
	t.Helper()
	SkipIfNotConfigured(t)

	if Client == nil {
		err := InitClient()
		require.NoError(t, err, "Failed to initialise Notary client")
	}
}

// NewContext returns a context bound to the configured request timeout.
func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), Config.RequestTimeout)
}

// Cleanup registers fn as a test cleanup function. When NOTARY_SKIP_CLEANUP=true
// the function is silently omitted so test data persists for debugging.
func Cleanup(t *testing.T, fn func()) {
	t.Helper()
	if !Config.SkipCleanup {
		t.Cleanup(fn)
	} else if Config.Verbose {
		t.Log("Skipping cleanup (NOTARY_SKIP_CLEANUP=true)")
	}
}

// LogTestStage logs a named stage with optional GitHub Actions ::notice annotation.
func LogTestStage(t *testing.T, stage, message string, args ...any) {
	t.Helper()
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	if isGitHubActions() {
		fmt.Printf("::notice title=%s::%s\n", stage, formatted)
	}
	if Config.Verbose {
		t.Logf("[%s] %s", stage, formatted)
	}
}

// LogTestSuccess logs a successful step.
func LogTestSuccess(t *testing.T, message string, args ...any) {
	t.Helper()
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	if isGitHubActions() {
		fmt.Printf("::notice title=Success::%s\n", formatted)
	}
	if Config.Verbose {
		t.Logf("OK: %s", formatted)
	}
}

// LogTestWarning logs a non-fatal warning.
func LogTestWarning(t *testing.T, message string, args ...any) {
	t.Helper()
	formatted := message
	if len(args) > 0 {
		formatted = fmt.Sprintf(message, args...)
	}
	if isGitHubActions() {
		fmt.Printf("::warning title=Warning::%s\n", formatted)
	}
	if Config.Verbose {
		t.Logf("WARNING: %s", formatted)
	}
}

// LogSubmissionError logs the result of a submission-related error. A 404 is
// treated as expected; other errors are non-fatal warnings.
func LogSubmissionError(t *testing.T, submissionID string, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if notary.IsNotFound(err) {
		LogTestStage(t, "Info", "submission ID=%s not found (404 — expected)", submissionID)
		return
	}
	LogTestWarning(t, "submission ID=%s error: %v", submissionID, err)
}

// PollUntil retries fn every interval until it returns true or timeout elapses.
func PollUntil(t *testing.T, timeout, interval time.Duration, fn func() bool) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

// RetryOnNotFound retries fn when it returns a 404 error, with exponential
// back-off. Returns nil on success or the last error when retries are exhausted.
func RetryOnNotFound(t *testing.T, maxRetries int, initialDelay time.Duration, fn func() error) error {
	t.Helper()
	var lastErr error
	delay := initialDelay

	for i := 0; i < maxRetries; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !notary.IsNotFound(lastErr) {
			return lastErr
		}
		if i < maxRetries-1 {
			if Config.Verbose {
				t.Logf("Resource not found (404), retry %d/%d — waiting %v", i+1, maxRetries, delay)
			}
			time.Sleep(delay)
			delay *= 2
		}
	}
	return lastErr
}

func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}
