package database

import (
	"fmt"

	"github.com/donnigundala/dg-core/contracts/foundation"
)

// DatabaseServiceProvider implements the PluginProvider interface.
// This provides a simple, plug-and-play integration for applications.
//
// For advanced use cases requiring custom configuration,
// use the library functions (NewManager) directly.
type DatabaseServiceProvider struct {
	// Config holds database configuration
	// Auto-injected by dg-core if using config:"database" tag
	Config Config `config:"database"`
}

// Name returns the provider name
func (p *DatabaseServiceProvider) Name() string {
	return "database"
}

// Version returns the provider version
func (p *DatabaseServiceProvider) Version() string {
	return "1.5.0"
}

// Dependencies returns the provider dependencies
func (p *DatabaseServiceProvider) Dependencies() []string {
	return []string{}
}

// Register registers the database services
func (p *DatabaseServiceProvider) Register(app foundation.Application) error {
	// Use provided config or default
	cfg := p.Config
	if cfg.Driver == "" {
		cfg = DefaultConfig()
	}

	// Register database manager
	app.Singleton("database", func() (interface{}, error) {
		// Create manager without logger
		// Logger will be injected separately if needed
		manager, err := NewManager(cfg, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create database manager: %w", err)
		}

		return manager, nil
	})

	// Auto-register named connections in container
	// This allows direct resolution via app.Make("database.analytics")
	for name := range cfg.Connections {
		connectionName := "database." + name

		// Use a function to properly capture the loop variable
		func(capturedName string) {
			app.Singleton(connectionName, func() (interface{}, error) {
				// Resolve manager to get connection
				managerInstance, err := app.Make("database")
				if err != nil {
					return nil, fmt.Errorf("failed to resolve database manager: %w", err)
				}
				mgr := managerInstance.(*Manager)
				return mgr.Connection(capturedName), nil
			})
		}(name)
	}

	return nil
}

// Boot boots the database services
func (p *DatabaseServiceProvider) Boot(app foundation.Application) error {
	// Database will be resolved when needed
	// No need to verify resolution here to avoid deadlock
	return nil
}

// Shutdown gracefully closes database connections.
func (p *DatabaseServiceProvider) Shutdown(app foundation.Application) error {
	dbInstance, err := app.Make("database")
	if err != nil {
		return nil // Database not initialized
	}

	manager := dbInstance.(*Manager)
	return manager.Close()
}

// loggerAdapter adapts a generic logger to database.Logger interface.
type loggerAdapter struct {
	logger interface {
		Info(msg string, keysAndValues ...interface{})
		Warn(msg string, keysAndValues ...interface{})
	}
}

func (l *loggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *loggerAdapter) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, keysAndValues...)
}
