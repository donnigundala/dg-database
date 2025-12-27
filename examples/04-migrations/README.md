# Database Migrations

This example demonstrates how to use the built-in migration system for versioned database schema changes.

## Features
- Declarative Up/Down migrations.
- Migration status tracking.
- Rollback support.

## Running the Example
```bash
go run main.go
```

## Highlights
```go
migrations := []dgdatabase.Migration{
    {
        ID: "001_initial",
        Up: func(db *gorm.DB) error { return db.AutoMigrate(&User{}) },
        Down: func(db *gorm.DB) error { return db.Migrator().DropTable(&User{}) },
    },
}

// Run migrations
manager.Migrate(migrations)

// Rollback latest
manager.Rollback(migrations)
```
