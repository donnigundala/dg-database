package database

import (
	"fmt"

	"github.com/donnigundala/dg-core/contracts/foundation"
	"gorm.io/gorm"
)

// Resolve resolves the main database manager from the application container.
// Returns the DB interface for flexibility and testability.
//
// Example:
//
//	db, err := database.Resolve(app)
//	if err != nil {
//	    return err
//	}
//	db.DB().Find(&users)
func Resolve(app foundation.Application) (DB, error) {
	instance, err := app.Make("database")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve database: %w", err)
	}
	return instance.(DB), nil
}

// MustResolve resolves the main database manager or panics.
// Use this when you're certain the database is available.
//
// Example:
//
//	db := database.MustResolve(app)
//	db.DB().Find(&users)
func MustResolve(app foundation.Application) DB {
	db, err := Resolve(app)
	if err != nil {
		panic(err)
	}
	return db
}

// ResolveConnection resolves a named database connection from the container.
// Named connections are auto-registered by the provider if defined in config.
//
// Example:
//
//	analyticsDB, err := database.ResolveConnection(app, "analytics")
//	if err != nil {
//	    return err
//	}
//	analyticsDB.Find(&events)
func ResolveConnection(app foundation.Application, name string) (*gorm.DB, error) {
	connectionName := "database." + name
	instance, err := app.Make(connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve connection '%s': %w", name, err)
	}
	return instance.(*gorm.DB), nil
}

// MustResolveConnection resolves a named connection or panics.
// Use this when you're certain the connection exists.
//
// Example:
//
//	analyticsDB := database.MustResolveConnection(app, "analytics")
//	analyticsDB.Find(&events)
func MustResolveConnection(app foundation.Application, name string) *gorm.DB {
	db, err := ResolveConnection(app, name)
	if err != nil {
		panic(err)
	}
	return db
}

// Injectable provides a helper for dependency injection.
// It wraps the application container to provide easy access to database services.
//
// Example:
//
//	type UserService struct {
//	    inject *database.Injectable
//	}
//
//	func NewUserService(app foundation.Application) *UserService {
//	    return &UserService{
//	        inject: database.NewInjectable(app),
//	    }
//	}
//
//	func (s *UserService) GetUsers() {
//	    s.inject.DB().DB().Find(&users)
//	    s.inject.Connection("analytics").Create(&event)
//	}
type Injectable struct {
	app foundation.Application
}

// NewInjectable creates a new injectable helper.
func NewInjectable(app foundation.Application) *Injectable {
	return &Injectable{app: app}
}

// DB returns the main database manager.
func (i *Injectable) DB() DB {
	return MustResolve(i.app)
}

// Connection returns a named database connection.
func (i *Injectable) Connection(name string) *gorm.DB {
	return MustResolveConnection(i.app, name)
}

// TryConnection attempts to resolve a named connection.
// Returns nil if the connection doesn't exist (useful for optional connections).
func (i *Injectable) TryConnection(name string) *gorm.DB {
	db, err := ResolveConnection(i.app, name)
	if err != nil {
		return nil
	}
	return db
}
