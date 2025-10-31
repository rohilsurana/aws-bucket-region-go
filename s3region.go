package s3region

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrRegionHeaderNotFound = errors.New("x-amz-bucket-region header not found in response")
var ErrBucketNotFound = errors.New("aws s3 bucket not found") // HEAD request returns 404

// GetBucketRegionByName takes a bucket name and returns its region by constructing
// the S3 URL and performing a HEAD request to extract the x-amz-bucket-region header.
func GetBucketRegionByName(bucketName string) (string, error) {
	url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucketName)

	resp, err := http.Head(url)
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
func GetBucketRegionFromARN(arn string) (string, error) {
	bucketName := strings.TrimPrefix(arn, "arn:aws:s3:::")
	// Remove any path after bucket name
	if idx := strings.Index(bucketName, "/"); idx != -1 {
		bucketName = bucketName[:idx]
	}
	return GetBucketRegionByName(bucketName)
}

// GetBucketRegionFromS3URI extracts the bucket name from an S3 URI and returns its region.
// Accepts S3 URI format: s3://bucket-name or s3://bucket-name/path/to/object
func GetBucketRegionFromS3URI(uri string) (string, error) {
	bucketName := strings.TrimPrefix(uri, "s3://")
	// Remove any path after bucket name
	if idx := strings.Index(bucketName, "/"); idx != -1 {
		bucketName = bucketName[:idx]
	}
	return GetBucketRegionByName(bucketName)
}

// GetBucketRegionFromHTTPURL extracts the bucket name from an HTTP/HTTPS URL and returns its region.
// Supports both virtual-hosted-style and path-style URLs:
// - Virtual-hosted: https://bucket-name.s3.amazonaws.com/path/to/object
// - Path-style: https://s3.amazonaws.com/bucket-name/path/to/object
// - Path-style with region: https://s3.us-west-2.amazonaws.com/bucket-name/path/to/object
func GetBucketRegionFromHTTPURL(url string) (string, error) {
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
			return GetBucketRegionByName(bucketName)
		}
	}

	// Path-style URL (s3.amazonaws.com/bucket-name or s3.region.amazonaws.com/bucket-name)
	// Extract bucket name from path (first segment)
	if path != "" {
		bucketName := path
		if idx := strings.Index(path, "/"); idx != -1 {
			bucketName = path[:idx]
		}
		return GetBucketRegionByName(bucketName)
	}

	// If we couldn't parse it, treat the host as bucket name
	return GetBucketRegionByName(host)
}

// GetBucketRegion is the main umbrella function that accepts any S3 identifier format
// and automatically detects the type to extract the bucket region. Supports:
// - Bucket name: my-bucket or my-bucket/path/to/object
// - S3 URI: s3://my-bucket or s3://my-bucket/path/to/object
// - AWS ARN: arn:aws:s3:::my-bucket or arn:aws:s3:::my-bucket/path
// - HTTP/HTTPS URL: https://my-bucket.s3.amazonaws.com or https://my-bucket.s3.amazonaws.com/path/to/object
func GetBucketRegion(input string) (string, error) {
	// Handle AWS ARN format
	if strings.HasPrefix(input, "arn:aws:s3:::") {
		return GetBucketRegionFromARN(input)
	}

	// Handle S3 URI format
	if strings.HasPrefix(input, "s3://") {
		return GetBucketRegionFromS3URI(input)
	}

	// Handle HTTP/HTTPS URL format
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return GetBucketRegionFromHTTPURL(input)
	}

	// Handle plain bucket name with or without path
	bucketName := input
	if idx := strings.Index(input, "/"); idx != -1 {
		bucketName = input[:idx]
	}
	return GetBucketRegionByName(bucketName)
}
