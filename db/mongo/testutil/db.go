package testutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB represents a test database
type TestDB struct {
	container testcontainers.Container
	uri       string
}

// NewTestDB creates a new test database
func NewTestDB(ctx context.Context) (*TestDB, error) {
	// Try to use existing database first
	if uri := os.Getenv("TEST_MONGO_URI"); uri != "" {
		return &TestDB{uri: uri}, nil
	}

	// Fallback to Docker container
	req := testcontainers.ContainerRequest{
		Image:        "mongo:5",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE": "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	return &TestDB{
		container: container,
		uri:       uri,
	}, nil
}

// URI returns the connection URI for the test database
func (db *TestDB) URI() string {
	return db.uri
}

// Close stops and removes the test database container if it was created
func (db *TestDB) Close(ctx context.Context) error {
	if db.container != nil {
		return db.container.Terminate(ctx)
	}
	return nil
}

// WaitForReady waits for the database to be ready
func (db *TestDB) WaitForReady(ctx context.Context) error {
	// Wait for a short time to ensure the database is ready
	time.Sleep(time.Second)
	return nil
}
