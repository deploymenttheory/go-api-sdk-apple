// Package upload uploads a submission's file to the S3 bucket Apple's notary
// service hands out, signed with AWS Signature Version 4 using the temporary
// credentials the submission response carries. It depends only on the
// standard library: there is no AWS SDK dependency.
//
// A file goes up in one PUT where S3 will take it, and in parts where it
// will not: S3 refuses a single request over 5 GiB, which a package can
// exceed.
//
// The region is not in Apple's documentation; us-west-2 is what every
// working client uses.
package upload

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// S3Region is where Apple's submission buckets live.
const S3Region = "us-west-2"

// maxSinglePut is S3's limit for one PUT. Anything larger has to go up in
// parts, which is what multipartUpload does.
const maxSinglePut = 5 << 30

// DefaultPartSize is how much of the file each part carries when a upload
// is split. S3 allows 5 MiB to 5 GiB per part and at most 10,000 parts, so
// 100 MiB covers a terabyte while keeping a failed part cheap to retry.
const DefaultPartSize = 100 << 20

// minPartSize is S3's floor for every part but the last.
const minPartSize = 5 << 20

// accelerateHost is S3's transfer acceleration endpoint, which routes the
// upload over Amazon's edge network rather than straight to the region.
//
// Every upload uses it, and there is no option not to. Apple's own
// documented example asks for it, building its S3 client with
// use_accelerate_endpoint: True, and notarytool turns it on as well. A
// setting whose right answer is always the same is a constant.
const accelerateHost = ".s3-accelerate.amazonaws.com"

// S3Credentials are the temporary credentials a submission returns.
type S3Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Bucket          string
	Object          string
}

// Uploader uploads a submission's file.
type Uploader interface {
	Upload(ctx context.Context, creds S3Credentials, path string, sha256Hex string, progress func(written, total int64)) error
}

// S3Uploader puts the file into the bucket directly.
type S3Uploader struct {
	Client *http.Client
	// Endpoint overrides the S3 endpoint (tests). Empty means AWS.
	Endpoint string
	// Now overrides the clock (tests).
	Now func() time.Time
	// PartSize is how much each part of a split upload carries. Zero
	// means DefaultPartSize.
	PartSize int64
}

// NewS3Uploader returns an uploader with a generous timeout: a package
// can be gigabytes.
func NewS3Uploader() *S3Uploader {
	return &S3Uploader{Client: &http.Client{Timeout: 4 * time.Hour}}
}

// Upload PUTs the file. sha256Hex is the file's SHA-256, which S3 wants
// as x-amz-content-sha256 and which the submission already computed.
func (u *S3Uploader) Upload(ctx context.Context, creds S3Credentials, path, sha256Hex string, progress func(written, total int64)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	target := u.targetURL(creds)
	if st.Size() > maxSinglePut {
		return u.multipartUpload(ctx, creds, f, st.Size(), target, progress)
	}

	const maxAttempts = 4
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return err
		}
		body := &progressReader{r: f, total: st.Size(), fn: progress}
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, target, body)
		if err != nil {
			return err
		}
		req.ContentLength = st.Size()
		req.Header.Set("Content-Type", "application/octet-stream")
		now := time.Now().UTC()
		if u.Now != nil {
			now = u.Now().UTC()
		}
		signV4(req, creds, S3Region, "s3", sha256Hex, now)

		client := u.Client
		if client == nil {
			client = http.DefaultClient
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
		} else {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("S3 returned %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
			if resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
				return fmt.Errorf("notary/upload: upload failed: %w", lastErr)
			}
		}
		if attempt < maxAttempts {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt*attempt) * 2 * time.Second):
			}
		}
	}
	return fmt.Errorf("notary/upload: upload failed after %d attempts: %w", maxAttempts, lastErr)
}

type progressReader struct {
	r       io.Reader
	total   int64
	written int64
	fn      func(written, total int64)
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.written += int64(n)
	if p.fn != nil && n > 0 {
		p.fn(p.written, p.total)
	}
	return n, err
}

// escapePath URI-encodes each path segment the way SigV4 canonicalizes
// it (RFC 3986 unreserved characters kept, everything else percent-encoded).
func escapePath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = awsEscape(part, false)
	}
	return strings.Join(parts, "/")
}

// awsEscape percent-encodes for SigV4: unreserved characters pass, and
// for a path segment "/" is kept out of the encoding by the caller.
func awsEscape(s string, encodeSlash bool) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9', c == '-', c == '_', c == '.', c == '~':
			b.WriteByte(c)
		case c == '/' && !encodeSlash:
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

