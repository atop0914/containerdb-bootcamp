// Package mysql provides MySQL container management for testing.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/atop0914/containerdb-bootcamp/internal/config"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// New starts a new MySQL container and returns a connected *sql.DB.
// Call the returned cleanup function to stop and remove the container.
//
// Example:
//
//	ctx := context.Background()
//	db, cleanup, err := mysql.New(ctx)
//	if err != nil {
//	    panic(err)
//	}
//	defer cleanup()
func New(ctx context.Context) (*sql.DB, func(), error) {
	cfg := config.DefaultMySQLConfig()
	return NewWithConfig(ctx, cfg)
}

// NewWithConfig starts a MySQL container with custom configuration.
func NewWithConfig(ctx context.Context, cfg *config.MySQLConfig) (*sql.DB, func(), error) {
	mysqlContainer, err := mysql.Run(ctx,
		cfg.Image,
		mysql.WithUsername(cfg.Username),
		mysql.WithPassword(cfg.Password),
		mysql.WithDatabase(cfg.Database),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start mysql container: %w", err)
	}

	// Get connection details
	hostPort, err := mysqlContainer.Port(ctx, 3306)
	if err != nil {
		mysqlContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get port: %w", err)
	}

	host, portStr, err := mysqlContainer.HostPort(ctx, 3306)
	if err != nil {
		mysqlContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get host port: %w", err)
	}

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Username, cfg.Password, host, portStr, cfg.Database)

	// Wait for MySQL to be ready (using hostPort which includes the mapped port)
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		mysqlContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Configure pool
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(time.Hour)

	// Verify connection with timeout
	ctx, cancel := context.WithTimeout(ctx, cfg.HealthCheckTimeout)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		mysqlContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("mysql not ready: %w", err)
	}

	cleanup := func() {
		pool.Close()
		mysqlContainer.Terminate(context.Background())
	}

	return pool, cleanup, nil
}

// Container exposes the underlying testcontainers MySQL module for advanced use.
func Container(ctx context.Context) (*mysql.MySQLContainer, error) {
	cfg := config.DefaultMySQLConfig()
	return mysql.Run(ctx, cfg.Image)
}
