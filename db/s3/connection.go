package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Connection represents a connection to S3.
type Connection struct {
	client *s3.S3
	tracer trace.Tracer
}

// connectionOptions holds configuration for S3 connection
type connectionOptions struct {
	region           string
	endpoint         string
	accessKey        string
	secretKey        string
	sessionToken     string
	httpTimeout      time.Duration
	maxRetries       int
	forcePathStyle   bool
	disableSSL       bool
	enableTracing    bool
	credentialsChain bool
}

// ConnectionOption is a function that configures connection options.
type ConnectionOption func(opts *connectionOptions)

// WithRegion sets the AWS region.
func WithRegion(region string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.region = region
	}
}

// WithEndpoint sets the S3 endpoint URL.
func WithEndpoint(endpoint string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.endpoint = endpoint
	}
}

// WithCredentials sets the AWS credentials.
func WithCredentials(accessKey, secretKey string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.accessKey = accessKey
		opts.secretKey = secretKey
	}
}

// WithSessionToken sets the AWS session token for temporary credentials.
func WithSessionToken(sessionToken string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.sessionToken = sessionToken
	}
}

// WithHTTPTimeout sets the HTTP client timeout.
func WithHTTPTimeout(timeout time.Duration) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.httpTimeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(maxRetries int) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.maxRetries = maxRetries
	}
}

// WithForcePathStyle forces path-style addressing.
func WithForcePathStyle(enable bool) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.forcePathStyle = enable
	}
}

// WithDisableSSL disables SSL for the connection.
func WithDisableSSL(disable bool) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.disableSSL = disable
	}
}

// WithTracing turns on/off tracing through OpenTelemetry.
func WithTracing(enable bool) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.enableTracing = enable
	}
}

// WithCredentialsChain uses the AWS credentials chain instead of static credentials.
func WithCredentialsChain(enable bool) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.credentialsChain = enable
	}
}

// WithMinIOEndpoint is a convenience function for MinIO endpoints.
// It automatically enables PathStyle and disables SSL if no scheme provided.
func WithMinIOEndpoint(endpoint string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.endpoint = endpoint
		opts.forcePathStyle = true

		// Auto-detect if SSL should be disabled for local MinIO
		if !strings.HasPrefix(endpoint, "https://") &&
			!strings.HasPrefix(endpoint, "http://") {
			opts.disableSSL = true
		}
	}
}

// NewConnection creates a new connection to S3.
func NewConnection(ctx context.Context, opts ...ConnectionOption) (ConnectionAPI, error) {
	// Apply default options
	connOpts := &connectionOptions{
		region:        DefaultRegion,
		httpTimeout:   DefaultHTTPTimeout,
		maxRetries:    DefaultMaxRetries,
		enableTracing: true, // default is true
	}

	for _, opt := range opts {
		if opt != nil {
			opt(connOpts)
		}
	}

	// Create AWS config
	cfg := aws.NewConfig().
		WithHTTPClient(&http.Client{
			Timeout: connOpts.httpTimeout,
		}).
		WithRegion(connOpts.region).
		WithMaxRetries(connOpts.maxRetries).
		WithS3ForcePathStyle(connOpts.forcePathStyle).
		WithDisableSSL(connOpts.disableSSL).
		WithCredentialsChainVerboseErrors(true)

	// Set endpoint if provided
	if connOpts.endpoint != "" {
		cfg = cfg.WithEndpoint(connOpts.endpoint)
	}

	// Auto-enable PathStyle for MinIO endpoints if not explicitly set
	if connOpts.endpoint != "" && IsMinIOEndpoint(connOpts.endpoint) && !connOpts.forcePathStyle {
		connOpts.forcePathStyle = true
	}

	// Set credentials
	if !connOpts.credentialsChain && connOpts.accessKey != "" && connOpts.secretKey != "" {
		cfg = cfg.WithCredentials(credentials.NewStaticCredentials(
			connOpts.accessKey,
			connOpts.secretKey,
			connOpts.sessionToken,
		))
	}

	// Create session
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	// Create S3 client
	client := s3.New(sess)

	conn := &Connection{
		client: client,
	}

	// Set up tracing
	if connOpts.enableTracing {
		conn.tracer = otel.Tracer("s3")
	}

	return conn, nil
}

// IsMinIOEndpoint detects if endpoint looks like MinIO
func IsMinIOEndpoint(endpoint string) bool {
	return strings.Contains(endpoint, "minio") ||
		strings.Contains(endpoint, ":9000") ||
		strings.Contains(endpoint, ":9001")
}

// Close closes the connection to S3.
func (c *Connection) Close() error {
	// S3 client doesn't require explicit closing
	return nil
}

// Client returns the S3 client.
func (c *Connection) Client() *s3.S3 {
	return c.client
}

