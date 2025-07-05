package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
)

// ConnectionCloser defines the interface for connection management.
type ConnectionCloser interface {
	// Close closes the connection.
	Close() error
	// Client returns the S3 client instance.
	Client() *s3.S3
	// Ping checks the connection to the S3 service.
	Ping(ctx context.Context) error
}

// ObjectAPI defines the interface for object operations.
type ObjectAPI interface {
	// PutObject uploads an object to S3.
	PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	// GetObject downloads an object from S3.
	GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	// DeleteObject deletes an object from S3.
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	// HeadObject retrieves metadata for an object without downloading it.
	HeadObject(ctx context.Context, input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
	// CopyObject copies an object from one location to another.
	CopyObject(ctx context.Context, input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error)
	// ListObjects lists objects in a bucket.
	ListObjects(ctx context.Context, input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	// ListObjectsV2 lists objects in a bucket using the V2 API.
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
}

// BucketAPI defines the interface for bucket operations.
type BucketAPI interface {
	// CreateBucket creates a new bucket.
	CreateBucket(ctx context.Context, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	// DeleteBucket deletes a bucket.
	DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)
	// ListBuckets lists all buckets.
	ListBuckets(ctx context.Context, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	// HeadBucket checks if a bucket exists and is accessible.
	HeadBucket(ctx context.Context, input *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	// GetBucketLocation retrieves the region where a bucket is located.
	GetBucketLocation(ctx context.Context, input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error)
}

// MultipartAPI defines the interface for multipart upload operations.
type MultipartAPI interface {
	// CreateMultipartUpload initiates a multipart upload.
	CreateMultipartUpload(ctx context.Context, input *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)
	// UploadPart uploads a part of a multipart upload.
	UploadPart(ctx context.Context, input *s3.UploadPartInput) (*s3.UploadPartOutput, error)
	// CompleteMultipartUpload completes a multipart upload.
	CompleteMultipartUpload(ctx context.Context, input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)
	// AbortMultipartUpload aborts a multipart upload.
	AbortMultipartUpload(ctx context.Context, input *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error)
	// ListMultipartUploads lists in-progress multipart uploads.
	ListMultipartUploads(ctx context.Context, input *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error)
	// ListParts lists parts of a multipart upload.
	ListParts(ctx context.Context, input *s3.ListPartsInput) (*s3.ListPartsOutput, error)
}

// PresignedAPI defines the interface for presigned URL operations.
type PresignedAPI interface {
	// GetObjectPresignedURL generates a presigned URL for GetObject operation.
	GetObjectPresignedURL(bucket, key string, expires int64) (string, error)
	// PutObjectPresignedURL generates a presigned URL for PutObject operation.
	PutObjectPresignedURL(bucket, key string, expires int64) (string, error)
}

// HelperAPI defines the interface for helper operations.
type HelperAPI interface {
	// PutObjectSimple uploads data to S3 with simple parameters.
	PutObjectSimple(ctx context.Context, bucket, key string, data io.Reader, acl string) error
	// GetObjectSimple downloads data from S3 with simple parameters.
	GetObjectSimple(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	// DeleteObjectSimple deletes an object from S3 with simple parameters.
	DeleteObjectSimple(ctx context.Context, bucket, key string) error
	// ObjectExists checks if an object exists in S3.
	ObjectExists(ctx context.Context, bucket, key string) (bool, error)
}

// ConnectionAPI defines the interface for all S3 operations.
type ConnectionAPI interface {
	ConnectionCloser
	ObjectAPI
	BucketAPI
	MultipartAPI
	PresignedAPI
	HelperAPI
}
