# aws-bucket-region-go

A simple Go package to detect AWS S3 bucket region without requiring AWS SDK or credentials.

## How it works

This package performs a HEAD request to an S3 bucket URL and extracts the region from the `x-amz-bucket-region` HTTP response header.

## Installation

```bash
go get github.com/rohilsurana/aws-bucket-region-go
```

## Usage

### Using bucket name

```go
package main

import (
    "fmt"
    "log"

    "github.com/rohilsurana/aws-bucket-region-go"
)

func main() {
    region, err := s3region.GetBucketRegionByName("my-bucket")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bucket region: %s\n", region)
}
```

### Using full URL

```go
region, err := s3region.GetBucketRegion("https://my-bucket.s3.amazonaws.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Bucket region: %s\n", region)
```

## Example

Run the example with bucket name, s3:// URI, AWS ARN, or full URL (with or without object paths):

```bash
cd example
go run main.go my-bucket
go run main.go my-bucket/path/to/object
go run main.go s3://my-bucket/path/to/object
go run main.go arn:aws:s3:::my-bucket/path
go run main.go https://my-bucket.s3.amazonaws.com/path/to/object
```

## API

### `GetBucketRegionByName(bucketName string) (string, error)`

Takes a bucket name and returns its region by constructing the S3 URL and performing a HEAD request.

**Parameters:**
- `bucketName`: S3 bucket name (e.g., `my-bucket`)

**Returns:**
- `string`: The AWS region code (e.g., `us-west-2`)
- `error`: Error if the request fails or the region header is missing

### `GetBucketRegion(url string) (string, error)`

Performs a HEAD request to the given S3 URL and returns the bucket region.

**Parameters:**
- `url`: Full S3 bucket URL (e.g., `https://bucket-name.s3.amazonaws.com`)

**Returns:**
- `string`: The AWS region code (e.g., `us-west-2`)
- `error`: Error if the request fails or the region header is missing

### Error Variables

- `ErrRegionHeaderNotFound`: Returned when the `x-amz-bucket-region` header is not found
- `ErrBucketNotFound`: Returned when the bucket does not exist (404 response)

## License

MIT