// Ping checks the connection to the S3 service.
func (c *Connection) Ping(ctx context.Context) error {
	_, err := c.client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to ping S3: %w", err)
	}
	return nil
}

// Object operations

// PutObject uploads an object to S3.
func (c *Connection) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return c.client.PutObjectWithContext(ctx, input)
}

// GetObject downloads an object from S3.
func (c *Connection) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.client.GetObjectWithContext(ctx, input)
}

// DeleteObject deletes an object from S3.
func (c *Connection) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return c.client.DeleteObjectWithContext(ctx, input)
}

// HeadObject retrieves metadata for an object without downloading it.
func (c *Connection) HeadObject(ctx context.Context, input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return c.client.HeadObjectWithContext(ctx, input)
}

// CopyObject copies an object from one location to another.
func (c *Connection) CopyObject(ctx context.Context, input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return c.client.CopyObjectWithContext(ctx, input)
}

// ListObjects lists objects in a bucket.
func (c *Connection) ListObjects(ctx context.Context, input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return c.client.ListObjectsWithContext(ctx, input)
}

// ListObjectsV2 lists objects in a bucket using the V2 API.
func (c *Connection) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return c.client.ListObjectsV2WithContext(ctx, input)
}

// Bucket operations

// CreateBucket creates a new bucket.
func (c *Connection) CreateBucket(ctx context.Context, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return c.client.CreateBucketWithContext(ctx, input)
}

// DeleteBucket deletes a bucket.
func (c *Connection) DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return c.client.DeleteBucketWithContext(ctx, input)
}

// ListBuckets lists all buckets.
func (c *Connection) ListBuckets(ctx context.Context, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return c.client.ListBucketsWithContext(ctx, input)
}

// HeadBucket checks if a bucket exists and is accessible.
func (c *Connection) HeadBucket(ctx context.Context, input *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return c.client.HeadBucketWithContext(ctx, input)
}

// GetBucketLocation retrieves the region where a bucket is located.
func (c *Connection) GetBucketLocation(ctx context.Context, input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	return c.client.GetBucketLocationWithContext(ctx, input)
}

// Multipart operations

// CreateMultipartUpload initiates a multipart upload.
func (c *Connection) CreateMultipartUpload(ctx context.Context, input *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	return c.client.CreateMultipartUploadWithContext(ctx, input)
}

// UploadPart uploads a part of a multipart upload.
func (c *Connection) UploadPart(ctx context.Context, input *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
	return c.client.UploadPartWithContext(ctx, input)
}

// CompleteMultipartUpload completes a multipart upload.
func (c *Connection) CompleteMultipartUpload(ctx context.Context, input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
	return c.client.CompleteMultipartUploadWithContext(ctx, input)
}

// AbortMultipartUpload aborts a multipart upload.
func (c *Connection) AbortMultipartUpload(ctx context.Context, input *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error) {
	return c.client.AbortMultipartUploadWithContext(ctx, input)
}

// ListMultipartUploads lists in-progress multipart uploads.
func (c *Connection) ListMultipartUploads(ctx context.Context, input *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
	return c.client.ListMultipartUploadsWithContext(ctx, input)
}

// ListParts lists parts of a multipart upload.
func (c *Connection) ListParts(ctx context.Context, input *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
	return c.client.ListPartsWithContext(ctx, input)
}

// Presigned URL operations

// GetObjectPresignedURL generates a presigned URL for GetObject operation.
func (c *Connection) GetObjectPresignedURL(bucket, key string, expires int64) (string, error) {
	req, _ := c.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(time.Duration(expires) * time.Second)
}

// PutObjectPresignedURL generates a presigned URL for PutObject operation.
func (c *Connection) PutObjectPresignedURL(bucket, key string, expires int64) (string, error) {
	req, _ := c.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(time.Duration(expires) * time.Second)
}

// Helper operations

// PutObjectSimple uploads data to S3 with simple parameters.
func (c *Connection) PutObjectSimple(ctx context.Context, bucket, key string, data io.Reader, acl string) error {
	if acl == "" {
		acl = DefaultACL
	}

	_, err := c.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(data),
		ACL:    aws.String(acl),
	})
	return err
}

// GetObjectSimple downloads data from S3 with simple parameters.
func (c *Connection) GetObjectSimple(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	result, err := c.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

// DeleteObjectSimple deletes an object from S3 with simple parameters.
func (c *Connection) DeleteObjectSimple(ctx context.Context, bucket, key string) error {
	_, err := c.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// ObjectExists checks if an object exists in S3.
func (c *Connection) ObjectExists(ctx context.Context, bucket, key string) (bool, error) {
	_, err := c.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey, "NotFound":
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
