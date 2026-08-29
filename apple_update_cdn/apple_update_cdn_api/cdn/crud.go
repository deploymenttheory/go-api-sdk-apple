package cdn

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/client"
	"resty.dev/v3"
)

// CDNService handles Apple CDN URL parsing and file metadata resolution.
//
// Apple's firmware CDN (updates.cdn-apple.com) is S3-backed. It does not expose
// a directory listing, but individual IPSW files can be inspected via HEAD
// requests which return SHA-1 and SHA-256 checksums alongside file size.
type CDNService struct {
	client client.Client
}

// NewService creates a new CDN service.
func NewService(c client.Client) *CDNService {
	return &CDNService{client: c}
}

// ParseURL parses an Apple CDN IPSW URL into its structural components.
// This is a pure parsing operation — no HTTP request is made.
//
// Expected URL format:
//
//	https://updates.cdn-apple.com/{catalogRelease}/{assetType}/{assetID}/{uuid}/{platform}_{version}_{build}_{restoreType}.ipsw
//
// Example:
//
//	https://updates.cdn-apple.com/2026WinterFCS/fullrestores/122-28781/DCB2FF13-06CB-44C2-BCA2-DFCAF3521D46/UniversalMac_26.4.1_25E253_Restore.ipsw
func (s *CDNService) ParseURL(rawURL string) (*CDNURLInfo, error) {
	return ParseURL(rawURL)
}

// ParseURL is the package-level URL parser, callable without a CDNService instance.
// See CDNService.ParseURL for full documentation.
func ParseURL(rawURL string) (*CDNURLInfo, error) {
	if rawURL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Host != CDNHost {
		return nil, fmt.Errorf("URL host %q is not the Apple CDN host %q", parsed.Host, CDNHost)
	}

	// Strip leading slash and split into segments.
	segments := strings.SplitN(strings.TrimPrefix(parsed.Path, "/"), "/", ExpectedPathSegments)
	if len(segments) != ExpectedPathSegments {
		return nil, fmt.Errorf("unexpected URL path structure: expected %d segments, got %d", ExpectedPathSegments, len(segments))
	}

	catalogRelease := segments[0]
	assetType := segments[1]
	assetID := segments[2]
	uuid := segments[3]
	filename := segments[4]

	if !strings.HasSuffix(filename, IPSWExtension) {
		return nil, fmt.Errorf("filename %q does not have %s extension", filename, IPSWExtension)
	}

	platform, version, build, restoreType, err := parseFilename(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filename %q: %w", filename, err)
	}

	return &CDNURLInfo{
		CatalogRelease: catalogRelease,
		AssetType:      assetType,
		AssetID:        assetID,
		UUID:           uuid,
		Filename:       filename,
		Platform:       platform,
		Version:        version,
		Build:          build,
		RestoreType:    restoreType,
	}, nil
}

// GetFileMetadataV1 issues a HEAD request to the Apple CDN and returns file
// metadata extracted from the response headers. The full IPSW file is not
// downloaded.
//
// Apple's CDN returns the following useful headers:
//   - Content-Length: file size in bytes
//   - x-amz-meta-digest-sh1: SHA-1 checksum
//   - x-amz-meta-digest-sha256: SHA-256 checksum
//   - Last-Modified: upload timestamp
//   - ETag: entity tag
func (s *CDNService) GetFileMetadataV1(ctx context.Context, rawURL string) (*CDNFileMetadata, *resty.Response, error) {
	if rawURL == "" {
		return nil, nil, fmt.Errorf("URL is required")
	}

	resp, err := s.client.NewRequest(ctx).Head(rawURL)
	if err != nil {
		return nil, resp, err
	}

	metadata := &CDNFileMetadata{
		URL:         rawURL,
		ContentType: resp.Header().Get("Content-Type"),
		ETag:        strings.Trim(resp.Header().Get("Etag"), `"`),
		SHA1:        resp.Header().Get(HeaderSHA1),
		SHA256:      resp.Header().Get(HeaderSHA256),
	}

	if cl := resp.Header().Get("Content-Length"); cl != "" {
		metadata.ContentLength, _ = strconv.ParseInt(cl, 10, 64)
	}

	if lm := resp.Header().Get("Last-Modified"); lm != "" {
		metadata.LastModified, _ = http.ParseTime(lm)
	}

	return metadata, resp, nil
}

