# aws-bucket-region-go

A simple Go package to detect AWS S3 bucket region without requiring AWS SDK or credentials.

## How it works

This package performs a HEAD request to an S3 bucket URL and extracts the region from the `x-amz-bucket-region` HTTP response header. This uses the [AWS S3 HeadBucket API](https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html) which returns the bucket's region in the response headers.

**Note:** Since this makes an actual HTTP HEAD request to AWS servers:
- There will be network latency (typically 100-500ms depending on your location and the bucket's region)
- Requires network access to `*.s3.amazonaws.com`
- No AWS credentials or SDK required

## Supported Formats

- **Bucket name**: `my-bucket` or `my-bucket/path/to/object`
- **S3 URI**: `s3://my-bucket` or `s3://my-bucket/path/to/object`
- **AWS ARN**: `arn:aws:s3:::my-bucket` or `arn:aws:s3:::my-bucket/path`
- **Virtual-hosted-style URL**: `https://my-bucket.s3.amazonaws.com/path/to/object`
- **Path-style URL**: `https://s3.amazonaws.com/my-bucket/path/to/object`
- **Path-style URL with region**: `https://s3.us-west-2.amazonaws.com/my-bucket/path`

## Installation

```bash
go get github.com/rohilsurana/aws-bucket-region-go
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/rohilsurana/aws-bucket-region-go"
)

func main() {
    // Works with any format - bucket name, S3 URI, ARN, or HTTP/HTTPS URL
    region, err := s3region.GetBucketRegion(context.Background(), "my-bucket")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bucket region: %s\n", region)
}
```

### Using Custom HTTP Client

You can provide a custom HTTP client for advanced use cases like custom timeouts, proxies, or TLS configuration:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/rohilsurana/aws-bucket-region-go"
)

func main() {
    // Create a custom HTTP client with 5-second timeout
    customClient := &http.Client{
        Timeout: 5 * time.Second,
    }

    region, err := s3region.GetBucketRegion(
        context.Background(),
        "my-bucket",
        s3region.WithHTTPClient(customClient),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bucket region: %s\n", region)
}
```

You can also implement the `HTTPClient` interface for custom behavior:

```go
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}
```

## Example

Run the example with any supported format:

```bash
cd example

# Bucket name
go run main.go my-bucket
go run main.go my-bucket/path/to/object

# S3 URI
go run main.go s3://my-bucket/path/to/object

# AWS ARN
go run main.go arn:aws:s3:::my-bucket/path

# Virtual-hosted-style URL
go run main.go https://my-bucket.s3.amazonaws.com/path/to/object

# Path-style URL
go run main.go https://s3.amazonaws.com/my-bucket/path/to/object
go run main.go https://s3.us-west-2.amazonaws.com/my-bucket/path
```

## API

### `GetBucketRegion(ctx context.Context, input string, opts ...Option) (string, error)`

Main function that automatically detects the input format and returns the bucket region. Supports all formats below.

**Parameters:**
- `ctx`: Context for timeout and cancellation control
- `input`: Any S3 identifier format (bucket name, S3 URI, ARN, or HTTP/HTTPS URL)
- `opts`: Optional configuration options (e.g., `WithHTTPClient`)

**Returns:**
- `string`: The AWS region code (e.g., `us-west-2`)
- `error`: Error if the request fails or the region header is missing

### Format-Specific Functions

Power users can call these directly if they know the input format:

#### `GetBucketRegionByName(ctx context.Context, bucketName string, opts ...Option) (string, error)`

Takes a bucket name and returns its region.

**Parameters:**
- `ctx`: Context for timeout and cancellation control
- `bucketName`: S3 bucket name (e.g., `my-bucket`)
- `opts`: Optional configuration options

#### `GetBucketRegionFromS3URI(ctx context.Context, uri string, opts ...Option) (string, error)`

Extracts bucket name from S3 URI and returns its region.

**Parameters:**
- `ctx`: Context for timeout and cancellation control
- `uri`: S3 URI (e.g., `s3://my-bucket` or `s3://my-bucket/path/to/object`)
- `opts`: Optional configuration options

#### `GetBucketRegionFromARN(ctx context.Context, arn string, opts ...Option) (string, error)`

Extracts bucket name from AWS ARN and returns its region.

**Parameters:**
- `ctx`: Context for timeout and cancellation control
- `arn`: AWS S3 ARN (e.g., `arn:aws:s3:::my-bucket` or `arn:aws:s3:::my-bucket/path`)
- `opts`: Optional configuration options

#### `GetBucketRegionFromHTTPURL(ctx context.Context, url string, opts ...Option) (string, error)`

Extracts bucket name from HTTP/HTTPS URL and returns its region. Supports both virtual-hosted-style and path-style URLs.

**Parameters:**
- `ctx`: Context for timeout and cancellation control
- `url`: HTTP/HTTPS URL in either format:
  - Virtual-hosted: `https://my-bucket.s3.amazonaws.com/path/to/object`
  - Path-style: `https://s3.amazonaws.com/my-bucket/path/to/object`
  - Path-style with region: `https://s3.us-west-2.amazonaws.com/my-bucket/path/to/object`
- `opts`: Optional configuration options

### Configuration Options

#### `WithHTTPClient(client HTTPClient) Option`

Sets a custom HTTP client for S3 requests. If not provided, `http.DefaultClient` is used.

**Example:**
```go
customClient := &http.Client{Timeout: 10 * time.Second}
region, err := s3region.GetBucketRegion(ctx, "my-bucket", s3region.WithHTTPClient(customClient))
```

### Error Variables

- `ErrInvalidBucketName`: Returned when the bucket name doesn't follow AWS S3 naming rules
- `ErrRegionHeaderNotFound`: Returned when the `x-amz-bucket-region` header is not found
- `ErrBucketNotFound`: Returned when the bucket does not exist (404 response)

## License

MIT
