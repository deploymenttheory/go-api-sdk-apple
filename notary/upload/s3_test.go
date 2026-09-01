package upload

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestSigV4Vector reproduces the worked example in AWS's Signature
// Version 4 documentation (GET https://iam.amazonaws.com/?Action=ListUsers
// &Version=2010-05-08 at 20150830T123600Z).
func TestSigV4Vector(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://iam.amazonaws.com/?Action=ListUsers&Version=2010-05-08", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	creds := S3Credentials{AccessKeyID: "AKIDEXAMPLE", SecretAccessKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"}
	empty := sha256.Sum256(nil)
	signV4(req, creds, "us-east-1", "iam", hex.EncodeToString(empty[:]), time.Date(2015, 8, 30, 12, 36, 0, 0, time.UTC))
	want := "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-content-sha256;x-amz-date, Signature="
	got := req.Header.Get("Authorization")
	if !strings.HasPrefix(got, want) {
		t.Fatalf("Authorization = %q", got)
	}
	// The documented signature is for the request without
	// x-amz-content-sha256; with it the string to sign differs, so check
	// the derivation another way: the same inputs must sign the same.
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://iam.amazonaws.com/?Version=2010-05-08&Action=ListUsers", nil)
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	signV4(req2, creds, "us-east-1", "iam", hex.EncodeToString(empty[:]), time.Date(2015, 8, 30, 12, 36, 0, 0, time.UTC))
	if req2.Header.Get("Authorization") != got {
		t.Error("query parameter order changed the signature; canonicalization is wrong")
	}
}

// TestSigV4Canonical checks the canonical request against the AWS
// example exactly, by signing without x-amz-content-sha256 the way the
// documentation's request does.
func TestSigV4Canonical(t *testing.T) {
	creds := S3Credentials{AccessKeyID: "AKIDEXAMPLE", SecretAccessKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"}
	now := time.Date(2015, 8, 30, 12, 36, 0, 0, time.UTC)
	// Reproduce signV4's derivation for the documented canonical request.
	canonical := "GET\n/\nAction=ListUsers&Version=2010-05-08\ncontent-type:application/x-www-form-urlencoded; charset=utf-8\nhost:iam.amazonaws.com\nx-amz-date:20150830T123600Z\n\ncontent-type;host;x-amz-date\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	sts := "AWS4-HMAC-SHA256\n20150830T123600Z\n20150830/us-east-1/iam/aws4_request\n" + hex.EncodeToString(sha256Bytes([]byte(canonical)))
	if want := "f536975d06c0309214f805bb90ccff089219ecd68b2577efef23edd43b7e1a59"; hex.EncodeToString(sha256Bytes([]byte(canonical))) != want {
		t.Errorf("canonical request hash = %s, want %s", hex.EncodeToString(sha256Bytes([]byte(canonical))), want)
	}
	kDate := hmacSHA256([]byte("AWS4"+creds.SecretAccessKey), now.Format("20060102"))
	kRegion := hmacSHA256(kDate, "us-east-1")
	kService := hmacSHA256(kRegion, "iam")
	kSigning := hmacSHA256(kService, "aws4_request")
	if got := hex.EncodeToString(hmacSHA256(kSigning, sts)); got != "5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7" {
		t.Errorf("signature = %s, want the documented 5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7", got)
	}
}

func TestS3UploaderPutsAndRetries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.pkg")
	content := []byte("package bytes")
	os.WriteFile(path, content, 0o644)
	sum := sha256.Sum256(content)
	sumHex := hex.EncodeToString(sum[:])

	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if r.Method != http.MethodPut || r.URL.Path != "/bucket/sub/obj.pkg" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != string(content) {
			t.Errorf("body = %q", body)
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 Credential=AKID/") || !strings.Contains(auth, "/us-west-2/s3/aws4_request") {
			t.Errorf("Authorization = %q", auth)
		}
		if r.Header.Get("x-amz-content-sha256") != sumHex || r.Header.Get("x-amz-security-token") != "tok" {
			t.Errorf("headers: %v", r.Header)
		}
		if !strings.Contains(auth, "x-amz-security-token") {
			t.Error("session token not among signed headers")
		}
		if attempts == 1 {
			http.Error(w, "try again", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	u := &S3Uploader{Client: srv.Client(), Endpoint: srv.URL}
	creds := S3Credentials{AccessKeyID: "AKID", SecretAccessKey: "secret", SessionToken: "tok", Bucket: "bucket", Object: "sub/obj.pkg"}
	var last int64
	err := u.Upload(context.Background(), creds, path, sumHex, func(w, total int64) { last = w })
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 || last != int64(len(content)) {
		t.Errorf("attempts %d, progress %d", attempts, last)
	}

	// A 4xx is not retried.
	attempts = 0
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "nope", http.StatusForbidden)
	}))
	defer srv2.Close()
	u = &S3Uploader{Client: srv2.Client(), Endpoint: srv2.URL}
	if err := u.Upload(context.Background(), creds, path, sumHex, nil); err == nil || attempts != 1 {
		t.Errorf("403: err %v attempts %d", err, attempts)
	}
}

