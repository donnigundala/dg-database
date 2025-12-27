# Read/Write Splitting

This example demonstrates how to configure and use master/slave database architecture for load balancing and high availability.

## Features
- Manual configuration for Master and multiple Slaves.
- Load balancing strategies (Round-robin, Random, Weighted).
- Automatic routing of reads to slaves and writes to master.
- Manual control to force master/slave usage.
- Integrated health checks for all nodes.

## Running the Example
```bash
go run main.go
```

## Highlights
```go
// Master-Slave configuration
config := dgdatabase.Config{
    ReadWriteSplitting: true,
    AutoRouting:        true,
    Master:             masterConfig,
    Slaves:             []dgdatabase.ConnectionConfig{slave1, slave2},
}

// Write to master (forced)
manager.Write().Create(&product)

// Read from slave (load-balanced)
manager.Read().Find(&products)
```
