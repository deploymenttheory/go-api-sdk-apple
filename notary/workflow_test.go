package notary

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary/client"
	"github.com/deploymenttheory/go-sdk-appleservices/notary/upload"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"resty.dev/v3"
)

const testSubmissionID = "2efe2717-52ef-43a5-96dc-0797e4ca1041"

// mockAuth satisfies the transport's auth provider without signing anything,
// so tests need no real App Store Connect key.
type mockAuth struct{}

func (mockAuth) ApplyAuth(*resty.Request) error { return nil }

// newTestClient builds a notary Client whose HTTP calls are served by
// httpmock rather than reaching Apple.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	c, err := NewClient(
		"test-key-id", "test-issuer-id", "dummy-key",
		client.WithAuth(mockAuth{}),
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(c.transport.GetHTTPClient().Client())
	t.Cleanup(httpmock.DeactivateAndReset)
	return c
}

// pointUploadAt redirects SubmitAndWait's uploader at a local endpoint for the
// duration of the test.
func pointUploadAt(t *testing.T, srv *httptest.Server) {
	t.Helper()
	orig := newS3Uploader
	newS3Uploader = func() *upload.S3Uploader {
		return &upload.S3Uploader{Client: srv.Client(), Endpoint: srv.URL}
	}
	t.Cleanup(func() { newS3Uploader = orig })
}

const submitResponse = `{
  "data": {
    "attributes": {
      "awsAccessKeyId": "ASIAIOSFODNN7EXAMPLE",
      "awsSecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
      "awsSessionToken": "AQoDYXdzEJr...",
      "bucket": "EXAMPLE-BUCKET",
      "object": "EXAMPLE-KEY-NAME"
    },
    "id": "` + testSubmissionID + `",
    "type": "submissionsPostResponse"
  },
  "meta": {}
}`

func statusResponse(status string) string {
	return `{
  "data": {
    "attributes": {"createdDate": "2022-06-08T01:38:09.498Z", "name": "app.pkg", "status": "` + status + `"},
    "id": "` + testSubmissionID + `",
    "type": "submissions"
  },
  "meta": {}
}`
}

// jsonResponder serves body with a JSON content type, which resty needs to
// unmarshal the response into the SDK's structs.
func jsonResponder(body string) httpmock.Responder {
	return func(*http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, body)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	}
}

// writeTempPackage writes a small file to stand in for the package to upload.
func writeTempPackage(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "app.pkg")
	require.NoError(t, os.WriteFile(path, []byte("package bytes for notarization"), 0o600))
	return path
}

func TestSubmitAndWait_Accepted(t *testing.T) {
	var uploaded atomic.Bool
	s3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && r.URL.Path == "/EXAMPLE-BUCKET/EXAMPLE-KEY-NAME" {
			uploaded.Store(true)
			w.WriteHeader(http.StatusOK)
			return
		}
		t.Errorf("unexpected S3 request %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer s3.Close()

	c := newTestClient(t)
	pointUploadAt(t, s3)

	base := "https://appstoreconnect.apple.com/notary/v2/submissions"
	httpmock.RegisterResponder(http.MethodPost, base,
		jsonResponder(submitResponse))
	httpmock.RegisterResponder(http.MethodGet, base+"/"+testSubmissionID,
		jsonResponder(statusResponse(StatusAccepted)))

	var uploadedTotal int64
	res, err := c.SubmitAndWait(context.Background(), SubmitInput{
		FilePath:       writeTempPackage(t),
		Name:           "app.pkg",
		UploadProgress: func(written, total int64) { uploadedTotal = total },
	}, WaitOptions{Interval: time.Millisecond, Timeout: 2 * time.Second})

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, testSubmissionID, res.SubmissionID)
	require.NotNil(t, res.Status)
	assert.Equal(t, StatusAccepted, res.Status.Status)
	assert.True(t, res.Status.Done())
	assert.True(t, uploaded.Load(), "the file was not uploaded to S3")
	assert.Positive(t, uploadedTotal, "upload progress was not reported")
	assert.Nil(t, res.Issues)
}

func TestSubmitAndWait_RejectedReturnsLogIssues(t *testing.T) {
	s3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer s3.Close()

	// The developer log lives behind a pre-signed URL, which FetchLog GETs
	// directly (not through the notary API). A local server stands in for it.
	logSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"issues":[{"severity":"error","message":"unsigned binary","path":"Contents/MacOS/app","docUrl":"https://developer.apple.com/x"}]}`))
	}))
	defer logSrv.Close()

	c := newTestClient(t)
	pointUploadAt(t, s3)

	base := "https://appstoreconnect.apple.com/notary/v2/submissions"
	httpmock.RegisterResponder(http.MethodPost, base,
		jsonResponder(submitResponse))
	httpmock.RegisterResponder(http.MethodGet, base+"/"+testSubmissionID,
		jsonResponder(statusResponse(StatusInvalid)))
	logURLResponse := `{"data":{"attributes":{"developerLogUrl":"` + logSrv.URL + `/log.json"},"id":"` + testSubmissionID + `","type":"submissionsLog"},"meta":{}}`
	httpmock.RegisterResponder(http.MethodGet, base+"/"+testSubmissionID+"/logs",
		jsonResponder(logURLResponse))

	res, err := c.SubmitAndWait(context.Background(), SubmitInput{
		FilePath: writeTempPackage(t),
		Name:     "app.pkg",
	}, WaitOptions{Interval: time.Millisecond, Timeout: 2 * time.Second})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRejected), "want ErrRejected, got %v", err)
	require.NotNil(t, res)
	require.NotNil(t, res.Status)
	assert.Equal(t, StatusInvalid, res.Status.Status)
	require.Len(t, res.Issues, 1)
	assert.Equal(t, "error", res.Issues[0].Severity)
	assert.Equal(t, "unsigned binary", res.Issues[0].Message)
	assert.NotEmpty(t, res.Log)
}
