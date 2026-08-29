package search

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/itunes/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates an iTunes transport with httpmock enabled.
func setupMockClient(t *testing.T) *SearchService {
	t.Helper()

	transport, err := client.NewTransport(
		client.WithLogger(zap.NewNop()),
		client.WithRetryCount(0), // disable retries for tests
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(transport.GetHTTPClient().Client())
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewService(transport)
}

// mockSearchBody returns a minimal SearchResponse JSON payload.
func mockSearchBody(resultCount int, results []Result) string {
	body, _ := json.Marshal(SearchResponse{ResultCount: resultCount, Results: results})
	return string(body)
}

// jsonResponder wraps a body string in an HTTP response with Content-Type: application/json.
// resty v3 only auto-unmarshals SetResult targets when the response content type is JSON.
func jsonResponder(status int, body string) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(status, body)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	}
}

// =============================================================================
// SearchV1
// =============================================================================

func TestSearchV1_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://itunes.apple.com/search",
		jsonResponder(200, mockSearchBody(2, []Result{
			{WrapperType: "track", Kind: "song", ArtistID: 909253, ArtistName: "Jack Johnson", TrackName: "Flake"},
			{WrapperType: "track", Kind: "song", ArtistID: 909253, ArtistName: "Jack Johnson", TrackName: "Banana Pancakes"},
		})))

	result, resp, err := svc.SearchV1(context.Background(), &SearchOptions{
		Term:  "Jack Johnson",
		Media: MediaMusic,
		Limit: 2,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ResultCount)
	assert.Len(t, result.Results, 2)
	assert.Equal(t, "Jack Johnson", result.Results[0].ArtistName)
	assert.Equal(t, "Flake", result.Results[0].TrackName)
}

func TestSearchV1_QueryParamsForwarded(t *testing.T) {
	svc := setupMockClient(t)

	var capturedQuery map[string]string
	httpmock.RegisterResponder("GET", `=~^https://itunes\.apple\.com/search`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = make(map[string]string)
			for k, v := range req.URL.Query() {
				if len(v) > 0 {
					capturedQuery[k] = v[0]
				}
			}
			resp := httpmock.NewStringResponse(200, mockSearchBody(0, nil))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	_, _, err := svc.SearchV1(context.Background(), &SearchOptions{
		Term:      "Taylor Swift",
		Media:     MediaMusic,
		Entity:    EntityMusicTrack,
		Attribute: AttributeArtistTerm,
		Country:   "US",
		Limit:     25,
		Lang:      "en_us",
		Explicit:  ExplicitNo,
		Version:   2,
	})

	require.NoError(t, err)
	assert.Equal(t, "Taylor Swift", capturedQuery["term"])
	assert.Equal(t, MediaMusic, capturedQuery["media"])
	assert.Equal(t, EntityMusicTrack, capturedQuery["entity"])
	assert.Equal(t, AttributeArtistTerm, capturedQuery["attribute"])
	assert.Equal(t, "US", capturedQuery["country"])
	assert.Equal(t, "25", capturedQuery["limit"])
	assert.Equal(t, "en_us", capturedQuery["lang"])
	assert.Equal(t, ExplicitNo, capturedQuery["explicit"])
	assert.Equal(t, "2", capturedQuery["version"])
}

func TestSearchV1_LimitCappedAtMax(t *testing.T) {
	svc := setupMockClient(t)

	var capturedLimit string
	httpmock.RegisterResponder("GET", `=~^https://itunes\.apple\.com/search`,
		func(req *http.Request) (*http.Response, error) {
			capturedLimit = req.URL.Query().Get("limit")
			resp := httpmock.NewStringResponse(200, mockSearchBody(0, nil))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	_, _, err := svc.SearchV1(context.Background(), &SearchOptions{Term: "test", Limit: 9999})

	require.NoError(t, err)
	assert.Equal(t, "200", capturedLimit, "limit should be capped at MaxLimit (200)")
}

func TestSearchV1_EmptyTermError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.SearchV1(context.Background(), &SearchOptions{Term: ""})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "search term is required")
}

func TestSearchV1_NilOptsError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.SearchV1(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "search term is required")
}

func TestSearchV1_NoResults(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://itunes.apple.com/search",
		jsonResponder(200, mockSearchBody(0, nil)))

	result, resp, err := svc.SearchV1(context.Background(), &SearchOptions{Term: "xyzzy_no_results_expected"})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, 0, result.ResultCount)
	assert.Empty(t, result.Results)
}

func TestSearchV1_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://itunes.apple.com/search",
		httpmock.NewStringResponder(503, "Service Unavailable"))

	_, resp, err := svc.SearchV1(context.Background(), &SearchOptions{Term: "test"})

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, err.Error(), "503")
}

