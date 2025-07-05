# S3 Library

This library provides a simple and consistent interface for working with AWS S3 and S3-compatible storage systems.

## Features

- **Connection Management**: Easy S3 client initialization with configurable options
- **Object Operations**: Upload, download, delete, and list objects
- **Bucket Operations**: Create, delete, and list buckets
- **Multipart Uploads**: Support for large file uploads
- **Presigned URLs**: Generate temporary URLs for secure access
- **Helper Functions**: Simplified operations for common use cases
- **OpenTelemetry Integration**: Built-in tracing support
- **Testing Utilities**: Easy testing with MinIO containers

## Installation

```bash
go get github.com/rshelekhov/golib/db/s3
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "log"

    s3lib "github.com/rshelekhov/golib/db/s3"
)

func main() {
    ctx := context.Background()

    // Create connection
    conn, err := s3lib.NewConnection(ctx,
        s3lib.WithRegion("us-east-1"),
        s3lib.WithCredentials("your-access-key", "your-secret-key"),
        s3lib.WithEndpoint("https://s3.amazonaws.com"), // optional for custom endpoints
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Use the connection...
}
```

### Configuration Options

```go
conn, err := s3lib.NewConnection(ctx,
    s3lib.WithRegion("us-west-2"),
    s3lib.WithCredentials("access-key", "secret-key"),
    s3lib.WithEndpoint("https://minio.example.com"),
    s3lib.WithHTTPTimeout(30*time.Second),
    s3lib.WithMaxRetries(5),
    s3lib.WithForcePathStyle(true),
    s3lib.WithDisableSSL(false),
    s3lib.WithTracing(true),
)
```

### Object Operations

```go
// Upload an object
err := conn.PutObjectSimple(ctx, "my-bucket", "my-key", strings.NewReader("content"), "private")

// Download an object
reader, err := conn.GetObjectSimple(ctx, "my-bucket", "my-key")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// Check if object exists
exists, err := conn.ObjectExists(ctx, "my-bucket", "my-key")

// Delete an object
err = conn.DeleteObjectSimple(ctx, "my-bucket", "my-key")
```

### Advanced Object Operations

```go
// Upload with full options
_, err := conn.PutObject(ctx, &s3.PutObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("my-key"),
    Body:   bytes.NewReader(data),
    ACL:    aws.String("private"),
    ContentType: aws.String("application/json"),
})

// Download with full options
result, err := conn.GetObject(ctx, &s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("my-key"),
})
```

### Bucket Operations

```go
// Create bucket
_, err := conn.CreateBucket(ctx, &s3.CreateBucketInput{
    Bucket: aws.String("my-new-bucket"),
})

// List buckets
result, err := conn.ListBuckets(ctx, &s3.ListBucketsInput{})

// Delete bucket
_, err := conn.DeleteBucket(ctx, &s3.DeleteBucketInput{
    Bucket: aws.String("my-bucket"),
})
```

### Presigned URLs

```go
// Generate presigned URL for download (expires in 1 hour)
url, err := conn.GetObjectPresignedURL("my-bucket", "my-key", 3600)

// Generate presigned URL for upload (expires in 1 hour)
url, err := conn.PutObjectPresignedURL("my-bucket", "my-key", 3600)
```

### Multipart Uploads

```go
// Create multipart upload
resp, err := conn.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("large-file"),
})

// Upload parts...
part, err := conn.UploadPart(ctx, &s3.UploadPartInput{
    Bucket:     aws.String("my-bucket"),
    Key:        aws.String("large-file"),
    UploadId:   resp.UploadId,
    PartNumber: aws.Int64(1),
    Body:       bytes.NewReader(partData),
})

// Complete multipart upload
_, err = conn.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
    Bucket:   aws.String("my-bucket"),
    Key:      aws.String("large-file"),
    UploadId: resp.UploadId,
    MultipartUpload: &s3.CompletedMultipartUpload{
        Parts: []*s3.CompletedPart{
            {
                ETag:       part.ETag,
                PartNumber: aws.Int64(1),
            },
        },
    },
})
```

## Testing

The library includes testing utilities for easy integration testing with MinIO:

```go
func TestS3Operations(t *testing.T) {
    ctx := context.Background()

    // Create test container
    container := testutil.NewTestContainer(ctx, t)
    defer container.Close(ctx)

    // Create connection
    conn, err := container.NewTestConnection(ctx)
    require.NoError(t, err)

    // Create test bucket
    err = container.CreateTestBucket(ctx, conn, "test-bucket")
    require.NoError(t, err)

    // Test operations...
}
```

## Environment Variables

The library supports AWS credential chain, which includes:

- `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
- `AWS_PROFILE`
- IAM roles for EC2 instances
- IAM roles for ECS tasks
- IAM roles for Lambda functions

To use credential chain instead of static credentials:

```go
conn, err := s3lib.NewConnection(ctx,
    s3lib.WithCredentialsChain(true),
    s3lib.WithRegion("us-east-1"),
)
```

## Error Handling

The library preserves AWS SDK error types, so you can handle specific errors:

```go
_, err := conn.GetObject(ctx, &s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("non-existent-key"),
})
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case s3.ErrCodeNoSuchKey:
            // Handle missing object
        case s3.ErrCodeNoSuchBucket:
            // Handle missing bucket
        }
    }
}
```

## Constants

- `DefaultHTTPTimeout`: 10 seconds
- `DefaultRegion`: "us-east-1"
- `DefaultACL`: "private"
- `DefaultMaxRetries`: 3
