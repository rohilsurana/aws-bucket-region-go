package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	s3region "github.com/rohilsurana/aws-bucket-region-go"
)

var (
	timeout = flag.Duration("timeout", 10*time.Second, "HTTP request timeout")
	version = flag.Bool("version", false, "Print version information")
	help    = flag.Bool("help", false, "Show help message")
)

const (
	appVersion = "1.0.0"
	appName    = "s3region"
)

func main() {
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: S3 bucket identifier required\n\n")
		printHelp()
		os.Exit(1)
	}

	input := flag.Arg(0)

	// Create custom HTTP client with timeout
	client := &http.Client{
		Timeout: *timeout,
	}

	// Get bucket region
	region, err := s3region.GetBucketRegion(
		context.Background(),
		input,
		s3region.WithHTTPClient(client),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(region)
}

func printHelp() {
	fmt.Printf(`%s - Get AWS S3 bucket region without credentials

Usage:
  %s [options] <s3-identifier>

Arguments:
  <s3-identifier>    S3 bucket identifier in any supported format:
                     - Bucket name: my-bucket
                     - S3 URI: s3://my-bucket or s3://my-bucket/path
                     - AWS ARN: arn:aws:s3:::my-bucket
                     - HTTP URL: https://my-bucket.s3.amazonaws.com

Options:
  -timeout duration  HTTP request timeout (default 10s)
  -version          Print version information
  -help             Show this help message

Examples:
  %s my-bucket
  %s s3://my-bucket/path/to/object
  %s arn:aws:s3:::my-bucket
  %s https://my-bucket.s3.amazonaws.com/object
  %s -timeout 5s my-bucket

`, appName, appName, appName, appName, appName, appName, appName)
}