func TestSearchV1_ContextCancellation(t *testing.T) {
	svc := setupMockClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond)

	_, _, err := svc.SearchV1(ctx, &SearchOptions{Term: "test"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// =============================================================================
// LookupV1
// =============================================================================

func TestLookupV1_ByID_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://itunes.apple.com/lookup",
		jsonResponder(200, mockSearchBody(1, []Result{
			{WrapperType: "artist", ArtistID: 909253, ArtistName: "Jack Johnson"},
		})))

	result, resp, err := svc.LookupV1(context.Background(), &LookupOptions{ID: 909253})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ResultCount)
	assert.Equal(t, "Jack Johnson", result.Results[0].ArtistName)
}

func TestLookupV1_ByIDs_CommaSeparated(t *testing.T) {
	svc := setupMockClient(t)

	var capturedID string
	httpmock.RegisterResponder("GET", `=~^https://itunes\.apple\.com/lookup`,
		func(req *http.Request) (*http.Response, error) {
			capturedID = req.URL.Query().Get("id")
			resp := httpmock.NewStringResponse(200, mockSearchBody(2, []Result{
				{ArtistID: 909253, ArtistName: "Jack Johnson"},
				{ArtistID: 26921078, ArtistName: "Weezer"},
			}))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	result, _, err := svc.LookupV1(context.Background(), &LookupOptions{IDs: []int{909253, 26921078}})

	require.NoError(t, err)
	assert.Equal(t, "909253,26921078", capturedID)
	assert.Equal(t, 2, result.ResultCount)
}

func TestLookupV1_ByUPC_Success(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://itunes.apple.com/lookup",
		jsonResponder(200, mockSearchBody(1, []Result{
			{WrapperType: "collection", CollectionName: "On and On"},
		})))

	result, resp, err := svc.LookupV1(context.Background(), &LookupOptions{UPC: "720642462928"})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "On and On", result.Results[0].CollectionName)
}

func TestLookupV1_ByAMGArtistIDs_QueryParams(t *testing.T) {
	svc := setupMockClient(t)

	var capturedQuery map[string]string
	httpmock.RegisterResponder("GET", `=~^https://itunes\.apple\.com/lookup`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = make(map[string]string)
			for k, v := range req.URL.Query() {
				if len(v) > 0 {
					capturedQuery[k] = v[0]
				}
			}
			resp := httpmock.NewStringResponse(200, mockSearchBody(0, nil))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	_, _, err := svc.LookupV1(context.Background(), &LookupOptions{
		AMGArtistIDs: []string{"468749", "5723"},
		Entity:       EntityAlbum,
		Limit:        5,
		Sort:         SortRecent,
		Country:      "US",
	})

	require.NoError(t, err)
	assert.Equal(t, "468749,5723", capturedQuery["amgArtistId"])
	assert.Equal(t, EntityAlbum, capturedQuery["entity"])
	assert.Equal(t, "5", capturedQuery["limit"])
	assert.Equal(t, SortRecent, capturedQuery["sort"])
	assert.Equal(t, "US", capturedQuery["country"])
}

func TestLookupV1_NilOptsError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.LookupV1(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lookup options are required")
}

func TestLookupV1_NoIdentifierError(t *testing.T) {
	svc := setupMockClient(t)

	_, _, err := svc.LookupV1(context.Background(), &LookupOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one lookup identifier is required")
}

func TestLookupV1_LimitCappedAtMax(t *testing.T) {
	svc := setupMockClient(t)

	var capturedLimit string
	httpmock.RegisterResponder("GET", `=~^https://itunes\.apple\.com/lookup`,
		func(req *http.Request) (*http.Response, error) {
			capturedLimit = req.URL.Query().Get("limit")
			resp := httpmock.NewStringResponse(200, mockSearchBody(0, nil))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	_, _, err := svc.LookupV1(context.Background(), &LookupOptions{ID: 909253, Limit: 9999})

	require.NoError(t, err)
	assert.Equal(t, "200", capturedLimit, "limit should be capped at MaxLimit (200)")
}

func TestLookupV1_HTTPError(t *testing.T) {
	svc := setupMockClient(t)

	httpmock.RegisterResponder("GET", "https://itunes.apple.com/lookup",
		httpmock.NewStringResponder(404, "Not Found"))

	_, resp, err := svc.LookupV1(context.Background(), &LookupOptions{ID: 999999999})

	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, err.Error(), "404")
}

// =============================================================================
// Constant smoke tests
// =============================================================================

func TestMediaConstants(t *testing.T) {
	assert.Equal(t, "music", MediaMusic)
	assert.Equal(t, "podcast", MediaPodcast)
	assert.Equal(t, "movie", MediaMovie)
	assert.Equal(t, "software", MediaSoftware)
	assert.Equal(t, "tvShow", MediaTVShow)
	assert.Equal(t, "audiobook", MediaAudiobook)
}

func TestExplicitConstants(t *testing.T) {
	assert.Equal(t, "Yes", ExplicitYes)
	assert.Equal(t, "No", ExplicitNo)
}

func TestSortConstants(t *testing.T) {
	assert.Equal(t, "recent", SortRecent)
	assert.Equal(t, "popular", SortPopular)
}
