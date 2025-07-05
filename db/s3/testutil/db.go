package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	s3lib "github.com/rshelekhov/golib/db/s3"
)

const (
	// MinIOImage is the Docker image for MinIO (S3-compatible storage)
	MinIOImage = "minio/minio:latest"
	// MinIOAccessKey is the default access key for MinIO
	MinIOAccessKey = "minioadmin"
	// MinIOSecretKey is the default secret key for MinIO
	MinIOSecretKey = "minioadmin"
	// MinIOPort is the default port for MinIO
	MinIOPort = "9000"
)

// TestContainer represents a test container for S3-compatible storage
type TestContainer struct {
	Container testcontainers.Container
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
}

// NewTestContainer creates a new test container with MinIO
func NewTestContainer(ctx context.Context, t *testing.T) *TestContainer {
	req := testcontainers.ContainerRequest{
		Image:        MinIOImage,
		ExposedPorts: []string{MinIOPort + "/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     MinIOAccessKey,
			"MINIO_ROOT_PASSWORD": MinIOSecretKey,
		},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForHTTP("/minio/health/ready").WithPort(MinIOPort),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, MinIOPort)
	require.NoError(t, err)

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	return &TestContainer{
		Container: container,
		Endpoint:  endpoint,
		AccessKey: MinIOAccessKey,
		SecretKey: MinIOSecretKey,
		Region:    "us-east-1",
	}
}

// Close terminates the test container
func (tc *TestContainer) Close(ctx context.Context) error {
	return tc.Container.Terminate(ctx)
}

// NewTestConnection creates a new S3 connection for testing
func (tc *TestContainer) NewTestConnection(ctx context.Context) (s3lib.ConnectionAPI, error) {
	return s3lib.NewConnection(ctx,
		s3lib.WithEndpoint(tc.Endpoint),
		s3lib.WithCredentials(tc.AccessKey, tc.SecretKey),
		s3lib.WithRegion(tc.Region),
		s3lib.WithForcePathStyle(true),
		s3lib.WithDisableSSL(true),
		s3lib.WithTracing(false),
	)
}

// CreateTestBucket creates a test bucket for testing
func (tc *TestContainer) CreateTestBucket(ctx context.Context, conn s3lib.ConnectionAPI, bucketName string) error {
	_, err := conn.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	return err
}
