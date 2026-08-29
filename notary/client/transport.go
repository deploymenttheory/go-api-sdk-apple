package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/notary/constants"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Transport represents the main Apple Notary API transport layer.
type Transport struct {
	httpClient   *resty.Client
	logger       *zap.Logger
	auth         AuthProvider
	errorHandler *ErrorHandler
	baseURL      string
}

// Ensure Transport implements Client interface.
var _ Client = (*Transport)(nil)

// APIResponse represents the standard API response structure.
type APIResponse[T any] struct {
	Data  []T   `json:"data"`
	Meta  Meta  `json:"meta"`
	Links Links `json:"links"`
}

// NewTransport creates a new HTTP transport for the Apple Notary API.
// This is an internal function - users should use notary.NewClient() instead.
func NewTransport(keyID, issuerID string, privateKey any, options ...ClientOption) (*Transport, error) {
	if keyID == "" {
		return nil, fmt.Errorf("keyID is required")
	}
	if issuerID == "" {
		return nil, fmt.Errorf("issuerID is required")
	}
	if privateKey == nil {
		return nil, fmt.Errorf("privateKey is required")
	}

	logger := zap.NewNop()

	auth := NewJWTAuth(JWTAuthConfig{
		KeyID:      keyID,
		IssuerID:   issuerID,
		PrivateKey: privateKey,
		Audience:   constants.DefaultJWTAudience,
	})

	httpClient := resty.New()
	httpClient.
		SetBaseURL(constants.DefaultBaseURL).
		SetTimeout(30*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetHeader("User-Agent", DefaultUserAgent)

	errorHandler := NewErrorHandler(logger)

	transport := &Transport{
		httpClient:   httpClient,
		logger:       logger,
		auth:         auth,
		errorHandler: errorHandler,
		baseURL:      constants.DefaultBaseURL,
	}

	for _, option := range options {
		if err := option(transport); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	httpClient.AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
		if err := transport.auth.ApplyAuth(req); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		transport.logger.Info("API request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)

		return nil
	})

	httpClient.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		transport.logger.Info("API response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.String("status", resp.Status()),
		)

		if resp.StatusCode() == 401 {
			if jwtAuth, ok := transport.auth.(*JWTAuth); ok {
				transport.logger.Info("Received 401 response, forcing JWT token refresh")
				jwtAuth.ForceRefresh()
			}
		}

		return nil
	})

	transport.logger.Info("Apple Notary API client created",
		zap.String("issuer_id", issuerID),
		zap.String("base_url", transport.baseURL))

	return transport, nil
}

// NewRequest returns a new RequestBuilder for constructing API requests.
func (t *Transport) NewRequest(ctx context.Context) *RequestBuilder {
	return &RequestBuilder{
		req:      t.httpClient.R().SetContext(ctx),
		executor: t,
	}
}

// QueryBuilder returns a new query builder instance.
func (t *Transport) QueryBuilder() *QueryBuilder {
	return NewQueryBuilder()
}

// GetLogger returns the configured logger.
func (t *Transport) GetLogger() *zap.Logger {
	return t.logger
}

// GetHTTPClient returns the underlying HTTP client for testing purposes.
func (t *Transport) GetHTTPClient() *resty.Client {
	return t.httpClient
}

// Close closes the HTTP client and cleans up resources.
func (t *Transport) Close() error {
	if t.httpClient != nil {
		t.httpClient.Close()
	}
	return nil
}

// execute implements requestExecutor — handles all HTTP method routing and error processing.
func (t *Transport) execute(req *resty.Request, method, path string, result any) (*resty.Response, error) {
	var apiErr ErrorResponse
	req.SetResultError(&apiErr)

	if result != nil {
		req.SetResult(result)
	}

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "PATCH":
		resp, err = req.Patch(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsStatusFailure() {
		return resp, t.errorHandler.HandleError(resp, &apiErr)
	}

	return resp, nil
}

// executeGetBytes implements requestExecutor — returns raw response bytes without JSON unmarshaling.
func (t *Transport) executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error) {
	resp, err := t.execute(req, "GET", path, nil)
	if err != nil {
		return resp, nil, err
	}
	return resp, resp.Bytes(), nil
}

