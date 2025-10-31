package s3region

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrRegionHeaderNotFound = errors.New("x-amz-bucket-region header not found in response")
var ErrBucketNotFound = errors.New("aws s3 bucket not found") // HEAD request returns 404

// GetBucketRegion performs a HEAD request to the given S3 URL and extracts the bucket region
// from the x-amz-bucket-region header. Returns an error if the request fails or header is missing.
func GetBucketRegion(url string) (string, error) {
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

// GetBucketRegionByName takes a bucket name and returns its region by constructing
// the S3 URL and performing a HEAD request to extract the x-amz-bucket-region header.
func GetBucketRegionByName(bucketName string) (string, error) {
	url := fmt.Sprintf("https://%s.s3.amazonaws.com", bucketName)
	return GetBucketRegion(url)
}
