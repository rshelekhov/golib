package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionManager defines the interface for transaction management.
type TransactionManager interface {
	// RunTransaction executes the given function within a transaction.
	RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// ConnectionCloser defines the interface for connection management.
type ConnectionCloser interface {
	// Close closes the connection.
	Close(ctx context.Context) error
	// Database returns the database instance.
	Database() *mongo.Database
	// Client returns the client instance.
	Client() *mongo.Client
	// Ping checks the connection to the database.
	Ping(ctx context.Context) error
}

// Inserter defines the interface for insert operations.
type Inserter interface {
	// InsertOne inserts a single document into the collection.
	InsertOne(ctx context.Context, collection string, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	// InsertMany inserts multiple documents into the collection.
	InsertMany(ctx context.Context, collection string, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
}

// Finder defines the interface for find operations.
type Finder interface {
	// FindOne finds a single document in the collection.
	FindOne(ctx context.Context, collection string, filter any, result any, opts ...*options.FindOneOptions) error
	// Find finds documents in the collection.
	Find(ctx context.Context, collection string, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

// Updater defines the interface for update operations.
type Updater interface {
	// UpdateOne updates a single document in the collection.
	UpdateOne(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// UpdateMany updates multiple documents in the collection.
	UpdateMany(ctx context.Context, collection string, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// Deleter defines the interface for delete operations.
type Deleter interface {
	// DeleteOne deletes a single document from the collection.
	DeleteOne(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	// DeleteMany deletes multiple documents from the collection.
	DeleteMany(ctx context.Context, collection string, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// Counter defines the interface for count operations.
type Counter interface {
	// CountDocuments counts the number of documents in the collection.
	CountDocuments(ctx context.Context, collection string, filter any, opts ...*options.CountOptions) (int64, error)
}

// Aggregator defines the interface for aggregation operations.
type Aggregator interface {
	// Aggregate performs an aggregation operation on the collection.
	Aggregate(ctx context.Context, collection string, pipeline any, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
}

// ConnectionManager defines the interface for all database operations.
type ConnectionManager interface {
	ConnectionCloser
	Inserter
	Finder
	Updater
	Deleter
	Counter
	Aggregator
}
