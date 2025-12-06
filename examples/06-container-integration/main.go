package main

import (
	"fmt"
	"log"

	"github.com/donnigundala/dg-core/foundation"
	database "github.com/donnigundala/dg-database"
)

// User model
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"size:100;uniqueIndex"`
}

// AnalyticsEvent model
type AnalyticsEvent struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"index"`
	Action string `gorm:"size:50"`
}

// AuditLog model
type AuditLog struct {
	ID      uint   `gorm:"primaryKey"`
	UserID  uint   `gorm:"index"`
	Message string `gorm:"size:255"`
}

func main() {
	fmt.Println("=== Container Integration Example ===\n")

	// 1. Setup application
	app := foundation.New(".")

	// 2. Configure database with named connections
	config := database.Config{
		Driver:   "sqlite",
		FilePath: "./primary.db",

		// Named connections - automatically registered in container!
		Connections: map[string]database.ConnectionConfig{
			"analytics": {
				Driver:   "sqlite",
				FilePath: "./analytics.db",
			},
			"audit": {
				Driver:   "sqlite",
				FilePath: "./audit.db",
			},
		},
	}

	// 3. Register provider
	provider := &database.DatabaseServiceProvider{Config: config}
	if err := app.Register(provider); err != nil {
		log.Fatal(err)
	}

	if err := app.Boot(); err != nil {
		log.Fatal(err)
	}

	// 4. Auto-migrate models
	db, _ := database.Resolve(app)
	db.DB().AutoMigrate(&User{})

	analyticsDB, _ := database.ResolveConnection(app, "analytics")
	analyticsDB.AutoMigrate(&AnalyticsEvent{})

	auditDB, _ := database.ResolveConnection(app, "audit")
	auditDB.AutoMigrate(&AuditLog{})

	fmt.Println("✅ Databases initialized and migrated\n")

	// ========================================
	// Pattern 1: Direct Resolution from Container
	// ========================================
	fmt.Println("--- Pattern 1: Direct Resolution ---")
	directAnalytics, err := app.Make("database.analytics")
	if err != nil {
		log.Fatal(err)
	}
	analyticsConn := directAnalytics.(*database.Manager).Connection("analytics")

	event := &AnalyticsEvent{
		UserID: 1,
		Action: "login",
	}
	analyticsConn.Create(event)
	fmt.Printf("✅ Created analytics event via direct resolution: %+v\n\n", event)

	// ========================================
	// Pattern 2: Helper Functions
	// ========================================
	fmt.Println("--- Pattern 2: Helper Functions ---")

	// Main database
	mainDB, err := database.Resolve(app)
	if err != nil {
		log.Fatal(err)
	}

	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mainDB.DB().Create(user)
	fmt.Printf("✅ Created user via Resolve(): %+v\n", user)

	// Named connection
	auditConnection, err := database.ResolveConnection(app, "audit")
	if err != nil {
		log.Fatal(err)
	}

	auditLog := &AuditLog{
		UserID:  user.ID,
		Message: "User created",
	}
	auditConnection.Create(auditLog)
	fmt.Printf("✅ Created audit log via ResolveConnection(): %+v\n\n", auditLog)

	// ========================================
	// Pattern 3: Injectable (Best for Services)
	// ========================================
	fmt.Println("--- Pattern 3: Injectable Pattern ---")

	userService := NewUserService(app)
	createdUser := userService.CreateUser("Jane Smith", "jane@example.com")
	fmt.Printf("✅ User created via service: %+v\n", createdUser)

	// ========================================
	// Verify Data
	// ========================================
	fmt.Println("\n--- Verification ---")

	var users []User
	mainDB.DB().Find(&users)
	fmt.Printf("Total users in primary DB: %d\n", len(users))

	var events []AnalyticsEvent
	analyticsConnection, _ := database.ResolveConnection(app, "analytics")
	analyticsConnection.Find(&events)
	fmt.Printf("Total events in analytics DB: %d\n", len(events))

	var logs []AuditLog
	auditConnection.Find(&logs)
	fmt.Printf("Total logs in audit DB: %d\n", len(logs))

	fmt.Println("\n✅ Example completed successfully!")
}

// ========================================
// UserService - Demonstrates Injectable Pattern
// ========================================

type UserService struct {
	inject *database.Injectable
}

func NewUserService(app *foundation.Application) *UserService {
	return &UserService{
		inject: database.NewInjectable(app),
	}
}

func (s *UserService) CreateUser(name, email string) *User {
	user := &User{
		Name:  name,
		Email: email,
	}

	// Use main database
	s.inject.DB().DB().Create(user)

	// Track in analytics
	s.inject.Connection("analytics").Create(&AnalyticsEvent{
		UserID: user.ID,
		Action: "user_created",
	})

	// Log in audit
	s.inject.Connection("audit").Create(&AuditLog{
		UserID:  user.ID,
		Message: fmt.Sprintf("User %s created", name),
	})

	return user
}