// executePaginated implements requestExecutor — cursor-based pagination loop.
func (t *Transport) executePaginated(req *resty.Request, path string, mergePage func([]byte) error) (*resty.Response, error) {
	currentParams := make(map[string]string)
	for k, v := range req.QueryParams {
		if len(v) > 0 {
			currentParams[k] = v[0]
		}
	}

	var lastResp *resty.Response

	for {
		pageReq := t.httpClient.R().SetContext(req.Context())
		for k, v := range req.Header {
			if len(v) > 0 {
				pageReq.SetHeader(k, v[0])
			}
		}
		for k, v := range currentParams {
			if v != "" {
				pageReq.SetQueryParam(k, v)
			}
		}

		var apiErr ErrorResponse
		pageReq.SetResultError(&apiErr)

		resp, err := pageReq.Get(path)
		if err != nil {
			return resp, fmt.Errorf("request failed: %w", err)
		}
		if resp.IsStatusFailure() {
			return resp, t.errorHandler.HandleError(resp, &apiErr)
		}

		lastResp = resp
		rawResponse := resp.Bytes()

		if err := mergePage(rawResponse); err != nil {
			return resp, err
		}

		var pageInfo struct {
			Links *Links `json:"links,omitempty"`
		}
		if err := parseJSON(rawResponse, &pageInfo); err != nil {
			return resp, fmt.Errorf("failed to parse pagination info: %w", err)
		}

		if !HasNextPage(pageInfo.Links) {
			break
		}

		nextParams, err := extractParamsFromURL(pageInfo.Links.Next)
		if err != nil {
			return resp, fmt.Errorf("failed to parse next URL: %w", err)
		}

		for k, v := range nextParams {
			currentParams[k] = v
		}
	}

	return lastResp, nil
}

// NewTransportFromEnv creates a transport using environment variables.
// Requires APPLE_KEY_ID and APPLE_ISSUER_ID plus exactly one of:
//   - APPLE_PRIVATE_KEY_PEM  — PEM-encoded private key supplied inline
//   - APPLE_PRIVATE_KEY_PATH — path to a PEM private key file
func NewTransportFromEnv(options ...ClientOption) (*Transport, error) {
	keyID := os.Getenv("APPLE_KEY_ID")
	issuerID := os.Getenv("APPLE_ISSUER_ID")
	privateKeyPEM := os.Getenv("APPLE_PRIVATE_KEY_PEM")
	privateKeyPath := os.Getenv("APPLE_PRIVATE_KEY_PATH")

	if keyID == "" {
		return nil, fmt.Errorf("APPLE_KEY_ID environment variable is required")
	}
	if issuerID == "" {
		return nil, fmt.Errorf("APPLE_ISSUER_ID environment variable is required")
	}

	var privateKey any
	var err error

	switch {
	case privateKeyPEM != "":
		privateKey, err = ParsePrivateKey([]byte(privateKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("failed to parse APPLE_PRIVATE_KEY_PEM: %w", err)
		}
	case privateKeyPath != "":
		privateKey, err = LoadPrivateKeyFromFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key from APPLE_PRIVATE_KEY_PATH: %w", err)
		}
	default:
		return nil, fmt.Errorf("either APPLE_PRIVATE_KEY_PEM or APPLE_PRIVATE_KEY_PATH environment variable is required")
	}

	return NewTransport(keyID, issuerID, privateKey, options...)
}

// NewTransportFromFile creates a transport using credentials from files.
func NewTransportFromFile(keyID, issuerID, privateKeyPath string, options ...ClientOption) (*Transport, error) {
	if keyID == "" {
		return nil, fmt.Errorf("keyID is required")
	}
	if issuerID == "" {
		return nil, fmt.Errorf("issuerID is required")
	}
	if privateKeyPath == "" {
		return nil, fmt.Errorf("privateKeyPath is required")
	}

	privateKey, err := LoadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return NewTransport(keyID, issuerID, privateKey, options...)
}
