// Copyright (c) 2025 A Bit of Help, Inc.

// Package interfaces provides database interfaces for the db package.
package interfaces

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

// SQLTxInterface defines the interface for sql.Tx operations
type SQLTxInterface interface {
	Commit() error
	Rollback() error
}

// MongoClientInterface defines the interface for mongo.Client operations
type MongoClientInterface interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Connect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

// PgxPoolInterface defines the interface for pgxpool.Pool operations
type PgxPoolInterface interface {
	Ping(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// PgxTxInterface defines the interface for pgx.Tx operations
type PgxTxInterface interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}