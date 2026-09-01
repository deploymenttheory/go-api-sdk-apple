package notary

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary/notary_api/submissions"
	"github.com/deploymenttheory/go-sdk-appleservices/notary/upload"
)

// Submission statuses, as Apple's notary service reports them.
const (
	StatusAccepted   = "Accepted"
	StatusInProgress = "In Progress"
	StatusInvalid    = "Invalid"
	StatusRejected   = "Rejected"
)

// WebhookChannel is the only notification channel Apple's notary service
// defines.
const WebhookChannel = "webhook"

// ErrRejected is returned when a submission reaches a terminal status other
// than Accepted.
var ErrRejected = errors.New("notary: submission was not accepted")

// ErrTimeout is returned when the wait deadline passes before Apple returns a
// verdict.
var ErrTimeout = errors.New("notary: timed out waiting for a verdict")

// Status is the state of a submission.
type Status struct {
	ID          string
	Name        string
	Status      string
	CreatedDate string
}

// Done reports whether Apple has finished with the submission.
func (s *Status) Done() bool { return s.Status != StatusInProgress && s.Status != "" }

// FileSHA256 hashes a file and returns the lowercase hexadecimal digest, which
// is what a submission and the S3 upload both need.
func FileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// WaitOptions configures Wait.
type WaitOptions struct {
	Interval time.Duration // default 30s
	Timeout  time.Duration // default 30m
	// Progress, when set, is called after each poll with the latest status.
	Progress func(*Status)
}

// Wait polls the submission until it is done or the timeout passes. A finished
// submission that Apple did not accept returns ErrRejected with the status; a
// deadline reached first returns ErrTimeout.
func Wait(ctx context.Context, c *Client, id string, o WaitOptions) (*Status, error) {
	if o.Interval <= 0 {
		o.Interval = 30 * time.Second
	}
	if o.Timeout <= 0 {
		o.Timeout = 30 * time.Minute
	}
	deadline := time.Now().Add(o.Timeout)
	var last *Status
	for {
		resp, _, err := c.NotaryAPI.Submissions.GetSubmissionStatusV2(ctx, id)
		if err != nil {
			return last, fmt.Errorf("notary: status: %w", err)
		}
		st := &Status{
			ID:          resp.Data.ID,
			Name:        resp.Data.Attributes.Name,
			Status:      resp.Data.Attributes.Status,
			CreatedDate: resp.Data.Attributes.CreatedDate,
		}
		last = st
		if o.Progress != nil {
			o.Progress(st)
		}
		if st.Done() {
			if st.Status != StatusAccepted {
				return st, fmt.Errorf("%w: %s", ErrRejected, st.Status)
			}
			return st, nil
		}
		if time.Now().After(deadline) {
			return st, ErrTimeout
		}
		select {
		case <-ctx.Done():
			return st, ctx.Err()
		case <-time.After(o.Interval):
		}
	}
}

// LogIssue is one entry from the issues list of a developer log.
type LogIssue struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Path     string `json:"path"`
	Code     any    `json:"code"`
	DocURL   string `json:"docUrl"`
}

// FetchLog downloads the developer log, a JSON document, for a submission.
//
// It asks the notary API for the temporary log URL, then GETs it with
// httpClient (a plain client, since the URL is a pre-signed S3 link that does
// not want the notary API's authentication). A nil httpClient uses a default
// with a 60-second timeout. The body is capped at 16 MiB and validated as
// JSON.
func FetchLog(ctx context.Context, c *Client, httpClient *http.Client, id string) (json.RawMessage, error) {
	resp, _, err := c.NotaryAPI.Submissions.GetSubmissionLogV2(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("notary: log: %w", err)
	}
	u := resp.Data.Attributes.DeveloperLogURL
	if u == "" {
		return nil, errors.New("notary: no log is available for the submission yet")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}
	logResp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notary: log download: %w", err)
	}
	defer logResp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(logResp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	if logResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notary: log download returned %s", logResp.Status)
	}
	if !json.Valid(data) {
		return nil, errors.New("notary: the log is not JSON")
	}
	return json.RawMessage(data), nil
}

