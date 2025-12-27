# Basic Database Operations

This example demonstrates the foundational CRUD (Create, Read, Update, Delete) operations using a single database connection.

## Features
- Manual configuration using `DefaultConfig()` and fluent setters.
- Basic CRUD using GORM's `*gorm.DB` instance.
- Transaction handling using `manager.WithTx()`.

## Running the Example
```bash
go run main.go
```

## Highlights
```go
// Create manager
manager, _ := dgdatabase.NewManager(config, nil)

// Get GORM instance
db := manager.DB()

// Basic CRUD
db.Create(&user)
db.Find(&users)
db.Model(&user).Update("email", "new@email.com")
db.Delete(&user)
```
