package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	region, err := s3region.GetBucketRegion(context.Background(), input)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Bucket region: %s\n", region)
}
