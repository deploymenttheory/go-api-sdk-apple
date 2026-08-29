package cdn

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a CDN transport with httpmock enabled.
func setupMockClient(t *testing.T) *CDNService {
	t.Helper()

	transport, err := client.NewTransport(
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewService(transport)
}

const (
	testIPSWURL     = "https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw"
	testSHA1        = "03078f4af82bff5473398ca49f99288c76253fe8"
	testSHA256      = "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90"
	testContentType = "application/x-apple-aspen-config"
)

const testContentLength = int64(19734779897)

// =============================================================================
// ParseURL (package-level function)
// =============================================================================

func TestParseURL_Success(t *testing.T) {
	info, err := ParseURL(testIPSWURL)

	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, "2026WinterFCS", info.CatalogRelease)
	assert.Equal(t, "fullrestores", info.AssetType)
	assert.Equal(t, "122-28781", info.AssetID)
	assert.Equal(t, "DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46", info.UUID)
	assert.Equal(t, "UniversalMac_26.4.1_25E253_Restore.ipsw", info.Filename)
	assert.Equal(t, "UniversalMac", info.Platform)
	assert.Equal(t, "26.4.1", info.Version)
	assert.Equal(t, "25E253", info.Build)
	assert.Equal(t, "Restore", info.RestoreType)
}

func TestParseURL_EmptyURLError(t *testing.T) {
	_, err := ParseURL("")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "URL is required")
}

func TestParseURL_WrongHostError(t *testing.T) {
	_, err := ParseURL("https://example.com/some/path/file.ipsw")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the Apple CDN host")
}

func TestParseURL_InvalidURLError(t *testing.T) {
	_, err := ParseURL("://invalid url")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid URL")
}

func TestParseURL_TooFewSegmentsError(t *testing.T) {
	_, err := ParseURL("https://updates.cdn-apple.com/onlyone")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected URL path structure")
}

func TestParseURL_NotIPSWExtensionError(t *testing.T) {
	_, err := ParseURL("https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.zip")

	require.Error(t, err)
	assert.Contains(t, err.Error(), ".ipsw")
}

func TestParseURL_InvalidFilenameFormatError(t *testing.T) {
	// Filename has only 2 underscore parts instead of 4.
	_, err := ParseURL("https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_Restore.ipsw")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse filename")
}

// CDNService.ParseURL should delegate to the package-level ParseURL.
func TestCDNServiceParseURL_DelegatesToPackageLevel(t *testing.T) {
	svc := setupMockClient(t)

	info, err := svc.ParseURL(testIPSWURL)

	require.NoError(t, err)
	assert.Equal(t, "26.4.1", info.Version)
}

// =============================================================================
// GetFileMetadataV1
// =============================================================================

func TestGetFileMetadataV1_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("HEAD", testIPSWURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Content-Type", testContentType)
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", testContentLength))
			resp.Header.Set("x-amz-meta-digest-sh1", testSHA1)
			resp.Header.Set("x-amz-meta-digest-sha256", testSHA256)
			resp.Header.Set("Last-Modified", "Mon, 27 Jan 2025 12:00:00 GMT")
			resp.Header.Set("Etag", `"abc123etag"`)
			return resp, nil
		})

	meta, resp, err := svc.GetFileMetadataV1(context.Background(), testIPSWURL)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, meta)
	assert.Equal(t, testIPSWURL, meta.URL)
	assert.Equal(t, testContentLength, meta.ContentLength)
	assert.Equal(t, testSHA1, meta.SHA1)
	assert.Equal(t, testSHA256, meta.SHA256)
	assert.Equal(t, testContentType, meta.ContentType)
	assert.Equal(t, "abc123etag", meta.ETag)
	assert.False(t, meta.LastModified.IsZero(), "LastModified should be parsed")
}