// signV4 adds the Authorization, x-amz-date, x-amz-content-sha256 and
// x-amz-security-token headers per AWS Signature Version 4.
func signV4(req *http.Request, creds S3Credentials, region, service, payloadHash string, now time.Time) {
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	if creds.SessionToken != "" {
		req.Header.Set("x-amz-security-token", creds.SessionToken)
	}
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}

	// Canonical headers: host plus every x-amz-* and content-type, lower
	// case, sorted.
	headers := map[string]string{"host": host}
	for k, v := range req.Header {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-amz-") || lk == "content-type" {
			headers[lk] = strings.TrimSpace(strings.Join(v, ","))
		}
	}
	var names []string
	for k := range headers {
		names = append(names, k)
	}
	sort.Strings(names)
	var canonHeaders strings.Builder
	for _, k := range names {
		canonHeaders.WriteString(k + ":" + headers[k] + "\n")
	}
	signedHeaders := strings.Join(names, ";")

	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalPath(req.URL),
		canonicalQuery(req.URL),
		canonHeaders.String(),
		signedHeaders,
		payloadHash,
	}, "\n")

	scope := dateStamp + "/" + region + "/" + service + "/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		hex.EncodeToString(sha256Bytes([]byte(canonicalRequest))),
	}, "\n")

	kDate := hmacSHA256([]byte("AWS4"+creds.SecretAccessKey), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	signature := hex.EncodeToString(hmacSHA256(kSigning, stringToSign))

	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		creds.AccessKeyID, scope, signedHeaders, signature))
}

func canonicalPath(u *urlpkg.URL) string {
	p := u.EscapedPath()
	if p == "" {
		return "/"
	}
	// S3 expects the path as sent, encoded once; awsEscape above already
	// produced that form, so re-escaping would double-encode.
	return p
}

func canonicalQuery(u *urlpkg.URL) string {
	if u.RawQuery == "" {
		return ""
	}
	values, _ := urlpkg.ParseQuery(u.RawQuery)
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		vs := values[k]
		sort.Strings(vs)
		for _, v := range vs {
			parts = append(parts, awsEscape(k, true)+"="+awsEscape(v, true))
		}
	}
	return strings.Join(parts, "&")
}

func sha256Bytes(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func hmacSHA256(key []byte, data string) []byte {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(data))
	return m.Sum(nil)
}

// targetURL is where the object goes.
//
// A real upload uses a virtual-hosted URL, where the bucket is part of the
// host: S3's transfer acceleration endpoint, for the reason accelerateHost
// gives. A test endpoint has no bucket in its host, so the bucket goes in
// the path instead.
func (u *S3Uploader) targetURL(creds S3Credentials) string {
	if u.Endpoint != "" {
		return u.Endpoint + "/" + creds.Bucket + "/" + escapePath(creds.Object)
	}
	return "https://" + creds.Bucket + accelerateHost + "/" + escapePath(creds.Object)
}

// partSize is how much of the file each part carries.
func (u *S3Uploader) partSize(size int64) int64 {
	part := u.PartSize
	if part <= 0 {
		part = DefaultPartSize
	}
	if part < minPartSize {
		part = minPartSize
	}
	// S3 takes at most 10,000 parts, so a big enough file needs bigger
	// parts however small the caller asked for.
	if needed := (size + maxParts - 1) / maxParts; needed > part {
		part = needed
	}
	return part
}

// maxParts is S3's limit on how many parts one upload may have.
const maxParts = 10000

