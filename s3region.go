package s3region

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrRegionHeaderNotFound = errors.New("x-amz-bucket-region header not found in response")
var ErrBucketNotFound = errors.New("aws s3 bucket not found") // HEAD request returns 404
var ErrInvalidBucketName = errors.New("invalid S3 bucket name")

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

// isValidBucketName validates an S3 bucket name according to AWS naming rules.
func isValidBucketName(name string) bool {
	// Check length: must be between 3 and 63 characters
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	// Must begin and end with a letter or number
	first := name[0]
	last := name[len(name)-1]
	if !((first >= 'a' && first <= 'z') || (first >= '0' && first <= '9')) {
		return false
	}
	if !((last >= 'a' && last <= 'z') || (last >= '0' && last <= '9')) {
		return false
	}

	// Check characters and consecutive periods
	prevDot := false
	dotCount := 0
	for i := 0; i < len(name); i++ {
		c := name[i]
		// Can only consist of lowercase letters, numbers, periods (.), and hyphens (-)
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '.' || c == '-') {
			return false
		}

		// Must not contain two adjacent periods
		if c == '.' {
			dotCount++
			if prevDot {
				return false
			}
			prevDot = true
		} else {
			prevDot = false
		}
	}

	// Must not be formatted as an IP address (e.g., 192.168.5.4)
	if dotCount == 3 {
		parts := strings.Split(name, ".")
		if len(parts) == 4 {
			allNumbers := true
			for _, part := range parts {
				if len(part) == 0 {
					allNumbers = false
					break
				}
				for j := 0; j < len(part); j++ {
					if part[j] < '0' || part[j] > '9' {
						allNumbers = false
						break
					}
				}
				if !allNumbers {
					break
				}
			}
			if allNumbers {
				return false
			}
		}
	}

	return true
}

// GetBucketRegionByName takes a bucket name and returns its region by constructing
// the S3 URL and performing a HEAD request to extract the x-amz-bucket-region header.
func GetBucketRegionByName(ctx context.Context, bucketName string, opts ...Option) (string, error) {
	cfg := &config{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	if !isValidBucketName(bucketName) {
		return "", fmt.Errorf("%w: %s", ErrInvalidBucketName, bucketName)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucketName)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := cfg.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform HEAD request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrBucketNotFound
	}
	region := resp.Header.Get("x-amz-bucket-region")
	if region == "" {
		return "", ErrRegionHeaderNotFound
	}

	return strings.TrimSpace(region), nil
}

// GetBucketRegionFromARN extracts the bucket name from an AWS S3 ARN and returns its region.
// Accepts ARN format: arn:aws:s3:::bucket-name or arn:aws:s3:::bucket-name/path/to/object
func GetBucketRegionFromARN(ctx context.Context, arn string, opts ...Option) (string, error) {
	bucketName := strings.TrimPrefix(arn, "arn:aws:s3:::")
	// Remove any path after bucket name
	if idx := strings.Index(bucketName, "/"); idx != -1 {
		bucketName = bucketName[:idx]
	}
	return GetBucketRegionByName(ctx, bucketName, opts...)
}

// GetBucketRegionFromS3URI extracts the bucket name from an S3 URI and returns its region.
// Accepts S3 URI format: s3://bucket-name or s3://bucket-name/path/to/object
func GetBucketRegionFromS3URI(ctx context.Context, uri string, opts ...Option) (string, error) {
	bucketName := strings.TrimPrefix(uri, "s3://")
	// Remove any path after bucket name
	if idx := strings.Index(bucketName, "/"); idx != -1 {
		bucketName = bucketName[:idx]
	}
	return GetBucketRegionByName(ctx, bucketName, opts...)
}

// GetBucketRegionFromHTTPURL extracts the bucket name from an HTTP/HTTPS URL and returns its region.
// Supports both virtual-hosted-style and path-style URLs:
// - Virtual-hosted: https://bucket-name.s3.amazonaws.com/path/to/object
// - Path-style: https://s3.amazonaws.com/bucket-name/path/to/object
// - Path-style with region: https://s3.us-west-2.amazonaws.com/bucket-name/path/to/object
func GetBucketRegionFromHTTPURL(ctx context.Context, url string, opts ...Option) (string, error) {
	// Remove protocol
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Get the host part (before first /)
	host := url
	path := ""
	if idx := strings.Index(url, "/"); idx != -1 {
		host = url[:idx]
		path = url[idx+1:]
	}

	// Check if this is a virtual-hosted-style URL (bucket-name.s3.amazonaws.com)
	if strings.Contains(host, ".s3.") || strings.Contains(host, ".s3-") {
		// Extract bucket name from host (before .s3.)
		if idx := strings.Index(host, ".s3"); idx != -1 {
			bucketName := host[:idx]
			return GetBucketRegionByName(ctx, bucketName, opts...)
		}
	}

	// Path-style URL (s3.amazonaws.com/bucket-name or s3.region.amazonaws.com/bucket-name)
	// Extract bucket name from path (first segment)
	if path != "" {
		bucketName := path
		if idx := strings.Index(path, "/"); idx != -1 {
			bucketName = path[:idx]
		}
		return GetBucketRegionByName(ctx, bucketName, opts...)
	}

	// If we couldn't parse it, treat the host as bucket name
	return GetBucketRegionByName(ctx, host, opts...)
}

// GetBucketRegion is the main umbrella function that accepts any S3 identifier format
// and automatically detects the type to extract the bucket region. Supports:
// - Bucket name: my-bucket or my-bucket/path/to/object
// - S3 URI: s3://my-bucket or s3://my-bucket/path/to/object
// - AWS ARN: arn:aws:s3:::my-bucket or arn:aws:s3:::my-bucket/path
// - HTTP/HTTPS URL: https://my-bucket.s3.amazonaws.com or https://my-bucket.s3.amazonaws.com/path/to/object
func GetBucketRegion(ctx context.Context, input string, opts ...Option) (string, error) {
	// Handle AWS ARN format
	if strings.HasPrefix(input, "arn:aws:s3:::") {
		return GetBucketRegionFromARN(ctx, input, opts...)
	}

	// Handle S3 URI format
	if strings.HasPrefix(input, "s3://") {
		return GetBucketRegionFromS3URI(ctx, input, opts...)
	}

	// Handle HTTP/HTTPS URL format
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return GetBucketRegionFromHTTPURL(ctx, input, opts...)
	}

	// Handle plain bucket name with or without path
	bucketName := input
	if idx := strings.Index(input, "/"); idx != -1 {
		bucketName = input[:idx]
	}
	return GetBucketRegionByName(ctx, bucketName, opts...)
}
