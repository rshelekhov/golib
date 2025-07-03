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
	connStr   string
}

// NewTestDB creates a new test database
func NewTestDB(ctx context.Context) (*TestDB, error) {
	// Try to use existing database first
	if connStr := os.Getenv("TEST_DB_CONN_STRING"); connStr != "" {
		return &TestDB{connStr: connStr}, nil
	}

	// Fallback to Docker container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
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

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	connStr := fmt.Sprintf("postgres://test:test@%s:%s/test?sslmode=disable", host, port.Port())

	return &TestDB{
		container: container,
		connStr:   connStr,
	}, nil
}

// ConnStr returns the connection string for the test database
func (db *TestDB) ConnStr() string {
	return db.connStr
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