// completedPart is one finished part, as the completion request names it.
type completedPart struct {
	PartNumber int    `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

// multipartUpload sends a file too large for one request, in parts.
//
// Every part is signed on its own, so each is hashed before it is sent: the
// file is read twice, once to hash and once to upload, rather than held in
// memory. A failure aborts the upload, so a half-finished one does not sit
// in Apple's bucket accruing storage nobody asked for.
func (u *S3Uploader) multipartUpload(ctx context.Context, creds S3Credentials, f *os.File, size int64, target string, progress func(written, total int64)) error {
	uploadID, err := u.startMultipart(ctx, creds, target)
	if err != nil {
		return err
	}
	parts, err := u.uploadParts(ctx, creds, f, size, target, uploadID, progress)
	if err != nil {
		// Best effort: the upload has already failed, and the abort
		// failing too is not what the caller needs to hear about.
		if aerr := u.abortMultipart(ctx, creds, target, uploadID); aerr != nil {
			// Nothing to do but leave a trace for a reader of the logs.
			err = fmt.Errorf("%w (the incomplete upload could not be aborted: %v)", err, aerr)
		}
		return err
	}
	return u.completeMultipart(ctx, creds, target, uploadID, parts)
}

// uploadParts sends each part in turn and collects what completion needs.
func (u *S3Uploader) uploadParts(ctx context.Context, creds S3Credentials, f *os.File, size int64, target, uploadID string, progress func(written, total int64)) ([]completedPart, error) {
	part := u.partSize(size)
	var parts []completedPart
	var sent int64
	for offset, number := int64(0), 1; offset < size; number++ {
		length := part
		if remaining := size - offset; remaining < length {
			length = remaining
		}
		sum, err := hashSection(f, offset, length)
		if err != nil {
			return nil, err
		}
		etag, err := u.uploadPart(ctx, creds, f, offset, length, target, uploadID, number, sum, func(written int64) {
			if progress != nil {
				progress(sent+written, size)
			}
		})
		if err != nil {
			return nil, fmt.Errorf("notary/upload: part %d: %w", number, err)
		}
		parts = append(parts, completedPart{PartNumber: number, ETag: etag})
		offset += length
		sent += length
	}
	return parts, nil
}

// hashSection returns the SHA-256 of a stretch of the file, which is what
// signing a part needs.
func hashSection(f *os.File, offset, length int64) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, io.NewSectionReader(f, offset, length)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// startMultipart asks S3 to open an upload and returns its identifier.
func (u *S3Uploader) startMultipart(ctx context.Context, creds S3Credentials, target string) (string, error) {
	body, err := u.do(ctx, creds, http.MethodPost, target+"?uploads=", nil, 0, emptyPayloadHash)
	if err != nil {
		return "", fmt.Errorf("notary/upload: starting a multipart upload: %w", err)
	}
	var out struct {
		UploadID string `xml:"UploadId"`
	}
	if err := xml.Unmarshal(body, &out); err != nil || out.UploadID == "" {
		return "", fmt.Errorf("notary/upload: S3 did not return an upload id: %s", strings.TrimSpace(string(body)))
	}
	return out.UploadID, nil
}

// uploadPart sends one part and returns the ETag completion will quote.
func (u *S3Uploader) uploadPart(ctx context.Context, creds S3Credentials, f *os.File, offset, length int64, target, uploadID string, number int, sum string, progress func(int64)) (string, error) {
	url := fmt.Sprintf("%s?partNumber=%d&uploadId=%s", target, number, urlpkg.QueryEscape(uploadID))
	section := io.NewSectionReader(f, offset, length)
	body := &progressReader{r: section, total: length, fn: func(written, _ int64) { progress(written) }}
	resp, err := u.doRequest(ctx, creds, http.MethodPut, url, body, length, sum)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10)); err != nil {
		return "", err
	}
	etag := resp.Header.Get("ETag")
	if etag == "" {
		return "", fmt.Errorf("S3 returned no ETag for part %d", number)
	}
	return etag, nil
}

// completeMultipart tells S3 the parts are all there.
func (u *S3Uploader) completeMultipart(ctx context.Context, creds S3Credentials, target, uploadID string, parts []completedPart) error {
	payload, err := xml.Marshal(struct {
		XMLName xml.Name        `xml:"CompleteMultipartUpload"`
		Parts   []completedPart `xml:"Part"`
	}{Parts: parts})
	if err != nil {
		return err
	}
	sum := sha256.Sum256(payload)
	url := target + "?uploadId=" + urlpkg.QueryEscape(uploadID)
	body, err := u.do(ctx, creds, http.MethodPost, url, bytes.NewReader(payload), int64(len(payload)), hex.EncodeToString(sum[:]))
	if err != nil {
		return fmt.Errorf("notary/upload: completing the upload: %w", err)
	}
	// S3 can report a failure inside a 200 response on this call alone,
	// so the body has to be read rather than trusted.
	if bytes.Contains(body, []byte("<Error")) {
		return fmt.Errorf("notary/upload: completing the upload: %s", strings.TrimSpace(string(body)))
	}
	return nil
}

// abortMultipart discards an upload that will not be completed.
func (u *S3Uploader) abortMultipart(ctx context.Context, creds S3Credentials, target, uploadID string) error {
	_, err := u.do(ctx, creds, http.MethodDelete, target+"?uploadId="+urlpkg.QueryEscape(uploadID), nil, 0, emptyPayloadHash)
	return err
}

// emptyPayloadHash is the SHA-256 of no bytes, which S3 wants on a request
// that carries no body.
const emptyPayloadHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// do performs a signed request and returns its body.
func (u *S3Uploader) do(ctx context.Context, creds S3Credentials, method, url string, body io.Reader, length int64, payloadHash string) ([]byte, error) {
	resp, err := u.doRequest(ctx, creds, method, url, body, length, payloadHash)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}

// doRequest signs and sends one request, and fails on any status that is
// not a success. The caller closes the body.
func (u *S3Uploader) doRequest(ctx context.Context, creds S3Credentials, method, url string, body io.Reader, length int64, payloadHash string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.ContentLength = length
	if length > 0 {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	now := time.Now().UTC()
	if u.Now != nil {
		now = u.Now().UTC()
	}
	signV4(req, creds, S3Region, "s3", payloadHash, now)

	client := u.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		resp.Body.Close()
		return nil, fmt.Errorf("S3 returned %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}
	return resp, nil
}
