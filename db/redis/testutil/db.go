package testutil

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB represents a test Redis database
type TestDB struct {
	container testcontainers.Container
	host      string
	port      int
	password  string
	db        int
}

// NewTestDB creates a new test Redis database
func NewTestDB(ctx context.Context) (*TestDB, error) {
	// Try to use existing database first
	if addr := os.Getenv("TEST_REDIS_ADDR"); addr != "" {
		host, port, err := parseRedisAddr(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TEST_REDIS_ADDR: %w", err)
		}

		password := os.Getenv("TEST_REDIS_PASSWORD")
		dbStr := os.Getenv("TEST_REDIS_DB")
		db := 0
		if dbStr != "" {
			if parsedDB, err := strconv.Atoi(dbStr); err == nil {
				db = parsedDB
			}
		}

		return &TestDB{
			host:     host,
			port:     port,
			password: password,
			db:       db,
		}, nil
	}

	// Fallback to Docker container
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections"),
			wait.ForListeningPort("6379/tcp"),
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

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		return nil, fmt.Errorf("failed to parse port: %w", err)
	}

	return &TestDB{
		container: container,
		host:      host,
		port:      portInt,
		db:        0,
	}, nil
}

// Host returns the Redis host
func (db *TestDB) Host() string {
	return db.host
}

// Port returns the Redis port
func (db *TestDB) Port() int {
	return db.port
}

// Password returns the Redis password
func (db *TestDB) Password() string {
	return db.password
}

// DB returns the Redis database number
func (db *TestDB) DB() int {
	return db.db
}

// Addr returns the Redis address in host:port format
func (db *TestDB) Addr() string {
	return fmt.Sprintf("%s:%d", db.host, db.port)
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

// parseRedisAddr parses Redis address in format host:port
func parseRedisAddr(addr string) (host string, port int, err error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address format, expected host:port")
	}

	host = parts[0]
	port, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %w", err)
	}

	return host, port, nil
}