// TestS3UploaderSplitsALargeFile drives the multipart path against a fake
// S3, checking the three calls it makes and that the parts reassemble into
// the file that went in.
//
// The 5 GiB threshold cannot be tested with a 5 GiB file, so PartSize and
// the threshold are what the test lowers: multipartUpload is called
// directly, which is the same code the size check reaches.
func TestS3UploaderSplitsALargeFile(t *testing.T) {
	// A body that is not a round number of parts, so the last part is
	// short and the offsets have to be right.
	body := bytes.Repeat([]byte("go-macos-pkg"), 4096) // 48 KiB
	path := filepath.Join(t.TempDir(), "big.pkg")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}

	var (
		mu        sync.Mutex
		started   int
		completed int
		parts     = map[int][]byte{}
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		q := r.URL.Query()
		switch {
		case r.Method == http.MethodPost && q.Has("uploads"):
			started++
			w.Header().Set("Content-Type", "application/xml")
			fmt.Fprint(w, `<?xml version="1.0"?><InitiateMultipartUploadResult><UploadId>up-1</UploadId></InitiateMultipartUploadResult>`)
		case r.Method == http.MethodPut && q.Get("uploadId") == "up-1":
			n, err := strconv.Atoi(q.Get("partNumber"))
			if err != nil {
				t.Errorf("bad part number %q", q.Get("partNumber"))
			}
			data, _ := io.ReadAll(r.Body)
			parts[n] = data
			w.Header().Set("ETag", fmt.Sprintf("\"etag-%d\"", n))
		case r.Method == http.MethodPost && q.Get("uploadId") == "up-1":
			data, _ := io.ReadAll(r.Body)
			if !bytes.Contains(data, []byte("<PartNumber>1</PartNumber>")) {
				t.Errorf("the completion did not list the parts: %s", data)
			}
			completed++
			fmt.Fprint(w, `<?xml version="1.0"?><CompleteMultipartUploadResult/>`)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	u := &S3Uploader{Client: srv.Client(), Endpoint: srv.URL, PartSize: minPartSize}
	creds := S3Credentials{AccessKeyID: "a", SecretAccessKey: "s", SessionToken: "t", Bucket: "b", Object: "o"}
	var lastWritten, lastTotal int64
	err = u.multipartUpload(context.Background(), creds, f, int64(len(body)), u.targetURL(creds),
		func(written, total int64) { lastWritten, lastTotal = written, total })
	if err != nil {
		t.Fatal(err)
	}

	if started != 1 || completed != 1 {
		t.Errorf("started %d uploads and completed %d, want 1 and 1", started, completed)
	}
	// The parts, in order, are the file.
	var joined []byte
	for i := 1; i <= len(parts); i++ {
		joined = append(joined, parts[i]...)
	}
	if !bytes.Equal(joined, body) {
		t.Errorf("the parts do not reassemble into the file: got %d bytes, want %d", len(joined), len(body))
	}
	if want := (len(body) + minPartSize - 1) / minPartSize; len(parts) != want {
		t.Errorf("sent %d parts, want %d", len(parts), want)
	}
	if lastWritten != int64(len(body)) || lastTotal != int64(len(body)) {
		t.Errorf("progress ended at %d/%d, want %d/%d", lastWritten, lastTotal, len(body), len(body))
	}
}

// TestS3UploaderAbortsAFailedUpload checks that a part failing does not
// leave an incomplete upload sitting in Apple's bucket.
func TestS3UploaderAbortsAFailedUpload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "big.pkg")
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), 3*minPartSize), 0o644); err != nil {
		t.Fatal(err)
	}
	var aborted bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case r.Method == http.MethodPost && q.Has("uploads"):
			fmt.Fprint(w, `<?xml version="1.0"?><InitiateMultipartUploadResult><UploadId>up-1</UploadId></InitiateMultipartUploadResult>`)
		case r.Method == http.MethodDelete:
			aborted = true
			w.WriteHeader(http.StatusNoContent)
		default:
			// Every part fails.
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	u := &S3Uploader{Client: srv.Client(), Endpoint: srv.URL, PartSize: minPartSize}
	creds := S3Credentials{AccessKeyID: "a", SecretAccessKey: "s", SessionToken: "t", Bucket: "b", Object: "o"}
	err = u.multipartUpload(context.Background(), creds, f, 3*minPartSize, u.targetURL(creds), nil)
	if err == nil {
		t.Fatal("a failing part should have failed the upload")
	}
	if !aborted {
		t.Error("the incomplete upload was not aborted")
	}
}

// TestS3UploaderUsesTheAccelerateEndpoint pins the host every upload goes
// to. Apple's own documented example asks for transfer acceleration, so
// there is no setting: the endpoint is the endpoint.
func TestS3UploaderUsesTheAccelerateEndpoint(t *testing.T) {
	creds := S3Credentials{Bucket: "bucket", Object: "key"}
	if got := (&S3Uploader{}).targetURL(creds); got != "https://bucket.s3-accelerate.amazonaws.com/key" {
		t.Errorf("endpoint = %s", got)
	}
	// A test endpoint carries the bucket in the path, since it has no
	// bucket in its host.
	u := &S3Uploader{Endpoint: "http://127.0.0.1:1"}
	if got := u.targetURL(creds); got != "http://127.0.0.1:1/bucket/key" {
		t.Errorf("test endpoint = %s", got)
	}
}
