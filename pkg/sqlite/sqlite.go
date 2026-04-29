// Package sqlite provides SQLite helpers for testing without containers.
package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// TempDB creates a temporary SQLite database file.
// The file is automatically removed when the database is closed.
func TempDB() (*sql.DB, func(), error) {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("testdb-%d.sqlite", os.Getpid()))
	
	db, err := sql.Open("sqlite3", tmpFile+"?mode=memory&cache=shared")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpFile)
	}

	return db, cleanup, nil
}

// InMemory creates an in-memory SQLite database.
// Useful for fast tests that don't need persistence.
func InMemory() (*sql.DB, func(), error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, nil
}