func TestGetFileMetadataV1_ETagStripsQuotes(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("HEAD", testIPSWURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Etag", `"quoted-etag-value"`)
			return resp, nil
		})

	meta, _, err := svc.GetFileMetadataV1(context.Background(), testIPSWURL)

	require.NoError(t, err)
	assert.Equal(t, "quoted-etag-value", meta.ETag, "ETag quotes should be stripped")
}

func TestGetFileMetadataV1_EmptyURLError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.GetFileMetadataV1(context.Background(), "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "URL is required")
}

func TestGetFileMetadataV1_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("HEAD", testIPSWURL,
		httpmock.NewStringResponder(403, "Forbidden"))

	_, resp, err := svc.GetFileMetadataV1(context.Background(), testIPSWURL)

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, err.Error(), "403")
}

func TestGetFileMetadataV1_MissingHeaders(t *testing.T) {
	svc := setupMockClient(t)

	// Respond with 200 but no metadata headers — should not panic.
	httpmock.RegisterResponder("HEAD", testIPSWURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})

	meta, _, err := svc.GetFileMetadataV1(context.Background(), testIPSWURL)

	require.NoError(t, err)
	assert.Equal(t, int64(0), meta.ContentLength)
	assert.Empty(t, meta.SHA1)
	assert.Empty(t, meta.SHA256)
	assert.True(t, meta.LastModified.IsZero())
}

// =============================================================================
// DownloadFileV1
// =============================================================================

// headAndGetResponder registers both HEAD and GET responders for url.
// The HEAD response carries the checksum headers; the GET response body
// contains content.
func headAndGetResponder(url string, content []byte) {
	sha1sum := sha1.Sum(content)
	sha256sum := sha256.Sum256(content)

	httpmock.RegisterResponder("HEAD", url,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(content)))
			resp.Header.Set("x-amz-meta-digest-sh1", hex.EncodeToString(sha1sum[:]))
			resp.Header.Set("x-amz-meta-digest-sha256", hex.EncodeToString(sha256sum[:]))
			resp.Header.Set("Content-Type", "application/x-apple-aspen-config")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, content)
			resp.Header.Set("Content-Type", "application/x-apple-aspen-config")
			return resp, nil
		})
}

func TestDownloadFileV1_Success(t *testing.T) {
	svc := setupMockClient(t)

	content := []byte("fake ipsw file content for testing")
	headAndGetResponder(testIPSWURL, content)

	destPath := t.TempDir() + "/test.ipsw"
	result, resp, err := svc.DownloadFileV1(context.Background(), testIPSWURL, destPath, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)
	assert.Equal(t, testIPSWURL, result.URL)
	assert.Equal(t, destPath, result.DestPath)
	assert.Equal(t, int64(len(content)), result.BytesWritten)
	assert.True(t, result.Verified, "checksums should match")
	assert.False(t, result.Duration == 0)

	// Verify file contents on disk match what was served.
	got, err := os.ReadFile(destPath)
	require.NoError(t, err)
	assert.Equal(t, content, got)
}

func TestDownloadFileV1_ChecksumsCorrect(t *testing.T) {
	svc := setupMockClient(t)

	content := []byte("checksum verification test content")
	headAndGetResponder(testIPSWURL, content)

	destPath := t.TempDir() + "/verify.ipsw"
	result, _, err := svc.DownloadFileV1(context.Background(), testIPSWURL, destPath, nil)

	require.NoError(t, err)

	expectedSHA1 := sha1.Sum(content)
	expectedSHA256 := sha256.Sum256(content)
	assert.Equal(t, hex.EncodeToString(expectedSHA1[:]), result.SHA1)
	assert.Equal(t, hex.EncodeToString(expectedSHA256[:]), result.SHA256)
}

func TestDownloadFileV1_ProgressCallbackInvoked(t *testing.T) {
	svc := setupMockClient(t)

	content := []byte("progress callback test content")
	headAndGetResponder(testIPSWURL, content)

	var lastWritten, lastTotal int64
	progressFn := func(written, total int64) {
		lastWritten = written
		lastTotal = total
	}

	destPath := t.TempDir() + "/progress.ipsw"
	_, _, err := svc.DownloadFileV1(context.Background(), testIPSWURL, destPath, progressFn)

	require.NoError(t, err)
	assert.Equal(t, int64(len(content)), lastWritten, "final progress should equal file size")
	assert.Equal(t, int64(len(content)), lastTotal, "total should match Content-Length from HEAD")
}

