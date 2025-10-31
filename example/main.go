package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	s3region "github.com/rohilsurana/aws-bucket-region-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <s3-bucket-name-or-url>")
		fmt.Println("Example: go run main.go testing-bucket")
		fmt.Println("Example: go run main.go testing-bucket/path/to/object")
		fmt.Println("Example: go run main.go s3://testing-bucket/path/to/object")
		fmt.Println("Example: go run main.go arn:aws:s3:::testing-bucket/path")
		fmt.Println("Example: go run main.go https://testing-bucket.s3.amazonaws.com/path/to/object")
		os.Exit(1)
	}

	input := os.Args[1]

	var region string
	var err error

	// Handle AWS ARN format: arn:aws:s3:::bucket-name
	if strings.HasPrefix(input, "arn:aws:s3:::") {
		bucketName := strings.TrimPrefix(input, "arn:aws:s3:::")
		// Remove any path after bucket name
		if idx := strings.Index(bucketName, "/"); idx != -1 {
			bucketName = bucketName[:idx]
		}
		region, err = s3region.GetBucketRegionByName(bucketName)
	} else if strings.HasPrefix(input, "s3://") {
		bucketName := strings.TrimPrefix(input, "s3://")
		// Remove any path after bucket name
		if idx := strings.Index(bucketName, "/"); idx != -1 {
			bucketName = bucketName[:idx]
		}
		region, err = s3region.GetBucketRegionByName(bucketName)
	} else if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		// Extract bucket name from URL with path
		// Format: https://bucket-name.s3.amazonaws.com/path/to/object
		url := input
		// Remove protocol
		url = strings.TrimPrefix(url, "https://")
		url = strings.TrimPrefix(url, "http://")
		// Get the host part (before first /)
		host := url
		if idx := strings.Index(url, "/"); idx != -1 {
			host = url[:idx]
		}
		// Extract bucket name from host (before .s3.amazonaws.com)
		bucketName := host
		if idx := strings.Index(host, ".s3."); idx != -1 {
			bucketName = host[:idx]
		}
		region, err = s3region.GetBucketRegionByName(bucketName)
	} else {
		// Handle plain bucket name with or without path
		bucketName := input
		if idx := strings.Index(input, "/"); idx != -1 {
			bucketName = input[:idx]
		}
		region, err = s3region.GetBucketRegionByName(bucketName)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Bucket region: %s\n", region)
}
