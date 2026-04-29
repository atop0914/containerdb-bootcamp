// Package postgres provides PostgreSQL container management for testing.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/atop0914/containerdb-bootcamp/internal/config"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// New starts a new PostgreSQL container and returns a connected *sql.DB.
// Call the returned cleanup function to stop and remove the container.
//
// Example:
//
//	ctx := context.Background()
//	db, cleanup, err := postgres.New(ctx)
//	if err != nil {
//	    panic(err)
//	}
//	defer cleanup()
func New(ctx context.Context) (*sql.DB, func(), error) {
	cfg := config.DefaultPostgresConfig()
	return NewWithConfig(ctx, cfg)
}

// NewWithConfig starts a PostgreSQL container with custom configuration.
func NewWithConfig(ctx context.Context, cfg *config.PostgresConfig) (*sql.DB, func(), error) {
	dbName := cfg.Database
	dbUser := cfg.Username
	dbPassword := cfg.Password

	pgContainer, err := postgres.Run(ctx,
		cfg.Image,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	host, portStr, err := pgContainer.HostPort(ctx, 5432)
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get host port: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, portStr, dbUser, dbPassword, dbName)

	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to open db: %w", err)
	}

	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(ctx, cfg.HealthCheckTimeout)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("postgres not ready: %w", err)
	}

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(context.Background())
	}

	return pool, cleanup, nil
}
