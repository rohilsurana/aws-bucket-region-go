package s3region

import (
	"errors"
	"fmt"
)

var ErrRegionHeaderNotFound = errors.New("x-amz-bucket-region header not found in response")
var ErrBucketNotFound = errors.New("aws s3 bucket not found") // HEAD request returns 404
var ErrInvalidBucketName = errors.New("invalid S3 bucket name")

// Error provides structured error information with context about the operation.
type Error struct {
	Op         string // Operation: "GetBucketRegion", "GetBucketRegionByName", etc.
	BucketName string // The bucket name being queried
	Input      string // Original input provided by user
	Err        error  // Underlying error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s(%q): %v", e.Op, e.Input, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// newError creates a new Error with the given parameters.
func newError(op, bucketName, input string, err error) error {
	return &Error{
		Op:         op,
		BucketName: bucketName,
		Input:      input,
		Err:        err,
	}
}