func TestDownloadFileV1_ChecksumMismatchDeletesFile(t *testing.T) {
	svc := setupMockClient(t)

	content := []byte("real content")
	badContent := []byte("tampered content")
	sha1sum := sha1.Sum(badContent)
	sha256sum := sha256.Sum256(badContent)

	// HEAD reports checksums for badContent, but GET delivers real content.
	httpmock.RegisterResponder("HEAD", testIPSWURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(content)))
			resp.Header.Set("x-amz-meta-digest-sh1", hex.EncodeToString(sha1sum[:]))
			resp.Header.Set("x-amz-meta-digest-sha256", hex.EncodeToString(sha256sum[:]))
			return resp, nil
		})
	httpmock.RegisterResponder("GET", testIPSWURL,
		httpmock.NewBytesResponder(200, content))

	destPath := t.TempDir() + "/mismatch.ipsw"
	_, _, err := svc.DownloadFileV1(context.Background(), testIPSWURL, destPath, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")

	// Partial file must be cleaned up.
	_, statErr := os.Stat(destPath)
	assert.True(t, os.IsNotExist(statErr), "file should be removed after checksum mismatch")
}

func TestDownloadFileV1_EmptyURLError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.DownloadFileV1(context.Background(), "", t.TempDir()+"/x.ipsw", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "URL is required")
}

func TestDownloadFileV1_EmptyDestPathError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.DownloadFileV1(context.Background(), testIPSWURL, "", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "destination path is required")
}

func TestDownloadFileV1_HTTPErrorOnGet(t *testing.T) {
	svc := setupMockClient(t)

	content := []byte("data")
	sha1sum := sha1.Sum(content)
	sha256sum := sha256.Sum256(content)

	httpmock.RegisterResponder("HEAD", testIPSWURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "")
			resp.Header.Set("x-amz-meta-digest-sh1", hex.EncodeToString(sha1sum[:]))
			resp.Header.Set("x-amz-meta-digest-sha256", hex.EncodeToString(sha256sum[:]))
			return resp, nil
		})
	httpmock.RegisterResponder("GET", testIPSWURL,
		httpmock.NewStringResponder(403, "Forbidden"))

	destPath := t.TempDir() + "/forbidden.ipsw"
	_, _, err := svc.DownloadFileV1(context.Background(), testIPSWURL, destPath, nil)

	require.Error(t, err)

	// File must not be left behind on server error.
	_, statErr := os.Stat(destPath)
	assert.True(t, os.IsNotExist(statErr), "file should be removed after failed download")
}

func TestDownloadFileV1_CreatesDestinationDirectory(t *testing.T) {
	svc := setupMockClient(t)

	content := []byte("mkdir test")
	headAndGetResponder(testIPSWURL, content)

	// Destination inside a subdirectory that does not yet exist.
	destPath := t.TempDir() + "/nested/dir/test.ipsw"
	_, _, err := svc.DownloadFileV1(context.Background(), testIPSWURL, destPath, nil)

	require.NoError(t, err)
	_, statErr := os.Stat(destPath)
	assert.NoError(t, statErr, "destination file should exist")
}

// =============================================================================
// CDN constants smoke tests
// =============================================================================

func TestCDNConstants(t *testing.T) {
	assert.Equal(t, "updates.cdn-apple.com", CDNHost)
	assert.Equal(t, "x-amz-meta-digest-sh1", HeaderSHA1)
	assert.Equal(t, "x-amz-meta-digest-sha256", HeaderSHA256)
	assert.Equal(t, ".ipsw", IPSWExtension)
	assert.Equal(t, 5, ExpectedPathSegments)
	assert.Equal(t, 4, ExpectedFilenameSegments)
}
