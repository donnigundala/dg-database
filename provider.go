package dgdatabase

import (
	"fmt"
	"reflect"

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

// NewDatabaseServiceProvider creates a new database service provider.
func NewDatabaseServiceProvider() *DatabaseServiceProvider {
	return &DatabaseServiceProvider{}
}

// Name returns the provider name
func (p *DatabaseServiceProvider) Name() string {
	return Binding
}

// Version returns the provider version
func (p *DatabaseServiceProvider) Version() string {
	return Version
}

// Dependencies returns the provider dependencies
func (p *DatabaseServiceProvider) Dependencies() []string {
	return []string{}
}

// Register registers the database services into the container.
func (p *DatabaseServiceProvider) Register(app foundation.Application) error {
	// Register database manager
	app.Singleton(Binding, func() (interface{}, error) {
		// Use provided config or default
		cfg := p.Config
		if cfg.Driver == "" {
			cfg = DefaultConfig()
		}

		var loggerInstance Logger
		// Try to resolve logger from container
		if log, err := app.Make("logger"); err == nil {
			// Adapt the logger to our Logger interface
			if adapted, ok := log.(interface {
				Debug(msg string, args ...interface{})
				Info(msg string, args ...interface{})
				Warn(msg string, args ...interface{})
				Error(msg string, args ...interface{})
			}); ok {
				loggerInstance = &loggerAdapter{logger: adapted}
			}
		}

		// Create manager with resolved logger (or nil)
		manager, err := NewManager(cfg, loggerInstance)
		if err != nil {
			return nil, fmt.Errorf("failed to create database manager: %w", err)
		}

		return manager, nil
	})

	// Auto-register named connections in container
	// This allows direct resolution via app.Make("database.analytics")
	for name := range p.Config.Connections {
		connectionName := Binding + "." + name

		// Use a function to properly capture the loop variable
		func(capturedName string) {
			app.Singleton(connectionName, func() (interface{}, error) {
				// Resolve manager to get connection
				managerInstance, err := app.Make(Binding)
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
	// Try to resolve manager and register metrics
	instance, err := app.Make(Binding)
	if err == nil {
		if manager, ok := instance.(*Manager); ok {
			if err := manager.RegisterMetrics(); err != nil {
				// We don't fail boot if metrics fail, just log it
				if log, err := app.Make("logger"); err == nil {
					if l, ok := log.(interface {
						Warn(msg string, args ...interface{})
					}); ok {
						l.Warn("Failed to register database metrics", "error", err)
					}
				}
			}
		}
	}
	return nil
}

// Shutdown gracefully closes database connections.
func (p *DatabaseServiceProvider) Shutdown(app foundation.Application) error {
	dbInstance, err := app.Make(Binding)
	if err != nil {
		return nil // Database not initialized
	}

	manager := dbInstance.(*Manager)
	return manager.Close()
}

// loggerAdapter adapts a generic logger to database.Logger interface.
type loggerAdapter struct {
	logger interface {
		Debug(msg string, args ...interface{})
		Info(msg string, args ...interface{})
		Warn(msg string, args ...interface{})
		Error(msg string, args ...interface{})
	}
}

func (l *loggerAdapter) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *loggerAdapter) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *loggerAdapter) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *loggerAdapter) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *loggerAdapter) With(args ...interface{}) Logger {
	// Try to call With(args...) via reflection to support different return types
	v := reflect.ValueOf(l.logger)
	m := v.MethodByName("With")
	if m.IsValid() {
		valArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			valArgs[i] = reflect.ValueOf(arg)
		}
		results := m.Call(valArgs)
		if len(results) == 1 {
			if nextLogger, ok := results[0].Interface().(interface {
				Debug(msg string, args ...interface{})
				Info(msg string, args ...interface{})
				Warn(msg string, args ...interface{})
				Error(msg string, args ...interface{})
			}); ok {
				return &loggerAdapter{logger: nextLogger}
			}
		}
	}
	return l
}
