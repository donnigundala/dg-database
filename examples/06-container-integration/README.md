# Container Integration Example

This example demonstrates the three patterns for accessing database connections using the v1.5.0 container integration features.

## Features Demonstrated

1. **Auto-Registration**: Named connections automatically registered in container
2. **Direct Resolution**: Access via `app.Make("database.analytics")`
3. **Helper Functions**: Use `Resolve()` and `ResolveConnection()`
4. **Injectable Pattern**: Clean dependency injection for services

## Running the Example

```bash
cd examples/06-container-integration
go run main.go
```

## What It Does

1. Creates three SQLite databases (primary, analytics, audit)
2. Demonstrates all three access patterns
3. Shows how to use Injectable in services
4. Verifies data across all databases

## Output

```
=== Container Integration Example ===

✅ Databases initialized and migrated

--- Pattern 1: Direct Resolution ---
✅ Created analytics event via direct resolution

--- Pattern 2: Helper Functions ---
✅ Created user via Resolve()
✅ Created audit log via ResolveConnection()

--- Pattern 3: Injectable Pattern ---
✅ User created via service

--- Verification ---
Total users in primary DB: 2
Total events in analytics DB: 2
Total logs in audit DB: 2

✅ Example completed successfully!
```

## Key Takeaways

### Pattern 1: Direct Resolution
```go
analyticsDB, _ := app.Make("database.analytics")
```
- Simple and direct
- Good for one-off usage

### Pattern 2: Helper Functions
```go
db, _ := database.Resolve(app)
analyticsDB, _ := database.ResolveConnection(app, "analytics")
```
- Type-safe
- Better error messages
- Recommended for most cases

### Pattern 3: Injectable (Best Practice)
```go
type UserService struct {
    inject *database.Injectable
}

func NewUserService(app foundation.Application) *UserService {
    return &UserService{
        inject: database.NewInjectable(app),
    }
}
```
- Clean dependency injection
- Easy to test
- Best for services and controllers
- Recommended for production code

## Cleanup

```bash
rm *.db
```
