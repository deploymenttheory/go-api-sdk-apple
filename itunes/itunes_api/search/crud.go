package search

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-appleservices/itunes/client"
	"github.com/deploymenttheory/go-sdk-appleservices/itunes/constants"
	"resty.dev/v3"
)

// SearchService handles communication with the iTunes Search and Lookup API.
//
// iTunes Search API docs: https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/
type SearchService struct {
	client client.Client
}

// NewService creates a new iTunes search service.
func NewService(c client.Client) *SearchService {
	return &SearchService{client: c}
}

// SearchV1 searches the iTunes Store using the Search endpoint.
// URL: GET https://itunes.apple.com/search
// https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/Searching.html
func (s *SearchService) SearchV1(ctx context.Context, opts *SearchOptions) (*SearchResponse, *resty.Response, error) {
	if opts == nil || opts.Term == "" {
		return nil, nil, fmt.Errorf("search term is required")
	}

	params := s.client.QueryBuilder()
	params.AddString("term", opts.Term)
	params.AddString("country", opts.Country)
	params.AddString("media", opts.Media)
	params.AddString("entity", opts.Entity)
	params.AddString("attribute", opts.Attribute)
	params.AddString("callback", opts.Callback)
	params.AddString("lang", opts.Lang)
	params.AddString("explicit", opts.Explicit)

	if opts.Limit > 0 {
		limit := opts.Limit
		if limit > constants.MaxLimit {
			limit = constants.MaxLimit
		}
		params.AddInt("limit", limit)
	}

	if opts.Version > 0 {
		params.AddInt("version", opts.Version)
	}

	var result SearchResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		SetResult(&result).
		Get(constants.EndpointSearch)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// LookupV1 looks up items in the iTunes Store by a known identifier.
// URL: GET https://itunes.apple.com/lookup
// https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/LookupExamples.html
func (s *SearchService) LookupV1(ctx context.Context, opts *LookupOptions) (*SearchResponse, *resty.Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("lookup options are required")
	}

	hasIdentifier := opts.ID > 0 || len(opts.IDs) > 0 ||
		opts.UPC != "" || opts.EAN != "" || opts.ISRC != "" || opts.ISBN != "" ||
		opts.AMGArtistID != "" || len(opts.AMGArtistIDs) > 0 ||
		opts.AMGAlbumID != "" || len(opts.AMGAlbumIDs) > 0 ||
		opts.AMGVideoID != ""

	if !hasIdentifier {
		return nil, nil, fmt.Errorf("at least one lookup identifier is required (ID, UPC, EAN, ISRC, ISBN, AMGArtistID, AMGAlbumID, or AMGVideoID)")
	}

	params := s.client.QueryBuilder()

	// IDs wins over single ID when both are supplied.
	if len(opts.IDs) > 0 {
		params.AddIntSlice("id", opts.IDs)
	} else if opts.ID > 0 {
		params.AddInt("id", opts.ID)
	}

	params.AddString("upc", opts.UPC)
	params.AddString("ean", opts.EAN)
	params.AddString("isrc", opts.ISRC)
	params.AddString("isbn", opts.ISBN)

	if len(opts.AMGArtistIDs) > 0 {
		params.AddStringSlice("amgArtistId", opts.AMGArtistIDs)
	} else {
		params.AddString("amgArtistId", opts.AMGArtistID)
	}

	if len(opts.AMGAlbumIDs) > 0 {
		params.AddStringSlice("amgAlbumId", opts.AMGAlbumIDs)
	} else {
		params.AddString("amgAlbumId", opts.AMGAlbumID)
	}

	params.AddString("amgVideoId", opts.AMGVideoID)
	params.AddString("entity", opts.Entity)
	params.AddString("sort", opts.Sort)
	params.AddString("country", opts.Country)

	if opts.Limit > 0 {
		limit := opts.Limit
		if limit > constants.MaxLimit {
			limit = constants.MaxLimit
		}
		params.AddInt("limit", limit)
	}

	var result SearchResponse

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetQueryParams(params.Build()).
		SetResult(&result).
		Get(constants.EndpointLookup)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
