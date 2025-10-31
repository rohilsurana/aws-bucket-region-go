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

```go
package main

import (
    "fmt"
    "log"

    "github.com/rohilsurana/aws-bucket-region-go"
)

func main() {
    // Works with any format - bucket name, S3 URI, ARN, or HTTP/HTTPS URL
    region, err := s3region.GetBucketRegion("my-bucket")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bucket region: %s\n", region)
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

### `GetBucketRegion(input string) (string, error)`

Main function that automatically detects the input format and returns the bucket region. Supports all formats below.

**Parameters:**
- `input`: Any S3 identifier format (bucket name, S3 URI, ARN, or HTTP/HTTPS URL)

**Returns:**
- `string`: The AWS region code (e.g., `us-west-2`)
- `error`: Error if the request fails or the region header is missing

### Format-Specific Functions

Power users can call these directly if they know the input format:

#### `GetBucketRegionByName(bucketName string) (string, error)`

Takes a bucket name and returns its region.

**Parameters:**
- `bucketName`: S3 bucket name (e.g., `my-bucket`)

#### `GetBucketRegionFromS3URI(uri string) (string, error)`

Extracts bucket name from S3 URI and returns its region.

**Parameters:**
- `uri`: S3 URI (e.g., `s3://my-bucket` or `s3://my-bucket/path/to/object`)

#### `GetBucketRegionFromARN(arn string) (string, error)`

Extracts bucket name from AWS ARN and returns its region.

**Parameters:**
- `arn`: AWS S3 ARN (e.g., `arn:aws:s3:::my-bucket` or `arn:aws:s3:::my-bucket/path`)

#### `GetBucketRegionFromHTTPURL(url string) (string, error)`

Extracts bucket name from HTTP/HTTPS URL and returns its region. Supports both virtual-hosted-style and path-style URLs.

**Parameters:**
- `url`: HTTP/HTTPS URL in either format:
  - Virtual-hosted: `https://my-bucket.s3.amazonaws.com/path/to/object`
  - Path-style: `https://s3.amazonaws.com/my-bucket/path/to/object`
  - Path-style with region: `https://s3.us-west-2.amazonaws.com/my-bucket/path/to/object`

### Error Variables

- `ErrInvalidBucketName`: Returned when the bucket name doesn't follow AWS S3 naming rules
- `ErrRegionHeaderNotFound`: Returned when the `x-amz-bucket-region` header is not found
- `ErrBucketNotFound`: Returned when the bucket does not exist (404 response)

## License

MIT
