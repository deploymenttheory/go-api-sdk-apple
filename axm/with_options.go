package axm

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"go.uber.org/zap"
)

// ClientOption configures the Apple Business Manager API transport at construction time.
// Pass one or more ClientOption values to NewClient, NewClientFromFile, or NewClientFromEnv.
type ClientOption = client.ClientOption

// WithBaseURL sets a custom base URL, overriding the default Apple Business Manager endpoint.
func WithBaseURL(baseURL string) ClientOption {
	return client.WithBaseURL(baseURL)
}

// WithLogger sets a custom zap logger. Returns an error if logger is nil.
func WithLogger(logger *zap.Logger) ClientOption {
	return client.WithLogger(logger)
}

// WithTimeout sets the timeout for all HTTP requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return client.WithTimeout(timeout)
}

// WithRetryCount sets the maximum number of retries for failed requests.
func WithRetryCount(count int) ClientOption {
	return client.WithRetryCount(count)
}

// WithRetryWaitTime sets the initial wait time between retry attempts.
func WithRetryWaitTime(waitTime time.Duration) ClientOption {
	return client.WithRetryWaitTime(waitTime)
}

// WithRetryMaxWaitTime sets the maximum wait time between retry attempts.
func WithRetryMaxWaitTime(maxWaitTime time.Duration) ClientOption {
	return client.WithRetryMaxWaitTime(maxWaitTime)
}

// WithUserAgent sets a custom user-agent string.
func WithUserAgent(userAgent string) ClientOption {
	return client.WithUserAgent(userAgent)
}

// WithCustomAgent appends a custom identifier to the default user agent.
func WithCustomAgent(customAgent string) ClientOption {
	return client.WithCustomAgent(customAgent)
}

// WithDebug enables resty's request/response debug logging.
func WithDebug() ClientOption {
	return client.WithDebug()
}

// WithGlobalHeader adds a single header to every outgoing request.
func WithGlobalHeader(key, value string) ClientOption {
	return client.WithGlobalHeader(key, value)
}

// WithGlobalHeaders adds multiple headers to every outgoing request.
func WithGlobalHeaders(headers map[string]string) ClientOption {
	return client.WithGlobalHeaders(headers)
}

// WithProxy sets an HTTP proxy for all requests.
func WithProxy(proxyURL string) ClientOption {
	return client.WithProxy(proxyURL)
}

// WithTLSClientConfig sets custom TLS configuration.
func WithTLSClientConfig(tlsConfig *tls.Config) ClientOption {
	return client.WithTLSClientConfig(tlsConfig)
}

// WithClientCertificate sets a client certificate for mutual TLS authentication.
func WithClientCertificate(certFile, keyFile string) ClientOption {
	return client.WithClientCertificate(certFile, keyFile)
}

// WithClientCertificateFromString sets a client certificate from PEM-encoded strings.
func WithClientCertificateFromString(certPEM, keyPEM string) ClientOption {
	return client.WithClientCertificateFromString(certPEM, keyPEM)
}

// WithRootCertificates adds custom root CA certificates for server validation.
func WithRootCertificates(pemFilePaths ...string) ClientOption {
	return client.WithRootCertificates(pemFilePaths...)
}

// WithRootCertificateFromString adds a custom root CA certificate from PEM string.
func WithRootCertificateFromString(pemContent string) ClientOption {
	return client.WithRootCertificateFromString(pemContent)
}

// WithTransport sets a custom HTTP transport (http.RoundTripper).
func WithTransport(transport http.RoundTripper) ClientOption {
	return client.WithTransport(transport)
}

// WithInsecureSkipVerify disables TLS certificate verification (use only for testing).
func WithInsecureSkipVerify() ClientOption {
	return client.WithInsecureSkipVerify()
}

// WithMinTLSVersion sets the minimum TLS version for connections.
func WithMinTLSVersion(minVersion uint16) ClientOption {
	return client.WithMinTLSVersion(minVersion)
}

// WithAudience sets a custom JWT audience (default: "appstoreconnect-v1").
func WithAudience(audience string) ClientOption {
	return client.WithAudience(audience)
}

// WithScope sets a custom JWT scope (default: "business.api").
func WithScope(scope string) ClientOption {
	return client.WithScope(scope)
}

// IsNotFound returns true when err is an API 404 response.
// Use this in cleanup functions to treat "already deleted" as non-fatal.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var apiErr *client.APIError
	return errors.As(err, &apiErr) && apiErr.Status == "404"
}

// ParsePrivateKey parses a PEM-encoded private key (ECDSA or RSA) from bytes.
func ParsePrivateKey(keyData []byte) (any, error) {
	return client.ParsePrivateKey(keyData)
}

// LoadPrivateKeyFromFile reads and parses a private key from a .p8 file path.
func LoadPrivateKeyFromFile(filePath string) (any, error) {
	return client.LoadPrivateKeyFromFile(filePath)
}
