package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// Connection represents a connection to MongoDB.
type Connection struct {
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration
}

// connectionOptions holds configuration for MongoDB connection
type connectionOptions struct {
	enableTracing bool
	timeout       *time.Duration
	serverAPI     *string
}

// ConnectionOption is a function that configures connection options.
type ConnectionOption func(opts *connectionOptions)

// WithTimeout sets the connection timeout.
func WithTimeout(d time.Duration) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.timeout = &d
	}
}

// WithServerAPI sets the server API version.
func WithServerAPI(version string) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.serverAPI = &version
	}
}

// WithTracing turns on/off tracing through otelmongo
func WithTracing(enable bool) ConnectionOption {
	return func(opts *connectionOptions) {
		opts.enableTracing = enable
	}
}

// NewConnection creates a new connection to MongoDB.
func NewConnection(ctx context.Context, uri string, dbName string, opts ...ConnectionOption) (ConnectionManager, error) {
	clientOpts := options.Client().ApplyURI(uri)

	// Apply default options
	connOpts := &connectionOptions{
		enableTracing: true, // default is true
	}

	for _, opt := range opts {
		if opt != nil {
			opt(connOpts)
		}
	}

	// Apply tracing if enabled
	if connOpts.enableTracing {
		clientOpts.SetMonitor(otelmongo.NewMonitor())
	}

	// Apply timeout
	if connOpts.timeout != nil {
		clientOpts.SetConnectTimeout(*connOpts.timeout)
	} else {
		clientOpts.SetConnectTimeout(DefaultConnectionTimeout)
	}

	// Apply server API
	if connOpts.serverAPI != nil {
		clientOpts.SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	conn := &Connection{
		client:   client,
		database: client.Database(dbName),
		timeout:  DefaultConnectionTimeout,
	}

	return conn, nil
}

// Close closes the connection to MongoDB.
func (c *Connection) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// Database returns the MongoDB database.
func (c *Connection) Database() *mongo.Database {
	return c.database
}

// Client returns the MongoDB client.
func (c *Connection) Client() *mongo.Client {
	return c.client
}

// InsertOne inserts a single document into the collection.
func (c *Connection) InsertOne(ctx context.Context, collection string, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := c.database.Collection(collection).InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return result, nil
}

// InsertMany inserts multiple documents into the collection.
func (c *Connection) InsertMany(ctx context.Context, collection string, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := c.database.Collection(collection).InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}
	return result, nil
}

// FindOne finds a single document in the collection.
func (c *Connection) FindOne(ctx context.Context, collection string, filter any, result any, opts ...*options.FindOneOptions) error {
	err := c.database.Collection(collection).FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}
	return nil
}

// Find finds documents in the collection.
func (c *Connection) Find(ctx context.Context, collection string, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := c.database.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	return cursor, nil
}

// UpdateOne updates a single document in the collection.
func (c *Connection) UpdateOne(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := c.database.Collection(collection).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

// UpdateMany updates multiple documents in the collection.
func (c *Connection) UpdateMany(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := c.database.Collection(collection).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %w", err)
	}
	return result, nil
}

// DeleteOne deletes a single document from the collection.
func (c *Connection) DeleteOne(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := c.database.Collection(collection).DeleteOne(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}

// DeleteMany deletes multiple documents from the collection.
func (c *Connection) DeleteMany(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := c.database.Collection(collection).DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to delete documents: %w", err)
	}
	return result, nil
}

// CountDocuments counts the number of documents in the collection.
func (c *Connection) CountDocuments(ctx context.Context, collection string, filter any, opts ...*options.CountOptions) (int64, error) {
	count, err := c.database.Collection(collection).CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Aggregate performs an aggregation operation on the collection.
func (c *Connection) Aggregate(ctx context.Context, collection string, pipeline any, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := c.database.Collection(collection).Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate documents: %w", err)
	}
	return cursor, nil
}

// Ping checks if the MongoDB server is available.
func (c *Connection) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, nil)
}
