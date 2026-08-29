package search

import (
	"context"
	"testing"

	acc "github.com/deploymenttheory/go-sdk-appleservices/itunes/acceptance"
	"github.com/deploymenttheory/go-sdk-appleservices/itunes/itunes_api/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TestAcceptance_Search_SearchV1
// Verifies the /search endpoint returns valid responses for common queries.
// =============================================================================

func TestAcceptance_Search_SearchV1(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.ItunesAPI.Search
	ctx := context.Background()

	// --- Basic search ---
	t.Run("BasicMusicSearch", func(t *testing.T) {
		acc.LogTestStage(t, "Search", "Searching for 'Jack Johnson' music")

		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.SearchV1(ctx1, &search.SearchOptions{
			Term:  "Jack Johnson",
			Media: search.MediaMusic,
			Limit: 5,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "search for 'Jack Johnson' should return results")

		acc.LogTestSuccess(t, "SearchV1 music: resultCount=%d returned=%d", result.ResultCount, len(result.Results))
	})

	// --- Podcast search ---
	t.Run("PodcastSearch", func(t *testing.T) {
		acc.LogTestStage(t, "Search", "Searching for podcasts")

		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		result, resp, err := svc.SearchV1(ctx2, &search.SearchOptions{
			Term:  "technology",
			Media: search.MediaPodcast,
			Limit: 5,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "search for 'technology' podcasts should return results")

		acc.LogTestSuccess(t, "SearchV1 podcast: returned=%d", len(result.Results))
	})

	// --- Software search ---
	t.Run("SoftwareSearch", func(t *testing.T) {
		acc.LogTestStage(t, "Search", "Searching for software apps")

		ctx3, cancel3 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel3()

		result, resp, err := svc.SearchV1(ctx3, &search.SearchOptions{
			Term:  "weather",
			Media: search.MediaSoftware,
			Limit: 5,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "search for 'weather' software should return results")

		acc.LogTestSuccess(t, "SearchV1 software: returned=%d", len(result.Results))
	})

	// --- Entity filter ---
	t.Run("WithEntityFilter", func(t *testing.T) {
		acc.LogTestStage(t, "Search", "Searching with entity=musicTrack")

		ctx4, cancel4 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel4()

		result, resp, err := svc.SearchV1(ctx4, &search.SearchOptions{
			Term:   "Adele",
			Media:  search.MediaMusic,
			Entity: search.EntityMusicTrack,
			Limit:  10,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results)

		// Every result should be a track-type wrapper
		for _, r := range result.Results {
			assert.Equal(t, "track", r.WrapperType, "entity=musicTrack should only return tracks")
		}

		acc.LogTestSuccess(t, "SearchV1 entity=musicTrack: returned=%d", len(result.Results))
	})

	// --- Country filter ---
	t.Run("WithCountryFilter", func(t *testing.T) {
		acc.LogTestStage(t, "Search", "Searching in GB store")

		ctx5, cancel5 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel5()

		result, resp, err := svc.SearchV1(ctx5, &search.SearchOptions{
			Term:    "Beatles",
			Media:   search.MediaMusic,
			Country: "GB",
			Limit:   5,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "search in GB store should return results")

		acc.LogTestSuccess(t, "SearchV1 country=GB: returned=%d", len(result.Results))
	})

	// --- Limit is respected (capped at MaxLimit) ---
	t.Run("LimitBoundary", func(t *testing.T) {
		acc.LogTestStage(t, "Search", "Verifying limit is honoured")

		ctx6, cancel6 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel6()

		const requestedLimit = 3
		result, resp, err := svc.SearchV1(ctx6, &search.SearchOptions{
			Term:  "Taylor Swift",
			Media: search.MediaMusic,
			Limit: requestedLimit,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Results), requestedLimit,
			"result count should not exceed the requested limit")

		acc.LogTestSuccess(t, "SearchV1 limit=%d: returned=%d", requestedLimit, len(result.Results))
	})
}

// =============================================================================
// TestAcceptance_Search_SearchV1_ValidationErrors
// Verifies client-side validation fires before any HTTP call is made.
// =============================================================================

func TestAcceptance_Search_SearchV1_ValidationErrors(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.ItunesAPI.Search
	ctx := context.Background()

	t.Run("EmptyTerm", func(t *testing.T) {
		_, _, err := svc.SearchV1(ctx, &search.SearchOptions{Term: ""})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search term is required")
	})

	t.Run("NilOptions", func(t *testing.T) {
		_, _, err := svc.SearchV1(ctx, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search term is required")
	})
}

// =============================================================================
// TestAcceptance_Search_LookupV1
// Verifies the /lookup endpoint returns valid responses for known identifiers.
// =============================================================================

func TestAcceptance_Search_LookupV1(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.ItunesAPI.Search
	ctx := context.Background()

	// --- Lookup by single iTunes ID ---
	t.Run("ByID_Artist", func(t *testing.T) {
		// 909253 = Jack Johnson (well-known, stable iTunes artist ID)
		acc.LogTestStage(t, "Lookup", "Looking up artist ID=909253 (Jack Johnson)")

		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.LookupV1(ctx1, &search.LookupOptions{ID: 909253})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "lookup by known artist ID should return a result")
		assert.Equal(t, int64(909253), result.Results[0].ArtistID)

		acc.LogTestSuccess(t, "LookupV1 ID=909253: artist=%q", result.Results[0].ArtistName)
	})

	// --- Lookup by multiple iTunes IDs ---
	t.Run("ByMultipleIDs", func(t *testing.T) {
		// 909253 = Jack Johnson, 26921078 = Weezer
		acc.LogTestStage(t, "Lookup", "Looking up multiple artist IDs: 909253, 26921078")

		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		result, resp, err := svc.LookupV1(ctx2, &search.LookupOptions{
			IDs: []int{909253, 26921078},
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "lookup by multiple IDs should return results")

		acc.LogTestSuccess(t, "LookupV1 IDs=[909253,26921078]: returned=%d", len(result.Results))
	})

	// --- Lookup by UPC ---
	t.Run("ByUPC", func(t *testing.T) {
		// Jack Johnson "On and On" UPC
		acc.LogTestStage(t, "Lookup", "Looking up by UPC=720642462928")

		ctx3, cancel3 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel3()

		result, resp, err := svc.LookupV1(ctx3, &search.LookupOptions{UPC: "720642462928"})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		// UPC lookup returns a result or an empty list — either is a valid API response
		acc.LogTestSuccess(t, "LookupV1 UPC: returned=%d", len(result.Results))
	})

	// --- Lookup with entity filter ---
	t.Run("WithEntityAndLimit", func(t *testing.T) {
		// Look up Jack Johnson's albums
		acc.LogTestStage(t, "Lookup", "Looking up Jack Johnson albums (ID=909253, entity=album)")

		ctx4, cancel4 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel4()

		result, resp, err := svc.LookupV1(ctx4, &search.LookupOptions{
			ID:     909253,
			Entity: search.EntityAlbum,
			Limit:  5,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Results, "Jack Johnson should have albums")
		assert.LessOrEqual(t, len(result.Results), 5+1, // +1 for the artist record itself
			"result count should not greatly exceed the requested limit")

		acc.LogTestSuccess(t, "LookupV1 ID=909253 entity=album limit=5: returned=%d", len(result.Results))
	})

	// --- Lookup with sort ---
	t.Run("WithSortRecent", func(t *testing.T) {
		acc.LogTestStage(t, "Lookup", "Looking up Jack Johnson albums sorted by recent")

		ctx5, cancel5 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel5()

		result, resp, err := svc.LookupV1(ctx5, &search.LookupOptions{
			ID:     909253,
			Entity: search.EntityAlbum,
			Sort:   search.SortRecent,
			Limit:  5,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)

		acc.LogTestSuccess(t, "LookupV1 sort=recent: returned=%d", len(result.Results))
	})
}

// =============================================================================
// TestAcceptance_Search_LookupV1_ValidationErrors
// Verifies client-side validation fires before any HTTP call is made.
// =============================================================================

func TestAcceptance_Search_LookupV1_ValidationErrors(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.ItunesAPI.Search
	ctx := context.Background()

	t.Run("NilOptions", func(t *testing.T) {
		_, _, err := svc.LookupV1(ctx, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "lookup options are required")
	})

	t.Run("NoIdentifier", func(t *testing.T) {
		_, _, err := svc.LookupV1(ctx, &search.LookupOptions{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one lookup identifier is required")
	})
}
