// Copyright (c) 2025 A Bit of Help, Inc.

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// SQLDBInterface defines the interface for sql.DB operations
type SQLDBInterface interface {
	PingContext(ctx context.Context) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// For type assertion to ensure sql.DB implements SQLDBInterface
var _ SQLDBInterface = (*sql.DB)(nil)

// SQLTxInterface defines the interface for sql.Tx operations
type SQLTxInterface interface {
	Commit() error
	Rollback() error
}

// For type assertion to ensure sql.Tx implements SQLTxInterface
var _ SQLTxInterface = (*sql.Tx)(nil)

// MongoClientInterface defines the interface for mongo.Client operations
type MongoClientInterface interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Connect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

// For type assertion to ensure mongo.Client implements MongoClientInterface
var _ MongoClientInterface = (*mongo.Client)(nil)

// PgxPoolInterface defines the interface for pgxpool.Pool operations
type PgxPoolInterface interface {
	Ping(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// For type assertion to ensure pgxpool.Pool implements PgxPoolInterface
var _ PgxPoolInterface = (*pgxpool.Pool)(nil)

// PgxTxInterface defines the interface for pgx.Tx operations
type PgxTxInterface interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// For type assertion to ensure pgx.Tx implements PgxTxInterface
var _ PgxTxInterface = (pgx.Tx)(nil)
