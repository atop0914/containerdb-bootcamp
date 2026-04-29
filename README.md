# ContainerDB

A lightweight containerized database toolkit for Go development and testing. Spin up real databases in containers with a single function call — no Docker Compose required.

## Features

- **MySQL, PostgreSQL, SQLite** support out of the box
- **One-line setup**: `db, cleanup, err := mysql.New(ctx)` — that's it
- **Auto-cleanup**: Containers are automatically cleaned up when the test ends
- **Random port allocation**: No port conflicts between parallel tests
- **Configurable**: Custom image tags, credentials, volumes, health checks
- **Zero external dependencies**: Pure Go, only depends on `testcontainers-go`

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/atop0914/containerdb-bootcamp/pkg/mysql"
)

func main() {
    ctx := context.Background()
    
    // Spin up a MySQL container
    db, cleanup, err := mysql.New(ctx)
    if err != nil {
        panic(err)
    }
    defer cleanup() // Always clean up
    
    // Use db.Pool as *sql.DB
    var version string
    db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)
    fmt.Println("MySQL version:", version)
}
```

## Architecture

```
pkg/
├── mysql/       # MySQL container wrapper
├── postgres/    # PostgreSQL container wrapper  
└── sqlite/      # SQLite helper (no container needed)
cmd/
└── containerdb/ # CLI tool
internal/
├── container/   # Base container management
├── database/    # Common database interfaces
└── config/      # Configuration types
```

## Motivation

Testcontainers (JVM) is too heavy for Go projects. This library provides the same experience with a native Go API.

## License

MIT
# ContainerDB