// DownloadFileV1 downloads an IPSW file from the Apple CDN to destPath.
//
// The method first issues a HEAD request to obtain the expected file size and
// checksums, then streams the GET response body directly to disk — the full
// file is never held in memory. After the download completes, both SHA-1 and
// SHA-256 are verified against the values returned by the CDN. DownloadResult
// reports whether verification passed.
//
// progressFn is called on each write with the cumulative bytes written and the
// expected total size. Pass nil to disable progress reporting.
//
// If destPath already exists it is overwritten. If the download or checksum
// verification fails the partially-written file is removed.
//
// Apple CDN IPSW files are typically 15–22 GB. Ensure sufficient disk space
// before calling this method.
func (s *CDNService) DownloadFileV1(ctx context.Context, rawURL, destPath string, progressFn ProgressFunc) (*DownloadResult, *resty.Response, error) {
	if rawURL == "" {
		return nil, nil, fmt.Errorf("URL is required")
	}
	if destPath == "" {
		return nil, nil, fmt.Errorf("destination path is required")
	}

	// HEAD first — obtain expected size and checksums before touching disk.
	meta, _, err := s.GetFileMetadataV1(ctx, rawURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve file metadata: %w", err)
	}

	// Ensure the destination directory exists.
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return nil, nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create destination file: %w", err)
	}

	// Build a multi-writer: file + SHA-1 + SHA-256 + optional progress.
	sha1h := sha1.New()
	sha256h := sha256.New()

	writers := []io.Writer{f, sha1h, sha256h}
	if progressFn != nil {
		writers = append(writers, &progressWriter{
			total:      meta.ContentLength,
			progressFn: progressFn,
		})
	}

	mw := io.MultiWriter(writers...)

	start := time.Now()
	resp, n, err := s.client.NewRequest(ctx).Download(rawURL, mw)

	// Always close the file; clean up on any error.
	f.Close()
	if err != nil {
		os.Remove(destPath)
		return nil, resp, fmt.Errorf("download failed: %w", err)
	}

	actualSHA1 := hex.EncodeToString(sha1h.Sum(nil))
	actualSHA256 := hex.EncodeToString(sha256h.Sum(nil))

	// Verified = true only when ≥1 checksum was provided and all matched.
	verified := meta.SHA1 != "" || meta.SHA256 != ""
	if meta.SHA1 != "" && !strings.EqualFold(meta.SHA1, actualSHA1) {
		verified = false
	}
	if meta.SHA256 != "" && !strings.EqualFold(meta.SHA256, actualSHA256) {
		verified = false
	}

	if !verified && (meta.SHA1 != "" || meta.SHA256 != "") {
		os.Remove(destPath)
		return nil, resp, fmt.Errorf("checksum mismatch: expected sha1=%s sha256=%s, got sha1=%s sha256=%s",
			meta.SHA1, meta.SHA256, actualSHA1, actualSHA256)
	}

	return &DownloadResult{
		URL:          rawURL,
		DestPath:     destPath,
		BytesWritten: n,
		SHA1:         actualSHA1,
		SHA256:       actualSHA256,
		Duration:     time.Since(start),
		Verified:     verified,
	}, resp, nil
}

// progressWriter wraps an optional ProgressFunc, forwarding write counts
// without buffering any data.
type progressWriter struct {
	written    int64
	total      int64
	progressFn ProgressFunc
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.written += int64(len(p))
	pw.progressFn(pw.written, pw.total)
	return len(p), nil
}

// parseFilename extracts platform, version, build, and restore type from an
// IPSW filename of the form:
//
//	{Platform}_{Version}_{Build}_{RestoreType}.ipsw
//
// Example: "UniversalMac_26.4.1_25E253_Restore.ipsw"
func parseFilename(filename string) (platform, version, build, restoreType string, err error) {
	base := strings.TrimSuffix(filename, IPSWExtension)
	parts := strings.SplitN(base, "_", ExpectedFilenameSegments)
	if len(parts) != ExpectedFilenameSegments {
		return "", "", "", "", fmt.Errorf("expected %d underscore-delimited parts, got %d", ExpectedFilenameSegments, len(parts))
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

