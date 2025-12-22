package dgdatabase

import (
	"context"

	"gorm.io/gorm"
)

// DB is the main database interface that defines the contract for database operations.
// This interface allows for easier testing and flexibility in implementation.
type DB interface {
	// Core database operations
	DB() *gorm.DB
	Close() error
	Ping() error

	// Multi-connection support
	Connection(name string) *gorm.DB
	HasConnection(name string) bool
	AddConnection(name string, config ConnectionConfig) error
	RemoveConnection(name string) error

	// Read/Write splitting
	Master() *gorm.DB
	Read() *gorm.DB
	Write() *gorm.DB
	Slave(index int) *gorm.DB

	// Transaction support
	WithTx(fn TransactionFunc) error
	WithTxContext(ctx context.Context, fn TransactionFunc) error
	TX() *TransactionHelper
	Transaction(fn func(*gorm.DB) error) error

	// Health monitoring
	HealthCheck() map[string]bool
	DetailedHealthCheck() map[string]ConnectionHealth
	IsHealthy() bool
	IsFullyHealthy() bool

	// Connection pool statistics
	Stats() PoolStats
	AllStats() map[string]PoolStats
	ConnectionStats(name string) PoolStats

	// Migration support
	AutoMigrate(models ...interface{}) error
	Migrate(migrations []Migration) error
	Rollback(migrations []Migration) error
	MigrationStatus() ([]string, error)
}

// Ensure Manager implements DB interface at compile time
var _ DB = (*Manager)(nil)
