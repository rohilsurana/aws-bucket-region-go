package s3region

import "net/http"

// HTTPClient interface allows custom HTTP client implementations.
// The standard *http.Client implements this interface.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// config holds configuration options for S3 region lookup.
type config struct {
	httpClient HTTPClient
}

// Option is a function that configures the internal config.
type Option func(*config)

// WithHTTPClient sets a custom HTTP client for S3 requests.
// If not provided, http.DefaultClient is used.
func WithHTTPClient(client HTTPClient) Option {
	return func(c *config) {
		c.httpClient = client
	}
}
