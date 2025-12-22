package dgdatabase

// logDebug logs a debug message using the logger.
func logDebug(logger Logger, msg string, args ...interface{}) {
	if logger != nil {
		logger.Debug(msg, args...)
	}
}

// logInfo logs an info message using the logger.
func logInfo(logger Logger, msg string, args ...interface{}) {
	if logger != nil {
		logger.Info(msg, args...)
	}
}

// logWarn logs a warning message using the logger.
func logWarn(logger Logger, msg string, args ...interface{}) {
	if logger != nil {
		logger.Warn(msg, args...)
	}
}

// logError logs an error message using the logger.
func logError(logger Logger, msg string, args ...interface{}) {
	if logger != nil {
		logger.Error(msg, args...)
	}
}
