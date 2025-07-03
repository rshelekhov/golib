package pgxv5

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type key string

const (
	txKey key = "tx"
)

const (
	maxConnIdleTimeDefault     = time.Minute
	maxConnLifeTimeDefault     = time.Hour
	minConnectionsCountDefault = 2
	maxConnectionsCountDefault = 10
)

// TxAccessMode is the transaction access mode (read write or read only)
type TxAccessMode = pgx.TxAccessMode

// Transaction access modes
const (
	ReadWrite = pgx.ReadWrite
	ReadOnly  = pgx.ReadOnly
)
