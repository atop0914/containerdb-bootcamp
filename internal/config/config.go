// Package config provides configuration types for containerized databases.
package config

import (
	"time"
)

// MySQLConfig holds configuration for MySQL containers.
type MySQLConfig struct {
	// Image is the Docker image to use. Defaults to mysql:8.0
	Image string
	// Username for MySQL. Defaults to root
	Username string
	// Password for MySQL. Defaults to rootpassword
	Password string
	// Database name. Defaults to testdb
	Database string
	// HealthCheckTimeout for container readiness. Defaults to 30s
	HealthCheckTimeout time.Duration
}

// DefaultMySQLConfig returns a config with sensible defaults.
func DefaultMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		Image:              "mysql:8.0",
		Username:           "root",
		Password:           "rootpassword",
		Database:           "testdb",
		HealthCheckTimeout: 30 * time.Second,
	}
}

// PostgresConfig holds configuration for PostgreSQL containers.
type PostgresConfig struct {
	// Image is the Docker image to use. Defaults to postgres:16-alpine
	Image string
	// Username. Defaults to postgres
	Username string
	// Password. Defaults to postgres
	Password string
	// Database name. Defaults to testdb
	Database string
	// HealthCheckTimeout. Defaults to 30s
	HealthCheckTimeout time.Duration
}

// DefaultPostgresConfig returns a config with sensible defaults.
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Image:              "postgres:16-alpine",
		Username:           "postgres",
		Password:           "postgres",
		Database:           "testdb",
		HealthCheckTimeout: 30 * time.Second,
	}
}
