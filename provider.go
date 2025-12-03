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
	return "dg-database"
}

// Version returns the provider version
func (p *DatabaseServiceProvider) Version() string {
	return "1.4.0"
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
		// Try to resolve logger (optional)
		var logger Logger
		if loggerInstance, err := app.Make("logger"); err == nil {
			// Adapt dg-core logger to database.Logger interface
			if l, ok := loggerInstance.(interface {
				Info(msg string, keysAndValues ...interface{})
				Warn(msg string, keysAndValues ...interface{})
			}); ok {
				logger = &loggerAdapter{logger: l}
			}
		}

		manager, err := NewManager(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create database manager: %w", err)
		}
		return manager, nil
	})

	return nil
}

// Boot boots the database services
func (p *DatabaseServiceProvider) Boot(app foundation.Application) error {
	// Get manager
	managerInstance, err := app.Make("database")
	if err != nil {
		return fmt.Errorf("failed to get database manager: %w", err)
	}
	manager := managerInstance.(*Manager)

	// Test connection
	if err := manager.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	// Log success if logger available
	if manager.logger != nil {
		manager.logInfo("Database connected successfully",
			"driver", p.Config.Driver,
			"database", p.Config.Database)

		if p.Config.ReadWriteSplitting {
			manager.logInfo("Read/write splitting enabled",
				"slaves", len(p.Config.Slaves),
				"strategy", p.Config.SlaveStrategy)
		}

		if len(p.Config.Connections) > 0 {
			manager.logInfo("Named connections established",
				"count", len(p.Config.Connections))
		}
	}

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
