package s3

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		opts    []ConnectionOption
		wantErr bool
	}{
		{
			name: "basic connection with credentials",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "connection with custom endpoint",
			opts: []ConnectionOption{
				WithRegion("us-west-2"),
				WithEndpoint("https://s3.custom.com"),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "connection with custom timeout",
			opts: []ConnectionOption{
				WithRegion("eu-west-1"),
				WithHTTPTimeout(30 * time.Second),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "connection with path style",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithForcePathStyle(true),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "connection with SSL disabled",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithDisableSSL(true),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "connection with tracing disabled",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithTracing(false),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
		{
			name: "connection with credential chain",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithCredentialsChain(true),
			},
			wantErr: false,
		},
		{
			name: "connection with session token",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithCredentials("test-key", "test-secret"),
				WithSessionToken("test-token"),
			},
			wantErr: false,
		},
		{
			name: "connection with max retries",
			opts: []ConnectionOption{
				WithRegion("us-east-1"),
				WithMaxRetries(5),
				WithCredentials("test-key", "test-secret"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := NewConnection(ctx, tt.opts...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)

				// Test that we can get the client
				client := conn.Client()
				assert.NotNil(t, client)

				// Test close
				err = conn.Close()
				assert.NoError(t, err)
			}
		})
	}
}

func TestConnectionOptions(t *testing.T) {
	ctx := context.Background()

	t.Run("default options", func(t *testing.T) {
		conn, err := NewConnection(ctx)
		require.NoError(t, err)
		require.NotNil(t, conn)

		assert.NotNil(t, conn.Client())

		err = conn.Close()
		assert.NoError(t, err)
	})

	t.Run("nil options", func(t *testing.T) {
		conn, err := NewConnection(ctx, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, conn)

		err = conn.Close()
		assert.NoError(t, err)
	})
}

func TestConnectionInterface(t *testing.T) {
	ctx := context.Background()

	conn, err := NewConnection(ctx,
		WithRegion("us-east-1"),
		WithCredentials("test-key", "test-secret"),
	)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	// Test that connection implements all required interfaces
	assert.Implements(t, (*ConnectionAPI)(nil), conn)
	assert.Implements(t, (*ConnectionCloser)(nil), conn)
	assert.Implements(t, (*ObjectAPI)(nil), conn)
	assert.Implements(t, (*BucketAPI)(nil), conn)
	assert.Implements(t, (*MultipartAPI)(nil), conn)
	assert.Implements(t, (*PresignedAPI)(nil), conn)
	assert.Implements(t, (*HelperAPI)(nil), conn)
}

func TestPresignedURLs(t *testing.T) {
	ctx := context.Background()

	conn, err := NewConnection(ctx,
		WithRegion("us-east-1"),
		WithCredentials("test-key", "test-secret"),
	)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	bucket := "test-bucket"
	key := "test-key"
	expires := int64(3600)

	t.Run("get object presigned URL", func(t *testing.T) {
		url, err := conn.GetObjectPresignedURL(bucket, key, expires)
		assert.NoError(t, err)
		assert.NotEmpty(t, url)
		assert.Contains(t, url, bucket)
		assert.Contains(t, url, key)
	})

	t.Run("put object presigned URL", func(t *testing.T) {
		url, err := conn.PutObjectPresignedURL(bucket, key, expires)
		assert.NoError(t, err)
		assert.NotEmpty(t, url)
		assert.Contains(t, url, bucket)
		assert.Contains(t, url, key)
	})
}

// BenchmarkNewConnection benchmarks connection creation
func BenchmarkNewConnection(b *testing.B) {
	ctx := context.Background()
	opts := []ConnectionOption{
		WithRegion("us-east-1"),
		WithCredentials("test-key", "test-secret"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := NewConnection(ctx, opts...)
		if err != nil {
			b.Fatal(err)
		}
		conn.Close()
	}
}
