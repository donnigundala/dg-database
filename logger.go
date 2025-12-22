package dgdatabase

// Logger defines the interface for database logging.
// Any logger implementation must provide these methods to be used with the database manager.
//
// This interface is designed to be compatible with dg-core's logging.Logger.
type Logger interface {
	// Debug logs a debug message with optional key-value pairs
	Debug(msg string, args ...interface{})

	// Info logs an informational message with optional key-value pairs
	Info(msg string, args ...interface{})

	// Warn logs a warning message with optional key-value pairs
	Warn(msg string, args ...interface{})

	// Error logs an error message with optional key-value pairs
	Error(msg string, args ...interface{})

	// With returns a new logger with the given fields added
	With(args ...interface{}) Logger
}
