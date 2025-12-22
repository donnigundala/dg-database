package dgdatabase

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

// GormLogger adapts our Logger interface to GORM's logger interface.
type GormLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
	Logger        Logger
}

// NewGormLogger creates a new GormLogger.
func NewGormLogger(log Logger, level logger.LogLevel, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		Logger:        log,
		LogLevel:      level,
		SlowThreshold: slowThreshold,
	}
}

// LogMode set log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info log info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.Logger != nil && l.LogLevel >= logger.Info {
		l.Logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn log warn
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.Logger != nil && l.LogLevel >= logger.Warn {
		l.Logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error log error
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.Logger != nil && l.LogLevel >= logger.Error {
		l.Logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace log sql
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Logger == nil || l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error:
		sql, rows := fc()
		l.Logger.Error("GORM Trace Error", "elapsed", elapsed, "rows", rows, "sql", sql, "error", err)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		l.Logger.Warn("GORM Slow Query", "elapsed", elapsed, "rows", rows, "sql", sql)
	case l.LogLevel >= logger.Info:
		sql, rows := fc()
		l.Logger.Info("GORM Trace", "elapsed", elapsed, "rows", rows, "sql", sql)
	}
}
