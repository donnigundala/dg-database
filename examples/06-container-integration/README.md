# Container Integration Example

This example demonstrates the core container integration patterns for accessing database connections in the DG Framework.

1. **Auto-Registration**: Named connections automatically registered in container
2. **Direct Resolution**: Access via `app.Make("database.analytics")`
3. **Helper Functions**: Use `Resolve()` and `ResolveConnection()`
4. **Injectable Pattern**: Clean dependency injection for services
5. **Phase 6 Service Provider**: Standardized registration pattern

## Running the Example

```bash
go run main.go
```

## Key Integration Patterns

### Registration
The provider automatically registers your primary and named connections. In this example, we set the config manually:

```go
config := dgdatabase.Config{...}
provider := dgdatabase.NewDatabaseServiceProvider()
provider.Config = config
app.Register(provider)
```

## Expected Output

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

### Pattern 1: Direct Resolution
```go
analyticsDB, _ := app.Make("database.analytics")
```
- Simple and direct
- Useful for dynamic resolution

### Pattern 2: Helper Functions
```go
db, _ := dgdatabase.Resolve(app)
analyticsDB, _ := dgdatabase.ResolveConnection(app, "analytics")
```
- Type-safe
- Better error messages
- Recommended for most cases

### Pattern 3: Injectable (Best Practice)
```go
type UserService struct {
    inject *dgdatabase.Injectable
}

func NewUserService(app *foundation.Application) *UserService {
    return &UserService{
        inject: dgdatabase.NewInjectable(app),
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
