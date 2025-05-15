package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoTransactionManager manages MongoDB transactions.
type MongoTransactionManager struct {
	conn *Connection
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(conn *Connection) TransactionManager {
	return &MongoTransactionManager{conn: conn}
}

// RunTransaction executes the given function within a transaction.
func (m *MongoTransactionManager) RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := m.conn.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err := fn(sessCtx); err != nil {
			return nil, fmt.Errorf("transaction execution failed: %w", err)
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}
