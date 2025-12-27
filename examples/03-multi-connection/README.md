# Multi-Connection Management

This example demonstrates how to manage multiple distinct database connections simultaneously, such as separate databases for "Primary", "Analytics", and "Logs".

## Features
- Named connections configuration.
- Runtime connection management (Add/Remove).
- Global health checks across all configured connections.

## Running the Example
```bash
go run main.go
```

## Highlights
```go
// Named connections
config := dgdatabase.Config{
    Connections: map[string]dgdatabase.ConnectionConfig{
        "analytics": analyticsConfig,
    },
}

// Access specific connection
analyticsDB := manager.Connection("analytics")

// Add connection at runtime
manager.AddConnection("storage", storageConfig)
```
