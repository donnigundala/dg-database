# PostgreSQL Schema Support

This example demonstrates how to use PostgreSQL schemas for multi-tenancy or logical database separation.

## Features
- Per-connection schema configuration.
- Support for multi-tenant schema isolation.
- Integration with read/write splitting.

## Running the Example
```bash
go run main.go
```

## Highlights
```go
// Per-tenant schema configuration
config := dgdatabase.ConnectionConfig{
    Driver: "postgres",
    Name:   "myapp",
    Schema: "tenant_1",
}
```
