package pgx

import (
	"context"

	"github.com/apex/log"
	"github.com/jackc/pgx/v4"
)

// Logger holds an Apex log instance
type Logger struct {
	logger       log.Interface
	ignoreErrors func(err error) bool
}

// NewLogger builds a new logger instance for pgx given an Apex log instance
func NewLogger(logger log.Interface, opts ...LoggerOpt) *Logger {
	l := &Logger{logger: logger}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Log a pgx log message to the underlying log instance, implements pgx.Logger
func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if data["err"] != nil && l.ignoreErrors != nil && l.ignoreErrors(data["err"].(error)) {
		return
	}

	e := l.logger.
		WithFields(toLogFields(data)).
		WithField("component", "db.driver").
		WithField("level", level.String())

	if data["err"] != nil {
		e = e.WithError(data["err"].(error))
	}

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

func toLogFields(data map[string]interface{}) log.Fields {
	fields := make(map[string]interface{}, len(data))
	for k, v := range data {
		if k == "err" {
			continue
		}
		fields[k] = v
	}
	return fields
}

// LoggerOpt sets options for the logger
type LoggerOpt func(*Logger)

func WithIgnoreErrors(matcher func(err error) bool) LoggerOpt {
	return func(l *Logger) {
		l.ignoreErrors = matcher
	}
}
