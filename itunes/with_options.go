package itunes

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/itunes/client"
	"go.uber.org/zap"
)

// ClientOption configures the iTunes Search API transport at construction time.
// Pass one or more ClientOption values to NewClient or NewDefaultClient.
type ClientOption = client.ClientOption

// WithBaseURL sets a custom base URL, overriding the default iTunes endpoint.
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

// WithRootCertificates adds custom root CA certificates for server validation.
func WithRootCertificates(pemFilePaths ...string) ClientOption {
	return client.WithRootCertificates(pemFilePaths...)
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