// ParseLogIssues pulls the issues out of a developer log for display.
func ParseLogIssues(log json.RawMessage) []LogIssue {
	var doc struct {
		Issues []LogIssue `json:"issues"`
	}
	_ = json.Unmarshal(log, &doc)
	return doc.Issues
}

// newS3Uploader builds the uploader SubmitAndWait uses. It is a variable so
// tests can point the upload at a local server; production always gets the
// real transfer-acceleration uploader.
var newS3Uploader = upload.NewS3Uploader

// SubmitInput describes a file to notarize.
type SubmitInput struct {
	// FilePath is the package (or zip, dmg) to notarize.
	FilePath string
	// Name is the submission name Apple records and shows in the log.
	Name string
	// Webhook, when set, is a public URL Apple posts the verdict to when
	// notarization finishes, so a job need not sit and poll. Apple's own
	// documentation warns it is best effort.
	Webhook string
	// UploadProgress, when set, is called as the file goes up to S3.
	UploadProgress func(written, total int64)
}

// Result is the outcome of SubmitAndWait.
type Result struct {
	// SubmissionID is the id Apple assigned, usable to fetch the log later.
	SubmissionID string
	// Status is the final status polling saw.
	Status *Status
	// Log is the developer log, populated only when the submission was not
	// accepted (and the log could be fetched).
	Log json.RawMessage
	// Issues are the parsed issues from Log, when there are any.
	Issues []LogIssue
}

// SubmitAndWait notarizes a file end to end: it hashes the file, registers the
// submission, uploads the file to the S3 bucket Apple hands out, and polls
// until Apple returns a verdict.
//
// On a verdict of Accepted it returns a Result with the final Status and a nil
// error. On any other terminal status it fetches and parses the developer log,
// populates the Result's Log and Issues, and returns ErrRejected (wrapped)
// with that Result. The Result is non-nil whenever the submission was
// registered, so a caller can always reach the SubmissionID.
func (c *Client) SubmitAndWait(ctx context.Context, in SubmitInput, wo WaitOptions) (*Result, error) {
	sum, err := FileSHA256(in.FilePath)
	if err != nil {
		return nil, fmt.Errorf("notary: hashing %s: %w", in.FilePath, err)
	}

	req := &submissions.NewSubmissionRequest{
		SHA256:         sum,
		SubmissionName: in.Name,
	}
	if in.Webhook != "" {
		req.Notifications = []submissions.NewSubmissionRequestNotification{
			{Channel: WebhookChannel, Target: in.Webhook},
		}
	}
	sub, _, err := c.NotaryAPI.Submissions.SubmitSoftwareV2(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("notary: submit: %w", err)
	}
	result := &Result{SubmissionID: sub.Data.ID}

	a := sub.Data.Attributes
	creds := upload.S3Credentials{
		AccessKeyID:     a.AWSAccessKeyID,
		SecretAccessKey: a.AWSSecretAccessKey,
		SessionToken:    a.AWSSessionToken,
		Bucket:          a.Bucket,
		Object:          a.Object,
	}
	if err := newS3Uploader().Upload(ctx, creds, in.FilePath, sum, in.UploadProgress); err != nil {
		return result, err
	}

	st, err := Wait(ctx, c, result.SubmissionID, wo)
	result.Status = st
	if err != nil {
		// A rejection is worth the log: fetch and parse it, but do not let a
		// log failure mask the verdict the caller came for.
		if errors.Is(err, ErrRejected) {
			if log, logErr := FetchLog(ctx, c, nil, result.SubmissionID); logErr == nil {
				result.Log = log
				result.Issues = ParseLogIssues(log)
			}
		}
		return result, err
	}
	return result, nil
}
