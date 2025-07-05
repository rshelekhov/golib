package s3

import "time"

const (
	// DefaultHTTPTimeout is the default timeout for S3 HTTP client
	DefaultHTTPTimeout = 10 * time.Second
	// DefaultRegion is the default AWS region
	DefaultRegion = "us-east-1"
	// DefaultACL is the default ACL for S3 objects
	DefaultACL = "private"
	// DefaultMaxRetries is the default number of retries for S3 operations
	DefaultMaxRetries = 3
)
