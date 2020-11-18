package pgx

import (
	"context"

	"github.com/apex/log"
	"github.com/jackc/pgx/v4"
)

// Logger holds an Apex log instance
type Logger struct {
	logger log.Interface
}

// NewLogger builds a new logger instance for pgx given an Apex log instance
func NewLogger(logger log.Interface) *Logger {
	return &Logger{logger: logger}
}

// Log a pgx log message to the underlying log instance, implements pgx.Logger
func (p *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	e := p.logger.
		WithFields(log.Fields(data)).
		WithField("component", "db.driver").
		WithField("level", level.String())

	switch level {
	case pgx.LogLevelTrace:
		e.Debug(msg)
	case pgx.LogLevelDebug:
		e.Debug(msg)
	case pgx.LogLevelInfo:
		// All queries are logged with level Info, so consider this a Debug output
		e.Debug(msg)
	case pgx.LogLevelWarn:
		e.Warn(msg)
	case pgx.LogLevelError:
		e.Error(msg)
	default:
		e.WithField("invalidPgxLogLevel", level).Error(msg)
	}
}
